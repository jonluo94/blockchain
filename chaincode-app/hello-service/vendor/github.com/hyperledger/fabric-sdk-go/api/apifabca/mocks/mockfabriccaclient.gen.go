// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hyperledger/fabric-sdk-go/api/apifabca (interfaces: FabricCAClient)

// Package mock_apifabca is a generated GoMock package.
package mock_apifabca

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	apicryptosuite "github.com/hyperledger/fabric-sdk-go/api/apicryptosuite"
	apifabca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	api "github.com/hyperledger/fabric-sdk-go/internal/github.com/hyperledger/fabric-ca/api"
)

// MockFabricCAClient is a mock of FabricCAClient interface
type MockFabricCAClient struct {
	ctrl     *gomock.Controller
	recorder *MockFabricCAClientMockRecorder
}

// MockFabricCAClientMockRecorder is the mock recorder for MockFabricCAClient
type MockFabricCAClientMockRecorder struct {
	mock *MockFabricCAClient
}

// NewMockFabricCAClient creates a new mock instance
func NewMockFabricCAClient(ctrl *gomock.Controller) *MockFabricCAClient {
	mock := &MockFabricCAClient{ctrl: ctrl}
	mock.recorder = &MockFabricCAClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFabricCAClient) EXPECT() *MockFabricCAClientMockRecorder {
	return m.recorder
}

// CAName mocks base method
func (m *MockFabricCAClient) CAName() string {
	ret := m.ctrl.Call(m, "CAName")
	ret0, _ := ret[0].(string)
	return ret0
}

// CAName indicates an expected call of CAName
func (mr *MockFabricCAClientMockRecorder) CAName() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CAName", reflect.TypeOf((*MockFabricCAClient)(nil).CAName))
}

// Enroll mocks base method
func (m *MockFabricCAClient) Enroll(arg0, arg1 string) (apicryptosuite.Key, []byte, error) {
	ret := m.ctrl.Call(m, "Enroll", arg0, arg1)
	ret0, _ := ret[0].(apicryptosuite.Key)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Enroll indicates an expected call of Enroll
func (mr *MockFabricCAClientMockRecorder) Enroll(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enroll", reflect.TypeOf((*MockFabricCAClient)(nil).Enroll), arg0, arg1)
}

// Reenroll mocks base method
func (m *MockFabricCAClient) Reenroll(arg0 apifabca.User) (apicryptosuite.Key, []byte, error) {
	ret := m.ctrl.Call(m, "Reenroll", arg0)
	ret0, _ := ret[0].(apicryptosuite.Key)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Reenroll indicates an expected call of Reenroll
func (mr *MockFabricCAClientMockRecorder) Reenroll(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reenroll", reflect.TypeOf((*MockFabricCAClient)(nil).Reenroll), arg0)
}

// Register mocks base method
func (m *MockFabricCAClient) Register(arg0 apifabca.User, arg1 *apifabca.RegistrationRequest) (string, error) {
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register
func (mr *MockFabricCAClientMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockFabricCAClient)(nil).Register), arg0, arg1)
}

// Revoke mocks base method
func (m *MockFabricCAClient) Revoke(arg0 apifabca.User, arg1 *apifabca.RevocationRequest) (*api.RevocationResponse, error) {
	ret := m.ctrl.Call(m, "Revoke", arg0, arg1)
	ret0, _ := ret[0].(*api.RevocationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Revoke indicates an expected call of Revoke
func (mr *MockFabricCAClientMockRecorder) Revoke(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revoke", reflect.TypeOf((*MockFabricCAClient)(nil).Revoke), arg0, arg1)
}
