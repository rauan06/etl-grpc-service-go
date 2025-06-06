// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: protos/price/v1/common.proto

package price

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListParamsSt struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int64                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int64                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Sort          []string               `protobuf:"bytes,3,rep,name=sort,proto3" json:"sort,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListParamsSt) Reset() {
	*x = ListParamsSt{}
	mi := &file_protos_price_v1_common_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListParamsSt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListParamsSt) ProtoMessage() {}

func (x *ListParamsSt) ProtoReflect() protoreflect.Message {
	mi := &file_protos_price_v1_common_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListParamsSt.ProtoReflect.Descriptor instead.
func (*ListParamsSt) Descriptor() ([]byte, []int) {
	return file_protos_price_v1_common_proto_rawDescGZIP(), []int{0}
}

func (x *ListParamsSt) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListParamsSt) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListParamsSt) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

type PaginationInfoSt struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int64                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int64                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PaginationInfoSt) Reset() {
	*x = PaginationInfoSt{}
	mi := &file_protos_price_v1_common_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaginationInfoSt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationInfoSt) ProtoMessage() {}

func (x *PaginationInfoSt) ProtoReflect() protoreflect.Message {
	mi := &file_protos_price_v1_common_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaginationInfoSt.ProtoReflect.Descriptor instead.
func (*PaginationInfoSt) Descriptor() ([]byte, []int) {
	return file_protos_price_v1_common_proto_rawDescGZIP(), []int{1}
}

func (x *PaginationInfoSt) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *PaginationInfoSt) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type ErrorRep struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Fields        map[string]string      `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ErrorRep) Reset() {
	*x = ErrorRep{}
	mi := &file_protos_price_v1_common_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ErrorRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorRep) ProtoMessage() {}

func (x *ErrorRep) ProtoReflect() protoreflect.Message {
	mi := &file_protos_price_v1_common_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorRep.ProtoReflect.Descriptor instead.
func (*ErrorRep) Descriptor() ([]byte, []int) {
	return file_protos_price_v1_common_proto_rawDescGZIP(), []int{2}
}

func (x *ErrorRep) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *ErrorRep) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ErrorRep) GetFields() map[string]string {
	if x != nil {
		return x.Fields
	}
	return nil
}

var File_protos_price_v1_common_proto protoreflect.FileDescriptor

const file_protos_price_v1_common_proto_rawDesc = "" +
	"\n" +
	"\x1cprotos/price/v1/common.proto\x12\x05price\"S\n" +
	"\fListParamsSt\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x03R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\x03R\bpageSize\x12\x12\n" +
	"\x04sort\x18\x03 \x03(\tR\x04sort\"C\n" +
	"\x10PaginationInfoSt\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x03R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\x03R\bpageSize\"\xa8\x01\n" +
	"\bErrorRep\x12\x12\n" +
	"\x04code\x18\x01 \x01(\tR\x04code\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x123\n" +
	"\x06fields\x18\x03 \x03(\v2\x1b.price.ErrorRep.FieldsEntryR\x06fields\x1a9\n" +
	"\vFieldsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01B\x11Z\x0f/price/v1;priceb\x06proto3"

var (
	file_protos_price_v1_common_proto_rawDescOnce sync.Once
	file_protos_price_v1_common_proto_rawDescData []byte
)

func file_protos_price_v1_common_proto_rawDescGZIP() []byte {
	file_protos_price_v1_common_proto_rawDescOnce.Do(func() {
		file_protos_price_v1_common_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protos_price_v1_common_proto_rawDesc), len(file_protos_price_v1_common_proto_rawDesc)))
	})
	return file_protos_price_v1_common_proto_rawDescData
}

var file_protos_price_v1_common_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protos_price_v1_common_proto_goTypes = []any{
	(*ListParamsSt)(nil),     // 0: price.ListParamsSt
	(*PaginationInfoSt)(nil), // 1: price.PaginationInfoSt
	(*ErrorRep)(nil),         // 2: price.ErrorRep
	nil,                      // 3: price.ErrorRep.FieldsEntry
}
var file_protos_price_v1_common_proto_depIdxs = []int32{
	3, // 0: price.ErrorRep.fields:type_name -> price.ErrorRep.FieldsEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protos_price_v1_common_proto_init() }
func file_protos_price_v1_common_proto_init() {
	if File_protos_price_v1_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protos_price_v1_common_proto_rawDesc), len(file_protos_price_v1_common_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_price_v1_common_proto_goTypes,
		DependencyIndexes: file_protos_price_v1_common_proto_depIdxs,
		MessageInfos:      file_protos_price_v1_common_proto_msgTypes,
	}.Build()
	File_protos_price_v1_common_proto = out.File
	file_protos_price_v1_common_proto_goTypes = nil
	file_protos_price_v1_common_proto_depIdxs = nil
}
