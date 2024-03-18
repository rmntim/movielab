// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/rmntim/movielab/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// MovieCreator is an autogenerated mock type for the MovieCreator type
type MovieCreator struct {
	mock.Mock
}

// CreateMovie provides a mock function with given fields: movie
func (_m *MovieCreator) CreateMovie(movie *entity.NewMovie) (int, error) {
	ret := _m.Called(movie)

	if len(ret) == 0 {
		panic("no return value specified for CreateMovie")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.NewMovie) (int, error)); ok {
		return rf(movie)
	}
	if rf, ok := ret.Get(0).(func(*entity.NewMovie) int); ok {
		r0 = rf(movie)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*entity.NewMovie) error); ok {
		r1 = rf(movie)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMovieCreator creates a new instance of MovieCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMovieCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MovieCreator {
	mock := &MovieCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}