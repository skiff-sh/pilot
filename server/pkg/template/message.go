package template

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FieldExpressions map containing protobuf field names to compiled
// templates.
type FieldExpressions map[string]Expression

type FieldTemplateOpts struct {
	Force bool
}

// WithForce compiles all available fields whether they contain
// an expression or not.
func WithForce() FieldTemplateOpt {
	return func(o *FieldTemplateOpts) {
		o.Force = true
	}
}

type FieldTemplateOpt func(o *FieldTemplateOpts)

func NewFieldTemplates(m proto.Message, opts ...FieldTemplateOpt) (FieldExpressions, error) {
	op := &FieldTemplateOpts{}
	for _, v := range opts {
		v(op)
	}
	ref := m.ProtoReflect()
	fields := ref.Descriptor().Fields()
	out := FieldExpressions{}
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.Kind() != protoreflect.StringKind {
			continue
		}

		fieldExpr := ref.Get(field).String()
		if op.Force || !ContainsExpression(fieldExpr) {
			continue
		}

		expr, err := CompileExpression(fieldExpr)
		if err != nil {
			return nil, fmt.Errorf(`field "%s" has invalid expression "%s": %w`, field.Name(), fieldExpr, err)
		}
		out[string(field.Name())] = expr
	}
	return out, nil
}

func (f FieldExpressions) Apply(m proto.Message, d Data) error {
	ref := m.ProtoReflect()
	fields := ref.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		temp, ok := f[string(field.Name())]
		if !ok {
			continue
		}

		out := temp.Eval(d)
		out = strings.ToLower(strings.TrimSpace(out))
		switch field.Kind() {
		case protoreflect.StringKind:
			ref.Set(field, protoreflect.ValueOfString(out))
		}
	}
	return nil
}
