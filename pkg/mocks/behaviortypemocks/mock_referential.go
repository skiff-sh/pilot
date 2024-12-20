// Code generated by mockery v2.50.0. DO NOT EDIT.

package behaviortypemocks

import mock "github.com/stretchr/testify/mock"

// Referential is an autogenerated mock type for the Referential type
type Referential struct {
	mock.Mock
}

type Referential_Expecter struct {
	mock *mock.Mock
}

func (_m *Referential) EXPECT() *Referential_Expecter {
	return &Referential_Expecter{mock: &_m.Mock}
}

// GetID provides a mock function with no fields
func (_m *Referential) GetID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Referential_GetID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetID'
type Referential_GetID_Call struct {
	*mock.Call
}

// GetID is a helper method to define mock.On call
func (_e *Referential_Expecter) GetID() *Referential_GetID_Call {
	return &Referential_GetID_Call{Call: _e.mock.On("GetID")}
}

func (_c *Referential_GetID_Call) Run(run func()) *Referential_GetID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Referential_GetID_Call) Return(_a0 string) *Referential_GetID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Referential_GetID_Call) RunAndReturn(run func() string) *Referential_GetID_Call {
	_c.Call.Return(run)
	return _c
}

// NewReferential creates a new instance of Referential. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReferential(t interface {
	mock.TestingT
	Cleanup(func())
}) *Referential {
	mock := &Referential{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
