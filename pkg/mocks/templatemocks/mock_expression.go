// Code generated by mockery v2.50.0. DO NOT EDIT.

package templatemocks

import (
	template "github.com/skiff-sh/pilot/pkg/template"
	mock "github.com/stretchr/testify/mock"
)

// Expression is an autogenerated mock type for the Expression type
type Expression struct {
	mock.Mock
}

type Expression_Expecter struct {
	mock *mock.Mock
}

func (_m *Expression) EXPECT() *Expression_Expecter {
	return &Expression_Expecter{mock: &_m.Mock}
}

// Eval provides a mock function with given fields: data
func (_m *Expression) Eval(data template.Data) string {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Eval")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(template.Data) string); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Expression_Eval_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Eval'
type Expression_Eval_Call struct {
	*mock.Call
}

// Eval is a helper method to define mock.On call
//   - data template.Data
func (_e *Expression_Expecter) Eval(data interface{}) *Expression_Eval_Call {
	return &Expression_Eval_Call{Call: _e.mock.On("Eval", data)}
}

func (_c *Expression_Eval_Call) Run(run func(data template.Data)) *Expression_Eval_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(template.Data))
	})
	return _c
}

func (_c *Expression_Eval_Call) Return(_a0 string) *Expression_Eval_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Expression_Eval_Call) RunAndReturn(run func(template.Data) string) *Expression_Eval_Call {
	_c.Call.Return(run)
	return _c
}

// NewExpression creates a new instance of Expression. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExpression(t interface {
	mock.TestingT
	Cleanup(func())
}) *Expression {
	mock := &Expression{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}