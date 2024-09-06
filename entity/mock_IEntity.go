// Code generated by mockery v2.43.2. DO NOT EDIT.

package entity

import (
	storage "github.com/AmadlaOrg/hery/storage"
	mock "github.com/stretchr/testify/mock"
)

// MockEntity is an autogenerated mock type for the IEntity type
type MockEntity struct {
	mock.Mock
}

type MockEntity_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEntity) EXPECT() *MockEntity_Expecter {
	return &MockEntity_Expecter{mock: &_m.Mock}
}

// CheckDuplicateEntity provides a mock function with given fields: entities, entityMeta
func (_m *MockEntity) CheckDuplicateEntity(entities []Entity, entityMeta Entity) error {
	ret := _m.Called(entities, entityMeta)

	if len(ret) == 0 {
		panic("no return value specified for CheckDuplicateEntity")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]Entity, Entity) error); ok {
		r0 = rf(entities, entityMeta)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEntity_CheckDuplicateEntity_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckDuplicateEntity'
type MockEntity_CheckDuplicateEntity_Call struct {
	*mock.Call
}

// CheckDuplicateEntity is a helper method to define mock.On call
//   - entities []Entity
//   - entityMeta Entity
func (_e *MockEntity_Expecter) CheckDuplicateEntity(entities interface{}, entityMeta interface{}) *MockEntity_CheckDuplicateEntity_Call {
	return &MockEntity_CheckDuplicateEntity_Call{Call: _e.mock.On("CheckDuplicateEntity", entities, entityMeta)}
}

func (_c *MockEntity_CheckDuplicateEntity_Call) Run(run func(entities []Entity, entityMeta Entity)) *MockEntity_CheckDuplicateEntity_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]Entity), args[1].(Entity))
	})
	return _c
}

func (_c *MockEntity_CheckDuplicateEntity_Call) Return(_a0 error) *MockEntity_CheckDuplicateEntity_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEntity_CheckDuplicateEntity_Call) RunAndReturn(run func([]Entity, Entity) error) *MockEntity_CheckDuplicateEntity_Call {
	_c.Call.Return(run)
	return _c
}

// FindEntityDir provides a mock function with given fields: paths, entityVals
func (_m *MockEntity) FindEntityDir(paths storage.AbsPaths, entityVals Entity) (string, error) {
	ret := _m.Called(paths, entityVals)

	if len(ret) == 0 {
		panic("no return value specified for FindEntityDir")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(storage.AbsPaths, Entity) (string, error)); ok {
		return rf(paths, entityVals)
	}
	if rf, ok := ret.Get(0).(func(storage.AbsPaths, Entity) string); ok {
		r0 = rf(paths, entityVals)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(storage.AbsPaths, Entity) error); ok {
		r1 = rf(paths, entityVals)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEntity_FindEntityDir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindEntityDir'
type MockEntity_FindEntityDir_Call struct {
	*mock.Call
}

// FindEntityDir is a helper method to define mock.On call
//   - paths storage.AbsPaths
//   - entityVals Entity
func (_e *MockEntity_Expecter) FindEntityDir(paths interface{}, entityVals interface{}) *MockEntity_FindEntityDir_Call {
	return &MockEntity_FindEntityDir_Call{Call: _e.mock.On("FindEntityDir", paths, entityVals)}
}

func (_c *MockEntity_FindEntityDir_Call) Run(run func(paths storage.AbsPaths, entityVals Entity)) *MockEntity_FindEntityDir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(storage.AbsPaths), args[1].(Entity))
	})
	return _c
}

func (_c *MockEntity_FindEntityDir_Call) Return(_a0 string, _a1 error) *MockEntity_FindEntityDir_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEntity_FindEntityDir_Call) RunAndReturn(run func(storage.AbsPaths, Entity) (string, error)) *MockEntity_FindEntityDir_Call {
	_c.Call.Return(run)
	return _c
}

// GeneratePseudoVersionPattern provides a mock function with given fields: name, version
func (_m *MockEntity) GeneratePseudoVersionPattern(name string, version string) string {
	ret := _m.Called(name, version)

	if len(ret) == 0 {
		panic("no return value specified for GeneratePseudoVersionPattern")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(name, version)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockEntity_GeneratePseudoVersionPattern_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GeneratePseudoVersionPattern'
type MockEntity_GeneratePseudoVersionPattern_Call struct {
	*mock.Call
}

// GeneratePseudoVersionPattern is a helper method to define mock.On call
//   - name string
//   - version string
func (_e *MockEntity_Expecter) GeneratePseudoVersionPattern(name interface{}, version interface{}) *MockEntity_GeneratePseudoVersionPattern_Call {
	return &MockEntity_GeneratePseudoVersionPattern_Call{Call: _e.mock.On("GeneratePseudoVersionPattern", name, version)}
}

func (_c *MockEntity_GeneratePseudoVersionPattern_Call) Run(run func(name string, version string)) *MockEntity_GeneratePseudoVersionPattern_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockEntity_GeneratePseudoVersionPattern_Call) Return(_a0 string) *MockEntity_GeneratePseudoVersionPattern_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEntity_GeneratePseudoVersionPattern_Call) RunAndReturn(run func(string, string) string) *MockEntity_GeneratePseudoVersionPattern_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEntity creates a new instance of MockEntity. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEntity(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEntity {
	mock := &MockEntity{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}