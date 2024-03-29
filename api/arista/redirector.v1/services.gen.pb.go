// Copyright (c) 2022 Arista Networks, Inc.  All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

//
// Code generated by boomtown. DO NOT EDIT.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.7
// source: arista/redirector.v1/services.gen.proto

package redirector

import (
	subscriptions "github.com/aristanetworks/cloudvision-go/api/arista/subscriptions"
	time "github.com/aristanetworks/cloudvision-go/api/arista/time"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AssignmentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies a Assignment instance to retrieve.
	// This value must be populated.
	Key *AssignmentKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Time indicates the time for which you are interested in the data.
	// If no time is given, the server will use the time at which it makes the request.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *AssignmentRequest) Reset() {
	*x = AssignmentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_redirector_v1_services_gen_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignmentRequest) ProtoMessage() {}

func (x *AssignmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arista_redirector_v1_services_gen_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignmentRequest.ProtoReflect.Descriptor instead.
func (*AssignmentRequest) Descriptor() ([]byte, []int) {
	return file_arista_redirector_v1_services_gen_proto_rawDescGZIP(), []int{0}
}

func (x *AssignmentRequest) GetKey() *AssignmentKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *AssignmentRequest) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type AssignmentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Value is the value requested.
	// This structure will be fully-populated as it exists in the datastore. If
	// optional fields were not given at creation, these fields will be empty or
	// set to default values.
	Value *Assignment `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Time carries the (UTC) timestamp of the last-modification of the
	// Assignment instance in this response.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *AssignmentResponse) Reset() {
	*x = AssignmentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_redirector_v1_services_gen_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignmentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignmentResponse) ProtoMessage() {}

func (x *AssignmentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_redirector_v1_services_gen_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignmentResponse.ProtoReflect.Descriptor instead.
func (*AssignmentResponse) Descriptor() ([]byte, []int) {
	return file_arista_redirector_v1_services_gen_proto_rawDescGZIP(), []int{1}
}

func (x *AssignmentResponse) GetValue() *Assignment {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *AssignmentResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type AssignmentStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// PartialEqFilter provides a way to server-side filter a GetAll/Subscribe.
	// This requires all provided fields to be equal to the response.
	//
	// While transparent to users, this field also allows services to optimize internal
	// subscriptions if filter(s) are sufficiently specific.
	PartialEqFilter []*Assignment `protobuf:"bytes,1,rep,name=partial_eq_filter,json=partialEqFilter,proto3" json:"partial_eq_filter,omitempty"`
	// TimeRange allows limiting response data to within a specified time window.
	// If this field is populated, at least one of the two time fields are required.
	//
	// This field is not allowed in the Subscribe RPC.
	Time *time.TimeBounds `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *AssignmentStreamRequest) Reset() {
	*x = AssignmentStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_redirector_v1_services_gen_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignmentStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignmentStreamRequest) ProtoMessage() {}

