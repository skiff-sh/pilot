package behavior

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"github.com/skiff-sh/pilot/pkg/template"

	pilot "github.com/skiff-sh/pilot/api/go"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type ActionTestSuite struct {
	suite.Suite
}

func (a *ActionTestSuite) TestActions() {
	type deps struct {
		Output behaviortype.Output
		Err    error
		Ctx    *behaviortype.Context
	}

	sleeperCalled := time.Duration(0)
	sleeperFunc = func(d time.Duration) {
		sleeperCalled = d
	}

	type test struct {
		Ctx                *behaviortype.Context
		Given              *pilot.Action
		ExpectedOutput     behaviortype.Output
		ExpectedOutputFunc func(d *deps)
		ID                 string
		ExpectedErr        string
		ExpectedCompileErr string
	}

	tests := map[string]test{
		"exec happy": {
			Given: &pilot.Action{
				Exec: &pilot.Action_Exec{
					Command: "/bin/sh",
					Args:    []string{"-c", `echo $derp`},
					EnvVars: map[string]string{"derp": "flerp"},
				},
			},
			ExpectedOutput: &ExecOutput{
				Output_ExecOutput: &pilot.Output_ExecOutput{
					Stdout: "flerp\n",
				},
			},
		},
		"exec invalid command": {
			Given: &pilot.Action{
				Exec: &pilot.Action_Exec{
					Command: "eddiecmd",
				},
			},
			ExpectedCompileErr: "exec: \"eddiecmd\": executable file not found in $PATH",
		},
		"exec invalid error code": {
			Given: &pilot.Action{
				Exec: &pilot.Action_Exec{
					Command: "curl",
					Args:    []string{"??"},
				},
			},
			ExpectedOutputFunc: func(d *deps) {
				a.Contains(d.Err.Error(), "exit status 3\ncurl: (3)")
			},
		},
		"http request happy": {
			Given: &pilot.Action{
				HttpRequest: &pilot.Action_HTTPRequest{
					Url:    "https://google.com",
					Method: http.MethodGet,
				},
			},
			ExpectedOutputFunc: func(d *deps) {
				resp := d.Output.(*HTTPResponseOutput)
				a.NotEmpty(resp.Headers)
				a.NotZero(resp.ContentLength)
				a.NotEmpty(resp.BodyRaw)
				resp.Headers = nil
				resp.BodyRaw = nil
				resp.ContentLength = 0
				exp := &pilot.Output_HTTPResponse{
					Status:     200,
					Proto:      "HTTP/2.0",
					ProtoMajor: 2,
				}
				a.Equal(exp.String(), resp.Output_HTTPResponse.String())
			},
		},
		"http request json": {
			Given: &pilot.Action{
				HttpRequest: &pilot.Action_HTTPRequest{
					Url:    "https://api.github.com/repos/skiff-sh/pilot",
					Method: http.MethodGet,
				},
			},
			ExpectedOutputFunc: func(d *deps) {
				resp := d.Output.(*HTTPResponseOutput)
				id := resp.Body.AsMap()["name"]
				a.Equal("pilot", id)
			},
		},
		"set status happy": {
			Given: &pilot.Action{
				SetStatus: &pilot.Action_SetStatus{
					Code:    1,
					Message: "message",
				},
			},
			ExpectedOutputFunc: func(d *deps) {
				a.Equal(status.New(1, "message"), d.Ctx.Response.Status)
			},
		},
		"set response happy": {
			Given: &pilot.Action{
				SetResponseField: &pilot.Action_SetResponseField{
					From: "derp",
					To:   "flerp",
				},
			},
			Ctx: &behaviortype.Context{
				Outputs:  template.Data{"derp": 1},
				Response: &behaviortype.Response{Body: make(template.Data)},
			},
			ExpectedOutputFunc: func(d *deps) {
				a.Equal(d.Ctx.Response.Body["flerp"], 1)
			},
		},
		"wait": {
			Given: &pilot.Action{
				Wait: durationpb.New(10 * time.Millisecond),
			},
			ExpectedOutputFunc: func(_ *deps) {
				a.Equal(10*time.Millisecond, sleeperCalled)
			},
		},
	}

	for desc, v := range tests {
		a.Run(desc, func() {
			act, err := CompileAction(v.ID, v.Given)
			if v.ExpectedCompileErr != "" || !a.NoError(err) {
				a.EqualError(err, v.ExpectedCompileErr)
				return
			}

			if v.Ctx == nil {
				v.Ctx = behaviortype.NewContext(context.TODO())
			}

			out, err := act.Act(v.Ctx)
			if v.ExpectedOutputFunc != nil {
				v.ExpectedOutputFunc(&deps{
					Ctx:    v.Ctx,
					Output: out,
					Err:    err,
				})
			} else {
				if v.ExpectedErr != "" || !a.NoError(err) {
					a.EqualError(err, v.ExpectedErr)
					return
				}
				a.Equal(v.ExpectedOutput, out)
			}
		})
	}
}

func TestActionTestSuite(t *testing.T) {
	suite.Run(t, new(ActionTestSuite))
}
