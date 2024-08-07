// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: proto/myservice.proto

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

type ProcessDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin      string `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"` // (!) origin is the original client and initiator of the request
	Source      string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"` // (!) source is the immediate sender of the request
	Destination string `protobuf:"bytes,3,opt,name=destination,proto3" json:"destination,omitempty"`
	Data        int64  `protobuf:"varint,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProcessDataRequest) Reset() {
	*x = ProcessDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_myservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessDataRequest) ProtoMessage() {}

func (x *ProcessDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_myservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessDataRequest.ProtoReflect.Descriptor instead.
func (*ProcessDataRequest) Descriptor() ([]byte, []int) {
	return file_proto_myservice_proto_rawDescGZIP(), []int{0}
}

func (x *ProcessDataRequest) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *ProcessDataRequest) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *ProcessDataRequest) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

func (x *ProcessDataRequest) GetData() int64 {
	if x != nil {
		return x.Data
	}
	return 0
}

type ProcessDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin      string `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"`
	Source      string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Destination string `protobuf:"bytes,3,opt,name=destination,proto3" json:"destination,omitempty"`
	Data        int64  `protobuf:"varint,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProcessDataResponse) Reset() {
	*x = ProcessDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_myservice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessDataResponse) ProtoMessage() {}

func (x *ProcessDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_myservice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessDataResponse.ProtoReflect.Descriptor instead.
func (*ProcessDataResponse) Descriptor() ([]byte, []int) {
	return file_proto_myservice_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessDataResponse) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *ProcessDataResponse) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *ProcessDataResponse) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

func (x *ProcessDataResponse) GetData() int64 {
	if x != nil {
		return x.Data
	}
	return 0
}

var File_proto_myservice_proto protoreflect.FileDescriptor

var file_proto_myservice_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x22, 0x7a, 0x0a, 0x12, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x74,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x7b,
	0x0a, 0x13, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x74,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x5b, 0x0a, 0x09, 0x4d,
	0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1d, 0x2e, 0x6d, 0x79, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6d, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_myservice_proto_rawDescOnce sync.Once
	file_proto_myservice_proto_rawDescData = file_proto_myservice_proto_rawDesc
)

func file_proto_myservice_proto_rawDescGZIP() []byte {
	file_proto_myservice_proto_rawDescOnce.Do(func() {
		file_proto_myservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_myservice_proto_rawDescData)
	})
	return file_proto_myservice_proto_rawDescData
}

var file_proto_myservice_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_myservice_proto_goTypes = []any{
	(*ProcessDataRequest)(nil),  // 0: myservice.ProcessDataRequest
	(*ProcessDataResponse)(nil), // 1: myservice.ProcessDataResponse
}
var file_proto_myservice_proto_depIdxs = []int32{
	0, // 0: myservice.MyService.ProcessData:input_type -> myservice.ProcessDataRequest
	1, // 1: myservice.MyService.ProcessData:output_type -> myservice.ProcessDataResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_myservice_proto_init() }
func file_proto_myservice_proto_init() {
	if File_proto_myservice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_myservice_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ProcessDataRequest); i {
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
		file_proto_myservice_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProcessDataResponse); i {
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
			RawDescriptor: file_proto_myservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_myservice_proto_goTypes,
		DependencyIndexes: file_proto_myservice_proto_depIdxs,
		MessageInfos:      file_proto_myservice_proto_msgTypes,
	}.Build()
	File_proto_myservice_proto = out.File
	file_proto_myservice_proto_rawDesc = nil
	file_proto_myservice_proto_goTypes = nil
	file_proto_myservice_proto_depIdxs = nil
}
