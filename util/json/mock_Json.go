// Code generated by mockery v2.42.1. DO NOT EDIT.

package json

import mock "github.com/stretchr/testify/mock"

// MockJson is an autogenerated mock type for the Json type
type MockJson struct {
	mock.Mock
}

type MockJson_Expecter struct {
	mock *mock.Mock
}

func (_m *MockJson) EXPECT() *MockJson_Expecter {
	return &MockJson_Expecter{mock: &_m.Mock}
}

// NewMockJson creates a new instance of MockJson. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockJson(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockJson {
	mock := &MockJson{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
