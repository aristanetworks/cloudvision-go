// Copyright (c) 2020 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

// Useful types that come from ietf-inet-types.yang

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: fmp/inet.proto

package fmp

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IPAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *IPAddress) Reset() {
	*x = IPAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_inet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPAddress) ProtoMessage() {}

func (x *IPAddress) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_inet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPAddress.ProtoReflect.Descriptor instead.
func (*IPAddress) Descriptor() ([]byte, []int) {
	return file_fmp_inet_proto_rawDescGZIP(), []int{0}
}

func (x *IPAddress) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type IPv4Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *IPv4Address) Reset() {
	*x = IPv4Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_inet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPv4Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPv4Address) ProtoMessage() {}

func (x *IPv4Address) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_inet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPv4Address.ProtoReflect.Descriptor instead.
func (*IPv4Address) Descriptor() ([]byte, []int) {
	return file_fmp_inet_proto_rawDescGZIP(), []int{1}
}

func (x *IPv4Address) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type IPv6Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *IPv6Address) Reset() {
	*x = IPv6Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_inet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPv6Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPv6Address) ProtoMessage() {}

func (x *IPv6Address) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_inet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPv6Address.ProtoReflect.Descriptor instead.
func (*IPv6Address) Descriptor() ([]byte, []int) {
	return file_fmp_inet_proto_rawDescGZIP(), []int{2}
}

func (x *IPv6Address) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type IPPrefix struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *IPPrefix) Reset() {
	*x = IPPrefix{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_inet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPPrefix) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPPrefix) ProtoMessage() {}

func (x *IPPrefix) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_inet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPPrefix.ProtoReflect.Descriptor instead.
func (*IPPrefix) Descriptor() ([]byte, []int) {
	return file_fmp_inet_proto_rawDescGZIP(), []int{3}
}

func (x *IPPrefix) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type IPv4Prefix struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *IPv4Prefix) Reset() {
	*x = IPv4Prefix{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_inet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPv4Prefix) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPv4Prefix) ProtoMessage() {}

func (x *IPv4Prefix) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_inet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPv4Prefix.ProtoReflect.Descriptor instead.
func (*IPv4Prefix) Descriptor() ([]byte, []int) {
	return file_fmp_inet_proto_rawDescGZIP(), []int{4}
}

func (x *IPv4Prefix) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type IPv6Prefix struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *IPv6Prefix) Reset() {
	*x = IPv6Prefix{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_inet_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPv6Prefix) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPv6Prefix) ProtoMessage() {}

func (x *IPv6Prefix) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_inet_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPv6Prefix.ProtoReflect.Descriptor instead.
func (*IPv6Prefix) Descriptor() ([]byte, []int) {
	return file_fmp_inet_proto_rawDescGZIP(), []int{5}
}

func (x *IPv6Prefix) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Port struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value uint32 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Port) Reset() {
	*x = Port{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_inet_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Port) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Port) ProtoMessage() {}

func (x *Port) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_inet_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Port.ProtoReflect.Descriptor instead.
func (*Port) Descriptor() ([]byte, []int) {
	return file_fmp_inet_proto_rawDescGZIP(), []int{6}
}

func (x *Port) GetValue() uint32 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_fmp_inet_proto protoreflect.FileDescriptor

var file_fmp_inet_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x66, 0x6d, 0x70, 0x2f, 0x69, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x03, 0x66, 0x6d, 0x70, 0x22, 0x21, 0x0a, 0x09, 0x49, 0x50, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x49, 0x50, 0x76, 0x34,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x23, 0x0a,
	0x0b, 0x49, 0x50, 0x76, 0x36, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x20, 0x0a, 0x08, 0x49, 0x50, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a, 0x0a, 0x49, 0x50, 0x76, 0x34, 0x50, 0x72, 0x65, 0x66,
	0x69, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a, 0x0a, 0x49, 0x50, 0x76, 0x36,
	0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x1c, 0x0a, 0x04,
	0x50, 0x6f, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x16, 0x5a, 0x14, 0x61, 0x72,
	0x69, 0x73, 0x74, 0x61, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x66,
	0x6d, 0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fmp_inet_proto_rawDescOnce sync.Once
	file_fmp_inet_proto_rawDescData = file_fmp_inet_proto_rawDesc
)

func file_fmp_inet_proto_rawDescGZIP() []byte {
	file_fmp_inet_proto_rawDescOnce.Do(func() {
		file_fmp_inet_proto_rawDescData = protoimpl.X.CompressGZIP(file_fmp_inet_proto_rawDescData)
	})
	return file_fmp_inet_proto_rawDescData
}

var file_fmp_inet_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_fmp_inet_proto_goTypes = []interface{}{
	(*IPAddress)(nil),   // 0: fmp.IPAddress
	(*IPv4Address)(nil), // 1: fmp.IPv4Address
	(*IPv6Address)(nil), // 2: fmp.IPv6Address
	(*IPPrefix)(nil),    // 3: fmp.IPPrefix
	(*IPv4Prefix)(nil),  // 4: fmp.IPv4Prefix
	(*IPv6Prefix)(nil),  // 5: fmp.IPv6Prefix
	(*Port)(nil),        // 6: fmp.Port
}
var file_fmp_inet_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_fmp_inet_proto_init() }
func file_fmp_inet_proto_init() {
	if File_fmp_inet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fmp_inet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPAddress); i {
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
		file_fmp_inet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPv4Address); i {
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
		file_fmp_inet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPv6Address); i {
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
		file_fmp_inet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPPrefix); i {
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
		file_fmp_inet_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPv4Prefix); i {
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
		file_fmp_inet_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPv6Prefix); i {
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
		file_fmp_inet_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Port); i {
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
			RawDescriptor: file_fmp_inet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fmp_inet_proto_goTypes,
		DependencyIndexes: file_fmp_inet_proto_depIdxs,
		MessageInfos:      file_fmp_inet_proto_msgTypes,
	}.Build()
	File_fmp_inet_proto = out.File
	file_fmp_inet_proto_rawDesc = nil
	file_fmp_inet_proto_goTypes = nil
	file_fmp_inet_proto_depIdxs = nil
}
