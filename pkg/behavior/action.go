package behavior

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/skiff-sh/pilot/api/go/pilot"

	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"github.com/skiff-sh/pilot/pkg/httptype"
	"github.com/skiff-sh/pilot/pkg/template"

	"github.com/goccy/go-json"
	"github.com/skiff-sh/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CompileAction(id string, b *pilot.Action) (behaviortype.Action, error) {
	switch {
	case b.Exec != nil:
		// validate that the command is valid.
		_, _, _, err := newCmd(context.TODO(), b.Exec)
		if err != nil {
			return nil, err
		}
		out := &Exec{
			ID:   id,
			Spec: b.Exec,
		}
		return out, nil
	case b.HttpRequest != nil:
		// validate that the request will be valid.
		_, err := newHTTPRequest(context.TODO(), b.HttpRequest)
		if err != nil {
			return nil, err
		}

		cl := *http.DefaultClient
		cl.Timeout = 5 * time.Second

		out := &HTTPRequest{
			ID:     id,
			Spec:   b.HttpRequest,
			Client: &cl,
		}
		return out, nil
	case b.SetStatus != nil:
		return &SetStatus{
			Spec: b.SetStatus,
		}, nil
	case b.Wait != nil:
		return &Wait{
			Dur: b.Wait.AsDuration(),
		}, nil
	case b.SetResponseField != nil:
		fe, err := template.NewFieldTemplates(b.SetResponseField, template.WithForce())
		if err != nil {
			return nil, err
		}
		return &SetResponseField{
			Expressions: fe,
			Spec:        b.SetResponseField,
		}, nil
	default:
		return nil, errors.New("unknown or missing action")
	}
}

var (
	_ behaviortype.Action      = &HTTPRequest{}
	_ behaviortype.Referential = &HTTPRequest{}
)

type HTTPRequest struct {
	ID     string
	Spec   *pilot.Action_HTTPRequest
	Client httptype.HttpDoer
}

func (h *HTTPRequest) GetID() string {
	return h.ID
}

func (h *HTTPRequest) Act(c *behaviortype.Context) (behaviortype.Output, error) {
	req, err := newHTTPRequest(c, h.Spec)
	if err != nil {
		return nil, err
	}
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	out := &pilot.Output_HTTPResponse{
		Status:        int32(resp.StatusCode),
		Proto:         resp.Proto,
		ProtoMajor:    int32(resp.ProtoMajor),
		ProtoMinor:    int32(resp.ProtoMinor),
		Headers:       make(map[string]string),
		ContentLength: resp.ContentLength,
	}
	for k, v := range resp.Header {
		if len(v) == 0 {
			continue
		}
		out.Headers[k] = v[0]
	}

	out.BodyRaw, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if out.ContentLength < 0 {
		out.ContentLength = int64(len(out.BodyRaw))
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		val := template.Data{}
		err = json.Unmarshal(out.BodyRaw, &val)
		if err != nil {
			return nil, err
		}
		out.Body = val.ToProto()
	}

	return &HTTPResponseOutput{out}, nil
}

var _ behaviortype.Action = &SetStatus{}

type SetStatus struct {
	Spec *pilot.Action_SetStatus
}

func (s *SetStatus) Act(c *behaviortype.Context) (behaviortype.Output, error) {
	c.Response.Status = status.New(codes.Code(s.Spec.Code), s.Spec.Message)
	return nil, nil
}

var _ behaviortype.Action = &SetResponseField{}

type SetResponseField struct {
	Expressions template.FieldExpressions
	Spec        *pilot.Action_SetResponseField
}

func (s *SetResponseField) Act(c *behaviortype.Context) (behaviortype.Output, error) {
	err := s.Expressions.Apply(s.Spec, c.Outputs)
	if err != nil {
		return nil, err
	}

	val, _ := c.Outputs.Get(s.Spec.From)
	c.Response.Body.Set(s.Spec.To, val)
	return nil, nil
}

var (
	_ behaviortype.Action      = &Exec{}
	_ behaviortype.Referential = &Exec{}
)

type Exec struct {
	ID   string
	Spec *pilot.Action_Exec
}

func (e *Exec) GetID() string {
	return e.ID
}

func (e *Exec) Act(c *behaviortype.Context) (behaviortype.Output, error) {
	cmd, stdout, stderr, err := newCmd(c, e.Spec)
	if err != nil {
		return nil, err
	}

	err = cmd.Run()
	out := &pilot.Output_ExecOutput{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	var v *exec.ExitError
	if errors.As(err, &v) {
		out.ExitCode = int32(v.ExitCode())
	}

	if err != nil {
		err = errors.Join(err, errors.New(stderr.String()))
	}

	return &ExecOutput{Output_ExecOutput: out}, err
}

var _ behaviortype.Action = &Wait{}

type Wait struct {
	Dur time.Duration
}

// Used for testing.
var sleeperFunc = time.Sleep

func (w *Wait) Act(_ *behaviortype.Context) (behaviortype.Output, error) {
	sleeperFunc(w.Dur)
	return nil, nil
}

func newHTTPRequest(ctx context.Context, spec *pilot.Action_HTTPRequest) (*http.Request, error) {
	var body io.Reader
	if len(spec.Body) > 0 {
		body = bytes.NewReader(spec.Body)
		if spec.Headers == nil {
			spec.Headers = make(map[string]string)
		}
		_, ok := spec.Headers["Content-Type"]
		if !ok {
			var contentType string
			if v := bytes.TrimSpace(spec.Body); len(v) > 0 && v[0] == '{' {
				contentType = "application/json"
			} else {
				contentType = http.DetectContentType(spec.Body)
			}
			spec.Headers["Content-Type"] = contentType
		}
	}

	req, err := http.NewRequestWithContext(ctx, spec.Method, spec.Url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range spec.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func newCmd(ctx context.Context, spec *pilot.Action_Exec) (cmd *exec.Cmd, stdout, stderr *bytes.Buffer, err error) {
	stdout, stderr = bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	//nolint: gosec // can add support to disable this action.
	cmd = exec.CommandContext(ctx, spec.Command, spec.Args...)
	if cmd.Err != nil {
		return nil, nil, nil, cmd.Err
	}

	cmd.Stderr, cmd.Stdout = stderr, stdout
	cmd.Dir = spec.WorkingDir
	cmd.Env = config.NewMap(spec.EnvVars).ToEnv()
	return cmd, stdout, stderr, nil
}
