// Copyright (c) 2023 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

//
// Code generated by boomtown. DO NOT EDIT.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.7
// source: arista/imagestatus.v1/services.gen.proto

package imagestatus

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

type SummaryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies a Summary instance to retrieve.
	// This value must be populated.
	Key *SummaryKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Time indicates the time for which you are interested in the data.
	// If no time is given, the server will use the time at which it makes the request.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *SummaryRequest) Reset() {
	*x = SummaryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_imagestatus_v1_services_gen_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SummaryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SummaryRequest) ProtoMessage() {}

func (x *SummaryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arista_imagestatus_v1_services_gen_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SummaryRequest.ProtoReflect.Descriptor instead.
func (*SummaryRequest) Descriptor() ([]byte, []int) {
	return file_arista_imagestatus_v1_services_gen_proto_rawDescGZIP(), []int{0}
}

func (x *SummaryRequest) GetKey() *SummaryKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *SummaryRequest) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type SummaryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Value is the value requested.
	// This structure will be fully-populated as it exists in the datastore. If
	// optional fields were not given at creation, these fields will be empty or
	// set to default values.
	Value *Summary `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Time carries the (UTC) timestamp of the last-modification of the
	// Summary instance in this response.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *SummaryResponse) Reset() {
	*x = SummaryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_imagestatus_v1_services_gen_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SummaryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SummaryResponse) ProtoMessage() {}

func (x *SummaryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_imagestatus_v1_services_gen_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SummaryResponse.ProtoReflect.Descriptor instead.
func (*SummaryResponse) Descriptor() ([]byte, []int) {
	return file_arista_imagestatus_v1_services_gen_proto_rawDescGZIP(), []int{1}
}

func (x *SummaryResponse) GetValue() *Summary {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SummaryResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type SummaryStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// PartialEqFilter provides a way to server-side filter a GetAll/Subscribe.
	// This requires all provided fields to be equal to the response.
	//
	// While transparent to users, this field also allows services to optimize internal
	// subscriptions if filter(s) are sufficiently specific.
	PartialEqFilter []*Summary `protobuf:"bytes,1,rep,name=partial_eq_filter,json=partialEqFilter,proto3" json:"partial_eq_filter,omitempty"`
	// TimeRange allows limiting response data to within a specified time window.
	// If this field is populated, at least one of the two time fields are required.
	//
	// For GetAll, the fields start and end can be used as follows:
	//
	//   * end: Returns the state of each Summary at end.
	//     * Each Summary response is fully-specified (all fields set).
	//   * start: Returns the state of each Summary at start, followed by updates until now.
	//     * Each Summary response at start is fully-specified, but updates may be partial.
	//   * start and end: Returns the state of each Summary at start, followed by updates
	//     until end.
	//     * Each Summary response at start is fully-specified, but updates until end may
	//       be partial.
	//
	// This field is not allowed in the Subscribe RPC.
	Time *time.TimeBounds `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *SummaryStreamRequest) Reset() {
	*x = SummaryStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_imagestatus_v1_services_gen_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SummaryStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SummaryStreamRequest) ProtoMessage() {}

