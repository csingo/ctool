// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: proto/app/enum.proto

package app

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

type Human int32

const (
	Human_NIL   Human = 0
	Human_MAN   Human = 1
	Human_WOMAN Human = 2
)

// Enum value maps for Human.
var (
	Human_name = map[int32]string{
		0: "NIL",
		1: "MAN",
		2: "WOMAN",
	}
	Human_value = map[string]int32{
		"NIL":   0,
		"MAN":   1,
		"WOMAN": 2,
	}
)

func (x Human) Enum() *Human {
	p := new(Human)
	*p = x
	return p
}

func (x Human) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Human) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_app_enum_proto_enumTypes[0].Descriptor()
}

func (Human) Type() protoreflect.EnumType {
	return &file_proto_app_enum_proto_enumTypes[0]
}

func (x Human) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Human.Descriptor instead.
func (Human) EnumDescriptor() ([]byte, []int) {
	return file_proto_app_enum_proto_rawDescGZIP(), []int{0}
}

var File_proto_app_enum_proto protoreflect.FileDescriptor

var file_proto_app_enum_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x65, 0x6e, 0x75, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x24, 0x0a, 0x05, 0x48, 0x75, 0x6d, 0x61, 0x6e, 0x12,
	0x07, 0x0a, 0x03, 0x4e, 0x49, 0x4c, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x4d, 0x41, 0x4e, 0x10,
	0x01, 0x12, 0x09, 0x0a, 0x05, 0x57, 0x4f, 0x4d, 0x41, 0x4e, 0x10, 0x02, 0x42, 0x08, 0x5a, 0x06,
	0x2e, 0x2f, 0x3b, 0x61, 0x70, 0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_app_enum_proto_rawDescOnce sync.Once
	file_proto_app_enum_proto_rawDescData = file_proto_app_enum_proto_rawDesc
)

func file_proto_app_enum_proto_rawDescGZIP() []byte {
	file_proto_app_enum_proto_rawDescOnce.Do(func() {
		file_proto_app_enum_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_app_enum_proto_rawDescData)
	})
	return file_proto_app_enum_proto_rawDescData
}

var file_proto_app_enum_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_app_enum_proto_goTypes = []interface{}{
	(Human)(0), // 0: Human
}
var file_proto_app_enum_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_app_enum_proto_init() }
func file_proto_app_enum_proto_init() {
	if File_proto_app_enum_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_app_enum_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_app_enum_proto_goTypes,
		DependencyIndexes: file_proto_app_enum_proto_depIdxs,
		EnumInfos:         file_proto_app_enum_proto_enumTypes,
	}.Build()
	File_proto_app_enum_proto = out.File
	file_proto_app_enum_proto_rawDesc = nil
	file_proto_app_enum_proto_goTypes = nil
	file_proto_app_enum_proto_depIdxs = nil
}
