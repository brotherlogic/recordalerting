// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: recordalerting.proto

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

type Problem_ProblemType int32

const (
	Problem_UNKNOWN            Problem_ProblemType = 0
	Problem_MISSING_WEIGHT     Problem_ProblemType = 1
	Problem_MISSING_FILED      Problem_ProblemType = 2
	Problem_MISSING_WIDTH      Problem_ProblemType = 3
	Problem_MISSING_CONDITION  Problem_ProblemType = 4
	Problem_MISSING_SLEEVE     Problem_ProblemType = 5
	Problem_NEEDS_KEEPER       Problem_ProblemType = 6
	Problem_NEEDS_DIGITAL      Problem_ProblemType = 7
	Problem_NEEDS_SALE_BUDGET  Problem_ProblemType = 8
	Problem_NEEDS_SOLD_DETAILS Problem_ProblemType = 9
	Problem_BAD_BANDCAMP       Problem_ProblemType = 10
	Problem_STALE_LIMBO        Problem_ProblemType = 11
)

// Enum value maps for Problem_ProblemType.
var (
	Problem_ProblemType_name = map[int32]string{
		0:  "UNKNOWN",
		1:  "MISSING_WEIGHT",
		2:  "MISSING_FILED",
		3:  "MISSING_WIDTH",
		4:  "MISSING_CONDITION",
		5:  "MISSING_SLEEVE",
		6:  "NEEDS_KEEPER",
		7:  "NEEDS_DIGITAL",
		8:  "NEEDS_SALE_BUDGET",
		9:  "NEEDS_SOLD_DETAILS",
		10: "BAD_BANDCAMP",
		11: "STALE_LIMBO",
	}
	Problem_ProblemType_value = map[string]int32{
		"UNKNOWN":            0,
		"MISSING_WEIGHT":     1,
		"MISSING_FILED":      2,
		"MISSING_WIDTH":      3,
		"MISSING_CONDITION":  4,
		"MISSING_SLEEVE":     5,
		"NEEDS_KEEPER":       6,
		"NEEDS_DIGITAL":      7,
		"NEEDS_SALE_BUDGET":  8,
		"NEEDS_SOLD_DETAILS": 9,
		"BAD_BANDCAMP":       10,
		"STALE_LIMBO":        11,
	}
)

func (x Problem_ProblemType) Enum() *Problem_ProblemType {
	p := new(Problem_ProblemType)
	*p = x
	return p
}

func (x Problem_ProblemType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Problem_ProblemType) Descriptor() protoreflect.EnumDescriptor {
	return file_recordalerting_proto_enumTypes[0].Descriptor()
}

func (Problem_ProblemType) Type() protoreflect.EnumType {
	return &file_recordalerting_proto_enumTypes[0]
}

