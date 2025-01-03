package pilot

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/skiff-sh/pilot/pkg/httptype"
	"github.com/skiff-sh/pilot/server/pkg/protoenc"

	"github.com/skiff-sh/pilot/api/go/pilot"

	"github.com/goccy/go-json"
	"github.com/skiff-sh/config/ptr"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
)

type Client interface {
	NewBehavior() NameBehaviorBuilder
	GRPC() Provoker
	HTTP() Provoker
}

type Provoker interface {
	Provoke(ctx context.Context, name string) (*structpb.Struct, error)
}

func New(cl pilot.PilotServiceClient, httpProv Provoker) Client {
	return &client{
		HTTPProvoker: httpProv,
		Cl:           cl,
	}
}

type client struct {
	HTTPProvoker Provoker
	Cl           pilot.PilotServiceClient
}

func (c *client) HTTP() Provoker {
	return c.HTTPProvoker
}

func (c *client) GRPC() Provoker {
	return NewGRPC(c.Cl)
}

func (c *client) NewBehavior() NameBehaviorBuilder {
	return &createBehaviorBuilder{
		Req: &pilot.Behavior{},
		Cl:  c,
	}
}

type NameBehaviorBuilder interface {
	Name(n string) TendencyBuilder
}

type TendencyBuilder interface {
	// Tendency adds a tendency to the behavior.
	Tendency() TendencyFieldBuilder
}

type TendencyFieldBuilder interface {
	Action() TendencyActionBuilder
	Condition(cond string) TendencyFieldBuilder
	ID(id string) TendencyFieldBuilder
}

type CreateBehaviorSender interface {
	TendencyBuilder
	Send(ctx context.Context) (*pilot.CreateBehavior_Response, error)
}

type TendencyActionBuilder interface {
	Wait(d time.Duration) CreateBehaviorSender
	HTTPRequest(url string, o ...HTTPRequestOpt) CreateBehaviorSender
	SetStatus(code codes.Code, msg string) CreateBehaviorSender
	SetResponseField(from, to string) CreateBehaviorSender
	Exec(cmd string, o ...ExecOpt) CreateBehaviorSender
}

type HTTPRequestOpt func(r *pilot.Action_HTTPRequest)

type ExecOpt func(r *pilot.Action_Exec)

func WithHTTPHeaders(m map[string]string) HTTPRequestOpt {
	return func(r *pilot.Action_HTTPRequest) {
		r.Headers = m
	}
}

func WithHTTPHeader(key, val string) HTTPRequestOpt {
	return func(r *pilot.Action_HTTPRequest) {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		r.Headers[key] = val
	}
}

func WithHTTPMethod(m string) HTTPRequestOpt {
	return func(r *pilot.Action_HTTPRequest) {
		r.Method = m
	}
}

func WithHTTPBodyRaw(b []byte) HTTPRequestOpt {
	return func(r *pilot.Action_HTTPRequest) {
		r.Body = b
	}
}

func WithHTTPJSONBody(a any) HTTPRequestOpt {
	b, _ := json.Marshal(a)
	return WithHTTPBodyRaw(b)
}

func WithExecArgs(a ...string) ExecOpt {
	return func(r *pilot.Action_Exec) {
		r.Args = append(r.Args, a...)
	}
}

func WithEnvVars(m map[string]string) ExecOpt {
	return func(r *pilot.Action_Exec) {
		r.EnvVars = m
	}
}

func WithExecDir(d string) ExecOpt {
	return func(r *pilot.Action_Exec) {
		r.WorkingDir = d
	}
}

var (
	_ NameBehaviorBuilder  = &createBehaviorBuilder{}
	_ TendencyBuilder      = &createBehaviorBuilder{}
	_ CreateBehaviorSender = &createBehaviorBuilder{}
)

type createBehaviorBuilder struct {
	Req *pilot.Behavior
	Cl  *client
}

func (c *createBehaviorBuilder) Send(ctx context.Context) (*pilot.CreateBehavior_Response, error) {
	return c.Cl.Cl.CreateBehavior(ctx, &pilot.CreateBehavior_Request{Behavior: c.Req})
}

