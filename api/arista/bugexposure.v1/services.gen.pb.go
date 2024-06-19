// Copyright (c) 2023 Arista Networks, Inc.  All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

//
// Code generated by boomtown. DO NOT EDIT.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: arista/bugexposure.v1/services.gen.proto

package bugexposure

import (
	subscriptions "github.com/aristanetworks/cloudvision-go/api/arista/subscriptions"
	time "github.com/aristanetworks/cloudvision-go/api/arista/time"
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

type MetaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Time holds the timestamp of the last item included in the metadata calculation.
	Time *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	// Operation indicates how the value in this response should be considered.
	// Under non-subscribe requests, this value should always be INITIAL. In a subscription,
	// once all initial data is streamed and the client begins to receive modification updates,
	// you should not see INITIAL again.
	Type subscriptions.Operation `protobuf:"varint,2,opt,name=type,proto3,enum=arista.subscriptions.Operation" json:"type,omitempty"`
	// Count is the number of items present under the conditions of the request.
	Count *wrapperspb.UInt32Value `protobuf:"bytes,3,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *MetaResponse) Reset() {
	*x = MetaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaResponse) ProtoMessage() {}

func (x *MetaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaResponse.ProtoReflect.Descriptor instead.
func (*MetaResponse) Descriptor() ([]byte, []int) {
	return file_arista_bugexposure_v1_services_gen_proto_rawDescGZIP(), []int{0}
}

func (x *MetaResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *MetaResponse) GetType() subscriptions.Operation {
	if x != nil {
		return x.Type
	}
	return subscriptions.Operation(0)
}

func (x *MetaResponse) GetCount() *wrapperspb.UInt32Value {
	if x != nil {
		return x.Count
	}
	return nil
}

type BugExposureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies a BugExposure instance to retrieve.
	// This value must be populated.
	Key *BugExposureKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Time indicates the time for which you are interested in the data.
	// If no time is given, the server will use the time at which it makes the request.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *BugExposureRequest) Reset() {
	*x = BugExposureRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BugExposureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BugExposureRequest) ProtoMessage() {}

func (x *BugExposureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BugExposureRequest.ProtoReflect.Descriptor instead.
func (*BugExposureRequest) Descriptor() ([]byte, []int) {
	return file_arista_bugexposure_v1_services_gen_proto_rawDescGZIP(), []int{1}
}

func (x *BugExposureRequest) GetKey() *BugExposureKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *BugExposureRequest) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type BugExposureResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Value is the value requested.
	// This structure will be fully-populated as it exists in the datastore. If
	// optional fields were not given at creation, these fields will be empty or
	// set to default values.
	Value *BugExposure `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Time carries the (UTC) timestamp of the last-modification of the
	// BugExposure instance in this response.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *BugExposureResponse) Reset() {
	*x = BugExposureResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BugExposureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BugExposureResponse) ProtoMessage() {}