func (x *SummaryStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arista_imagestatus_v1_services_gen_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SummaryStreamRequest.ProtoReflect.Descriptor instead.
func (*SummaryStreamRequest) Descriptor() ([]byte, []int) {
	return file_arista_imagestatus_v1_services_gen_proto_rawDescGZIP(), []int{2}
}

func (x *SummaryStreamRequest) GetPartialEqFilter() []*Summary {
	if x != nil {
		return x.PartialEqFilter
	}
	return nil
}

func (x *SummaryStreamRequest) GetTime() *time.TimeBounds {
	if x != nil {
		return x.Time
	}
	return nil
}

type SummaryStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Value is a value deemed relevant to the initiating request.
	// This structure will always have its key-field populated. Which other fields are
	// populated, and why, depends on the value of Operation and what triggered this notification.
	Value *Summary `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Time holds the timestamp of this Summary's last modification.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	// Operation indicates how the Summary value in this response should be considered.
	// Under non-subscribe requests, this value should always be INITIAL. In a subscription,
	// once all initial data is streamed and the client begins to receive modification updates,
	// you should not see INITIAL again.
	Type subscriptions.Operation `protobuf:"varint,3,opt,name=type,proto3,enum=arista.subscriptions.Operation" json:"type,omitempty"`
}

func (x *SummaryStreamResponse) Reset() {
	*x = SummaryStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_imagestatus_v1_services_gen_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SummaryStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SummaryStreamResponse) ProtoMessage() {}

func (x *SummaryStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_imagestatus_v1_services_gen_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SummaryStreamResponse.ProtoReflect.Descriptor instead.
func (*SummaryStreamResponse) Descriptor() ([]byte, []int) {
	return file_arista_imagestatus_v1_services_gen_proto_rawDescGZIP(), []int{3}
}

func (x *SummaryStreamResponse) GetValue() *Summary {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SummaryStreamResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *SummaryStreamResponse) GetType() subscriptions.Operation {
	if x != nil {
		return x.Type
	}
	return subscriptions.Operation_UNSPECIFIED
}

var File_arista_imagestatus_v1_services_gen_proto protoreflect.FileDescriptor

var file_arista_imagestatus_v1_services_gen_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76,
	0x31, 0x1a, 0x27, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x28, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x75, 0x0a,
	0x0e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x33, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04,
	0x74, 0x69, 0x6d, 0x65, 0x22, 0x77, 0x0a, 0x0f, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x8f, 0x01,
	0x0a, 0x14, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x4a, 0x0a, 0x11, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61,
	0x6c, 0x5f, 0x65, 0x71, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x52, 0x0f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x45, 0x71, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22,
	0xb2, 0x01, 0x0a, 0x15, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12,
	0x33, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e,
	0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x32, 0xba, 0x02, 0x0a, 0x0e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x4f, 0x6e,
	0x65, 0x12, 0x25, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x65, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x2b, 0x2e, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x68, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x12, 0x2b, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30,
	0x01, 0x42, 0x66, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x13,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x50, 0x01, 0x5a, 0x32, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x3b, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_arista_imagestatus_v1_services_gen_proto_rawDescOnce sync.Once
	file_arista_imagestatus_v1_services_gen_proto_rawDescData = file_arista_imagestatus_v1_services_gen_proto_rawDesc
)

func file_arista_imagestatus_v1_services_gen_proto_rawDescGZIP() []byte {
	file_arista_imagestatus_v1_services_gen_proto_rawDescOnce.Do(func() {
		file_arista_imagestatus_v1_services_gen_proto_rawDescData = protoimpl.X.CompressGZIP(file_arista_imagestatus_v1_services_gen_proto_rawDescData)
	})
	return file_arista_imagestatus_v1_services_gen_proto_rawDescData
}

var file_arista_imagestatus_v1_services_gen_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_arista_imagestatus_v1_services_gen_proto_goTypes = []interface{}{
	(*SummaryRequest)(nil),        // 0: arista.imagestatus.v1.SummaryRequest
	(*SummaryResponse)(nil),       // 1: arista.imagestatus.v1.SummaryResponse
	(*SummaryStreamRequest)(nil),  // 2: arista.imagestatus.v1.SummaryStreamRequest
	(*SummaryStreamResponse)(nil), // 3: arista.imagestatus.v1.SummaryStreamResponse
	(*SummaryKey)(nil),            // 4: arista.imagestatus.v1.SummaryKey
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
	(*Summary)(nil),               // 6: arista.imagestatus.v1.Summary
	(*time.TimeBounds)(nil),       // 7: arista.time.TimeBounds
	(subscriptions.Operation)(0),  // 8: arista.subscriptions.Operation
}
var file_arista_imagestatus_v1_services_gen_proto_depIdxs = []int32{
	4,  // 0: arista.imagestatus.v1.SummaryRequest.key:type_name -> arista.imagestatus.v1.SummaryKey
	5,  // 1: arista.imagestatus.v1.SummaryRequest.time:type_name -> google.protobuf.Timestamp
	6,  // 2: arista.imagestatus.v1.SummaryResponse.value:type_name -> arista.imagestatus.v1.Summary
	5,  // 3: arista.imagestatus.v1.SummaryResponse.time:type_name -> google.protobuf.Timestamp
	6,  // 4: arista.imagestatus.v1.SummaryStreamRequest.partial_eq_filter:type_name -> arista.imagestatus.v1.Summary
	7,  // 5: arista.imagestatus.v1.SummaryStreamRequest.time:type_name -> arista.time.TimeBounds
	6,  // 6: arista.imagestatus.v1.SummaryStreamResponse.value:type_name -> arista.imagestatus.v1.Summary
	5,  // 7: arista.imagestatus.v1.SummaryStreamResponse.time:type_name -> google.protobuf.Timestamp
	8,  // 8: arista.imagestatus.v1.SummaryStreamResponse.type:type_name -> arista.subscriptions.Operation
	0,  // 9: arista.imagestatus.v1.SummaryService.GetOne:input_type -> arista.imagestatus.v1.SummaryRequest
	2,  // 10: arista.imagestatus.v1.SummaryService.GetAll:input_type -> arista.imagestatus.v1.SummaryStreamRequest
	2,  // 11: arista.imagestatus.v1.SummaryService.Subscribe:input_type -> arista.imagestatus.v1.SummaryStreamRequest
	1,  // 12: arista.imagestatus.v1.SummaryService.GetOne:output_type -> arista.imagestatus.v1.SummaryResponse
	3,  // 13: arista.imagestatus.v1.SummaryService.GetAll:output_type -> arista.imagestatus.v1.SummaryStreamResponse
	3,  // 14: arista.imagestatus.v1.SummaryService.Subscribe:output_type -> arista.imagestatus.v1.SummaryStreamResponse
	12, // [12:15] is the sub-list for method output_type
	9,  // [9:12] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_arista_imagestatus_v1_services_gen_proto_init() }
func file_arista_imagestatus_v1_services_gen_proto_init() {
	if File_arista_imagestatus_v1_services_gen_proto != nil {
		return
	}
	file_arista_imagestatus_v1_imagestatus_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_arista_imagestatus_v1_services_gen_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SummaryRequest); i {
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
		file_arista_imagestatus_v1_services_gen_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SummaryResponse); i {
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
		file_arista_imagestatus_v1_services_gen_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SummaryStreamRequest); i {
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
		file_arista_imagestatus_v1_services_gen_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SummaryStreamResponse); i {
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
			RawDescriptor: file_arista_imagestatus_v1_services_gen_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_arista_imagestatus_v1_services_gen_proto_goTypes,
		DependencyIndexes: file_arista_imagestatus_v1_services_gen_proto_depIdxs,
		MessageInfos:      file_arista_imagestatus_v1_services_gen_proto_msgTypes,
	}.Build()
	File_arista_imagestatus_v1_services_gen_proto = out.File
	file_arista_imagestatus_v1_services_gen_proto_rawDesc = nil
	file_arista_imagestatus_v1_services_gen_proto_goTypes = nil
	file_arista_imagestatus_v1_services_gen_proto_depIdxs = nil
}
