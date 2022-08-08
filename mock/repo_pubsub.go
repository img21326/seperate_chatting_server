// Code generated by MockGen. DO NOT EDIT.
// Source: repo/pubsub/interface.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPubSubRepoInterface is a mock of PubSubRepoInterface interface.
type MockPubSubRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPubSubRepoInterfaceMockRecorder
}

// MockPubSubRepoInterfaceMockRecorder is the mock recorder for MockPubSubRepoInterface.
type MockPubSubRepoInterfaceMockRecorder struct {
	mock *MockPubSubRepoInterface
}

// NewMockPubSubRepoInterface creates a new mock instance.
func NewMockPubSubRepoInterface(ctrl *gomock.Controller) *MockPubSubRepoInterface {
	mock := &MockPubSubRepoInterface{ctrl: ctrl}
	mock.recorder = &MockPubSubRepoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubSubRepoInterface) EXPECT() *MockPubSubRepoInterfaceMockRecorder {
	return m.recorder
}

// Pub mocks base method.
func (m *MockPubSubRepoInterface) Pub(ctx context.Context, topic string, message []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pub", ctx, topic, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Pub indicates an expected call of Pub.
func (mr *MockPubSubRepoInterfaceMockRecorder) Pub(ctx, topic, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pub", reflect.TypeOf((*MockPubSubRepoInterface)(nil).Pub), ctx, topic, message)
}

// Sub mocks base method.
func (m *MockPubSubRepoInterface) Sub(ctx context.Context, topic string) func() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sub", ctx, topic)
	ret0, _ := ret[0].(func() ([]byte, error))
	return ret0
}

// Sub indicates an expected call of Sub.
func (mr *MockPubSubRepoInterfaceMockRecorder) Sub(ctx, topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sub", reflect.TypeOf((*MockPubSubRepoInterface)(nil).Sub), ctx, topic)
}