func (x Problem_ProblemType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Problem_ProblemType.Descriptor instead.
func (Problem_ProblemType) EnumDescriptor() ([]byte, []int) {
	return file_recordalerting_proto_rawDescGZIP(), []int{1, 0}
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Problems []*Problem `protobuf:"bytes,1,rep,name=problems,proto3" json:"problems,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordalerting_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_recordalerting_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_recordalerting_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetProblems() []*Problem {
	if x != nil {
		return x.Problems
	}
	return nil
}

type Problem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type                 Problem_ProblemType `protobuf:"varint,1,opt,name=type,proto3,enum=recordalerting.Problem_ProblemType" json:"type,omitempty"`
	InstanceId           int32               `protobuf:"varint,2,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	IssueNumber          int32               `protobuf:"varint,3,opt,name=issue_number,json=issueNumber,proto3" json:"issue_number,omitempty"`
	IssueOpenedTimestamp int64               `protobuf:"varint,4,opt,name=issue_opened_timestamp,json=issueOpenedTimestamp,proto3" json:"issue_opened_timestamp,omitempty"`
}

func (x *Problem) Reset() {
	*x = Problem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordalerting_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Problem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Problem) ProtoMessage() {}

func (x *Problem) ProtoReflect() protoreflect.Message {
	mi := &file_recordalerting_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Problem.ProtoReflect.Descriptor instead.
func (*Problem) Descriptor() ([]byte, []int) {
	return file_recordalerting_proto_rawDescGZIP(), []int{1}
}

func (x *Problem) GetType() Problem_ProblemType {
	if x != nil {
		return x.Type
	}
	return Problem_UNKNOWN
}

func (x *Problem) GetInstanceId() int32 {
	if x != nil {
		return x.InstanceId
	}
	return 0
}

func (x *Problem) GetIssueNumber() int32 {
	if x != nil {
		return x.IssueNumber
	}
	return 0
}

func (x *Problem) GetIssueOpenedTimestamp() int64 {
	if x != nil {
		return x.IssueOpenedTimestamp
	}
	return 0
}

var File_recordalerting_proto protoreflect.FileDescriptor

var file_recordalerting_proto_rawDesc = []byte{
	0x0a, 0x14, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x6c,
	0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x3d, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x33, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x6c, 0x65, 0x72, 0x74,
	0x69, 0x6e, 0x67, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x52, 0x08, 0x70, 0x72, 0x6f,
	0x62, 0x6c, 0x65, 0x6d, 0x73, 0x22, 0xb5, 0x03, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65,
	0x6d, 0x12, 0x37, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x23, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67,
	0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x69,
	0x73, 0x73, 0x75, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0b, 0x69, 0x73, 0x73, 0x75, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x34,
	0x0a, 0x16, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f, 0x6f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x14,
	0x69, 0x73, 0x73, 0x75, 0x65, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x22, 0xf6, 0x01, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x12, 0x0a, 0x0e, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x5f, 0x57, 0x45, 0x49,
	0x47, 0x48, 0x54, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4e, 0x47,
	0x5f, 0x46, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x4d, 0x49, 0x53, 0x53,
	0x49, 0x4e, 0x47, 0x5f, 0x57, 0x49, 0x44, 0x54, 0x48, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11, 0x4d,
	0x49, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x5f, 0x43, 0x4f, 0x4e, 0x44, 0x49, 0x54, 0x49, 0x4f, 0x4e,
	0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x5f, 0x53, 0x4c,
	0x45, 0x45, 0x56, 0x45, 0x10, 0x05, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x45, 0x45, 0x44, 0x53, 0x5f,
	0x4b, 0x45, 0x45, 0x50, 0x45, 0x52, 0x10, 0x06, 0x12, 0x11, 0x0a, 0x0d, 0x4e, 0x45, 0x45, 0x44,
	0x53, 0x5f, 0x44, 0x49, 0x47, 0x49, 0x54, 0x41, 0x4c, 0x10, 0x07, 0x12, 0x15, 0x0a, 0x11, 0x4e,
	0x45, 0x45, 0x44, 0x53, 0x5f, 0x53, 0x41, 0x4c, 0x45, 0x5f, 0x42, 0x55, 0x44, 0x47, 0x45, 0x54,
	0x10, 0x08, 0x12, 0x16, 0x0a, 0x12, 0x4e, 0x45, 0x45, 0x44, 0x53, 0x5f, 0x53, 0x4f, 0x4c, 0x44,
	0x5f, 0x44, 0x45, 0x54, 0x41, 0x49, 0x4c, 0x53, 0x10, 0x09, 0x12, 0x10, 0x0a, 0x0c, 0x42, 0x41,
	0x44, 0x5f, 0x42, 0x41, 0x4e, 0x44, 0x43, 0x41, 0x4d, 0x50, 0x10, 0x0a, 0x12, 0x0f, 0x0a, 0x0b,
	0x53, 0x54, 0x41, 0x4c, 0x45, 0x5f, 0x4c, 0x49, 0x4d, 0x42, 0x4f, 0x10, 0x0b, 0x42, 0x2e, 0x5a,
	0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x72, 0x6f, 0x74,
	0x68, 0x65, 0x72, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x61,
	0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_recordalerting_proto_rawDescOnce sync.Once
	file_recordalerting_proto_rawDescData = file_recordalerting_proto_rawDesc
)

func file_recordalerting_proto_rawDescGZIP() []byte {
	file_recordalerting_proto_rawDescOnce.Do(func() {
		file_recordalerting_proto_rawDescData = protoimpl.X.CompressGZIP(file_recordalerting_proto_rawDescData)
	})
	return file_recordalerting_proto_rawDescData
}

var file_recordalerting_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_recordalerting_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_recordalerting_proto_goTypes = []interface{}{
	(Problem_ProblemType)(0), // 0: recordalerting.Problem.ProblemType
	(*Config)(nil),           // 1: recordalerting.Config
	(*Problem)(nil),          // 2: recordalerting.Problem
}
var file_recordalerting_proto_depIdxs = []int32{
	2, // 0: recordalerting.Config.problems:type_name -> recordalerting.Problem
	0, // 1: recordalerting.Problem.type:type_name -> recordalerting.Problem.ProblemType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_recordalerting_proto_init() }
func file_recordalerting_proto_init() {
	if File_recordalerting_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_recordalerting_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
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
		file_recordalerting_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Problem); i {
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
			RawDescriptor: file_recordalerting_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_recordalerting_proto_goTypes,
		DependencyIndexes: file_recordalerting_proto_depIdxs,
		EnumInfos:         file_recordalerting_proto_enumTypes,
		MessageInfos:      file_recordalerting_proto_msgTypes,
	}.Build()
	File_recordalerting_proto = out.File
	file_recordalerting_proto_rawDesc = nil
	file_recordalerting_proto_goTypes = nil
	file_recordalerting_proto_depIdxs = nil
}
