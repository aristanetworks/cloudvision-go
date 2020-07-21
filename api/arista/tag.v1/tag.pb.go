// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: arista/tag.v1/tag.proto

package tag

import (
	_ "github.com/aristanetworks/cloudvision-go/api/fmp"
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
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

// CreatorType specifies an entity that creates something.
type CreatorType int32

const (
	CreatorType_CREATOR_TYPE_UNSPECIFIED CreatorType = 0
	// CREATOR_TYPE_SYSTEM is the type for something created by the system.
	CreatorType_CREATOR_TYPE_SYSTEM CreatorType = 1
	// CREATOR_TYPE_USER is the type for something created by a user.
	CreatorType_CREATOR_TYPE_USER CreatorType = 2
)

// Enum value maps for CreatorType.
var (
	CreatorType_name = map[int32]string{
		0: "CREATOR_TYPE_UNSPECIFIED",
		1: "CREATOR_TYPE_SYSTEM",
		2: "CREATOR_TYPE_USER",
	}
	CreatorType_value = map[string]int32{
		"CREATOR_TYPE_UNSPECIFIED": 0,
		"CREATOR_TYPE_SYSTEM":      1,
		"CREATOR_TYPE_USER":        2,
	}
)

func (x CreatorType) Enum() *CreatorType {
	p := new(CreatorType)
	*p = x
	return p
}

func (x CreatorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CreatorType) Descriptor() protoreflect.EnumDescriptor {
	return file_arista_tag_v1_tag_proto_enumTypes[0].Descriptor()
}

func (CreatorType) Type() protoreflect.EnumType {
	return &file_arista_tag_v1_tag_proto_enumTypes[0]
}

func (x CreatorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CreatorType.Descriptor instead.
func (CreatorType) EnumDescriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{0}
}

// TagKey uniquely identifies a tag for a network element.
type TagKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Label is the label of the tag.
	Label *wrappers.StringValue `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	// Value is the value of the tag.
	Value *wrappers.StringValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *TagKey) Reset() {
	*x = TagKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TagKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TagKey) ProtoMessage() {}

func (x *TagKey) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TagKey.ProtoReflect.Descriptor instead.
func (*TagKey) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{0}
}

func (x *TagKey) GetLabel() *wrappers.StringValue {
	if x != nil {
		return x.Label
	}
	return nil
}

func (x *TagKey) GetValue() *wrappers.StringValue {
	if x != nil {
		return x.Value
	}
	return nil
}

// InterfaceTagConfig is a label-value pair that may or may
// not be assigned to an interface.
type InterfaceTagConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies the interface tag.
	Key *TagKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *InterfaceTagConfig) Reset() {
	*x = InterfaceTagConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InterfaceTagConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InterfaceTagConfig) ProtoMessage() {}

func (x *InterfaceTagConfig) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InterfaceTagConfig.ProtoReflect.Descriptor instead.
func (*InterfaceTagConfig) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{1}
}

func (x *InterfaceTagConfig) GetKey() *TagKey {
	if x != nil {
		return x.Key
	}
	return nil
}

// InterfaceTag is a label-value pair that may or may
// not be assigned to an interface.
type InterfaceTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies the interface tag.
	Key *TagKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// CreatorType is the creator type of the tag.
	CreatorType CreatorType `protobuf:"varint,2,opt,name=creator_type,json=creatorType,proto3,enum=arista.tag.v1.CreatorType" json:"creator_type,omitempty"`
}

func (x *InterfaceTag) Reset() {
	*x = InterfaceTag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InterfaceTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InterfaceTag) ProtoMessage() {}

func (x *InterfaceTag) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InterfaceTag.ProtoReflect.Descriptor instead.
func (*InterfaceTag) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{2}
}

func (x *InterfaceTag) GetKey() *TagKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *InterfaceTag) GetCreatorType() CreatorType {
	if x != nil {
		return x.CreatorType
	}
	return CreatorType_CREATOR_TYPE_UNSPECIFIED
}

// InterfaceTagAssignmentKey uniquely identifies an interface
// tag assignment.
type InterfaceTagAssignmentKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Label is the label of the tag.
	Label *wrappers.StringValue `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	// Value is the value of the tag.
	Value *wrappers.StringValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// DeviceId is the ID of the interface's device.
	DeviceId *wrappers.StringValue `protobuf:"bytes,3,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	// InterfaceId is the ID of the interface.
	InterfaceId *wrappers.StringValue `protobuf:"bytes,4,opt,name=interface_id,json=interfaceId,proto3" json:"interface_id,omitempty"`
}

func (x *InterfaceTagAssignmentKey) Reset() {
	*x = InterfaceTagAssignmentKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InterfaceTagAssignmentKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InterfaceTagAssignmentKey) ProtoMessage() {}

func (x *InterfaceTagAssignmentKey) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InterfaceTagAssignmentKey.ProtoReflect.Descriptor instead.
func (*InterfaceTagAssignmentKey) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{3}
}

func (x *InterfaceTagAssignmentKey) GetLabel() *wrappers.StringValue {
	if x != nil {
		return x.Label
	}
	return nil
}

func (x *InterfaceTagAssignmentKey) GetValue() *wrappers.StringValue {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *InterfaceTagAssignmentKey) GetDeviceId() *wrappers.StringValue {
	if x != nil {
		return x.DeviceId
	}
	return nil
}

func (x *InterfaceTagAssignmentKey) GetInterfaceId() *wrappers.StringValue {
	if x != nil {
		return x.InterfaceId
	}
	return nil
}

// InterfaceTagAssignmentConfig is the assignment of an interface tag
// to a specific interface.
type InterfaceTagAssignmentConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies the interface tag assignment.
	Key *InterfaceTagAssignmentKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *InterfaceTagAssignmentConfig) Reset() {
	*x = InterfaceTagAssignmentConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InterfaceTagAssignmentConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InterfaceTagAssignmentConfig) ProtoMessage() {}

func (x *InterfaceTagAssignmentConfig) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InterfaceTagAssignmentConfig.ProtoReflect.Descriptor instead.
func (*InterfaceTagAssignmentConfig) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{4}
}

func (x *InterfaceTagAssignmentConfig) GetKey() *InterfaceTagAssignmentKey {
	if x != nil {
		return x.Key
	}
	return nil
}

// DeviceTagConfig is a label-value pair that may or may not
// be assigned to a device.
type DeviceTagConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies the device tag.
	Key *TagKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeviceTagConfig) Reset() {
	*x = DeviceTagConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceTagConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceTagConfig) ProtoMessage() {}

func (x *DeviceTagConfig) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceTagConfig.ProtoReflect.Descriptor instead.
func (*DeviceTagConfig) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{5}
}

func (x *DeviceTagConfig) GetKey() *TagKey {
	if x != nil {
		return x.Key
	}
	return nil
}

// DeviceTag is a label-value pair that may or may not
// be assigned to a device.
type DeviceTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies the device tag.
	Key *TagKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// CreatorType is the creator type of the tag.
	CreatorType CreatorType `protobuf:"varint,2,opt,name=creator_type,json=creatorType,proto3,enum=arista.tag.v1.CreatorType" json:"creator_type,omitempty"`
}

func (x *DeviceTag) Reset() {
	*x = DeviceTag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceTag) ProtoMessage() {}

func (x *DeviceTag) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceTag.ProtoReflect.Descriptor instead.
func (*DeviceTag) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{6}
}

func (x *DeviceTag) GetKey() *TagKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *DeviceTag) GetCreatorType() CreatorType {
	if x != nil {
		return x.CreatorType
	}
	return CreatorType_CREATOR_TYPE_UNSPECIFIED
}

// DeviceTagAssignmentKey uniquely identifies a device tag
// assignment.
type DeviceTagAssignmentKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Label is the label of the tag.
	Label *wrappers.StringValue `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	// Value is the value of the tag.
	Value *wrappers.StringValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// DeviceId is the ID of the device.
	DeviceId *wrappers.StringValue `protobuf:"bytes,3,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
}

func (x *DeviceTagAssignmentKey) Reset() {
	*x = DeviceTagAssignmentKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceTagAssignmentKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceTagAssignmentKey) ProtoMessage() {}

func (x *DeviceTagAssignmentKey) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceTagAssignmentKey.ProtoReflect.Descriptor instead.
func (*DeviceTagAssignmentKey) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{7}
}

func (x *DeviceTagAssignmentKey) GetLabel() *wrappers.StringValue {
	if x != nil {
		return x.Label
	}
	return nil
}

func (x *DeviceTagAssignmentKey) GetValue() *wrappers.StringValue {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *DeviceTagAssignmentKey) GetDeviceId() *wrappers.StringValue {
	if x != nil {
		return x.DeviceId
	}
	return nil
}

// DeviceTagAssignmentConfig is the assignment of a device tag to a
// specific device.
type DeviceTagAssignmentConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies the device tag assignment.
	Key *DeviceTagAssignmentKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeviceTagAssignmentConfig) Reset() {
	*x = DeviceTagAssignmentConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_tag_v1_tag_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceTagAssignmentConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceTagAssignmentConfig) ProtoMessage() {}

func (x *DeviceTagAssignmentConfig) ProtoReflect() protoreflect.Message {
	mi := &file_arista_tag_v1_tag_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceTagAssignmentConfig.ProtoReflect.Descriptor instead.
func (*DeviceTagAssignmentConfig) Descriptor() ([]byte, []int) {
	return file_arista_tag_v1_tag_proto_rawDescGZIP(), []int{8}
}

func (x *DeviceTagAssignmentConfig) GetKey() *DeviceTagAssignmentKey {
	if x != nil {
		return x.Key
	}
	return nil
}

var File_arista_tag_v1_tag_proto protoreflect.FileDescriptor

var file_arista_tag_v1_tag_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x2f,
	0x74, 0x61, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x66, 0x6d, 0x70, 0x2f, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x76,
	0x0a, 0x06, 0x54, 0x61, 0x67, 0x4b, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x32, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x04, 0x80, 0x8e, 0x19, 0x01, 0x22, 0x45, 0x0a, 0x12, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x54, 0x61, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x27, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x67, 0x4b, 0x65, 0x79,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x77, 0x22, 0x7e, 0x0a,
	0x0c, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x54, 0x61, 0x67, 0x12, 0x27, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x67, 0x4b, 0x65,
	0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3d, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x54, 0x79, 0x70, 0x65, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x6f, 0x22, 0x85, 0x02,
	0x0a, 0x19, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x54, 0x61, 0x67, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12,
	0x32, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x39, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x3f,
	0x0a, 0x0c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x0b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x49, 0x64, 0x3a,
	0x04, 0x80, 0x8e, 0x19, 0x01, 0x22, 0x62, 0x0a, 0x1c, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61,
	0x63, 0x65, 0x54, 0x61, 0x67, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3a, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x74, 0x61, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x54, 0x61, 0x67, 0x41,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x77, 0x22, 0x42, 0x0a, 0x0f, 0x44, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x54, 0x61, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x27, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x67, 0x4b, 0x65, 0x79,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x77, 0x22, 0x7b, 0x0a,
	0x09, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x61, 0x67, 0x12, 0x27, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x67, 0x4b, 0x65, 0x79, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x3d, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x54, 0x79,
	0x70, 0x65, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x6f, 0x22, 0xc1, 0x01, 0x0a, 0x16, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x61, 0x67, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x32, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x39, 0x0a,
	0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x3a, 0x04, 0x80, 0x8e, 0x19, 0x01, 0x22, 0x5c,
	0x0a, 0x19, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x61, 0x67, 0x41, 0x73, 0x73, 0x69, 0x67,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x37, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54,
	0x61, 0x67, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x77, 0x2a, 0x5b, 0x0a, 0x0b,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x18, 0x43,
	0x52, 0x45, 0x41, 0x54, 0x4f, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x52, 0x45,
	0x41, 0x54, 0x4f, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d,
	0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x52, 0x45, 0x41, 0x54, 0x4f, 0x52, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x10, 0x02, 0x42, 0x24, 0x5a, 0x22, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x72,
	0x69, 0x73, 0x74, 0x61, 0x2f, 0x74, 0x61, 0x67, 0x2e, 0x76, 0x31, 0x3b, 0x74, 0x61, 0x67, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arista_tag_v1_tag_proto_rawDescOnce sync.Once
	file_arista_tag_v1_tag_proto_rawDescData = file_arista_tag_v1_tag_proto_rawDesc
)

func file_arista_tag_v1_tag_proto_rawDescGZIP() []byte {
	file_arista_tag_v1_tag_proto_rawDescOnce.Do(func() {
		file_arista_tag_v1_tag_proto_rawDescData = protoimpl.X.CompressGZIP(file_arista_tag_v1_tag_proto_rawDescData)
	})
	return file_arista_tag_v1_tag_proto_rawDescData
}

var file_arista_tag_v1_tag_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_arista_tag_v1_tag_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_arista_tag_v1_tag_proto_goTypes = []interface{}{
	(CreatorType)(0),                     // 0: arista.tag.v1.CreatorType
	(*TagKey)(nil),                       // 1: arista.tag.v1.TagKey
	(*InterfaceTagConfig)(nil),           // 2: arista.tag.v1.InterfaceTagConfig
	(*InterfaceTag)(nil),                 // 3: arista.tag.v1.InterfaceTag
	(*InterfaceTagAssignmentKey)(nil),    // 4: arista.tag.v1.InterfaceTagAssignmentKey
	(*InterfaceTagAssignmentConfig)(nil), // 5: arista.tag.v1.InterfaceTagAssignmentConfig
	(*DeviceTagConfig)(nil),              // 6: arista.tag.v1.DeviceTagConfig
	(*DeviceTag)(nil),                    // 7: arista.tag.v1.DeviceTag
	(*DeviceTagAssignmentKey)(nil),       // 8: arista.tag.v1.DeviceTagAssignmentKey
	(*DeviceTagAssignmentConfig)(nil),    // 9: arista.tag.v1.DeviceTagAssignmentConfig
	(*wrappers.StringValue)(nil),         // 10: google.protobuf.StringValue
}
var file_arista_tag_v1_tag_proto_depIdxs = []int32{
	10, // 0: arista.tag.v1.TagKey.label:type_name -> google.protobuf.StringValue
	10, // 1: arista.tag.v1.TagKey.value:type_name -> google.protobuf.StringValue
	1,  // 2: arista.tag.v1.InterfaceTagConfig.key:type_name -> arista.tag.v1.TagKey
	1,  // 3: arista.tag.v1.InterfaceTag.key:type_name -> arista.tag.v1.TagKey
	0,  // 4: arista.tag.v1.InterfaceTag.creator_type:type_name -> arista.tag.v1.CreatorType
	10, // 5: arista.tag.v1.InterfaceTagAssignmentKey.label:type_name -> google.protobuf.StringValue
	10, // 6: arista.tag.v1.InterfaceTagAssignmentKey.value:type_name -> google.protobuf.StringValue
	10, // 7: arista.tag.v1.InterfaceTagAssignmentKey.device_id:type_name -> google.protobuf.StringValue
	10, // 8: arista.tag.v1.InterfaceTagAssignmentKey.interface_id:type_name -> google.protobuf.StringValue
	4,  // 9: arista.tag.v1.InterfaceTagAssignmentConfig.key:type_name -> arista.tag.v1.InterfaceTagAssignmentKey
	1,  // 10: arista.tag.v1.DeviceTagConfig.key:type_name -> arista.tag.v1.TagKey
	1,  // 11: arista.tag.v1.DeviceTag.key:type_name -> arista.tag.v1.TagKey
	0,  // 12: arista.tag.v1.DeviceTag.creator_type:type_name -> arista.tag.v1.CreatorType
	10, // 13: arista.tag.v1.DeviceTagAssignmentKey.label:type_name -> google.protobuf.StringValue
	10, // 14: arista.tag.v1.DeviceTagAssignmentKey.value:type_name -> google.protobuf.StringValue
	10, // 15: arista.tag.v1.DeviceTagAssignmentKey.device_id:type_name -> google.protobuf.StringValue
	8,  // 16: arista.tag.v1.DeviceTagAssignmentConfig.key:type_name -> arista.tag.v1.DeviceTagAssignmentKey
	17, // [17:17] is the sub-list for method output_type
	17, // [17:17] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_arista_tag_v1_tag_proto_init() }
func file_arista_tag_v1_tag_proto_init() {
	if File_arista_tag_v1_tag_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arista_tag_v1_tag_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TagKey); i {
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
		file_arista_tag_v1_tag_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InterfaceTagConfig); i {
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
		file_arista_tag_v1_tag_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InterfaceTag); i {
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
		file_arista_tag_v1_tag_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InterfaceTagAssignmentKey); i {
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
		file_arista_tag_v1_tag_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InterfaceTagAssignmentConfig); i {
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
		file_arista_tag_v1_tag_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceTagConfig); i {
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
		file_arista_tag_v1_tag_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceTag); i {
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
		file_arista_tag_v1_tag_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceTagAssignmentKey); i {
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
		file_arista_tag_v1_tag_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceTagAssignmentConfig); i {
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
			RawDescriptor: file_arista_tag_v1_tag_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arista_tag_v1_tag_proto_goTypes,
		DependencyIndexes: file_arista_tag_v1_tag_proto_depIdxs,
		EnumInfos:         file_arista_tag_v1_tag_proto_enumTypes,
		MessageInfos:      file_arista_tag_v1_tag_proto_msgTypes,
	}.Build()
	File_arista_tag_v1_tag_proto = out.File
	file_arista_tag_v1_tag_proto_rawDesc = nil
	file_arista_tag_v1_tag_proto_goTypes = nil
	file_arista_tag_v1_tag_proto_depIdxs = nil
}
