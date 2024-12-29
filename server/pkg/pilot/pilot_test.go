package pilot

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/skiff-sh/config/ptr"
	"github.com/skiff-sh/pilot/api/go/pilot"
	"github.com/skiff-sh/pilot/pkg/mocks/pilotmocks"
	"github.com/skiff-sh/pilot/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/durationpb"
)

type PilotTestSuite struct {
	suite.Suite
}

func (p *PilotTestSuite) TestBuilder() {
	type deps struct {
		Cl *pilotmocks.PilotServiceClient
	}

	type test struct {
		Constructor func(d *deps) test
		Given       CreateBehaviorSender
		Expected    *pilot.CreateBehavior_Request
		ExpectedErr string
	}

	cl := &client{}

	tests := map[string]test{
		"http and exec": {
			Expected: &pilot.CreateBehavior_Request{
				Behavior: &pilot.Behavior{
					Tendencies: []*pilot.Tendency{
						{
							Action: &pilot.Action{
								HttpRequest: &pilot.Action_HTTPRequest{
									Url:    "google.com",
									Method: http.MethodPost,
									Headers: map[string]string{
										"hi": "there",
										"h":  "t",
									},
									Body: []byte(`{"derp":"flerp"}`),
								},
							},
							Id: ptr.Ptr("id"),
						},
						{
							Action: &pilot.Action{
								Exec: &pilot.Action_Exec{
									Command:    "echo",
									Args:       []string{"hi"},
									EnvVars:    map[string]string{"KEY": "VAL"},
									WorkingDir: "/root",
								},
							},
							If: ptr.Ptr("cond"),
							Id: ptr.Ptr("exec"),
						},
					},
					Name: "derp",
				},
			},
			Given: cl.NewBehavior().Name("derp").Tendency().ID("id").Action().
				HTTPRequest("google.com",
					WithHTTPMethod(http.MethodPost),
					WithHTTPHeaders(map[string]string{"hi": "there"}),
					WithHTTPHeader("h", "t"),
					WithHTTPJSONBody(map[string]string{"derp": "flerp"}),
				).Tendency().Condition("cond").ID("exec").Action().
				Exec("echo", WithExecArgs("hi"), WithExecDir("/root"), WithEnvVars(map[string]string{"KEY": "VAL"})),
		},
		"set response field": {
			Given: cl.NewBehavior().Name("derp").Tendency().Action().SetResponseField("from", "to"),
			Expected: &pilot.CreateBehavior_Request{
				Behavior: &pilot.Behavior{
					Tendencies: []*pilot.Tendency{
						{
							Action: &pilot.Action{
								SetResponseField: &pilot.Action_SetResponseField{
									From: "from",
									To:   "to",
								},
							},
						},
					},
					Name: "derp",
				},
			},
		},
		"wait": {
			Given: cl.NewBehavior().Name("derp").Tendency().Action().Wait(time.Second),
			Expected: &pilot.CreateBehavior_Request{
				Behavior: &pilot.Behavior{
					Tendencies: []*pilot.Tendency{
						{
							Action: &pilot.Action{
								Wait: durationpb.New(time.Second),
							},
						},
					},
					Name: "derp",
				},
			},
		},
		"set status": {
			Given: cl.NewBehavior().Name("derp").Tendency().Action().SetStatus(codes.AlreadyExists, "exists"),
			Expected: &pilot.CreateBehavior_Request{
				Behavior: &pilot.Behavior{
					Tendencies: []*pilot.Tendency{
						{
							Action: &pilot.Action{
								SetStatus: &pilot.Action_SetStatus{
									Code:    uint32(codes.AlreadyExists),
									Message: "exists",
								},
							},
						},
					},
					Name: "derp",
				},
			},
		},
	}

	for desc, v := range tests {
		p.Run(desc, func() {
			d := &deps{
				Cl: new(pilotmocks.PilotServiceClient),
			}
			cl.Cl = d.Cl

			if v.Constructor != nil {
				v = v.Constructor(d)
			}

			ctx := context.TODO()

			d.Cl.EXPECT().CreateBehavior(ctx, mock.MatchedBy(func(req *pilot.CreateBehavior_Request) bool {
				return p.Empty(testutil.DiffProto(v.Expected, req))
			})).Return(nil, nil)
			_, err := v.Given.Send(ctx)
			if v.ExpectedErr != "" || !p.NoError(err) {
				p.EqualError(err, v.ExpectedErr)
				return
			}
		})
	}
}

func TestPilotTestSuite(t *testing.T) {
	suite.Run(t, new(PilotTestSuite))
}
