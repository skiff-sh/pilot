package behaviortype

import (
	"context"
	"github.com/skiff-sh/pilot/api/go/pilot"

	"github.com/skiff-sh/pilot/pkg/template"
	"google.golang.org/grpc/status"
)

type Action interface {
	Act(c *Context) (out Output, err error)
}

type Interface interface {
	Provoke(ctx context.Context) (*Response, error)
	GetName() string
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
		Outputs: template.Data{},
		Response: &Response{
			Body: template.Data{},
		},
	}
}

type Context struct {
	context.Context
	Outputs  template.Data
	Response *Response
}

type Response struct {
	Body   template.Data
	Status *status.Status
}

type Output interface {
	ToAPI() *pilot.Output
	ToRaw() template.Data
}

type Referential interface {
	GetID() string
}
