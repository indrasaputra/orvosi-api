// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/medical_record_creator.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/orvosi/api/entity"
)

// MockCreateMedicalRecord is a mock of CreateMedicalRecord interface
type MockCreateMedicalRecord struct {
	ctrl     *gomock.Controller
	recorder *MockCreateMedicalRecordMockRecorder
}

// MockCreateMedicalRecordMockRecorder is the mock recorder for MockCreateMedicalRecord
type MockCreateMedicalRecordMockRecorder struct {
	mock *MockCreateMedicalRecord
}

// NewMockCreateMedicalRecord creates a new mock instance
func NewMockCreateMedicalRecord(ctrl *gomock.Controller) *MockCreateMedicalRecord {
	mock := &MockCreateMedicalRecord{ctrl: ctrl}
	mock.recorder = &MockCreateMedicalRecordMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCreateMedicalRecord) EXPECT() *MockCreateMedicalRecordMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockCreateMedicalRecord) Create(ctx context.Context, record *entity.MedicalRecord) *entity.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, record)
	ret0, _ := ret[0].(*entity.Error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockCreateMedicalRecordMockRecorder) Create(ctx, record interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCreateMedicalRecord)(nil).Create), ctx, record)
}