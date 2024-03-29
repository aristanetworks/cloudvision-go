// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aristanetworks/cloudvision-go/device/cvclient (interfaces: CVClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	cvclient "github.com/aristanetworks/cloudvision-go/device/cvclient"
	provider "github.com/aristanetworks/cloudvision-go/provider"
	gomock "github.com/golang/mock/gomock"
	gnmi "github.com/openconfig/gnmi/proto/gnmi"
	grpc "google.golang.org/grpc"
)

// MockCVClient is a mock of CVClient interface.
type MockCVClient struct {
	ctrl     *gomock.Controller
	recorder *MockCVClientMockRecorder
}

// MockCVClientMockRecorder is the mock recorder for MockCVClient.
type MockCVClientMockRecorder struct {
	mock *MockCVClient
}

// NewMockCVClient creates a new mock instance.
func NewMockCVClient(ctrl *gomock.Controller) *MockCVClient {
	mock := &MockCVClient{ctrl: ctrl}
	mock.recorder = &MockCVClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCVClient) EXPECT() *MockCVClientMockRecorder {
	return m.recorder
}

// Capabilities mocks base method.
func (m *MockCVClient) Capabilities(arg0 context.Context, arg1 *gnmi.CapabilityRequest, arg2 ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Capabilities", varargs...)
	ret0, _ := ret[0].(*gnmi.CapabilityResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Capabilities indicates an expected call of Capabilities.
func (mr *MockCVClientMockRecorder) Capabilities(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capabilities", reflect.TypeOf((*MockCVClient)(nil).Capabilities), varargs...)
}

// ForProvider mocks base method.
func (m *MockCVClient) ForProvider(arg0 provider.GNMIProvider) cvclient.CVClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForProvider", arg0)
	ret0, _ := ret[0].(cvclient.CVClient)
	return ret0
}

// ForProvider indicates an expected call of ForProvider.
func (mr *MockCVClientMockRecorder) ForProvider(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForProvider", reflect.TypeOf((*MockCVClient)(nil).ForProvider), arg0)
}

// Get mocks base method.
func (m *MockCVClient) Get(arg0 context.Context, arg1 *gnmi.GetRequest, arg2 ...grpc.CallOption) (*gnmi.GetResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*gnmi.GetResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCVClientMockRecorder) Get(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCVClient)(nil).Get), varargs...)
}

// SendDeviceMetadata mocks base method.
func (m *MockCVClient) SendDeviceMetadata(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendDeviceMetadata", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendDeviceMetadata indicates an expected call of SendDeviceMetadata.
func (mr *MockCVClientMockRecorder) SendDeviceMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendDeviceMetadata", reflect.TypeOf((*MockCVClient)(nil).SendDeviceMetadata), arg0)
}

// SendHeartbeat mocks base method.
func (m *MockCVClient) SendHeartbeat(arg0 context.Context, arg1 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeartbeat", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeartbeat indicates an expected call of SendHeartbeat.
func (mr *MockCVClientMockRecorder) SendHeartbeat(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeartbeat", reflect.TypeOf((*MockCVClient)(nil).SendHeartbeat), arg0, arg1)
}

// Set mocks base method.
func (m *MockCVClient) Set(arg0 context.Context, arg1 *gnmi.SetRequest, arg2 ...grpc.CallOption) (*gnmi.SetResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Set", varargs...)
	ret0, _ := ret[0].(*gnmi.SetResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockCVClientMockRecorder) Set(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCVClient)(nil).Set), varargs...)
}

// SetManagedDevices mocks base method.
func (m *MockCVClient) SetManagedDevices(arg0 []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetManagedDevices", arg0)
}

// SetManagedDevices indicates an expected call of SetManagedDevices.
func (mr *MockCVClientMockRecorder) SetManagedDevices(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetManagedDevices", reflect.TypeOf((*MockCVClient)(nil).SetManagedDevices), arg0)
}

// Subscribe mocks base method.
func (m *MockCVClient) Subscribe(arg0 context.Context, arg1 ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(gnmi.GNMI_SubscribeClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockCVClientMockRecorder) Subscribe(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockCVClient)(nil).Subscribe), varargs...)
}
