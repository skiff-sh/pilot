package controller

import (
	"context"

	"github.com/skiff-sh/pilot/pkg/behavior"
	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"github.com/skiff-sh/pilot/pkg/protovalidatetype"

	"github.com/orcaman/concurrent-map/v2"
	"github.com/skiff-sh/pilot/api/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pilot.PilotServiceServer = &Controller{}

func NewController(val protovalidatetype.Validator) *Controller {
	out := &Controller{
		Behaviors: cmap.New[behaviortype.Interface](),
		Validator: val,
	}
	return out
}

type Controller struct {
	Behaviors cmap.ConcurrentMap[string, behaviortype.Interface]
	Validator protovalidatetype.Validator
}

func (c *Controller) CreateBehavior(_ context.Context, request *pilot.CreateBehavior_Request) (*pilot.CreateBehavior_Response, error) {
	err := c.Validator.Validate(request)
	if err != nil {
		return nil, err
	}
	beh, err := behavior.Compile(request.GetBehavior())
	if err != nil {
		return nil, err
	}

	c.Behaviors.Set(request.Behavior.Name, beh)

	return &pilot.CreateBehavior_Response{}, nil
}

func (c *Controller) ProvokeBehavior(ctx context.Context, request *pilot.ProvokeBehavior_Request) (*pilot.ProvokeBehavior_Response, error) {
	err := c.Validator.Validate(request)
	if err != nil {
		return nil, err
	}

	beh, ok := c.Behaviors.Get(request.GetName())
	if !ok {
		return nil, status.Newf(codes.NotFound, "behavior %s not found", request.GetName()).Err()
	}

	resp, err := beh.Provoke(ctx)
	if err != nil {
		return nil, err
	}

	if resp.Status != nil && resp.Status.Code() != codes.OK {
		return nil, resp.Status.Err()
	}

	return &pilot.ProvokeBehavior_Response{Body: resp.Body.ToProto()}, nil
}
