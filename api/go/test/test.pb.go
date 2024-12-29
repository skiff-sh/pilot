// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: test/test.proto

// For the purpose of testing various protobuf interactions.

package test

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

type Primitives struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	I32           int32                  `protobuf:"varint,1,opt,name=i32,proto3" json:"i32,omitempty"`
	I64           int64                  `protobuf:"varint,2,opt,name=i64,proto3" json:"i64,omitempty"`
	U32           uint32                 `protobuf:"varint,3,opt,name=u32,proto3" json:"u32,omitempty"`
	U64           uint64                 `protobuf:"varint,4,opt,name=u64,proto3" json:"u64,omitempty"`
	Str           string                 `protobuf:"bytes,5,opt,name=str,proto3" json:"str,omitempty"`
	Boolean       bool                   `protobuf:"varint,6,opt,name=boolean,proto3" json:"boolean,omitempty"`
	Dub           float64                `protobuf:"fixed64,7,opt,name=dub,proto3" json:"dub,omitempty"`
	F32           float32                `protobuf:"fixed32,8,opt,name=f32,proto3" json:"f32,omitempty"`
	Byt           []byte                 `protobuf:"bytes,9,opt,name=byt,proto3" json:"byt,omitempty"`
	I64Ptr        *int64                 `protobuf:"varint,10,opt,name=i64_ptr,json=i64Ptr,proto3,oneof" json:"i64_ptr,omitempty"`
	U32Ptr        *uint32                `protobuf:"varint,11,opt,name=u32_ptr,json=u32Ptr,proto3,oneof" json:"u32_ptr,omitempty"`
	U64Ptr        *uint64                `protobuf:"varint,12,opt,name=u64_ptr,json=u64Ptr,proto3,oneof" json:"u64_ptr,omitempty"`
	StrPtr        *string                `protobuf:"bytes,13,opt,name=str_ptr,json=strPtr,proto3,oneof" json:"str_ptr,omitempty"`
	BooleanPtr    *bool                  `protobuf:"varint,14,opt,name=boolean_ptr,json=booleanPtr,proto3,oneof" json:"boolean_ptr,omitempty"`
	DubPtr        *float64               `protobuf:"fixed64,15,opt,name=dub_ptr,json=dubPtr,proto3,oneof" json:"dub_ptr,omitempty"`
	F32Ptr        *float32               `protobuf:"fixed32,16,opt,name=f32_ptr,json=f32Ptr,proto3,oneof" json:"f32_ptr,omitempty"`
	BytPtr        []byte                 `protobuf:"bytes,17,opt,name=byt_ptr,json=bytPtr,proto3,oneof" json:"byt_ptr,omitempty"`
	I32Ptr        *int32                 `protobuf:"varint,18,opt,name=i32_ptr,json=i32Ptr,proto3,oneof" json:"i32_ptr,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Primitives) Reset() {
	*x = Primitives{}
	mi := &file_test_test_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Primitives) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Primitives) ProtoMessage() {}

func (x *Primitives) ProtoReflect() protoreflect.Message {
	mi := &file_test_test_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Primitives.ProtoReflect.Descriptor instead.
func (*Primitives) Descriptor() ([]byte, []int) {
	return file_test_test_proto_rawDescGZIP(), []int{0}
}

func (x *Primitives) GetI32() int32 {
	if x != nil {
		return x.I32
	}
	return 0
}

func (x *Primitives) GetI64() int64 {
	if x != nil {
		return x.I64
	}
	return 0
}

func (x *Primitives) GetU32() uint32 {
	if x != nil {
		return x.U32
	}
	return 0
}

func (x *Primitives) GetU64() uint64 {
	if x != nil {
		return x.U64
	}
	return 0
}

func (x *Primitives) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

func (x *Primitives) GetBoolean() bool {
	if x != nil {
		return x.Boolean
	}
	return false
}

func (x *Primitives) GetDub() float64 {
	if x != nil {
		return x.Dub
	}
	return 0
}

func (x *Primitives) GetF32() float32 {
	if x != nil {
		return x.F32
	}
	return 0
}

func (x *Primitives) GetByt() []byte {
	if x != nil {
		return x.Byt
	}
	return nil
}

func (x *Primitives) GetI64Ptr() int64 {
	if x != nil && x.I64Ptr != nil {
		return *x.I64Ptr
	}
	return 0
}

func (x *Primitives) GetU32Ptr() uint32 {
	if x != nil && x.U32Ptr != nil {
		return *x.U32Ptr
	}
	return 0
}

func (x *Primitives) GetU64Ptr() uint64 {
	if x != nil && x.U64Ptr != nil {
		return *x.U64Ptr
	}
	return 0
}

func (x *Primitives) GetStrPtr() string {
	if x != nil && x.StrPtr != nil {
		return *x.StrPtr
	}
	return ""
}

func (x *Primitives) GetBooleanPtr() bool {
	if x != nil && x.BooleanPtr != nil {
		return *x.BooleanPtr
	}
	return false
}

func (x *Primitives) GetDubPtr() float64 {
	if x != nil && x.DubPtr != nil {
		return *x.DubPtr
	}
	return 0
}

func (x *Primitives) GetF32Ptr() float32 {
	if x != nil && x.F32Ptr != nil {
		return *x.F32Ptr
	}
	return 0
}

func (x *Primitives) GetBytPtr() []byte {
	if x != nil {
		return x.BytPtr
	}
	return nil
}

func (x *Primitives) GetI32Ptr() int32 {
	if x != nil && x.I32Ptr != nil {
		return *x.I32Ptr
	}
	return 0
}

var File_test_test_proto protoreflect.FileDescriptor

var file_test_test_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x10, 0x73, 0x6b, 0x69, 0x66, 0x66, 0x2e, 0x70, 0x69, 0x6c, 0x6f, 0x74, 0x2e, 0x74,
	0x65, 0x73, 0x74, 0x22, 0xbc, 0x04, 0x0a, 0x0a, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76,
	0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x33, 0x32, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x69, 0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x36, 0x34, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x69, 0x36, 0x34, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x33, 0x32, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x36, 0x34, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x36, 0x34, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74,
	0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x12, 0x18, 0x0a, 0x07,
	0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x62,
	0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x75, 0x62, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x03, 0x64, 0x75, 0x62, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x33, 0x32, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x66, 0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x79,
	0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x62, 0x79, 0x74, 0x12, 0x1c, 0x0a, 0x07,
	0x69, 0x36, 0x34, 0x5f, 0x70, 0x74, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52,
	0x06, 0x69, 0x36, 0x34, 0x50, 0x74, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x75, 0x33,
	0x32, 0x5f, 0x70, 0x74, 0x72, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x01, 0x52, 0x06, 0x75,
	0x33, 0x32, 0x50, 0x74, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x75, 0x36, 0x34, 0x5f,
	0x70, 0x74, 0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04, 0x48, 0x02, 0x52, 0x06, 0x75, 0x36, 0x34,
	0x50, 0x74, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x73, 0x74, 0x72, 0x5f, 0x70, 0x74,
	0x72, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x06, 0x73, 0x74, 0x72, 0x50, 0x74,
	0x72, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a, 0x0b, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x5f,
	0x70, 0x74, 0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x48, 0x04, 0x52, 0x0a, 0x62, 0x6f, 0x6f,
	0x6c, 0x65, 0x61, 0x6e, 0x50, 0x74, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x64, 0x75,
	0x62, 0x5f, 0x70, 0x74, 0x72, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x01, 0x48, 0x05, 0x52, 0x06, 0x64,
	0x75, 0x62, 0x50, 0x74, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x66, 0x33, 0x32, 0x5f,
	0x70, 0x74, 0x72, 0x18, 0x10, 0x20, 0x01, 0x28, 0x02, 0x48, 0x06, 0x52, 0x06, 0x66, 0x33, 0x32,
	0x50, 0x74, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x62, 0x79, 0x74, 0x5f, 0x70, 0x74,
	0x72, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x07, 0x52, 0x06, 0x62, 0x79, 0x74, 0x50, 0x74,
	0x72, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x69, 0x33, 0x32, 0x5f, 0x70, 0x74, 0x72, 0x18,
	0x12, 0x20, 0x01, 0x28, 0x05, 0x48, 0x08, 0x52, 0x06, 0x69, 0x33, 0x32, 0x50, 0x74, 0x72, 0x88,
	0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x69, 0x36, 0x34, 0x5f, 0x70, 0x74, 0x72, 0x42, 0x0a,
	0x0a, 0x08, 0x5f, 0x75, 0x33, 0x32, 0x5f, 0x70, 0x74, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75,
	0x36, 0x34, 0x5f, 0x70, 0x74, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x73, 0x74, 0x72, 0x5f, 0x70,
	0x74, 0x72, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x5f, 0x70,
	0x74, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x64, 0x75, 0x62, 0x5f, 0x70, 0x74, 0x72, 0x42, 0x0a,
	0x0a, 0x08, 0x5f, 0x66, 0x33, 0x32, 0x5f, 0x70, 0x74, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x62,
	0x79, 0x74, 0x5f, 0x70, 0x74, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x69, 0x33, 0x32, 0x5f, 0x70,
	0x74, 0x72, 0x42, 0xb0, 0x01, 0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x6b, 0x69, 0x66, 0x66,
	0x2e, 0x70, 0x69, 0x6c, 0x6f, 0x74, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x42, 0x09, 0x54, 0x65, 0x73,
	0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6b, 0x69, 0x66, 0x66, 0x2d, 0x73, 0x68, 0x2f, 0x70, 0x69,
	0x6c, 0x6f, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x69, 0x6c, 0x6f, 0x74,
	0x2f, 0x74, 0x65, 0x73, 0x74, 0xa2, 0x02, 0x03, 0x53, 0x50, 0x54, 0xaa, 0x02, 0x10, 0x53, 0x6b,
	0x69, 0x66, 0x66, 0x2e, 0x50, 0x69, 0x6c, 0x6f, 0x74, 0x2e, 0x54, 0x65, 0x73, 0x74, 0xca, 0x02,
	0x10, 0x53, 0x6b, 0x69, 0x66, 0x66, 0x5c, 0x50, 0x69, 0x6c, 0x6f, 0x74, 0x5c, 0x54, 0x65, 0x73,
	0x74, 0xe2, 0x02, 0x1c, 0x53, 0x6b, 0x69, 0x66, 0x66, 0x5c, 0x50, 0x69, 0x6c, 0x6f, 0x74, 0x5c,
	0x54, 0x65, 0x73, 0x74, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x12, 0x53, 0x6b, 0x69, 0x66, 0x66, 0x3a, 0x3a, 0x50, 0x69, 0x6c, 0x6f, 0x74, 0x3a,
	0x3a, 0x54, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_test_test_proto_rawDescOnce sync.Once
	file_test_test_proto_rawDescData = file_test_test_proto_rawDesc
)

func file_test_test_proto_rawDescGZIP() []byte {
	file_test_test_proto_rawDescOnce.Do(func() {
		file_test_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_test_proto_rawDescData)
	})
	return file_test_test_proto_rawDescData
}

var file_test_test_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_test_test_proto_goTypes = []any{
	(*Primitives)(nil), // 0: skiff.pilot.test.Primitives
}
var file_test_test_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_test_proto_init() }
func file_test_test_proto_init() {
	if File_test_test_proto != nil {
		return
	}
	file_test_test_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_test_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_test_proto_goTypes,
		DependencyIndexes: file_test_test_proto_depIdxs,
		MessageInfos:      file_test_test_proto_msgTypes,
	}.Build()
	File_test_test_proto = out.File
	file_test_test_proto_rawDesc = nil
	file_test_test_proto_goTypes = nil
	file_test_test_proto_depIdxs = nil
}
