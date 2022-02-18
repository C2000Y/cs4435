// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: proto/greeter.proto

package proto

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

// The request message containing the action
type BankUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *BankUpdateRequest) Reset() {
	*x = BankUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_greeter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BankUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BankUpdateRequest) ProtoMessage() {}

func (x *BankUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_greeter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BankUpdateRequest.ProtoReflect.Descriptor instead.
func (*BankUpdateRequest) Descriptor() ([]byte, []int) {
	return file_proto_greeter_proto_rawDescGZIP(), []int{0}
}

func (x *BankUpdateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The response message containing the action
type BankUpdateReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *BankUpdateReply) Reset() {
	*x = BankUpdateReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_greeter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BankUpdateReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BankUpdateReply) ProtoMessage() {}

func (x *BankUpdateReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_greeter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BankUpdateReply.ProtoReflect.Descriptor instead.
func (*BankUpdateReply) Descriptor() ([]byte, []int) {
	return file_proto_greeter_proto_rawDescGZIP(), []int{1}
}

func (x *BankUpdateReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_greeter_proto protoreflect.FileDescriptor

var file_proto_greeter_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x62, 0x61, 0x6e, 0x6b, 0x22, 0x27, 0x0a, 0x11, 0x42,
	0x61, 0x6e, 0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2b, 0x0a, 0x0f, 0x42, 0x61, 0x6e, 0x6b, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x32, 0x46, 0x0a, 0x04, 0x42, 0x61, 0x6e, 0x6b, 0x12, 0x3e, 0x0a, 0x0a, 0x42, 0x61, 0x6e,
	0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x62, 0x61, 0x6e, 0x6b, 0x2e, 0x42,
	0x61, 0x6e, 0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x15, 0x2e, 0x62, 0x61, 0x6e, 0x6b, 0x2e, 0x42, 0x61, 0x6e, 0x6b, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x62, 0x61, 0x6e,
	0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_greeter_proto_rawDescOnce sync.Once
	file_proto_greeter_proto_rawDescData = file_proto_greeter_proto_rawDesc
)

func file_proto_greeter_proto_rawDescGZIP() []byte {
	file_proto_greeter_proto_rawDescOnce.Do(func() {
		file_proto_greeter_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_greeter_proto_rawDescData)
	})
	return file_proto_greeter_proto_rawDescData
}

var file_proto_greeter_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_greeter_proto_goTypes = []interface{}{
	(*BankUpdateRequest)(nil), // 0: bank.BankUpdateRequest
	(*BankUpdateReply)(nil),   // 1: bank.BankUpdateReply
}
var file_proto_greeter_proto_depIdxs = []int32{
	0, // 0: bank.Bank.BankUpdate:input_type -> bank.BankUpdateRequest
	1, // 1: bank.Bank.BankUpdate:output_type -> bank.BankUpdateReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_greeter_proto_init() }
func file_proto_greeter_proto_init() {
	if File_proto_greeter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_greeter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BankUpdateRequest); i {
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
		file_proto_greeter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BankUpdateReply); i {
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
			RawDescriptor: file_proto_greeter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_greeter_proto_goTypes,
		DependencyIndexes: file_proto_greeter_proto_depIdxs,
		MessageInfos:      file_proto_greeter_proto_msgTypes,
	}.Build()
	File_proto_greeter_proto = out.File
	file_proto_greeter_proto_rawDesc = nil
	file_proto_greeter_proto_goTypes = nil
	file_proto_greeter_proto_depIdxs = nil
}