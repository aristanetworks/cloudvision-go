// Copyright (c) 2023 Arista Networks, Inc.  All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: arista/lifecycle.v1/lifecycle.proto

package lifecycle

import (
	fmp "github.com/aristanetworks/cloudvision-go/api/fmp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// DeviceLifecycleSummaryKey is the key type for
// DeviceLifecycleSummary model
type DeviceLifecycleSummaryKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// device_id is the device ID
	DeviceId *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
}

func (x *DeviceLifecycleSummaryKey) Reset() {
	*x = DeviceLifecycleSummaryKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceLifecycleSummaryKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceLifecycleSummaryKey) ProtoMessage() {}

func (x *DeviceLifecycleSummaryKey) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceLifecycleSummaryKey.ProtoReflect.Descriptor instead.
func (*DeviceLifecycleSummaryKey) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_lifecycle_proto_rawDescGZIP(), []int{0}
}

func (x *DeviceLifecycleSummaryKey) GetDeviceId() *wrapperspb.StringValue {
	if x != nil {
		return x.DeviceId
	}
	return nil
}

// SoftwareEOL represents a software end of life
type SoftwareEOL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// version of a SoftwareEOL
	Version *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	// end_of_support of a SoftwareEOL
	EndOfSupport *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end_of_support,json=endOfSupport,proto3" json:"end_of_support,omitempty"`
}

func (x *SoftwareEOL) Reset() {
	*x = SoftwareEOL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SoftwareEOL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SoftwareEOL) ProtoMessage() {}

func (x *SoftwareEOL) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SoftwareEOL.ProtoReflect.Descriptor instead.
func (*SoftwareEOL) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_lifecycle_proto_rawDescGZIP(), []int{1}
}

func (x *SoftwareEOL) GetVersion() *wrapperspb.StringValue {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *SoftwareEOL) GetEndOfSupport() *timestamppb.Timestamp {
	if x != nil {
		return x.EndOfSupport
	}
	return nil
}

// DateAndModels has an "end of" date along with
// the models that has this exact "end of" date
type DateAndModels struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// "end of" date
	Date *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	// models with this exact "end of" date
	// mapped to its count
	Models *fmp.MapStringInt32 `protobuf:"bytes,2,opt,name=models,proto3" json:"models,omitempty"`
}

func (x *DateAndModels) Reset() {
	*x = DateAndModels{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DateAndModels) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DateAndModels) ProtoMessage() {}

func (x *DateAndModels) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DateAndModels.ProtoReflect.Descriptor instead.
func (*DateAndModels) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_lifecycle_proto_rawDescGZIP(), []int{2}
}

func (x *DateAndModels) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *DateAndModels) GetModels() *fmp.MapStringInt32 {
	if x != nil {
		return x.Models
	}
	return nil
}

// HardwareLifecycleSummary represents a hardware lifecycle summary
type HardwareLifecycleSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// end_of_life of a HardwareLifecycleSummary
	EndOfLife *DateAndModels `protobuf:"bytes,1,opt,name=end_of_life,json=endOfLife,proto3" json:"end_of_life,omitempty"`
	// end_of_sale of a HardwareLifecycleSummary
	EndOfSale *DateAndModels `protobuf:"bytes,2,opt,name=end_of_sale,json=endOfSale,proto3" json:"end_of_sale,omitempty"`
	// end_of_tac_support of a HardwareLifecycleSummary
	EndOfTacSupport *DateAndModels `protobuf:"bytes,3,opt,name=end_of_tac_support,json=endOfTacSupport,proto3" json:"end_of_tac_support,omitempty"`
	// end_of_hardware_rma_requests of a HardwareLifecycleSummary
	EndOfHardwareRmaRequests *DateAndModels `protobuf:"bytes,4,opt,name=end_of_hardware_rma_requests,json=endOfHardwareRmaRequests,proto3" json:"end_of_hardware_rma_requests,omitempty"`
}

func (x *HardwareLifecycleSummary) Reset() {
	*x = HardwareLifecycleSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HardwareLifecycleSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HardwareLifecycleSummary) ProtoMessage() {}

func (x *HardwareLifecycleSummary) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HardwareLifecycleSummary.ProtoReflect.Descriptor instead.
func (*HardwareLifecycleSummary) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_lifecycle_proto_rawDescGZIP(), []int{3}
}

func (x *HardwareLifecycleSummary) GetEndOfLife() *DateAndModels {
	if x != nil {
		return x.EndOfLife
	}
	return nil
}

func (x *HardwareLifecycleSummary) GetEndOfSale() *DateAndModels {
	if x != nil {
		return x.EndOfSale
	}
	return nil
}

func (x *HardwareLifecycleSummary) GetEndOfTacSupport() *DateAndModels {
	if x != nil {
		return x.EndOfTacSupport
	}
	return nil
}

func (x *HardwareLifecycleSummary) GetEndOfHardwareRmaRequests() *DateAndModels {
	if x != nil {
		return x.EndOfHardwareRmaRequests
	}
	return nil
}

