// Copyright (c) 2020 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

// Useful types that come from ietf-yang-types.yang

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.7
// source: fmp/yang.proto

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

type MACAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MACAddress) Reset() {
	*x = MACAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_yang_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MACAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MACAddress) ProtoMessage() {}

func (x *MACAddress) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_yang_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MACAddress.ProtoReflect.Descriptor instead.
func (*MACAddress) Descriptor() ([]byte, []int) {
	return file_fmp_yang_proto_rawDescGZIP(), []int{0}
}

func (x *MACAddress) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type RepeatedMACAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []*MACAddress `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *RepeatedMACAddress) Reset() {
	*x = RepeatedMACAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fmp_yang_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepeatedMACAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepeatedMACAddress) ProtoMessage() {}

func (x *RepeatedMACAddress) ProtoReflect() protoreflect.Message {
	mi := &file_fmp_yang_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepeatedMACAddress.ProtoReflect.Descriptor instead.
func (*RepeatedMACAddress) Descriptor() ([]byte, []int) {
	return file_fmp_yang_proto_rawDescGZIP(), []int{1}
}

func (x *RepeatedMACAddress) GetValues() []*MACAddress {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_fmp_yang_proto protoreflect.FileDescriptor

var file_fmp_yang_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x66, 0x6d, 0x70, 0x2f, 0x79, 0x61, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x03, 0x66, 0x6d, 0x70, 0x22, 0x22, 0x0a, 0x0a, 0x4d, 0x41, 0x43, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3d, 0x0a, 0x12, 0x52, 0x65, 0x70,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x4d, 0x41, 0x43, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x27, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x66, 0x6d, 0x70, 0x2e, 0x4d, 0x41, 0x43, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x42, 0x16, 0x5a, 0x14, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x66, 0x6d, 0x70,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fmp_yang_proto_rawDescOnce sync.Once
	file_fmp_yang_proto_rawDescData = file_fmp_yang_proto_rawDesc
)

func file_fmp_yang_proto_rawDescGZIP() []byte {
	file_fmp_yang_proto_rawDescOnce.Do(func() {
		file_fmp_yang_proto_rawDescData = protoimpl.X.CompressGZIP(file_fmp_yang_proto_rawDescData)
	})
	return file_fmp_yang_proto_rawDescData
}

var file_fmp_yang_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_fmp_yang_proto_goTypes = []interface{}{
	(*MACAddress)(nil),         // 0: fmp.MACAddress
	(*RepeatedMACAddress)(nil), // 1: fmp.RepeatedMACAddress
}
var file_fmp_yang_proto_depIdxs = []int32{
	0, // 0: fmp.RepeatedMACAddress.values:type_name -> fmp.MACAddress
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_fmp_yang_proto_init() }
func file_fmp_yang_proto_init() {
	if File_fmp_yang_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fmp_yang_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MACAddress); i {
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
		file_fmp_yang_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepeatedMACAddress); i {
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
			RawDescriptor: file_fmp_yang_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fmp_yang_proto_goTypes,
		DependencyIndexes: file_fmp_yang_proto_depIdxs,
		MessageInfos:      file_fmp_yang_proto_msgTypes,
	}.Build()
	File_fmp_yang_proto = out.File
	file_fmp_yang_proto_rawDesc = nil
	file_fmp_yang_proto_goTypes = nil
	file_fmp_yang_proto_depIdxs = nil
}
