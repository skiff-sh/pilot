package pilot

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/skiff-sh/config/ptr"
	pilot "github.com/skiff-sh/pilot/api/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

type Client interface {
	NewBehavior() NameBehaviorBuilder
	Provoke(ctx context.Context, name string) (*structpb.Struct, error)
}

func New(cl pilot.PilotServiceClient) Client {
	return &client{
		Cl: cl,
	}
}

type client struct {
	Cl pilot.PilotServiceClient
}

func (c *client) NewBehavior() NameBehaviorBuilder {
	return &createBehaviorBuilder{
		Req: &pilot.Behavior{},
		Cl:  c.Cl,
	}
}

func (c *client) Provoke(ctx context.Context, name string) (*structpb.Struct, error) {
	out, err := c.Cl.ProvokeBehavior(ctx, &pilot.ProvokeBehavior_Request{Name: name})
	if err != nil {
		return nil, err
	}

	return out.Body, nil
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

type TendencyAdder interface {
	Add() CreateBehaviorSender
}

type TendencyActionBuilder interface {
	Wait(d time.Duration) TendencyAdder
	HTTPRequest(url string, o ...HTTPRequestOpt) TendencyAdder
	SetStatus(code codes.Code, msg string) TendencyAdder
	SetResponseField(from, to string) TendencyAdder
	Exec(cmd string, o ...ExecOpt) TendencyAdder
}

type HTTPRequestOpt func(r *pilot.Action_HTTPRequest)

type ExecOpt func(r *pilot.Action_Exec)

func WithHTTPHeaders(m map[string]string) HTTPRequestOpt {
	return func(r *pilot.Action_HTTPRequest) {
		r.Headers = m
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

var _ NameBehaviorBuilder = &createBehaviorBuilder{}
var _ TendencyBuilder = &createBehaviorBuilder{}
var _ CreateBehaviorSender = &createBehaviorBuilder{}

type createBehaviorBuilder struct {
	Req *pilot.Behavior
	Cl  pilot.PilotServiceClient
}

func (c *createBehaviorBuilder) Send(ctx context.Context) (*pilot.CreateBehavior_Response, error) {
	return c.Cl.CreateBehavior(ctx, &pilot.CreateBehavior_Request{Behavior: c.Req})
}

func (c *createBehaviorBuilder) Tendency() TendencyFieldBuilder {
	return &tendencyBuilder{Parent: c, Tendency: &pilot.Tendency{}}
}

func (c *createBehaviorBuilder) Name(n string) TendencyBuilder {
	c.Req.Name = n
	return c
}

var _ TendencyActionBuilder = &tendencyBuilder{}
var _ TendencyAdder = &tendencyBuilder{}
var _ TendencyFieldBuilder = &tendencyBuilder{}

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

func (t *tendencyBuilder) Wait(d time.Duration) TendencyAdder {
	t.Tendency.Action.Wait = durationpb.New(d)
	return t
}

func (t *tendencyBuilder) HTTPRequest(url string, o ...HTTPRequestOpt) TendencyAdder {
	t.Tendency.Action.HttpRequest = &pilot.Action_HTTPRequest{
		Url: url,
	}
	for _, v := range o {
		v(t.Tendency.Action.HttpRequest)
	}
	return t
}

func (t *tendencyBuilder) SetStatus(code codes.Code, msg string) TendencyAdder {
	t.Tendency.Action.SetStatus = &pilot.Action_SetStatus{
		Code:    uint32(code),
		Message: msg,
	}
	return t
}

func (t *tendencyBuilder) SetResponseField(from, to string) TendencyAdder {
	t.Tendency.Action.SetResponseField = &pilot.Action_SetResponseField{
		From: from,
		To:   to,
	}
	return t
}

func (t *tendencyBuilder) Exec(cmd string, o ...ExecOpt) TendencyAdder {
	t.Tendency.Action.Exec = &pilot.Action_Exec{
		Command: cmd,
	}
	for _, v := range o {
		v(t.Tendency.Action.Exec)
	}
	return t
}