// DeviceLifecycleSummary is the state model that represents
// the lifecycle summary of a device
type DeviceLifecycleSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// DeviceLifecycleSummaryKey is the key of
	// DeviceLifecycleSummary
	Key *DeviceLifecycleSummaryKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// software_eol is the software end of life of
	// a device
	SoftwareEol *SoftwareEOL `protobuf:"bytes,2,opt,name=software_eol,json=softwareEol,proto3" json:"software_eol,omitempty"`
	// hardware_lifecycle_summary is the hardware lifecycle summary
	// of a device
	HardwareLifecycleSummary *HardwareLifecycleSummary `protobuf:"bytes,3,opt,name=hardware_lifecycle_summary,json=hardwareLifecycleSummary,proto3" json:"hardware_lifecycle_summary,omitempty"`
}

func (x *DeviceLifecycleSummary) Reset() {
	*x = DeviceLifecycleSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceLifecycleSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceLifecycleSummary) ProtoMessage() {}

func (x *DeviceLifecycleSummary) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_lifecycle_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceLifecycleSummary.ProtoReflect.Descriptor instead.
func (*DeviceLifecycleSummary) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_lifecycle_proto_rawDescGZIP(), []int{4}
}

func (x *DeviceLifecycleSummary) GetKey() *DeviceLifecycleSummaryKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *DeviceLifecycleSummary) GetSoftwareEol() *SoftwareEOL {
	if x != nil {
		return x.SoftwareEol
	}
	return nil
}

func (x *DeviceLifecycleSummary) GetHardwareLifecycleSummary() *HardwareLifecycleSummary {
	if x != nil {
		return x.HardwareLifecycleSummary
	}
	return nil
}

var File_arista_lifecycle_v1_lifecycle_proto protoreflect.FileDescriptor

var file_arista_lifecycle_v1_lifecycle_proto_rawDesc = []byte{
	0x0a, 0x23, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69,
	0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x66, 0x6d, 0x70,
	0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x12, 0x66, 0x6d, 0x70, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x19, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c,
	0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x4b,
	0x65, 0x79, 0x12, 0x39, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x3a, 0x04, 0x80,
	0x8e, 0x19, 0x01, 0x22, 0x87, 0x01, 0x0a, 0x0b, 0x53, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65,
	0x45, 0x4f, 0x4c, 0x12, 0x36, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a, 0x0e, 0x65,
	0x6e, 0x64, 0x5f, 0x6f, 0x66, 0x5f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0c, 0x65, 0x6e, 0x64, 0x4f, 0x66, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x6c, 0x0a,
	0x0d, 0x44, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x12, 0x2e,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x2b,
	0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x66, 0x6d, 0x70, 0x2e, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x49, 0x6e,
	0x74, 0x33, 0x32, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x22, 0xd7, 0x02, 0x0a, 0x18,
	0x48, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c,
	0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x42, 0x0a, 0x0b, 0x65, 0x6e, 0x64, 0x5f,
	0x6f, 0x66, 0x5f, 0x6c, 0x69, 0x66, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e,
	0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x4f, 0x66, 0x4c, 0x69, 0x66, 0x65, 0x12, 0x42, 0x0a, 0x0b,
	0x65, 0x6e, 0x64, 0x5f, 0x6f, 0x66, 0x5f, 0x73, 0x61, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63,
	0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x4f, 0x66, 0x53, 0x61, 0x6c, 0x65,
	0x12, 0x4f, 0x0a, 0x12, 0x65, 0x6e, 0x64, 0x5f, 0x6f, 0x66, 0x5f, 0x74, 0x61, 0x63, 0x5f, 0x73,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x52, 0x0f, 0x65, 0x6e, 0x64, 0x4f, 0x66, 0x54, 0x61, 0x63, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x62, 0x0a, 0x1c, 0x65, 0x6e, 0x64, 0x5f, 0x6f, 0x66, 0x5f, 0x68, 0x61, 0x72, 0x64,
	0x77, 0x61, 0x72, 0x65, 0x5f, 0x72, 0x6d, 0x61, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61,
	0x74, 0x65, 0x41, 0x6e, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x52, 0x18, 0x65, 0x6e, 0x64,
	0x4f, 0x66, 0x48, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x52, 0x6d, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x73, 0x22, 0x94, 0x02, 0x0a, 0x16, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x12, 0x40, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e,
	0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79,
	0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x43, 0x0a, 0x0c, 0x73, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65, 0x5f, 0x65,
	0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65, 0x45, 0x4f, 0x4c, 0x52, 0x0b, 0x73, 0x6f, 0x66, 0x74,
	0x77, 0x61, 0x72, 0x65, 0x45, 0x6f, 0x6c, 0x12, 0x6b, 0x0a, 0x1a, 0x68, 0x61, 0x72, 0x64, 0x77,
	0x61, 0x72, 0x65, 0x5f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x5f, 0x73, 0x75,
	0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x61, 0x72,
	0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x48, 0x61, 0x72, 0x64, 0x77, 0x61, 0x72, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79,
	0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x18, 0x68, 0x61, 0x72, 0x64,
	0x77, 0x61, 0x72, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x6f, 0x42, 0x72, 0x0a, 0x17,
	0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63,
	0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x50, 0x01, 0x5a, 0x4a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2f,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x2d, 0x67, 0x6f, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79,
	0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x3b, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arista_lifecycle_v1_lifecycle_proto_rawDescOnce sync.Once
	file_arista_lifecycle_v1_lifecycle_proto_rawDescData = file_arista_lifecycle_v1_lifecycle_proto_rawDesc
)

