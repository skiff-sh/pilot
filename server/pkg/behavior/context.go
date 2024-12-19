package behavior

import (
	"context"

	"github.com/skiff-sh/pilot/server/pkg/template"
	"google.golang.org/grpc/status"
)

func newContext(ctx context.Context) *Context {
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
