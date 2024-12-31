package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"

	"github.com/skiff-sh/pilot/api/go/pilot"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/bufbuild/protovalidate-go"
	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"github.com/skiff-sh/pilot/pkg/mocks/behaviortypemocks"
	"github.com/skiff-sh/pilot/pkg/protovalidatetype"
	"github.com/skiff-sh/pilot/pkg/template"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type ControllerTestSuite struct {
	suite.Suite
}

func (c *ControllerTestSuite) TestCreateBehavior() {
	type deps struct {
		Val protovalidatetype.Validator
	}

	type test struct {
		Given        *pilot.CreateBehavior_Request
		ExpectedErr  string
		ExpectedFunc func(c *PilotGRPC)
		Constructor  func(d *deps) test
	}

	tests := map[string]test{
		"happy": {
			Constructor: func(d *deps) test {
				return test{
					Given: &pilot.CreateBehavior_Request{
						Behavior: &pilot.Behavior{
							Tendencies: []*pilot.Tendency{
								{
									Action: &pilot.Action{
										Wait: durationpb.New(1 * time.Second),
									},
								},
							},
							Name: "derp",
						},
					},
					ExpectedFunc: func(con *PilotGRPC) {
						c.Equal(1, len(con.Behaviors.Keys()))
					},
				}
			},
		},
		"validation failed": {
			Constructor: func(d *deps) test {
				return test{
					ExpectedErr: "validation error:\n - behavior: value is required [required]",
				}
			},
		},
	}

	for desc, v := range tests {
		c.Run(desc, func() {
			val, _ := protovalidate.New()
			d := &deps{Val: val}
			if v.Constructor != nil {
				v = v.Constructor(d)
			}

			con := NewPilotGRPC(d.Val, cmap.New[behaviortype.Interface]())
			ctx := context.TODO()
			_, err := con.CreateBehavior(ctx, v.Given)
			if v.ExpectedErr != "" || !c.NoError(err) {
				c.EqualError(err, v.ExpectedErr)
				return
			}

			v.ExpectedFunc(con)
		})
	}
}

func (c *ControllerTestSuite) TestProvokeBehavior() {
	type deps struct {
		Val protovalidatetype.Validator
	}

	type test struct {
		Given       behaviortype.Interface
		GivenName   string
		Expected    map[string]any
		ExpectedErr string
		Constructor func(d *deps) test
	}

	tests := map[string]test{
		"happy": {
			Constructor: func(d *deps) test {
				beh := new(behaviortypemocks.Interface)
				name := gofakeit.LoremIpsumWord()
				beh.EXPECT().Provoke(mock.Anything).Return(&behaviortype.Response{Body: template.Data{"derp": "flerp"}}, nil)
				beh.EXPECT().GetName().Return(name)
				return test{
					Given:     beh,
					GivenName: name,
					Expected:  map[string]any{"derp": "flerp"},
				}
			},
		},
		"provoke error": {
			Constructor: func(d *deps) test {
				beh := new(behaviortypemocks.Interface)
				name := gofakeit.LoremIpsumWord()
				beh.EXPECT().Provoke(mock.Anything).Return(nil, errors.New("err"))
				beh.EXPECT().GetName().Return(name)
				return test{
					Given:       beh,
					GivenName:   name,
					ExpectedErr: "err",
				}
			},
		},
		"provoke body empty": {
			Constructor: func(d *deps) test {
				beh := new(behaviortypemocks.Interface)
				name := gofakeit.LoremIpsumWord()
				beh.EXPECT().Provoke(mock.Anything).Return(&behaviortype.Response{}, nil)
				beh.EXPECT().GetName().Return(name)
				return test{
					Given:     beh,
					GivenName: name,
					Expected:  map[string]any{},
				}
			},
		},
		"behavior not found": {
			Constructor: func(d *deps) test {
				beh := new(behaviortypemocks.Interface)
				name := "beh"
				beh.EXPECT().Provoke(mock.Anything).Return(&behaviortype.Response{}, nil)
				beh.EXPECT().GetName().Return(name + "asd")
				return test{
					Given:       beh,
					GivenName:   name,
					ExpectedErr: "rpc error: code = NotFound desc = behavior beh not found",
				}
			},
		},
		"provoke returns status": {
			Constructor: func(d *deps) test {
				beh := new(behaviortypemocks.Interface)
				name := gofakeit.LoremIpsumWord()
				beh.EXPECT().Provoke(mock.Anything).Return(&behaviortype.Response{Status: status.New(codes.AlreadyExists, "")}, nil)
				beh.EXPECT().GetName().Return(name)
				return test{
					Given:       beh,
					GivenName:   name,
					ExpectedErr: "rpc error: code = AlreadyExists desc = ",
				}
			},
		},
	}

	for desc, v := range tests {
		c.Run(desc, func() {
			val, _ := protovalidate.New()
			d := &deps{Val: val}
			if v.Constructor != nil {
				v = v.Constructor(d)
			}

			con := NewPilotGRPC(d.Val, cmap.New[behaviortype.Interface]())
			con.Behaviors.Set(v.Given.GetName(), v.Given)
			ctx := context.TODO()
			actual, err := con.ProvokeBehavior(ctx, &pilot.ProvokeBehavior_Request{Name: v.GivenName})
			if v.ExpectedErr != "" || !c.NoError(err) {
				c.EqualError(err, v.ExpectedErr)
				return
			}

			c.Equal(v.Expected, actual.Body.AsMap())
		})
	}
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}
