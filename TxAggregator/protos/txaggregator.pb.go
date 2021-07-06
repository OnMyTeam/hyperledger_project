// Copyright the Hyperledger Fabric contributors. All rights reserved.
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.1
// source: protos/txaggregator.proto

package protos

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

type TaggedTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// chaincode function name
	Functionname string `protobuf:"bytes,1,opt,name=functionname,proto3" json:"functionname,omitempty"`
	// chaincode key
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	// chaincode fieldname for aggregate
	Fieldname string `protobuf:"bytes,3,opt,name=fieldname,proto3" json:"fieldname,omitempty"`
	// operator
	Operator int32 `protobuf:"varint,4,opt,name=operator,proto3" json:"operator,omitempty"`
	// operand only numeric
	Operand int32 `protobuf:"varint,5,opt,name=operand,proto3" json:"operand,omitempty"`
	// precondition
	Precondition int32 `protobuf:"varint,6,opt,name=precondition,proto3" json:"precondition,omitempty"`
	// postcondition
	Postcondition int32 `protobuf:"varint,7,opt,name=postcondition,proto3" json:"postcondition,omitempty"`
}

func (x *TaggedTransaction) Reset() {
	*x = TaggedTransaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_txaggregator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaggedTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaggedTransaction) ProtoMessage() {}

func (x *TaggedTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_protos_txaggregator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaggedTransaction.ProtoReflect.Descriptor instead.
func (*TaggedTransaction) Descriptor() ([]byte, []int) {
	return file_protos_txaggregator_proto_rawDescGZIP(), []int{0}
}

func (x *TaggedTransaction) GetFunctionname() string {
	if x != nil {
		return x.Functionname
	}
	return ""
}

func (x *TaggedTransaction) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *TaggedTransaction) GetFieldname() string {
	if x != nil {
		return x.Fieldname
	}
	return ""
}

func (x *TaggedTransaction) GetOperator() int32 {
	if x != nil {
		return x.Operator
	}
	return 0
}

func (x *TaggedTransaction) GetOperand() int32 {
	if x != nil {
		return x.Operand
	}
	return 0
}

func (x *TaggedTransaction) GetPrecondition() int32 {
	if x != nil {
		return x.Precondition
	}
	return 0
}

func (x *TaggedTransaction) GetPostcondition() int32 {
	if x != nil {
		return x.Postcondition
	}
	return 0
}

type TaggedTransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A response message indicating whether the
	// endorsement of the action was successful
	Response int32 `protobuf:"varint,1,opt,name=response,proto3" json:"response,omitempty"`
	// The payload of response. It is the bytes of ProposalResponsePayload
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *TaggedTransactionResponse) Reset() {
	*x = TaggedTransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_txaggregator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaggedTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaggedTransactionResponse) ProtoMessage() {}

func (x *TaggedTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_txaggregator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaggedTransactionResponse.ProtoReflect.Descriptor instead.
func (*TaggedTransactionResponse) Descriptor() ([]byte, []int) {
	return file_protos_txaggregator_proto_rawDescGZIP(), []int{1}
}

func (x *TaggedTransactionResponse) GetResponse() int32 {
	if x != nil {
		return x.Response
	}
	return 0
}

func (x *TaggedTransactionResponse) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_protos_txaggregator_proto protoreflect.FileDescriptor

var file_protos_txaggregator_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x74, 0x78, 0x61, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x22, 0xe7, 0x01, 0x0a, 0x11, 0x54, 0x61, 0x67, 0x67, 0x65, 0x64, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x66, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x1c, 0x0a, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x6e, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x6e, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x70, 0x6f, 0x73, 0x74, 0x63,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d,
	0x70, 0x6f, 0x73, 0x74, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x51, 0x0a,
	0x19, 0x54, 0x61, 0x67, 0x67, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x32, 0x66, 0x0a, 0x0a, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x58,
	0x0a, 0x18, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x54, 0x61, 0x67, 0x67, 0x65, 0x64, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2e, 0x54, 0x61, 0x67, 0x67, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x54,
	0x61, 0x67, 0x67, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x29, 0x5a, 0x27, 0x68, 0x79, 0x70, 0x65,
	0x72, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x54, 0x78, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_txaggregator_proto_rawDescOnce sync.Once
	file_protos_txaggregator_proto_rawDescData = file_protos_txaggregator_proto_rawDesc
)

func file_protos_txaggregator_proto_rawDescGZIP() []byte {
	file_protos_txaggregator_proto_rawDescOnce.Do(func() {
		file_protos_txaggregator_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_txaggregator_proto_rawDescData)
	})
	return file_protos_txaggregator_proto_rawDescData
}

var file_protos_txaggregator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protos_txaggregator_proto_goTypes = []interface{}{
	(*TaggedTransaction)(nil),         // 0: protos.TaggedTransaction
	(*TaggedTransactionResponse)(nil), // 1: protos.TaggedTransactionResponse
}
var file_protos_txaggregator_proto_depIdxs = []int32{
	0, // 0: protos.Aggregator.ReceiveTaggedTransaction:input_type -> protos.TaggedTransaction
	1, // 1: protos.Aggregator.ReceiveTaggedTransaction:output_type -> protos.TaggedTransactionResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_txaggregator_proto_init() }
func file_protos_txaggregator_proto_init() {
	if File_protos_txaggregator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_txaggregator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaggedTransaction); i {
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
		file_protos_txaggregator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaggedTransactionResponse); i {
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
			RawDescriptor: file_protos_txaggregator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_txaggregator_proto_goTypes,
		DependencyIndexes: file_protos_txaggregator_proto_depIdxs,
		MessageInfos:      file_protos_txaggregator_proto_msgTypes,
	}.Build()
	File_protos_txaggregator_proto = out.File
	file_protos_txaggregator_proto_rawDesc = nil
	file_protos_txaggregator_proto_goTypes = nil
	file_protos_txaggregator_proto_depIdxs = nil
}