package behavior

import (
	pilot "github.com/skiff-sh/pilot/api/go"
	"github.com/skiff-sh/pilot/server/pkg/template"
)

type Output interface {
	ToAPI() *pilot.Output
	ToRaw() template.Data
}

var _ Output = &HTTPResponseOutput{}

type HTTPResponseOutput struct {
	*pilot.Output_HTTPResponse
}

func (h *HTTPResponseOutput) ToRaw() template.Data {
	d, _ := template.Unmarshal(h.Output_HTTPResponse)
	return d
}

func (h *HTTPResponseOutput) ToAPI() *pilot.Output {
	return &pilot.Output{HttpResponse: h.Output_HTTPResponse}
}

var _ Output = &ExecOutput{}

type ExecOutput struct {
	*pilot.Output_ExecOutput
}

func (e *ExecOutput) ToRaw() template.Data {
	d, _ := template.Unmarshal(e.Output_ExecOutput)
	return d
}

func (e *ExecOutput) ToAPI() *pilot.Output {
	return &pilot.Output{ExecOutput: e.Output_ExecOutput}
}
