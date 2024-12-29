package behavior

import (
	"context"
	"errors"
	"fmt"
	"github.com/skiff-sh/pilot/api/go/pilot"

	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"github.com/skiff-sh/pilot/pkg/template"
)

func Compile(beh *pilot.Behavior) (*Behavior, error) {
	out := &Behavior{
		Name:       beh.Name,
		Tendencies: make([]*Tendency, 0, len(beh.Tendencies)),
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

		act, err := CompileAction(tend.GetId(), tend.GetAction())
		if err != nil {
			return nil, errors.Join(fmt.Errorf("tendency %d has an invalid action", d), err)
		}
		t.Action = act
		out.Tendencies = append(out.Tendencies, t)
	}

	return out, nil
}

type Behavior struct {
	Name       string
	Tendencies []*Tendency
}

func (b *Behavior) GetName() string {
	return b.Name
}

func (b *Behavior) Provoke(ctx context.Context) (*behaviortype.Response, error) {
	c := behaviortype.NewContext(ctx)
	for _, v := range b.Tendencies {
		out, err := v.Act(c)
		if err != nil {
			return nil, err
		}
		if out == nil {
			continue
		}

		raw := out.ToRaw()
		if ref, ok := v.Action.(behaviortype.Referential); ok && ref.GetID() != "" {
			c.Outputs.Set(ref.GetID(), raw)
		}
	}
	return c.Response, nil
}

var _ behaviortype.Action = &Tendency{}

type Tendency struct {
	Action behaviortype.Action
	Cond   template.Expression
}

func (t *Tendency) Act(c *behaviortype.Context) (out behaviortype.Output, err error) {
	if t.Cond != nil && !template.IsTruthy(t.Cond.Eval(c.Outputs)) {
		return nil, nil
	}
	return t.Action.Act(c)
}
