// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	domain "courses/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockCourses is a mock of Courses interface.
type MockCourses struct {
	ctrl     *gomock.Controller
	recorder *MockCoursesMockRecorder
}

// MockCoursesMockRecorder is the mock recorder for MockCourses.
type MockCoursesMockRecorder struct {
	mock *MockCourses
}

// NewMockCourses creates a new mock instance.
func NewMockCourses(ctrl *gomock.Controller) *MockCourses {
	mock := &MockCourses{ctrl: ctrl}
	mock.recorder = &MockCoursesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCourses) EXPECT() *MockCoursesMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCourses) Create(ctx context.Context, course domain.Courses) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, course)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCoursesMockRecorder) Create(ctx, course interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCourses)(nil).Create), ctx, course)
}

// Delele mocks base method.
func (m *MockCourses) Delele(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delele", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delele indicates an expected call of Delele.
func (mr *MockCoursesMockRecorder) Delele(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delele", reflect.TypeOf((*MockCourses)(nil).Delele), ctx, id)
}

// GetByID mocks base method.
func (m *MockCourses) GetByID(ctx context.Context, id int) (domain.Courses, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(domain.Courses)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCoursesMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCourses)(nil).GetByID), ctx, id)
}

// GetCoursesByIdStudent mocks base method.
func (m *MockCourses) GetCoursesByIdStudent(ctx context.Context, studentId string) ([]domain.Courses, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoursesByIdStudent", ctx, studentId)
	ret0, _ := ret[0].([]domain.Courses)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoursesByIdStudent indicates an expected call of GetCoursesByIdStudent.
func (mr *MockCoursesMockRecorder) GetCoursesByIdStudent(ctx, studentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoursesByIdStudent", reflect.TypeOf((*MockCourses)(nil).GetCoursesByIdStudent), ctx, studentId)
}

// Update mocks base method.
func (m *MockCourses) Update(ctx context.Context, course domain.Courses) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, course)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCoursesMockRecorder) Update(ctx, course interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCourses)(nil).Update), ctx, course)
}