func file_arista_lifecycle_v1_lifecycle_proto_rawDescGZIP() []byte {
	file_arista_lifecycle_v1_lifecycle_proto_rawDescOnce.Do(func() {
		file_arista_lifecycle_v1_lifecycle_proto_rawDescData = protoimpl.X.CompressGZIP(file_arista_lifecycle_v1_lifecycle_proto_rawDescData)
	})
	return file_arista_lifecycle_v1_lifecycle_proto_rawDescData
}

var file_arista_lifecycle_v1_lifecycle_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_arista_lifecycle_v1_lifecycle_proto_goTypes = []interface{}{
	(*DeviceLifecycleSummaryKey)(nil), // 0: arista.lifecycle.v1.DeviceLifecycleSummaryKey
	(*SoftwareEOL)(nil),               // 1: arista.lifecycle.v1.SoftwareEOL
	(*DateAndModels)(nil),             // 2: arista.lifecycle.v1.DateAndModels
	(*HardwareLifecycleSummary)(nil),  // 3: arista.lifecycle.v1.HardwareLifecycleSummary
	(*DeviceLifecycleSummary)(nil),    // 4: arista.lifecycle.v1.DeviceLifecycleSummary
	(*wrapperspb.StringValue)(nil),    // 5: google.protobuf.StringValue
	(*timestamppb.Timestamp)(nil),     // 6: google.protobuf.Timestamp
	(*fmp.MapStringInt32)(nil),        // 7: fmp.MapStringInt32
}
var file_arista_lifecycle_v1_lifecycle_proto_depIdxs = []int32{
	5,  // 0: arista.lifecycle.v1.DeviceLifecycleSummaryKey.device_id:type_name -> google.protobuf.StringValue
	5,  // 1: arista.lifecycle.v1.SoftwareEOL.version:type_name -> google.protobuf.StringValue
	6,  // 2: arista.lifecycle.v1.SoftwareEOL.end_of_support:type_name -> google.protobuf.Timestamp
	6,  // 3: arista.lifecycle.v1.DateAndModels.date:type_name -> google.protobuf.Timestamp
	7,  // 4: arista.lifecycle.v1.DateAndModels.models:type_name -> fmp.MapStringInt32
	2,  // 5: arista.lifecycle.v1.HardwareLifecycleSummary.end_of_life:type_name -> arista.lifecycle.v1.DateAndModels
	2,  // 6: arista.lifecycle.v1.HardwareLifecycleSummary.end_of_sale:type_name -> arista.lifecycle.v1.DateAndModels
	2,  // 7: arista.lifecycle.v1.HardwareLifecycleSummary.end_of_tac_support:type_name -> arista.lifecycle.v1.DateAndModels
	2,  // 8: arista.lifecycle.v1.HardwareLifecycleSummary.end_of_hardware_rma_requests:type_name -> arista.lifecycle.v1.DateAndModels
	0,  // 9: arista.lifecycle.v1.DeviceLifecycleSummary.key:type_name -> arista.lifecycle.v1.DeviceLifecycleSummaryKey
	1,  // 10: arista.lifecycle.v1.DeviceLifecycleSummary.software_eol:type_name -> arista.lifecycle.v1.SoftwareEOL
	3,  // 11: arista.lifecycle.v1.DeviceLifecycleSummary.hardware_lifecycle_summary:type_name -> arista.lifecycle.v1.HardwareLifecycleSummary
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_arista_lifecycle_v1_lifecycle_proto_init() }
func file_arista_lifecycle_v1_lifecycle_proto_init() {
	if File_arista_lifecycle_v1_lifecycle_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arista_lifecycle_v1_lifecycle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceLifecycleSummaryKey); i {
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
		file_arista_lifecycle_v1_lifecycle_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SoftwareEOL); i {
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
		file_arista_lifecycle_v1_lifecycle_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DateAndModels); i {
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
		file_arista_lifecycle_v1_lifecycle_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HardwareLifecycleSummary); i {
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
		file_arista_lifecycle_v1_lifecycle_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceLifecycleSummary); i {
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
			RawDescriptor: file_arista_lifecycle_v1_lifecycle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arista_lifecycle_v1_lifecycle_proto_goTypes,
		DependencyIndexes: file_arista_lifecycle_v1_lifecycle_proto_depIdxs,
		MessageInfos:      file_arista_lifecycle_v1_lifecycle_proto_msgTypes,
	}.Build()
	File_arista_lifecycle_v1_lifecycle_proto = out.File
	file_arista_lifecycle_v1_lifecycle_proto_rawDesc = nil
	file_arista_lifecycle_v1_lifecycle_proto_goTypes = nil
	file_arista_lifecycle_v1_lifecycle_proto_depIdxs = nil
}
