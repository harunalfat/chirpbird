// Code generated by MockGen. DO NOT EDIT.
// Source: presentation/persistence/mongodb_channel_repository.go

// Package mock_persistence is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/harunalfat/chirpbird/backend/entities"
)

// MockChannelRepository is a mock of ChannelRepository interface.
type MockChannelRepository struct {
	ctrl     *gomock.Controller
	recorder *MockChannelRepositoryMockRecorder
}

// MockChannelRepositoryMockRecorder is the mock recorder for MockChannelRepository.
type MockChannelRepositoryMockRecorder struct {
	mock *MockChannelRepository
}

// NewMockChannelRepository creates a new mock instance.
func NewMockChannelRepository(ctrl *gomock.Controller) *MockChannelRepository {
	mock := &MockChannelRepository{ctrl: ctrl}
	mock.recorder = &MockChannelRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChannelRepository) EXPECT() *MockChannelRepositoryMockRecorder {
	return m.recorder
}

// Fetch mocks base method.
func (m *MockChannelRepository) Fetch(ctx context.Context, channelID string) (entities.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", ctx, channelID)
	ret0, _ := ret[0].(entities.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockChannelRepositoryMockRecorder) Fetch(ctx, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockChannelRepository)(nil).Fetch), ctx, channelID)
}

// FetchByName mocks base method.
func (m *MockChannelRepository) FetchByName(ctx context.Context, channelName string) (entities.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchByName", ctx, channelName)
	ret0, _ := ret[0].(entities.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchByName indicates an expected call of FetchByName.
func (mr *MockChannelRepositoryMockRecorder) FetchByName(ctx, channelName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchByName", reflect.TypeOf((*MockChannelRepository)(nil).FetchByName), ctx, channelName)
}

// Insert mocks base method.
func (m *MockChannelRepository) Insert(arg0 context.Context, arg1 entities.Channel) (entities.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(entities.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockChannelRepositoryMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockChannelRepository)(nil).Insert), arg0, arg1)
}

// Update mocks base method.
func (m *MockChannelRepository) Update(ctx context.Context, channelID string, updated entities.Channel) (entities.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, channelID, updated)
	ret0, _ := ret[0].(entities.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockChannelRepositoryMockRecorder) Update(ctx, channelID, updated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockChannelRepository)(nil).Update), ctx, channelID, updated)
}
