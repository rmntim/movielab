// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/rmntim/movielab/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// ActorGetter is an autogenerated mock type for the ActorGetter type
type ActorGetter struct {
	mock.Mock
}

// GetActors provides a mock function with given fields: limit, offset
func (_m *ActorGetter) GetActors(limit int, offset int) ([]entity.Actor, error) {
	ret := _m.Called(limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetActors")
	}

	var r0 []entity.Actor
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) ([]entity.Actor, error)); ok {
		return rf(limit, offset)
	}
	if rf, ok := ret.Get(0).(func(int, int) []entity.Actor); ok {
		r0 = rf(limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Actor)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewActorGetter creates a new instance of ActorGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActorGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActorGetter {
	mock := &ActorGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}