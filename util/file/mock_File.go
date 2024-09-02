// Code generated by mockery v2.43.2. DO NOT EDIT.

package file

import mock "github.com/stretchr/testify/mock"

// MockUtilFile is an autogenerated mock type for the File type
type MockUtilFile struct {
	mock.Mock
}

type MockUtilFile_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUtilFile) EXPECT() *MockUtilFile_Expecter {
	return &MockUtilFile_Expecter{mock: &_m.Mock}
}

// Exists provides a mock function with given fields: path
func (_m *MockUtilFile) Exists(path string) bool {
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

// MockUtilFile_Exists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exists'
type MockUtilFile_Exists_Call struct {
	*mock.Call
}

// Exists is a helper method to define mock.On call
//   - path string
func (_e *MockUtilFile_Expecter) Exists(path interface{}) *MockUtilFile_Exists_Call {
	return &MockUtilFile_Exists_Call{Call: _e.mock.On("Exists", path)}
}

func (_c *MockUtilFile_Exists_Call) Run(run func(path string)) *MockUtilFile_Exists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockUtilFile_Exists_Call) Return(_a0 bool) *MockUtilFile_Exists_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUtilFile_Exists_Call) RunAndReturn(run func(string) bool) *MockUtilFile_Exists_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUtilFile creates a new instance of MockUtilFile. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUtilFile(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUtilFile {
	mock := &MockUtilFile{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
