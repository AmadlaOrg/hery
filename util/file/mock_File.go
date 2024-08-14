// Code generated by mockery v2.43.2. DO NOT EDIT.

package file

import mock "github.com/stretchr/testify/mock"

// MockFile is an autogenerated mock type for the File type
type MockFile struct {
	mock.Mock
}

type MockFile_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFile) EXPECT() *MockFile_Expecter {
	return &MockFile_Expecter{mock: &_m.Mock}
}

// Exists provides a mock function with given fields: path
func (_m *MockFile) Exists(path string) bool {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for Exists")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockFile_Exists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exists'
type MockFile_Exists_Call struct {
	*mock.Call
}

// Exists is a helper method to define mock.On call
//   - path string
func (_e *MockFile_Expecter) Exists(path interface{}) *MockFile_Exists_Call {
	return &MockFile_Exists_Call{Call: _e.mock.On("Exists", path)}
}

func (_c *MockFile_Exists_Call) Run(run func(path string)) *MockFile_Exists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockFile_Exists_Call) Return(_a0 bool) *MockFile_Exists_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFile_Exists_Call) RunAndReturn(run func(string) bool) *MockFile_Exists_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFile creates a new instance of MockFile. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFile(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFile {
	mock := &MockFile{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
