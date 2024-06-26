// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: fmp/extensions.proto

package fmp

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_fmp_extensions_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51423,
		Name:          "fmp.model",
		Tag:           "bytes,51423,opt,name=model",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51424,
		Name:          "fmp.model_key",
		Tag:           "varint,51424,opt,name=model_key",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51425,
		Name:          "fmp.custom_filter",
		Tag:           "bytes,51425,opt,name=custom_filter",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51426,
		Name:          "fmp.no_default_filter",
		Tag:           "varint,51426,opt,name=no_default_filter",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51427,
		Name:          "fmp.require_set_key",
		Tag:           "varint,51427,opt,name=require_set_key",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51428,
		Name:          "fmp.unkeyed_model",
		Tag:           "bytes,51428,opt,name=unkeyed_model",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51429,
		Name:          "fmp.paginated",
		Tag:           "varint,51429,opt,name=paginated",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51449,
		Name:          "fmp.child_resource",
		Tag:           "bytes,51449,opt,name=child_resource",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51450,
		Name:          "fmp.sortable",
		Tag:           "bytes,51450,opt,name=sortable",
		Filename:      "fmp/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51623,
		Name:          "fmp.disable_yang",
		Tag:           "bytes,51623,opt,name=disable_yang",
		Filename:      "fmp/extensions.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// TODO: will need an official number from Google, just like gNMI extensions
	//
	//	this works for now, though.
	//
	// optional string model = 51423;
	E_Model = &file_fmp_extensions_proto_extTypes[0]
	// optional bool model_key = 51424;
	E_ModelKey = &file_fmp_extensions_proto_extTypes[1]
	// optional string custom_filter = 51425;
	E_CustomFilter = &file_fmp_extensions_proto_extTypes[2]
	// optional bool no_default_filter = 51426;
	E_NoDefaultFilter = &file_fmp_extensions_proto_extTypes[3]
	// optional bool require_set_key = 51427;
	E_RequireSetKey = &file_fmp_extensions_proto_extTypes[4]
	// optional string unkeyed_model = 51428;
	E_UnkeyedModel = &file_fmp_extensions_proto_extTypes[5]
	// optional bool paginated = 51429;
	E_Paginated = &file_fmp_extensions_proto_extTypes[6]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional string child_resource = 51449;
	E_ChildResource = &file_fmp_extensions_proto_extTypes[7]
	// optional string sortable = 51450;
	E_Sortable = &file_fmp_extensions_proto_extTypes[8]
)

// Extension fields to descriptorpb.FileOptions.
var (
	// optional string disable_yang = 51623;
	E_DisableYang = &file_fmp_extensions_proto_extTypes[9]
)

var File_fmp_extensions_proto protoreflect.FileDescriptor

var file_fmp_extensions_proto_rawDesc = []byte{
	0x0a, 0x14, 0x66, 0x6d, 0x70, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x66, 0x6d, 0x70, 0x1a, 0x20, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x37, 0x0a,
	0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xdf, 0x91, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x3a, 0x3e, 0x0a, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f,
	0x6b, 0x65, 0x79, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe0, 0x91, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x4b, 0x65, 0x79, 0x3a, 0x46, 0x0a, 0x0d, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe1, 0x91, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x3a, 0x4d,
	0x0a, 0x11, 0x6e, 0x6f, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe2, 0x91, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x6e, 0x6f,
	0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x3a, 0x49, 0x0a,
	0x0f, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x5f, 0x73, 0x65, 0x74, 0x5f, 0x6b, 0x65, 0x79,
	0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xe3, 0x91, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x53, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x3a, 0x46, 0x0a, 0x0d, 0x75, 0x6e, 0x6b, 0x65,
	0x79, 0x65, 0x64, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe4, 0x91, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x75, 0x6e, 0x6b, 0x65, 0x79, 0x65, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x3a, 0x3f, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe5,
	0x91, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x65,
	0x64, 0x3a, 0x46, 0x0a, 0x0e, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xf9, 0x91, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x68, 0x69, 0x6c,
	0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x3a, 0x3b, 0x0a, 0x08, 0x73, 0x6f, 0x72,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xfa, 0x91, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6f,
	0x72, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x3a, 0x41, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c,
	0x65, 0x5f, 0x79, 0x61, 0x6e, 0x67, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa7, 0x93, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69,
	0x73, 0x61, 0x62, 0x6c, 0x65, 0x59, 0x61, 0x6e, 0x67, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x2d, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x6d, 0x70, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_fmp_extensions_proto_goTypes = []interface{}{
	(*descriptorpb.MessageOptions)(nil), // 0: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 1: google.protobuf.FieldOptions
	(*descriptorpb.FileOptions)(nil),    // 2: google.protobuf.FileOptions
}
var file_fmp_extensions_proto_depIdxs = []int32{
	0,  // 0: fmp.model:extendee -> google.protobuf.MessageOptions
	0,  // 1: fmp.model_key:extendee -> google.protobuf.MessageOptions
	0,  // 2: fmp.custom_filter:extendee -> google.protobuf.MessageOptions
	0,  // 3: fmp.no_default_filter:extendee -> google.protobuf.MessageOptions
	0,  // 4: fmp.require_set_key:extendee -> google.protobuf.MessageOptions
	0,  // 5: fmp.unkeyed_model:extendee -> google.protobuf.MessageOptions
	0,  // 6: fmp.paginated:extendee -> google.protobuf.MessageOptions
	1,  // 7: fmp.child_resource:extendee -> google.protobuf.FieldOptions
	1,  // 8: fmp.sortable:extendee -> google.protobuf.FieldOptions
	2,  // 9: fmp.disable_yang:extendee -> google.protobuf.FileOptions
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	0,  // [0:10] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_fmp_extensions_proto_init() }
func file_fmp_extensions_proto_init() {
	if File_fmp_extensions_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_fmp_extensions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 10,
			NumServices:   0,
		},
		GoTypes:           file_fmp_extensions_proto_goTypes,
		DependencyIndexes: file_fmp_extensions_proto_depIdxs,
		ExtensionInfos:    file_fmp_extensions_proto_extTypes,
	}.Build()
	File_fmp_extensions_proto = out.File
	file_fmp_extensions_proto_rawDesc = nil
	file_fmp_extensions_proto_goTypes = nil
	file_fmp_extensions_proto_depIdxs = nil
}
