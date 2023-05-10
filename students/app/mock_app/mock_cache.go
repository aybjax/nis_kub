// Code generated by MockGen. DO NOT EDIT.
// Source: ./cache.go

// Package mock_app is a generated GoMock package.
package mock_app

import (
	reflect "reflect"

	pbdto "github.com/aybjax/nis_lib/pbdto"
	gomock "github.com/golang/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// InvalidateCreated mocks base method.
func (m *MockCache) InvalidateCreated() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateCreated")
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateCreated indicates an expected call of InvalidateCreated.
func (mr *MockCacheMockRecorder) InvalidateCreated() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateCreated", reflect.TypeOf((*MockCache)(nil).InvalidateCreated))
}

// InvalidateDeleted mocks base method.
func (m *MockCache) InvalidateDeleted(c_id string, oldCourseIds []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateDeleted", c_id, oldCourseIds)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateDeleted indicates an expected call of InvalidateDeleted.
func (mr *MockCacheMockRecorder) InvalidateDeleted(c_id, oldCourseIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateDeleted", reflect.TypeOf((*MockCache)(nil).InvalidateDeleted), c_id, oldCourseIds)
}

// InvalidateUpdated mocks base method.
func (m *MockCache) InvalidateUpdated(c_id string, newCourseIds, oldCourseIds []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateUpdated", c_id, newCourseIds, oldCourseIds)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateUpdated indicates an expected call of InvalidateUpdated.
func (mr *MockCacheMockRecorder) InvalidateUpdated(c_id, newCourseIds, oldCourseIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateUpdated", reflect.TypeOf((*MockCache)(nil).InvalidateUpdated), c_id, newCourseIds, oldCourseIds)
}

// RetrieveByCourseId mocks base method.
func (m *MockCache) RetrieveByCourseId(course_id string) ([]*pbdto.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveByCourseId", course_id)
	ret0, _ := ret[0].([]*pbdto.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveByCourseId indicates an expected call of RetrieveByCourseId.
func (mr *MockCacheMockRecorder) RetrieveByCourseId(course_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveByCourseId", reflect.TypeOf((*MockCache)(nil).RetrieveByCourseId), course_id)
}

// RetriveAll mocks base method.
func (m *MockCache) RetriveAll() ([]*pbdto.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetriveAll")
	ret0, _ := ret[0].([]*pbdto.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetriveAll indicates an expected call of RetriveAll.
func (mr *MockCacheMockRecorder) RetriveAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetriveAll", reflect.TypeOf((*MockCache)(nil).RetriveAll))
}

// RetriveOneById mocks base method.
func (m *MockCache) RetriveOneById(id string) (*pbdto.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetriveOneById", id)
	ret0, _ := ret[0].(*pbdto.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetriveOneById indicates an expected call of RetriveOneById.
func (mr *MockCacheMockRecorder) RetriveOneById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetriveOneById", reflect.TypeOf((*MockCache)(nil).RetriveOneById), id)
}

// WriteAll mocks base method.
func (m *MockCache) WriteAll(data []*pbdto.Student) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteAll", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteAll indicates an expected call of WriteAll.
func (mr *MockCacheMockRecorder) WriteAll(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAll", reflect.TypeOf((*MockCache)(nil).WriteAll), data)
}

// WriteByCourseId mocks base method.
func (m *MockCache) WriteByCourseId(course_id string, data []*pbdto.Student) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteByCourseId", course_id, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteByCourseId indicates an expected call of WriteByCourseId.
func (mr *MockCacheMockRecorder) WriteByCourseId(course_id, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteByCourseId", reflect.TypeOf((*MockCache)(nil).WriteByCourseId), course_id, data)
}

// WriteOneById mocks base method.
func (m *MockCache) WriteOneById(id string, data *pbdto.Student) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteOneById", id, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteOneById indicates an expected call of WriteOneById.
func (mr *MockCacheMockRecorder) WriteOneById(id, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteOneById", reflect.TypeOf((*MockCache)(nil).WriteOneById), id, data)
}
