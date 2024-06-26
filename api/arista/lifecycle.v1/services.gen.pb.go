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
// source: arista/lifecycle.v1/services.gen.proto

package lifecycle

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
		mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaResponse) ProtoMessage() {}

func (x *MetaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[0]
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
	return file_arista_lifecycle_v1_services_gen_proto_rawDescGZIP(), []int{0}
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

type DeviceLifecycleSummaryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key uniquely identifies a DeviceLifecycleSummary instance to retrieve.
	// This value must be populated.
	Key *DeviceLifecycleSummaryKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Time indicates the time for which you are interested in the data.
	// If no time is given, the server will use the time at which it makes the request.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *DeviceLifecycleSummaryRequest) Reset() {
	*x = DeviceLifecycleSummaryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceLifecycleSummaryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceLifecycleSummaryRequest) ProtoMessage() {}

func (x *DeviceLifecycleSummaryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceLifecycleSummaryRequest.ProtoReflect.Descriptor instead.
func (*DeviceLifecycleSummaryRequest) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_services_gen_proto_rawDescGZIP(), []int{1}
}

func (x *DeviceLifecycleSummaryRequest) GetKey() *DeviceLifecycleSummaryKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *DeviceLifecycleSummaryRequest) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type DeviceLifecycleSummaryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Value is the value requested.
	// This structure will be fully-populated as it exists in the datastore. If
	// optional fields were not given at creation, these fields will be empty or
	// set to default values.
	Value *DeviceLifecycleSummary `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Time carries the (UTC) timestamp of the last-modification of the
	// DeviceLifecycleSummary instance in this response.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *DeviceLifecycleSummaryResponse) Reset() {
	*x = DeviceLifecycleSummaryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceLifecycleSummaryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceLifecycleSummaryResponse) ProtoMessage() {}

func (x *DeviceLifecycleSummaryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceLifecycleSummaryResponse.ProtoReflect.Descriptor instead.
func (*DeviceLifecycleSummaryResponse) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_services_gen_proto_rawDescGZIP(), []int{2}
}

func (x *DeviceLifecycleSummaryResponse) GetValue() *DeviceLifecycleSummary {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *DeviceLifecycleSummaryResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type DeviceLifecycleSummaryStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// PartialEqFilter provides a way to server-side filter a GetAll/Subscribe.
	// This requires all provided fields to be equal to the response.
	//
	// While transparent to users, this field also allows services to optimize internal
	// subscriptions if filter(s) are sufficiently specific.
	PartialEqFilter []*DeviceLifecycleSummary `protobuf:"bytes,1,rep,name=partial_eq_filter,json=partialEqFilter,proto3" json:"partial_eq_filter,omitempty"`
	// TimeRange allows limiting response data to within a specified time window.
	// If this field is populated, at least one of the two time fields are required.
	//
	// For GetAll, the fields start and end can be used as follows:
	//
	//   - end: Returns the state of each DeviceLifecycleSummary at end.
	//   - Each DeviceLifecycleSummary response is fully-specified (all fields set).
	//   - start: Returns the state of each DeviceLifecycleSummary at start, followed by updates until now.
	//   - Each DeviceLifecycleSummary response at start is fully-specified, but updates may be partial.
	//   - start and end: Returns the state of each DeviceLifecycleSummary at start, followed by updates
	//     until end.
	//   - Each DeviceLifecycleSummary response at start is fully-specified, but updates until end may
	//     be partial.
	//
	// This field is not allowed in the Subscribe RPC.
	Time *time.TimeBounds `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *DeviceLifecycleSummaryStreamRequest) Reset() {
	*x = DeviceLifecycleSummaryStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceLifecycleSummaryStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceLifecycleSummaryStreamRequest) ProtoMessage() {}

func (x *DeviceLifecycleSummaryStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceLifecycleSummaryStreamRequest.ProtoReflect.Descriptor instead.
func (*DeviceLifecycleSummaryStreamRequest) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_services_gen_proto_rawDescGZIP(), []int{3}
}

func (x *DeviceLifecycleSummaryStreamRequest) GetPartialEqFilter() []*DeviceLifecycleSummary {
	if x != nil {
		return x.PartialEqFilter
	}
	return nil
}

func (x *DeviceLifecycleSummaryStreamRequest) GetTime() *time.TimeBounds {
	if x != nil {
		return x.Time
	}
	return nil
}

type DeviceLifecycleSummaryStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Value is a value deemed relevant to the initiating request.
	// This structure will always have its key-field populated. Which other fields are
	// populated, and why, depends on the value of Operation and what triggered this notification.
	Value *DeviceLifecycleSummary `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// Time holds the timestamp of this DeviceLifecycleSummary's last modification.
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	// Operation indicates how the DeviceLifecycleSummary value in this response should be considered.
	// Under non-subscribe requests, this value should always be INITIAL. In a subscription,
	// once all initial data is streamed and the client begins to receive modification updates,
	// you should not see INITIAL again.
	Type subscriptions.Operation `protobuf:"varint,3,opt,name=type,proto3,enum=arista.subscriptions.Operation" json:"type,omitempty"`
}

func (x *DeviceLifecycleSummaryStreamResponse) Reset() {
	*x = DeviceLifecycleSummaryStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceLifecycleSummaryStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceLifecycleSummaryStreamResponse) ProtoMessage() {}

func (x *DeviceLifecycleSummaryStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arista_lifecycle_v1_services_gen_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceLifecycleSummaryStreamResponse.ProtoReflect.Descriptor instead.
func (*DeviceLifecycleSummaryStreamResponse) Descriptor() ([]byte, []int) {
	return file_arista_lifecycle_v1_services_gen_proto_rawDescGZIP(), []int{4}
}

func (x *DeviceLifecycleSummaryStreamResponse) GetValue() *DeviceLifecycleSummary {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *DeviceLifecycleSummaryStreamResponse) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *DeviceLifecycleSummaryStreamResponse) GetType() subscriptions.Operation {
	if x != nil {
		return x.Type
	}
	return subscriptions.Operation(0)
}

var File_arista_lifecycle_v1_services_gen_proto protoreflect.FileDescriptor

var file_arista_lifecycle_v1_services_gen_proto_rawDesc = []byte{
	0x0a, 0x26, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x67,
	0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x23, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e,
	0x76, 0x31, 0x2f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x16, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x28, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa7, 0x01, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x73, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e,
	0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x91, 0x01, 0x0a, 0x1d, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79,
	0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x40, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e,
	0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63,
	0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x22, 0x93, 0x01, 0x0a, 0x1e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69,
	0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c,
	0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0xab, 0x01, 0x0a, 0x23, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x57, 0x0a, 0x11, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x65, 0x71, 0x5f,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x0f, 0x70, 0x61, 0x72, 0x74, 0x69,
	0x61, 0x6c, 0x45, 0x71, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x6f, 0x75, 0x6e, 0x64,
	0x73, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0xce, 0x01, 0x0a, 0x24, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x41, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2b, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65,
	0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1f, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x73, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x32, 0xf0, 0x04, 0x0a, 0x1d, 0x44, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x71, 0x0a, 0x06, 0x47, 0x65,
	0x74, 0x4f, 0x6e, 0x65, 0x12, 0x32, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69,
	0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75,
	0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7f, 0x0a,
	0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x38, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x39, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63,
	0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69,
	0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x82,
	0x01, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x38, 0x2e, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x39, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e,
	0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x30, 0x01, 0x12, 0x66, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x38,
	0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63,
	0x79, 0x63, 0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d,
	0x65, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6e, 0x0a, 0x0d, 0x53,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x38, 0x2e, 0x61,
	0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e,
	0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x7a, 0x0a, 0x17, 0x63,
	0x6f, 0x6d, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79,
	0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x11, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c,
	0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x50, 0x01, 0x5a, 0x4a, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x2d, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x3b, 0x6c, 0x69,
	0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arista_lifecycle_v1_services_gen_proto_rawDescOnce sync.Once
	file_arista_lifecycle_v1_services_gen_proto_rawDescData = file_arista_lifecycle_v1_services_gen_proto_rawDesc
)

func file_arista_lifecycle_v1_services_gen_proto_rawDescGZIP() []byte {
	file_arista_lifecycle_v1_services_gen_proto_rawDescOnce.Do(func() {
		file_arista_lifecycle_v1_services_gen_proto_rawDescData = protoimpl.X.CompressGZIP(file_arista_lifecycle_v1_services_gen_proto_rawDescData)
	})
	return file_arista_lifecycle_v1_services_gen_proto_rawDescData
}

var file_arista_lifecycle_v1_services_gen_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_arista_lifecycle_v1_services_gen_proto_goTypes = []interface{}{
	(*MetaResponse)(nil),                         // 0: arista.lifecycle.v1.MetaResponse
	(*DeviceLifecycleSummaryRequest)(nil),        // 1: arista.lifecycle.v1.DeviceLifecycleSummaryRequest
	(*DeviceLifecycleSummaryResponse)(nil),       // 2: arista.lifecycle.v1.DeviceLifecycleSummaryResponse
	(*DeviceLifecycleSummaryStreamRequest)(nil),  // 3: arista.lifecycle.v1.DeviceLifecycleSummaryStreamRequest
	(*DeviceLifecycleSummaryStreamResponse)(nil), // 4: arista.lifecycle.v1.DeviceLifecycleSummaryStreamResponse
	(*timestamppb.Timestamp)(nil),                // 5: google.protobuf.Timestamp
	(subscriptions.Operation)(0),                 // 6: arista.subscriptions.Operation
	(*wrapperspb.UInt32Value)(nil),               // 7: google.protobuf.UInt32Value
	(*DeviceLifecycleSummaryKey)(nil),            // 8: arista.lifecycle.v1.DeviceLifecycleSummaryKey
	(*DeviceLifecycleSummary)(nil),               // 9: arista.lifecycle.v1.DeviceLifecycleSummary
	(*time.TimeBounds)(nil),                      // 10: arista.time.TimeBounds
}
var file_arista_lifecycle_v1_services_gen_proto_depIdxs = []int32{
	5,  // 0: arista.lifecycle.v1.MetaResponse.time:type_name -> google.protobuf.Timestamp
	6,  // 1: arista.lifecycle.v1.MetaResponse.type:type_name -> arista.subscriptions.Operation
	7,  // 2: arista.lifecycle.v1.MetaResponse.count:type_name -> google.protobuf.UInt32Value
	8,  // 3: arista.lifecycle.v1.DeviceLifecycleSummaryRequest.key:type_name -> arista.lifecycle.v1.DeviceLifecycleSummaryKey
	5,  // 4: arista.lifecycle.v1.DeviceLifecycleSummaryRequest.time:type_name -> google.protobuf.Timestamp
	9,  // 5: arista.lifecycle.v1.DeviceLifecycleSummaryResponse.value:type_name -> arista.lifecycle.v1.DeviceLifecycleSummary
	5,  // 6: arista.lifecycle.v1.DeviceLifecycleSummaryResponse.time:type_name -> google.protobuf.Timestamp
	9,  // 7: arista.lifecycle.v1.DeviceLifecycleSummaryStreamRequest.partial_eq_filter:type_name -> arista.lifecycle.v1.DeviceLifecycleSummary
	10, // 8: arista.lifecycle.v1.DeviceLifecycleSummaryStreamRequest.time:type_name -> arista.time.TimeBounds
	9,  // 9: arista.lifecycle.v1.DeviceLifecycleSummaryStreamResponse.value:type_name -> arista.lifecycle.v1.DeviceLifecycleSummary
	5,  // 10: arista.lifecycle.v1.DeviceLifecycleSummaryStreamResponse.time:type_name -> google.protobuf.Timestamp
	6,  // 11: arista.lifecycle.v1.DeviceLifecycleSummaryStreamResponse.type:type_name -> arista.subscriptions.Operation
	1,  // 12: arista.lifecycle.v1.DeviceLifecycleSummaryService.GetOne:input_type -> arista.lifecycle.v1.DeviceLifecycleSummaryRequest
	3,  // 13: arista.lifecycle.v1.DeviceLifecycleSummaryService.GetAll:input_type -> arista.lifecycle.v1.DeviceLifecycleSummaryStreamRequest
	3,  // 14: arista.lifecycle.v1.DeviceLifecycleSummaryService.Subscribe:input_type -> arista.lifecycle.v1.DeviceLifecycleSummaryStreamRequest
	3,  // 15: arista.lifecycle.v1.DeviceLifecycleSummaryService.GetMeta:input_type -> arista.lifecycle.v1.DeviceLifecycleSummaryStreamRequest
	3,  // 16: arista.lifecycle.v1.DeviceLifecycleSummaryService.SubscribeMeta:input_type -> arista.lifecycle.v1.DeviceLifecycleSummaryStreamRequest
	2,  // 17: arista.lifecycle.v1.DeviceLifecycleSummaryService.GetOne:output_type -> arista.lifecycle.v1.DeviceLifecycleSummaryResponse
	4,  // 18: arista.lifecycle.v1.DeviceLifecycleSummaryService.GetAll:output_type -> arista.lifecycle.v1.DeviceLifecycleSummaryStreamResponse
	4,  // 19: arista.lifecycle.v1.DeviceLifecycleSummaryService.Subscribe:output_type -> arista.lifecycle.v1.DeviceLifecycleSummaryStreamResponse
	0,  // 20: arista.lifecycle.v1.DeviceLifecycleSummaryService.GetMeta:output_type -> arista.lifecycle.v1.MetaResponse
	0,  // 21: arista.lifecycle.v1.DeviceLifecycleSummaryService.SubscribeMeta:output_type -> arista.lifecycle.v1.MetaResponse
	17, // [17:22] is the sub-list for method output_type
	12, // [12:17] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_arista_lifecycle_v1_services_gen_proto_init() }
func file_arista_lifecycle_v1_services_gen_proto_init() {
	if File_arista_lifecycle_v1_services_gen_proto != nil {
		return
	}
	file_arista_lifecycle_v1_lifecycle_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_arista_lifecycle_v1_services_gen_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_arista_lifecycle_v1_services_gen_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceLifecycleSummaryRequest); i {
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
		file_arista_lifecycle_v1_services_gen_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceLifecycleSummaryResponse); i {
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
		file_arista_lifecycle_v1_services_gen_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceLifecycleSummaryStreamRequest); i {
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
		file_arista_lifecycle_v1_services_gen_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceLifecycleSummaryStreamResponse); i {
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
			RawDescriptor: file_arista_lifecycle_v1_services_gen_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_arista_lifecycle_v1_services_gen_proto_goTypes,
		DependencyIndexes: file_arista_lifecycle_v1_services_gen_proto_depIdxs,
		MessageInfos:      file_arista_lifecycle_v1_services_gen_proto_msgTypes,
	}.Build()
	File_arista_lifecycle_v1_services_gen_proto = out.File
	file_arista_lifecycle_v1_services_gen_proto_rawDesc = nil
	file_arista_lifecycle_v1_services_gen_proto_goTypes = nil
	file_arista_lifecycle_v1_services_gen_proto_depIdxs = nil
}
