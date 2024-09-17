// Code generated by mockery v2.43.2. DO NOT EDIT.

package validation

import mock "github.com/stretchr/testify/mock"

// MockEntitySchemaValidation is an autogenerated mock type for the IValidation type
type MockEntitySchemaValidation struct {
	mock.Mock
}

type MockEntitySchemaValidation_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEntitySchemaValidation) EXPECT() *MockEntitySchemaValidation_Expecter {
	return &MockEntitySchemaValidation_Expecter{mock: &_m.Mock}
}

// Id provides a mock function with given fields: id, collectionName, entityUri
func (_m *MockEntitySchemaValidation) Id(id string, collectionName string, entityUri string) error {
	ret := _m.Called(id, collectionName, entityUri)

	if len(ret) == 0 {
		panic("no return value specified for Id")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(id, collectionName, entityUri)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEntitySchemaValidation_Id_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Id'
type MockEntitySchemaValidation_Id_Call struct {
	*mock.Call
}

// Id is a helper method to define mock.On call
//   - id string
//   - collectionName string
//   - entityUri string
func (_e *MockEntitySchemaValidation_Expecter) Id(id interface{}, collectionName interface{}, entityUri interface{}) *MockEntitySchemaValidation_Id_Call {
	return &MockEntitySchemaValidation_Id_Call{Call: _e.mock.On("Id", id, collectionName, entityUri)}
}

func (_c *MockEntitySchemaValidation_Id_Call) Run(run func(id string, collectionName string, entityUri string)) *MockEntitySchemaValidation_Id_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockEntitySchemaValidation_Id_Call) Return(_a0 error) *MockEntitySchemaValidation_Id_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEntitySchemaValidation_Id_Call) RunAndReturn(run func(string, string, string) error) *MockEntitySchemaValidation_Id_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEntitySchemaValidation creates a new instance of MockEntitySchemaValidation. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEntitySchemaValidation(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEntitySchemaValidation {
	mock := &MockEntitySchemaValidation{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}