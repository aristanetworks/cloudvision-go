// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: arista/inventory.v1/inventory.proto

package inventory

import (
	_ "github.com/aristanetworks/cloudvision-go/api/fmp"
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// StreamingStatus the status of streaming telemetry for this device.
type StreamingStatus int32

const (
	// Unspecified is the uninitialized state of this enum.
	StreamingStatus_STREAMING_STATUS_UNSPECIFIED StreamingStatus = 0
	// Inactive indicates the device is not streaming telemetry.
	StreamingStatus_STREAMING_STATUS_INACTIVE StreamingStatus = 1
	// Active indicates the device is streaming telemetry.
	StreamingStatus_STREAMING_STATUS_ACTIVE StreamingStatus = 2
)

// Enum value maps for StreamingStatus.
var (
	StreamingStatus_name = map[int32]string{
		0: "STREAMING_STATUS_UNSPECIFIED",
		1: "STREAMING_STATUS_INACTIVE",
		2: "STREAMING_STATUS_ACTIVE",
	}
	StreamingStatus_value = map[string]int32{
		"STREAMING_STATUS_UNSPECIFIED": 0,
		"STREAMING_STATUS_INACTIVE":    1,
		"STREAMING_STATUS_ACTIVE":      2,
	}
)

func (x StreamingStatus) Enum() *StreamingStatus {
	p := new(StreamingStatus)
	*p = x
	return p
}

func (x StreamingStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StreamingStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_arista_inventory_v1_inventory_proto_enumTypes[0].Descriptor()
}

func (StreamingStatus) Type() protoreflect.EnumType {
	return &file_arista_inventory_v1_inventory_proto_enumTypes[0]
}

func (x StreamingStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StreamingStatus.Descriptor instead.
func (StreamingStatus) EnumDescriptor() ([]byte, []int) {
	return file_arista_inventory_v1_inventory_proto_rawDescGZIP(), []int{0}
}

// ExtendedAttributes wraps any additional, potentially non-standard, features
// or attributes the device reports.
type ExtendedAttributes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// FeatureEnabled is a map of feature name to enabled status.
	// If a feature is missing from this map it can be assumed off.
	FeatureEnabled map[string]bool `protobuf:"bytes,1,rep,name=feature_enabled,json=featureEnabled,proto3" json:"feature_enabled,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *ExtendedAttributes) Reset() {
	*x = ExtendedAttributes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_inventory_v1_inventory_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExtendedAttributes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtendedAttributes) ProtoMessage() {}

func (x *ExtendedAttributes) ProtoReflect() protoreflect.Message {
	mi := &file_arista_inventory_v1_inventory_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtendedAttributes.ProtoReflect.Descriptor instead.
func (*ExtendedAttributes) Descriptor() ([]byte, []int) {
	return file_arista_inventory_v1_inventory_proto_rawDescGZIP(), []int{0}
}

func (x *ExtendedAttributes) GetFeatureEnabled() map[string]bool {
	if x != nil {
		return x.FeatureEnabled
	}
	return nil
}

// DeviceKey uniquely identifies a single device.
type DeviceKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId *wrappers.StringValue `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
}

func (x *DeviceKey) Reset() {
	*x = DeviceKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_inventory_v1_inventory_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceKey) ProtoMessage() {}

func (x *DeviceKey) ProtoReflect() protoreflect.Message {
	mi := &file_arista_inventory_v1_inventory_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceKey.ProtoReflect.Descriptor instead.
func (*DeviceKey) Descriptor() ([]byte, []int) {
	return file_arista_inventory_v1_inventory_proto_rawDescGZIP(), []int{1}
}

func (x *DeviceKey) GetDeviceId() *wrappers.StringValue {
	if x != nil {
		return x.DeviceId
	}
	return nil
}

// Device is the primary model for this service.
type Device struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The key that uniquely identifies this device.
	Key *DeviceKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// SoftwareVersion gives the currently running device software version.
	SoftwareVersion *wrappers.StringValue `protobuf:"bytes,2,opt,name=software_version,json=softwareVersion,proto3" json:"software_version,omitempty"`
	// ModelName describes the hardware model of this device.
	ModelName *wrappers.StringValue `protobuf:"bytes,3,opt,name=model_name,json=modelName,proto3" json:"model_name,omitempty"`
	// HardwareREvision describes any revisional data to the model name.
	HardwareRevision *wrappers.StringValue `protobuf:"bytes,4,opt,name=hardware_revision,json=hardwareRevision,proto3" json:"hardware_revision,omitempty"`
	// FQDN gives the fully qualified hostname to reach the device.
	Fqdn *wrappers.StringValue `protobuf:"bytes,10,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
	// Hostname is the hostname as reported on the device.
	Hostname *wrappers.StringValue `protobuf:"bytes,11,opt,name=hostname,proto3" json:"hostname,omitempty"`
	// DomainName provides the domain name the device is registered on.
	DomainName *wrappers.StringValue `protobuf:"bytes,12,opt,name=domain_name,json=domainName,proto3" json:"domain_name,omitempty"`
	// SystemMacAddress provides the MAC address of the management port.
	SystemMacAddress *wrappers.StringValue `protobuf:"bytes,13,opt,name=system_mac_address,json=systemMacAddress,proto3" json:"system_mac_address,omitempty"`
	// BootTime indicates when the device was last booted.
	BootTime *timestamp.Timestamp `protobuf:"bytes,20,opt,name=boot_time,json=bootTime,proto3" json:"boot_time,omitempty"`
	// StreamingStatus the status of streaming telemetry for this device.
	StreamingStatus StreamingStatus `protobuf:"varint,30,opt,name=streaming_status,json=streamingStatus,proto3,enum=arista.inventory.v1.StreamingStatus" json:"streaming_status,omitempty"`
	// ExtendedAttributes wraps any additional, potentially non-standard, features
	// or attributes the device reports.
	ExtendedAttributes *ExtendedAttributes `protobuf:"bytes,31,opt,name=extended_attributes,json=extendedAttributes,proto3" json:"extended_attributes,omitempty"`
}

func (x *Device) Reset() {
	*x = Device{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_inventory_v1_inventory_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Device) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Device) ProtoMessage() {}

func (x *Device) ProtoReflect() protoreflect.Message {
	mi := &file_arista_inventory_v1_inventory_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Device.ProtoReflect.Descriptor instead.
func (*Device) Descriptor() ([]byte, []int) {
	return file_arista_inventory_v1_inventory_proto_rawDescGZIP(), []int{2}
}

func (x *Device) GetKey() *DeviceKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *Device) GetSoftwareVersion() *wrappers.StringValue {
	if x != nil {
		return x.SoftwareVersion
	}
	return nil
}

func (x *Device) GetModelName() *wrappers.StringValue {
	if x != nil {
		return x.ModelName
	}
	return nil
}

func (x *Device) GetHardwareRevision() *wrappers.StringValue {
	if x != nil {
		return x.HardwareRevision
	}
	return nil
}

func (x *Device) GetFqdn() *wrappers.StringValue {
	if x != nil {
		return x.Fqdn
	}
	return nil
}

func (x *Device) GetHostname() *wrappers.StringValue {
	if x != nil {
		return x.Hostname
	}
	return nil
}

func (x *Device) GetDomainName() *wrappers.StringValue {
	if x != nil {
		return x.DomainName
	}
	return nil
}

func (x *Device) GetSystemMacAddress() *wrappers.StringValue {
	if x != nil {
		return x.SystemMacAddress
	}
	return nil
}

func (x *Device) GetBootTime() *timestamp.Timestamp {
	if x != nil {
		return x.BootTime
	}
	return nil
}

func (x *Device) GetStreamingStatus() StreamingStatus {
	if x != nil {
		return x.StreamingStatus
	}
	return StreamingStatus_STREAMING_STATUS_UNSPECIFIED
}

func (x *Device) GetExtendedAttributes() *ExtendedAttributes {
	if x != nil {
		return x.ExtendedAttributes
	}
	return nil
}

var File_arista_inventory_v1_inventory_proto protoreflect.FileDescriptor

var file_arista_inventory_v1_inventory_proto_rawDesc = []byte{
	0x0a, 0x23, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x6e,
	0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x66, 0x6d, 0x70,
	0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xbd, 0x01, 0x0a, 0x12, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x64, 0x0a, 0x0f, 0x66, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x3b, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0e,
	0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x1a, 0x41,
	0x0a, 0x13, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x4c, 0x0a, 0x09, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x39,
	0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x3a, 0x04, 0x80, 0x8e, 0x19, 0x01, 0x22,
	0xee, 0x05, 0x0a, 0x06, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x47, 0x0a, 0x10,
	0x73, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x0f, 0x73, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3b, 0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x49, 0x0a, 0x11, 0x68, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x5f, 0x72,
	0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x10, 0x68, 0x61, 0x72,
	0x64, 0x77, 0x61, 0x72, 0x65, 0x52, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x0a,
	0x04, 0x66, 0x71, 0x64, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x66, 0x71, 0x64, 0x6e, 0x12,
	0x38, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x0b, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x4a, 0x0a, 0x12, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x5f, 0x6d, 0x61, 0x63, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x10, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4d, 0x61, 0x63, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x37, 0x0a, 0x09, 0x62, 0x6f, 0x6f, 0x74, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x08, 0x62, 0x6f, 0x6f, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x4f, 0x0a,
	0x10, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0f, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x58,
	0x0a, 0x13, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x1f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61, 0x72,
	0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x52, 0x12, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x6f,
	0x2a, 0x6f, 0x0a, 0x0f, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x20, 0x0a, 0x1c, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x49, 0x4e, 0x47,
	0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x49,
	0x4e, 0x47, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x41, 0x43, 0x54, 0x49,
	0x56, 0x45, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x49, 0x4e,
	0x47, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10,
	0x02, 0x42, 0x30, 0x5a, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x69, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x3b, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arista_inventory_v1_inventory_proto_rawDescOnce sync.Once
	file_arista_inventory_v1_inventory_proto_rawDescData = file_arista_inventory_v1_inventory_proto_rawDesc
)

func file_arista_inventory_v1_inventory_proto_rawDescGZIP() []byte {
	file_arista_inventory_v1_inventory_proto_rawDescOnce.Do(func() {
		file_arista_inventory_v1_inventory_proto_rawDescData = protoimpl.X.CompressGZIP(file_arista_inventory_v1_inventory_proto_rawDescData)
	})
	return file_arista_inventory_v1_inventory_proto_rawDescData
}

var file_arista_inventory_v1_inventory_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_arista_inventory_v1_inventory_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_arista_inventory_v1_inventory_proto_goTypes = []interface{}{
	(StreamingStatus)(0),         // 0: arista.inventory.v1.StreamingStatus
	(*ExtendedAttributes)(nil),   // 1: arista.inventory.v1.ExtendedAttributes
	(*DeviceKey)(nil),            // 2: arista.inventory.v1.DeviceKey
	(*Device)(nil),               // 3: arista.inventory.v1.Device
	nil,                          // 4: arista.inventory.v1.ExtendedAttributes.FeatureEnabledEntry
	(*wrappers.StringValue)(nil), // 5: google.protobuf.StringValue
	(*timestamp.Timestamp)(nil),  // 6: google.protobuf.Timestamp
}
var file_arista_inventory_v1_inventory_proto_depIdxs = []int32{
	4,  // 0: arista.inventory.v1.ExtendedAttributes.feature_enabled:type_name -> arista.inventory.v1.ExtendedAttributes.FeatureEnabledEntry
	5,  // 1: arista.inventory.v1.DeviceKey.device_id:type_name -> google.protobuf.StringValue
	2,  // 2: arista.inventory.v1.Device.key:type_name -> arista.inventory.v1.DeviceKey
	5,  // 3: arista.inventory.v1.Device.software_version:type_name -> google.protobuf.StringValue
	5,  // 4: arista.inventory.v1.Device.model_name:type_name -> google.protobuf.StringValue
	5,  // 5: arista.inventory.v1.Device.hardware_revision:type_name -> google.protobuf.StringValue
	5,  // 6: arista.inventory.v1.Device.fqdn:type_name -> google.protobuf.StringValue
	5,  // 7: arista.inventory.v1.Device.hostname:type_name -> google.protobuf.StringValue
	5,  // 8: arista.inventory.v1.Device.domain_name:type_name -> google.protobuf.StringValue
	5,  // 9: arista.inventory.v1.Device.system_mac_address:type_name -> google.protobuf.StringValue
	6,  // 10: arista.inventory.v1.Device.boot_time:type_name -> google.protobuf.Timestamp
	0,  // 11: arista.inventory.v1.Device.streaming_status:type_name -> arista.inventory.v1.StreamingStatus
	1,  // 12: arista.inventory.v1.Device.extended_attributes:type_name -> arista.inventory.v1.ExtendedAttributes
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_arista_inventory_v1_inventory_proto_init() }
func file_arista_inventory_v1_inventory_proto_init() {
	if File_arista_inventory_v1_inventory_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arista_inventory_v1_inventory_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExtendedAttributes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_arista_inventory_v1_inventory_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_arista_inventory_v1_inventory_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Device); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_arista_inventory_v1_inventory_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arista_inventory_v1_inventory_proto_goTypes,
		DependencyIndexes: file_arista_inventory_v1_inventory_proto_depIdxs,
		EnumInfos:         file_arista_inventory_v1_inventory_proto_enumTypes,
		MessageInfos:      file_arista_inventory_v1_inventory_proto_msgTypes,
	}.Build()
	File_arista_inventory_v1_inventory_proto = out.File
	file_arista_inventory_v1_inventory_proto_rawDesc = nil
	file_arista_inventory_v1_inventory_proto_goTypes = nil
	file_arista_inventory_v1_inventory_proto_depIdxs = nil
}