func (x *BugExposureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BugExposureResponse.ProtoReflect.Descriptor instead.
func (*BugExposureResponse) Descriptor() ([]byte, []int) {
	return file_arista_bugexposure_v1_services_gen_proto_rawDescGZIP(), []int{2}
}

func (x *BugExposureResponse) GetValue() *BugExposure {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *BugExposureResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type BugExposureStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// PartialEqFilter provides a way to server-side filter a GetAll/Subscribe.
	// This requires all provided fields to be equal to the response.
	//
	// While transparent to users, this field also allows services to optimize internal
	// subscriptions if filter(s) are sufficiently specific.
	PartialEqFilter []*BugExposure `protobuf:"bytes,1,rep,name=partial_eq_filter,json=partialEqFilter,proto3" json:"partial_eq_filter,omitempty"`
	// TimeRange allows limiting response data to within a specified time window.
	// If this field is populated, at least one of the two time fields are required.
	//
	// For GetAll, the fields start and end can be used as follows:
	//
	//   - end: Returns the state of each BugExposure at end.
	//   - Each BugExposure response is fully-specified (all fields set).
	//   - start: Returns the state of each BugExposure at start, followed by updates until now.
	//   - Each BugExposure response at start is fully-specified, but updates may be partial.
	//   - start and end: Returns the state of each BugExposure at start, followed by updates
	//     until end.
	//   - Each BugExposure response at start is fully-specified, but updates until end may
	//     be partial.
	//
	// This field is not allowed in the Subscribe RPC.
	Time *time.TimeBounds `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *BugExposureStreamRequest) Reset() {
	*x = BugExposureStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BugExposureStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BugExposureStreamRequest) ProtoMessage() {}

func (x *BugExposureStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BugExposureStreamRequest.ProtoReflect.Descriptor instead.
func (*BugExposureStreamRequest) Descriptor() ([]byte, []int) {
	return file_arista_bugexposure_v1_services_gen_proto_rawDescGZIP(), []int{3}
}

func (x *BugExposureStreamRequest) GetPartialEqFilter() []*BugExposure {
	if x != nil {
		return x.PartialEqFilter
	}
	return nil
}

func (x *BugExposureStreamRequest) GetTime() *time.TimeBounds {
	if x != nil {
		return x.Time
	}
	return nil
}

type BugExposureStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Value is a value deemed relevant to the initiating request.
	// This structure will always have its key-field populated. Which other fields are
	// populated, and why, depends on the value of Operation and what triggered this notification.
	Value *BugExposure `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Time holds the timestamp of this BugExposure's last modification.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	// Operation indicates how the BugExposure value in this response should be considered.
	// Under non-subscribe requests, this value should always be INITIAL. In a subscription,
	// once all initial data is streamed and the client begins to receive modification updates,
	// you should not see INITIAL again.
	Type subscriptions.Operation `protobuf:"varint,3,opt,name=type,proto3,enum=arista.subscriptions.Operation" json:"type,omitempty"`
}

func (x *BugExposureStreamResponse) Reset() {
	*x = BugExposureStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BugExposureStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BugExposureStreamResponse) ProtoMessage() {}

func (x *BugExposureStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_bugexposure_v1_services_gen_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BugExposureStreamResponse.ProtoReflect.Descriptor instead.
func (*BugExposureStreamResponse) Descriptor() ([]byte, []int) {
	return file_arista_bugexposure_v1_services_gen_proto_rawDescGZIP(), []int{4}
}

func (x *BugExposureStreamResponse) GetValue() *BugExposure {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *BugExposureStreamResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *BugExposureStreamResponse) GetType() subscriptions.Operation {
	if x != nil {
		return x.Type
	}
	return subscriptions.Operation(0)
}

var File_arista_bugexposure_v1_services_gen_proto protoreflect.FileDescriptor

var file_arista_bugexposure_v1_services_gen_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f,
	0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x1a, 0x27, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70,
	0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2f, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f,
	0x73, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x28, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa7, 0x01,
	0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e,
	0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x33,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x7d, 0x0a, 0x12, 0x42, 0x75, 0x67, 0x45, 0x78,
	0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x37, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x4b, 0x65,
	0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x7f, 0x0a, 0x13, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70,
	0x6f, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x97, 0x01, 0x0a, 0x18, 0x42, 0x75, 0x67, 0x45,
	0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x4e, 0x0a, 0x11, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x5f,
	0x65, 0x71, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f,
	0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73,
	0x75, 0x72, 0x65, 0x52, 0x0f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x45, 0x71, 0x46, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x74, 0x69, 0x6d, 0x65,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x52, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x22, 0xba, 0x01, 0x0a, 0x19, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72,
	0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x38, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f, 0x73,
	0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x75,
	0x72, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x32, 0xa0,
	0x04, 0x0a, 0x12, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5f, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x65, 0x12,
	0x29, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f,
	0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73,
	0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6d, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c,
	0x12, 0x2f, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70,
	0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f,
	0x73, 0x75, 0x72, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x30, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78,
	0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70,
	0x6f, 0x73, 0x75, 0x72, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x70, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x12, 0x2f, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65,
	0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78,
	0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67,
	0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45,
	0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x5f, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4d, 0x65,
	0x74, 0x61, 0x12, 0x2f, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65,
	0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78,
	0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67,
	0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x67, 0x0a, 0x0d, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x2f, 0x2e, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30,
	0x01, 0x42, 0x82, 0x01, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x62, 0x75, 0x67, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x42,
	0x13, 0x42, 0x75, 0x67, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x50, 0x01, 0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x2d, 0x67, 0x6f,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x62, 0x75, 0x67, 0x65,
	0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x3b, 0x62, 0x75, 0x67, 0x65, 0x78,
	0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arista_bugexposure_v1_services_gen_proto_rawDescOnce sync.Once
	file_arista_bugexposure_v1_services_gen_proto_rawDescData = file_arista_bugexposure_v1_services_gen_proto_rawDesc
)

func file_arista_bugexposure_v1_services_gen_proto_rawDescGZIP() []byte {
	file_arista_bugexposure_v1_services_gen_proto_rawDescOnce.Do(func() {
		file_arista_bugexposure_v1_services_gen_proto_rawDescData = protoimpl.X.CompressGZIP(file_arista_bugexposure_v1_services_gen_proto_rawDescData)
	})
	return file_arista_bugexposure_v1_services_gen_proto_rawDescData
}

var file_arista_bugexposure_v1_services_gen_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_arista_bugexposure_v1_services_gen_proto_goTypes = []interface{}{
	(*MetaResponse)(nil),              // 0: arista.bugexposure.v1.MetaResponse
	(*BugExposureRequest)(nil),        // 1: arista.bugexposure.v1.BugExposureRequest
	(*BugExposureResponse)(nil),       // 2: arista.bugexposure.v1.BugExposureResponse
	(*BugExposureStreamRequest)(nil),  // 3: arista.bugexposure.v1.BugExposureStreamRequest
	(*BugExposureStreamResponse)(nil), // 4: arista.bugexposure.v1.BugExposureStreamResponse
	(*timestamppb.Timestamp)(nil),     // 5: google.protobuf.Timestamp
	(subscriptions.Operation)(0),      // 6: arista.subscriptions.Operation
	(*wrapperspb.UInt32Value)(nil),    // 7: google.protobuf.UInt32Value
	(*BugExposureKey)(nil),            // 8: arista.bugexposure.v1.BugExposureKey
	(*BugExposure)(nil),               // 9: arista.bugexposure.v1.BugExposure
	(*time.TimeBounds)(nil),           // 10: arista.time.TimeBounds
}
var file_arista_bugexposure_v1_services_gen_proto_depIdxs = []int32{
	5,  // 0: arista.bugexposure.v1.MetaResponse.time:type_name -> google.protobuf.Timestamp
	6,  // 1: arista.bugexposure.v1.MetaResponse.type:type_name -> arista.subscriptions.Operation
	7,  // 2: arista.bugexposure.v1.MetaResponse.count:type_name -> google.protobuf.UInt32Value
	8,  // 3: arista.bugexposure.v1.BugExposureRequest.key:type_name -> arista.bugexposure.v1.BugExposureKey
	5,  // 4: arista.bugexposure.v1.BugExposureRequest.time:type_name -> google.protobuf.Timestamp
	9,  // 5: arista.bugexposure.v1.BugExposureResponse.value:type_name -> arista.bugexposure.v1.BugExposure
	5,  // 6: arista.bugexposure.v1.BugExposureResponse.time:type_name -> google.protobuf.Timestamp
	9,  // 7: arista.bugexposure.v1.BugExposureStreamRequest.partial_eq_filter:type_name -> arista.bugexposure.v1.BugExposure
	10, // 8: arista.bugexposure.v1.BugExposureStreamRequest.time:type_name -> arista.time.TimeBounds
	9,  // 9: arista.bugexposure.v1.BugExposureStreamResponse.value:type_name -> arista.bugexposure.v1.BugExposure
	5,  // 10: arista.bugexposure.v1.BugExposureStreamResponse.time:type_name -> google.protobuf.Timestamp
	6,  // 11: arista.bugexposure.v1.BugExposureStreamResponse.type:type_name -> arista.subscriptions.Operation
	1,  // 12: arista.bugexposure.v1.BugExposureService.GetOne:input_type -> arista.bugexposure.v1.BugExposureRequest
	3,  // 13: arista.bugexposure.v1.BugExposureService.GetAll:input_type -> arista.bugexposure.v1.BugExposureStreamRequest
	3,  // 14: arista.bugexposure.v1.BugExposureService.Subscribe:input_type -> arista.bugexposure.v1.BugExposureStreamRequest
	3,  // 15: arista.bugexposure.v1.BugExposureService.GetMeta:input_type -> arista.bugexposure.v1.BugExposureStreamRequest
	3,  // 16: arista.bugexposure.v1.BugExposureService.SubscribeMeta:input_type -> arista.bugexposure.v1.BugExposureStreamRequest
	2,  // 17: arista.bugexposure.v1.BugExposureService.GetOne:output_type -> arista.bugexposure.v1.BugExposureResponse
	4,  // 18: arista.bugexposure.v1.BugExposureService.GetAll:output_type -> arista.bugexposure.v1.BugExposureStreamResponse
	4,  // 19: arista.bugexposure.v1.BugExposureService.Subscribe:output_type -> arista.bugexposure.v1.BugExposureStreamResponse
	0,  // 20: arista.bugexposure.v1.BugExposureService.GetMeta:output_type -> arista.bugexposure.v1.MetaResponse
	0,  // 21: arista.bugexposure.v1.BugExposureService.SubscribeMeta:output_type -> arista.bugexposure.v1.MetaResponse
	17, // [17:22] is the sub-list for method output_type
	12, // [12:17] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_arista_bugexposure_v1_services_gen_proto_init() }
func file_arista_bugexposure_v1_services_gen_proto_init() {
	if File_arista_bugexposure_v1_services_gen_proto != nil {
		return
	}
	file_arista_bugexposure_v1_bugexposure_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_arista_bugexposure_v1_services_gen_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetaResponse); i {
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
		file_arista_bugexposure_v1_services_gen_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BugExposureRequest); i {
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
		file_arista_bugexposure_v1_services_gen_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BugExposureResponse); i {
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
		file_arista_bugexposure_v1_services_gen_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BugExposureStreamRequest); i {
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
		file_arista_bugexposure_v1_services_gen_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BugExposureStreamResponse); i {
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
			RawDescriptor: file_arista_bugexposure_v1_services_gen_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_arista_bugexposure_v1_services_gen_proto_goTypes,
		DependencyIndexes: file_arista_bugexposure_v1_services_gen_proto_depIdxs,
		MessageInfos:      file_arista_bugexposure_v1_services_gen_proto_msgTypes,
	}.Build()
	File_arista_bugexposure_v1_services_gen_proto = out.File
	file_arista_bugexposure_v1_services_gen_proto_rawDesc = nil
	file_arista_bugexposure_v1_services_gen_proto_goTypes = nil
	file_arista_bugexposure_v1_services_gen_proto_depIdxs = nil
}
