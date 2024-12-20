package behavior

import (
	"context"
	"testing"
	"time"

	"github.com/skiff-sh/config/ptr"
	pilot "github.com/skiff-sh/pilot/api/go"
	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"github.com/skiff-sh/pilot/pkg/mocks/behaviortypemocks"
	"github.com/skiff-sh/pilot/pkg/mocks/templatemocks"
	"github.com/skiff-sh/pilot/pkg/template"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/durationpb"
)

type BehaviorTestSuite struct {
	suite.Suite
}

func (b *BehaviorTestSuite) TestCompile() {
	type test struct {
		Given        *pilot.Behavior
		Expected     *Behavior
		ExpectedFunc func(b *Behavior)
		ExpectedErr  string
	}

	tests := map[string]test{
		"happy": {
			Given: &pilot.Behavior{Name: "derp", Tendencies: []*pilot.Tendency{
				{
					Action: &pilot.Action{Wait: durationpb.New(1 * time.Second)},
				},
			}},
			Expected: &Behavior{
				Name: "derp",
				Tendencies: []*Tendency{
					{
						Action: &Wait{Dur: 1 * time.Second},
					},
				},
			},
		},
		"cond": {
			Given: &pilot.Behavior{Name: "derp", Tendencies: []*pilot.Tendency{
				{
					If:     ptr.Ptr("{{derp}}"),
					Action: &pilot.Action{Wait: durationpb.New(1 * time.Second)},
				},
			}},
			ExpectedFunc: func(beh *Behavior) {
				b.True(template.IsTruthy(beh.Tendencies[0].Cond.Eval(template.Data{"derp": true})))
			},
		},
		"cond fails": {
			Given: &pilot.Behavior{Name: "derp", Tendencies: []*pilot.Tendency{
				{
					If:     ptr.Ptr("{{derp}}"),
					Action: &pilot.Action{Wait: durationpb.New(1 * time.Second)},
				},
			}},
			ExpectedFunc: func(beh *Behavior) {
				b.False(template.IsTruthy(beh.Tendencies[0].Cond.Eval(template.Data{"derp": false})))
			},
		},
		"empty action": {
			Given: &pilot.Behavior{Name: "derp", Tendencies: []*pilot.Tendency{
				{Action: &pilot.Action{}},
			}},
			ExpectedErr: "tendency 0 has an invalid action\nunknown or missing action",
		},
		"invalid condition": {
			Given: &pilot.Behavior{Name: "derp", Tendencies: []*pilot.Tendency{
				{
					If:     ptr.Ptr("{{derp}"),
					Action: &pilot.Action{Wait: durationpb.New(1 * time.Second)},
				},
			}},
			ExpectedErr: "tendency 0 has an invalid if field\n[Error (where: parser) in <string> | Line 1 Col 7 near '}'] '}}' expected",
		},
	}

	for desc, v := range tests {
		b.Run(desc, func() {
			actual, err := Compile(v.Given)
			if v.ExpectedErr != "" || !b.NoError(err) {
				b.EqualError(err, v.ExpectedErr)
				return
			}
			if v.ExpectedFunc != nil {
				v.ExpectedFunc(actual)
			} else {
				b.Equal(v.Expected, actual)
			}
		})
	}
}

func (b *BehaviorTestSuite) TestProvoke() {
	type test struct {
		Constructor func() test
		Given       *Behavior
		Expected    *behaviortype.Response
		ExpectedErr string
	}

	tests := map[string]test{
		"happy": {
			Constructor: func() test {
				act := new(behaviortypemocks.Action)
				act.EXPECT().Act(mock.Anything).Return(nil, nil)
				return test{
					Given: &Behavior{
						Name: "derp",
						Tendencies: []*Tendency{
							{Action: act},
						},
					},
					Expected: &behaviortype.Response{Body: template.Data{}},
				}
			},
		},
		"cond failed": {
			Constructor: func() test {
				expr := new(templatemocks.Expression)
				expr.EXPECT().Eval(mock.Anything).Return("false")
				return test{
					Given: &Behavior{
						Name: "derp",
						Tendencies: []*Tendency{
							{Cond: expr},
						},
					},
					Expected: &behaviortype.Response{Body: template.Data{}},
				}
			},
		},
		"referential action": {
			Constructor: func() test {
				act := &RefAction{
					Action:      new(behaviortypemocks.Action),
					Referential: new(behaviortypemocks.Referential),
				}
				act1 := new(behaviortypemocks.Action)
				expected := template.Data{
					"id": template.Data{"derp": "flerp"},
				}

				act1.EXPECT().Act(mock.Anything).RunAndReturn(func(ctx *behaviortype.Context) (behaviortype.Output, error) {
					b.Equal(expected, ctx.Outputs)
					return nil, nil
				})
				out := new(behaviortypemocks.Output)
				out.EXPECT().ToRaw().Return(template.Data{"derp": "flerp"})
				act.Action.EXPECT().Act(mock.Anything).Return(out, nil)
				act.Referential.EXPECT().GetID().Return("id")

				return test{
					Given: &Behavior{
						Name: "derp",
						Tendencies: []*Tendency{
							{Action: act},
							{Action: act1},
						},
					},
					Expected: &behaviortype.Response{
						Body: template.Data{},
					},
				}
			},
		},
		"outputs action": {
			Constructor: func() test {
				act := new(behaviortypemocks.Action)
				expected := template.Data{"derp": "flerp"}
				act.EXPECT().Act(mock.Anything).RunAndReturn(func(ctx *behaviortype.Context) (behaviortype.Output, error) {
					ctx.Response.Body.Set("id", expected)
					return nil, nil
				})

				return test{
					Given: &Behavior{
						Name: "derp",
						Tendencies: []*Tendency{
							{Action: act},
						},
					},
					Expected: &behaviortype.Response{
						Body: template.Data{
							"id": expected,
						},
					},
				}
			},
		},
	}

	for desc, v := range tests {
		b.Run(desc, func() {
			ctx := context.TODO()
			if v.Constructor != nil {
				v = v.Constructor()
			}

			actual, err := v.Given.Provoke(ctx)
			if v.ExpectedErr != "" || !b.NoError(err) {
				b.EqualError(err, v.ExpectedErr)
				return
			}

			b.Equal(v.Expected, actual)
		})
	}
}

type RefAction struct {
	*behaviortypemocks.Action
	*behaviortypemocks.Referential
}

func TestBehaviorTestSuite(t *testing.T) {
	suite.Run(t, new(BehaviorTestSuite))
}
