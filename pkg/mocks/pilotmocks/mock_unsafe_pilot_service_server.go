// Code generated by mockery v2.50.0. DO NOT EDIT.

package pilotmocks

import mock "github.com/stretchr/testify/mock"

// UnsafePilotServiceServer is an autogenerated mock type for the UnsafePilotServiceServer type
type UnsafePilotServiceServer struct {
	mock.Mock
}

type UnsafePilotServiceServer_Expecter struct {
	mock *mock.Mock
}

func (_m *UnsafePilotServiceServer) EXPECT() *UnsafePilotServiceServer_Expecter {
	return &UnsafePilotServiceServer_Expecter{mock: &_m.Mock}
}

// mustEmbedUnimplementedPilotServiceServer provides a mock function with no fields
func (_m *UnsafePilotServiceServer) mustEmbedUnimplementedPilotServiceServer() {
	_m.Called()
}

// UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'mustEmbedUnimplementedPilotServiceServer'
type UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call struct {
	*mock.Call
}

// mustEmbedUnimplementedPilotServiceServer is a helper method to define mock.On call
func (_e *UnsafePilotServiceServer_Expecter) mustEmbedUnimplementedPilotServiceServer() *UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call {
	return &UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call{Call: _e.mock.On("mustEmbedUnimplementedPilotServiceServer")}
}

func (_c *UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call) Run(run func()) *UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call) Return() *UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call {
	_c.Call.Return()
	return _c
}

func (_c *UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call) RunAndReturn(run func()) *UnsafePilotServiceServer_mustEmbedUnimplementedPilotServiceServer_Call {
	_c.Run(run)
	return _c
}

// NewUnsafePilotServiceServer creates a new instance of UnsafePilotServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUnsafePilotServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UnsafePilotServiceServer {
	mock := &UnsafePilotServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