func (x *AssignmentStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arista_redirector_v1_services_gen_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignmentStreamRequest.ProtoReflect.Descriptor instead.
func (*AssignmentStreamRequest) Descriptor() ([]byte, []int) {
	return file_arista_redirector_v1_services_gen_proto_rawDescGZIP(), []int{2}
}

func (x *AssignmentStreamRequest) GetPartialEqFilter() []*Assignment {
	if x != nil {
		return x.PartialEqFilter
	}
	return nil
}

func (x *AssignmentStreamRequest) GetTime() *time.TimeBounds {
	if x != nil {
		return x.Time
	}
	return nil
}

type AssignmentStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Value is a value deemed relevant to the initiating request.
	// This structure will always have its key-field populated. Which other fields are
	// populated, and why, depends on the value of Operation and what triggered this notification.
	Value *Assignment `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Time holds the timestamp of this Assignment's last modification.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	// Operation indicates how the Assignment value in this response should be considered.
	// Under non-subscribe requests, this value should always be INITIAL. In a subscription,
	// once all initial data is streamed and the client begins to receive modification updates,
	// you should not see INITIAL again.
	Type subscriptions.Operation `protobuf:"varint,3,opt,name=type,proto3,enum=arista.subscriptions.Operation" json:"type,omitempty"`
}

func (x *AssignmentStreamResponse) Reset() {
	*x = AssignmentStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_redirector_v1_services_gen_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignmentStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignmentStreamResponse) ProtoMessage() {}

func (x *AssignmentStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_redirector_v1_services_gen_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignmentStreamResponse.ProtoReflect.Descriptor instead.
func (*AssignmentStreamResponse) Descriptor() ([]byte, []int) {
	return file_arista_redirector_v1_services_gen_proto_rawDescGZIP(), []int{3}
}

func (x *AssignmentStreamResponse) GetValue() *Assignment {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *AssignmentStreamResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *AssignmentStreamResponse) GetType() subscriptions.Operation {
	if x != nil {
		return x.Type
	}
	return subscriptions.Operation_UNSPECIFIED
}

var File_arista_redirector_v1_services_gen_proto protoreflect.FileDescriptor

var file_arista_redirector_v1_services_gen_proto_rawDesc = []byte{
	0x0a, 0x27, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x67, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x1a,
	0x25, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x28,
	0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7a, 0x0a, 0x11, 0x41, 0x73, 0x73,
	0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x61, 0x72,
	0x69, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x7c, 0x0a, 0x12, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x22, 0x94, 0x01, 0x0a, 0x17, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x4c, 0x0a, 0x11, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x65, 0x71, 0x5f, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0f, 0x70, 0x61,
	0x72, 0x74, 0x69, 0x61, 0x6c, 0x45, 0x71, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x2b, 0x0a,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x72,
	0x69, 0x73, 0x74, 0x61, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x6f,
	0x75, 0x6e, 0x64, 0x73, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0xb7, 0x01, 0x0a, 0x18, 0x41,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e,
	0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12,
	0x33, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e,
	0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x32, 0xc9, 0x02, 0x0a, 0x11, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x06, 0x47, 0x65,
	0x74, 0x4f, 0x6e, 0x65, 0x12, 0x27, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e,
	0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x69, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c,
	0x6c, 0x12, 0x2d, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2e, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x30, 0x01, 0x12, 0x6c, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12,
	0x2d, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e,
	0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e,
	0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01,
	0x42, 0x32, 0x5a, 0x30, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x72, 0x65, 0x64, 0x69,
	0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x3b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arista_redirector_v1_services_gen_proto_rawDescOnce sync.Once
	file_arista_redirector_v1_services_gen_proto_rawDescData = file_arista_redirector_v1_services_gen_proto_rawDesc
)

func file_arista_redirector_v1_services_gen_proto_rawDescGZIP() []byte {
	file_arista_redirector_v1_services_gen_proto_rawDescOnce.Do(func() {
		file_arista_redirector_v1_services_gen_proto_rawDescData = protoimpl.X.CompressGZIP(file_arista_redirector_v1_services_gen_proto_rawDescData)
	})
	return file_arista_redirector_v1_services_gen_proto_rawDescData
}

var file_arista_redirector_v1_services_gen_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_arista_redirector_v1_services_gen_proto_goTypes = []interface{}{
	(*AssignmentRequest)(nil),        // 0: arista.redirector.v1.AssignmentRequest
	(*AssignmentResponse)(nil),       // 1: arista.redirector.v1.AssignmentResponse
	(*AssignmentStreamRequest)(nil),  // 2: arista.redirector.v1.AssignmentStreamRequest
	(*AssignmentStreamResponse)(nil), // 3: arista.redirector.v1.AssignmentStreamResponse
	(*AssignmentKey)(nil),            // 4: arista.redirector.v1.AssignmentKey
	(*timestamppb.Timestamp)(nil),    // 5: google.protobuf.Timestamp
	(*Assignment)(nil),               // 6: arista.redirector.v1.Assignment
	(*time.TimeBounds)(nil),          // 7: arista.time.TimeBounds
	(subscriptions.Operation)(0),     // 8: arista.subscriptions.Operation
}
var file_arista_redirector_v1_services_gen_proto_depIdxs = []int32{
	4,  // 0: arista.redirector.v1.AssignmentRequest.key:type_name -> arista.redirector.v1.AssignmentKey
	5,  // 1: arista.redirector.v1.AssignmentRequest.time:type_name -> google.protobuf.Timestamp
	6,  // 2: arista.redirector.v1.AssignmentResponse.value:type_name -> arista.redirector.v1.Assignment
	5,  // 3: arista.redirector.v1.AssignmentResponse.time:type_name -> google.protobuf.Timestamp
	6,  // 4: arista.redirector.v1.AssignmentStreamRequest.partial_eq_filter:type_name -> arista.redirector.v1.Assignment
	7,  // 5: arista.redirector.v1.AssignmentStreamRequest.time:type_name -> arista.time.TimeBounds
	6,  // 6: arista.redirector.v1.AssignmentStreamResponse.value:type_name -> arista.redirector.v1.Assignment
	5,  // 7: arista.redirector.v1.AssignmentStreamResponse.time:type_name -> google.protobuf.Timestamp
	8,  // 8: arista.redirector.v1.AssignmentStreamResponse.type:type_name -> arista.subscriptions.Operation
	0,  // 9: arista.redirector.v1.AssignmentService.GetOne:input_type -> arista.redirector.v1.AssignmentRequest
	2,  // 10: arista.redirector.v1.AssignmentService.GetAll:input_type -> arista.redirector.v1.AssignmentStreamRequest
	2,  // 11: arista.redirector.v1.AssignmentService.Subscribe:input_type -> arista.redirector.v1.AssignmentStreamRequest
	1,  // 12: arista.redirector.v1.AssignmentService.GetOne:output_type -> arista.redirector.v1.AssignmentResponse
	3,  // 13: arista.redirector.v1.AssignmentService.GetAll:output_type -> arista.redirector.v1.AssignmentStreamResponse
	3,  // 14: arista.redirector.v1.AssignmentService.Subscribe:output_type -> arista.redirector.v1.AssignmentStreamResponse
	12, // [12:15] is the sub-list for method output_type
	9,  // [9:12] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_arista_redirector_v1_services_gen_proto_init() }
func file_arista_redirector_v1_services_gen_proto_init() {
	if File_arista_redirector_v1_services_gen_proto != nil {
		return
	}
	file_arista_redirector_v1_redirector_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_arista_redirector_v1_services_gen_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignmentRequest); i {
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
		file_arista_redirector_v1_services_gen_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignmentResponse); i {
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
		file_arista_redirector_v1_services_gen_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignmentStreamRequest); i {
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
		file_arista_redirector_v1_services_gen_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignmentStreamResponse); i {
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
			RawDescriptor: file_arista_redirector_v1_services_gen_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_arista_redirector_v1_services_gen_proto_goTypes,
		DependencyIndexes: file_arista_redirector_v1_services_gen_proto_depIdxs,
		MessageInfos:      file_arista_redirector_v1_services_gen_proto_msgTypes,
	}.Build()
	File_arista_redirector_v1_services_gen_proto = out.File
	file_arista_redirector_v1_services_gen_proto_rawDesc = nil
	file_arista_redirector_v1_services_gen_proto_goTypes = nil
	file_arista_redirector_v1_services_gen_proto_depIdxs = nil
}
