package behavior

import (
	"context"
	"errors"
	"fmt"

	pilot "github.com/skiff-sh/pilot/api/go"
	"github.com/skiff-sh/pilot/server/pkg/template"
)

func Compile(beh *pilot.Behavior) (*Behavior, error) {
	out := &Behavior{
		Name:       beh.Name,
		Tendencies: make([]Tendency, 0, len(beh.Tendencies)),
	}
	for d, tend := range beh.Tendencies {
		t := new(Tendency)
		if v := tend.GetIf(); v != "" {
			expr, err := template.CompileExpression(v)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("tendency %d has an invalid if field", d), err)
			}
			t.Cond = expr
		}

		act, err := CompileAction(tend.GetId(), tend.Action)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("tendency %d has an invalid action", d), err)
		}
		t.Action = act
	}

	return out, nil
}

type Behavior struct {
	Name       string
	Tendencies []Tendency
}

func (b *Behavior) Provoke(ctx context.Context) (*Response, error) {
	c := newContext(ctx)
	for _, v := range b.Tendencies {
		out, err := v.Act(c)
		if err != nil {
			return nil, err
		}

		if out == nil {
			continue
		}

		raw := out.ToRaw()
		if ref, ok := v.Action.(Referential); ok && ref.GetID() != "" {
			c.Outputs.Set(ref.GetID(), raw)
		}
	}
	return c.Response, nil
}

var _ Action = &Tendency{}

type Tendency struct {
	Action Action
	Cond   template.Expression
}

func (t *Tendency) Act(c *Context) (out Output, err error) {
	if t.Cond != nil && !template.IsTruthy(t.Cond.Eval(c.Outputs)) {
		return nil, nil
	}
	return t.Action.Act(c)
}
