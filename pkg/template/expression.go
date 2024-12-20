package template

import (
	"strings"

	"github.com/flosch/pongo2/v6"
)

func CompileExpression(expr string) (Expression, error) {
	tmp, err := pongo2.FromString(expr)
	if err != nil {
		return nil, err
	}
	return &expression{Template: tmp}, nil
}

func ContainsExpression(s string) bool {
	return strings.Contains(s, "{{")
}

type Expression interface {
	Eval(data Data) string
}

type expression struct {
	Template *pongo2.Template
}

func (e *expression) Eval(data Data) string {
	v, _ := e.Template.Execute(data.toPongo())
	return v
}
