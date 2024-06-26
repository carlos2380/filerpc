// Code generated by MockGen. DO NOT EDIT.
// Source: internal/datastore/datastore.go

// Package datastore is a generated GoMock package.
package datastore

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileDataStore is a mock of FileDataStore interface.
type MockFileDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockFileDataStoreMockRecorder
}

// MockFileDataStoreMockRecorder is the mock recorder for MockFileDataStore.
type MockFileDataStoreMockRecorder struct {
	mock *MockFileDataStore
}

// NewMockFileDataStore creates a new mock instance.
func NewMockFileDataStore(ctrl *gomock.Controller) *MockFileDataStore {
	mock := &MockFileDataStore{ctrl: ctrl}
	mock.recorder = &MockFileDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileDataStore) EXPECT() *MockFileDataStoreMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockFileDataStore) Get(ctx context.Context, key string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockFileDataStoreMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockFileDataStore)(nil).Get), ctx, key)
}

// Save mocks base method.
func (m *MockFileDataStore) Save(ctx context.Context, key string, content []byte, hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, key, content, hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockFileDataStoreMockRecorder) Save(ctx, key, content, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockFileDataStore)(nil).Save), ctx, key, content, hash)
}
