// Code generated by MockGen. DO NOT EDIT.
// Source: ./db.go

// Package mock_app_db is a generated GoMock package.
package mock_app_db

import (
	reflect "reflect"

	pbdto "github.com/aybjax/nis_lib/pbdto"
	gomock "github.com/golang/mock/gomock"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// AddCourseIdTo mocks base method.
func (m *MockDB) AddCourseIdTo(id, courseId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCourseIdTo", id, courseId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCourseIdTo indicates an expected call of AddCourseIdTo.
func (mr *MockDBMockRecorder) AddCourseIdTo(id, courseId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCourseIdTo", reflect.TypeOf((*MockDB)(nil).AddCourseIdTo), id, courseId)
}

// Close mocks base method.
func (m *MockDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDB)(nil).Close))
}

// Create mocks base method.
func (m *MockDB) Create(payload *pbdto.Student) (string, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", payload)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockDBMockRecorder) Create(payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDB)(nil).Create), payload)
}

// Delete mocks base method.
func (m *MockDB) Delete(id string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockDBMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDB)(nil).Delete), id)
}

// DeleteCourseIdFrom mocks base method.
func (m *MockDB) DeleteCourseIdFrom(id, courseId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCourseIdFrom", id, courseId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCourseIdFrom indicates an expected call of DeleteCourseIdFrom.
func (mr *MockDBMockRecorder) DeleteCourseIdFrom(id, courseId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCourseIdFrom", reflect.TypeOf((*MockDB)(nil).DeleteCourseIdFrom), id, courseId)
}

// GetCourseIds mocks base method.
func (m *MockDB) GetCourseIds(id string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCourseIds", id)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCourseIds indicates an expected call of GetCourseIds.
func (mr *MockDBMockRecorder) GetCourseIds(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCourseIds", reflect.TypeOf((*MockDB)(nil).GetCourseIds), id)
}

// ReadAll mocks base method.
func (m *MockDB) ReadAll() ([]*pbdto.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAll")
	ret0, _ := ret[0].([]*pbdto.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAll indicates an expected call of ReadAll.
func (mr *MockDBMockRecorder) ReadAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAll", reflect.TypeOf((*MockDB)(nil).ReadAll))
}

// ReadByCourseId mocks base method.
func (m *MockDB) ReadByCourseId(course_id string) ([]*pbdto.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadByCourseId", course_id)
	ret0, _ := ret[0].([]*pbdto.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadByCourseId indicates an expected call of ReadByCourseId.
func (mr *MockDBMockRecorder) ReadByCourseId(course_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadByCourseId", reflect.TypeOf((*MockDB)(nil).ReadByCourseId), course_id)
}

// ReadById mocks base method.
func (m *MockDB) ReadById(id string) (*pbdto.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadById", id)
	ret0, _ := ret[0].(*pbdto.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadById indicates an expected call of ReadById.
func (mr *MockDBMockRecorder) ReadById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadById", reflect.TypeOf((*MockDB)(nil).ReadById), id)
}

// Update mocks base method.
func (m *MockDB) Update(id string, payload *pbdto.Student) ([]string, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, payload)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockDBMockRecorder) Update(id, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDB)(nil).Update), id, payload)
}
