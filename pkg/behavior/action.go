package behavior

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"github.com/skiff-sh/pilot/pkg/httptype"
	"github.com/skiff-sh/pilot/pkg/template"

	"github.com/goccy/go-json"
	"github.com/skiff-sh/config"
	pilot "github.com/skiff-sh/pilot/api/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CompileAction(id string, b *pilot.Action) (behaviortype.Action, error) {
	switch {
	case b.Exec != nil:
		out := &Exec{
			ID:     id,
			Stderr: bytes.NewBuffer(nil),
			Stdout: bytes.NewBuffer(nil),
		}
		//nolint: gosec // can add support to disable this action.
		out.Cmd = exec.Command(b.Exec.Command, b.Exec.Args...)
		if out.Cmd.Err != nil {
			return nil, out.Cmd.Err
		}

		out.Cmd.Stderr, out.Cmd.Stdout = out.Stderr, out.Stdout
		out.Cmd.Dir = b.Exec.WorkingDir
		out.Cmd.Env = config.NewMap(b.Exec.EnvVars).ToEnv()
		return out, nil
	case b.HttpRequest != nil:
		var body io.Reader
		if len(b.HttpRequest.Body) > 0 {
			body = bytes.NewReader(b.HttpRequest.Body)
		}
		req, err := http.NewRequest(b.HttpRequest.Method, b.HttpRequest.Url, body)
		if err != nil {
			return nil, err
		}
		for k, v := range b.HttpRequest.Headers {
			req.Header.Set(k, v)
		}

		out := &HTTPRequest{
			ID:     id,
			Spec:   b.HttpRequest,
			Req:    req,
			Client: http.DefaultClient,
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
	Req    *http.Request
	Client httptype.HttpDoer
}

func (h *HTTPRequest) GetID() string {
	return h.ID
}

func (h *HTTPRequest) Act(_ *behaviortype.Context) (behaviortype.Output, error) {
	resp, err := h.Client.Do(h.Req)
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
	ID     string
	Cmd    *exec.Cmd
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

func (e *Exec) GetID() string {
	return e.ID
}

func (e *Exec) Act(_ *behaviortype.Context) (behaviortype.Output, error) {
	err := e.Cmd.Run()
	out := &pilot.Output_ExecOutput{
		Stdout: e.Stdout.String(),
		Stderr: e.Stderr.String(),
	}
	var v *exec.ExitError
	if errors.As(err, &v) {
		out.ExitCode = int32(v.ExitCode())
	}

	if err != nil {
		err = errors.Join(err, errors.New(e.Stderr.String()))
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
