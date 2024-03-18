// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ActorDeleter is an autogenerated mock type for the ActorDeleter type
type ActorDeleter struct {
	mock.Mock
}

// DeleteActor provides a mock function with given fields: id
func (_m *ActorDeleter) DeleteActor(id int) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteActor")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewActorDeleter creates a new instance of ActorDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActorDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActorDeleter {
	mock := &ActorDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