func (c *createBehaviorBuilder) Tendency() TendencyFieldBuilder {
	return &tendencyBuilder{Parent: c, Tendency: &pilot.Tendency{}}
}

func (c *createBehaviorBuilder) Name(n string) TendencyBuilder {
	c.Req.Name = n
	return c
}

var (
	_ TendencyActionBuilder = &tendencyBuilder{}
	_ TendencyFieldBuilder  = &tendencyBuilder{}
)

type tendencyBuilder struct {
	Parent   *createBehaviorBuilder
	Tendency *pilot.Tendency
}

func (t *tendencyBuilder) Action() TendencyActionBuilder {
	t.Tendency.Action = &pilot.Action{}
	return t
}

func (t *tendencyBuilder) Condition(cond string) TendencyFieldBuilder {
	t.Tendency.If = ptr.Ptr(cond)
	return t
}

func (t *tendencyBuilder) ID(id string) TendencyFieldBuilder {
	t.Tendency.Id = ptr.Ptr(id)
	return t
}

func (t *tendencyBuilder) Add() CreateBehaviorSender {
	t.Parent.Req.Tendencies = append(t.Parent.Req.Tendencies, t.Tendency)
	return t.Parent
}

func (t *tendencyBuilder) Wait(d time.Duration) CreateBehaviorSender {
	t.Tendency.Action.Wait = durationpb.New(d)
	t.Parent.Req.Tendencies = append(t.Parent.Req.Tendencies, t.Tendency)
	return t.Parent
}

func (t *tendencyBuilder) HTTPRequest(url string, o ...HTTPRequestOpt) CreateBehaviorSender {
	t.Tendency.Action.HttpRequest = &pilot.Action_HTTPRequest{
		Url:    url,
		Method: http.MethodGet,
	}
	for _, v := range o {
		v(t.Tendency.Action.HttpRequest)
	}
	t.Parent.Req.Tendencies = append(t.Parent.Req.Tendencies, t.Tendency)
	return t.Parent
}

func (t *tendencyBuilder) SetStatus(code codes.Code, msg string) CreateBehaviorSender {
	t.Tendency.Action.SetStatus = &pilot.Action_SetStatus{
		Code:    uint32(code),
		Message: msg,
	}
	t.Parent.Req.Tendencies = append(t.Parent.Req.Tendencies, t.Tendency)
	return t.Parent
}

func (t *tendencyBuilder) SetResponseField(from, to string) CreateBehaviorSender {
	t.Tendency.Action.SetResponseField = &pilot.Action_SetResponseField{
		From: from,
		To:   to,
	}
	t.Parent.Req.Tendencies = append(t.Parent.Req.Tendencies, t.Tendency)
	return t.Parent
}

func (t *tendencyBuilder) Exec(cmd string, o ...ExecOpt) CreateBehaviorSender {
	t.Tendency.Action.Exec = &pilot.Action_Exec{
		Command: cmd,
	}
	for _, v := range o {
		v(t.Tendency.Action.Exec)
	}
	t.Parent.Req.Tendencies = append(t.Parent.Req.Tendencies, t.Tendency)
	return t.Parent
}

func NewGRPC(cl pilot.PilotServiceClient) Provoker {
	out := &grpcProvoker{
		Cl: cl,
	}

	return out
}

type grpcProvoker struct {
	Cl pilot.PilotServiceClient
}

func (g *grpcProvoker) Provoke(ctx context.Context, name string) (*structpb.Struct, error) {
	resp, err := g.Cl.ProvokeBehavior(ctx, &pilot.ProvokeBehavior_Request{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetBody(), nil
}

func NewHTTP(addr string, doer httptype.HttpDoer) Provoker {
	out := &httpProvoker{
		Cl:          doer,
		ProvokerURL: addr + "/api/v1/provoke",
	}
	return out
}

type httpProvoker struct {
	Cl          httptype.HttpDoer
	ProvokerURL string
}

func (h *httpProvoker) Provoke(ctx context.Context, name string) (*structpb.Struct, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.ProvokerURL+"/"+name, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.Cl.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	payload := new(pilot.ProvokeBehavior_Response)
	err = protoenc.ProtoUnmarshaller.Unmarshal(b, payload)
	if err != nil {
		return nil, err
	}

	return payload.GetBody(), nil
}
