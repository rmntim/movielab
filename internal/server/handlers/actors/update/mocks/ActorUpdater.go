// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/rmntim/movielab/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// ActorUpdater is an autogenerated mock type for the ActorUpdater type
type ActorUpdater struct {
	mock.Mock
}

// GetActorById provides a mock function with given fields: id
func (_m *ActorUpdater) GetActorById(id int) (*entity.Actor, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetActorById")
	}

	var r0 *entity.Actor
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*entity.Actor, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *entity.Actor); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Actor)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateActor provides a mock function with given fields: id, actor
func (_m *ActorUpdater) UpdateActor(id int, actor *entity.Actor) error {
	ret := _m.Called(id, actor)

	if len(ret) == 0 {
		panic("no return value specified for UpdateActor")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, *entity.Actor) error); ok {
		r0 = rf(id, actor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewActorUpdater creates a new instance of ActorUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActorUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActorUpdater {
	mock := &ActorUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
