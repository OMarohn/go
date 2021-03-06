// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: coaster.proto

package grpcCoaster

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

type CoasterMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Manufacture string `protobuf:"bytes,2,opt,name=manufacture,proto3" json:"manufacture,omitempty"`
	Id          string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Height      uint32 `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *CoasterMessage) Reset() {
	*x = CoasterMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coaster_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoasterMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoasterMessage) ProtoMessage() {}

func (x *CoasterMessage) ProtoReflect() protoreflect.Message {
	mi := &file_coaster_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoasterMessage.ProtoReflect.Descriptor instead.
func (*CoasterMessage) Descriptor() ([]byte, []int) {
	return file_coaster_proto_rawDescGZIP(), []int{0}
}

func (x *CoasterMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CoasterMessage) GetManufacture() string {
	if x != nil {
		return x.Manufacture
	}
	return ""
}

func (x *CoasterMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CoasterMessage) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

type CoastersMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Coasters []*CoasterMessage `protobuf:"bytes,1,rep,name=coasters,proto3" json:"coasters,omitempty"`
}

func (x *CoastersMessage) Reset() {
	*x = CoastersMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coaster_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoastersMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoastersMessage) ProtoMessage() {}

func (x *CoastersMessage) ProtoReflect() protoreflect.Message {
	mi := &file_coaster_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoastersMessage.ProtoReflect.Descriptor instead.
func (*CoastersMessage) Descriptor() ([]byte, []int) {
	return file_coaster_proto_rawDescGZIP(), []int{1}
}

func (x *CoastersMessage) GetCoasters() []*CoasterMessage {
	if x != nil {
		return x.Coasters
	}
	return nil
}

type CoasterIDMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CoasterIDMessage) Reset() {
	*x = CoasterIDMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coaster_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoasterIDMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoasterIDMessage) ProtoMessage() {}

func (x *CoasterIDMessage) ProtoReflect() protoreflect.Message {
	mi := &file_coaster_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoasterIDMessage.ProtoReflect.Descriptor instead.
func (*CoasterIDMessage) Descriptor() ([]byte, []int) {
	return file_coaster_proto_rawDescGZIP(), []int{2}
}

func (x *CoasterIDMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coaster_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_coaster_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_coaster_proto_rawDescGZIP(), []int{3}
}

var File_coaster_proto protoreflect.FileDescriptor

var file_coaster_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x63, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x22, 0x6e, 0x0a, 0x0e, 0x43, 0x6f, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x46, 0x0a, 0x0f, 0x43, 0x6f, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x63,
	0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x63, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x63, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x73,
	0x22, 0x22, 0x0a, 0x10, 0x43, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x49, 0x44, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xcb, 0x01,
	0x0a, 0x0e, 0x43, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x39, 0x0a, 0x0b, 0x67, 0x65, 0x74, 0x43, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x73, 0x12,
	0x0e, 0x2e, 0x63, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x18, 0x2e, 0x63, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x0a, 0x67,
	0x65, 0x74, 0x43, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x63, 0x6f, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x49, 0x44, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x17, 0x2e, 0x63, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43,
	0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12,
	0x3a, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x12, 0x17, 0x2e, 0x63, 0x6f, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0e, 0x2e, 0x63, 0x6f, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x26, 0x5a, 0x24, 0x6b,
	0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x2f, 0x76, 0x32, 0x2f, 0x73, 0x72, 0x63, 0x2f,
	0x6b, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x43, 0x6f, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coaster_proto_rawDescOnce sync.Once
	file_coaster_proto_rawDescData = file_coaster_proto_rawDesc
)

func file_coaster_proto_rawDescGZIP() []byte {
	file_coaster_proto_rawDescOnce.Do(func() {
		file_coaster_proto_rawDescData = protoimpl.X.CompressGZIP(file_coaster_proto_rawDescData)
	})
	return file_coaster_proto_rawDescData
}

var file_coaster_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_coaster_proto_goTypes = []interface{}{
	(*CoasterMessage)(nil),   // 0: coaster.CoasterMessage
	(*CoastersMessage)(nil),  // 1: coaster.CoastersMessage
	(*CoasterIDMessage)(nil), // 2: coaster.CoasterIDMessage
	(*Empty)(nil),            // 3: coaster.Empty
}
var file_coaster_proto_depIdxs = []int32{
	0, // 0: coaster.CoastersMessage.coasters:type_name -> coaster.CoasterMessage
	3, // 1: coaster.CoasterService.getCoasters:input_type -> coaster.Empty
	2, // 2: coaster.CoasterService.getCoaster:input_type -> coaster.CoasterIDMessage
	0, // 3: coaster.CoasterService.createCoaster:input_type -> coaster.CoasterMessage
	1, // 4: coaster.CoasterService.getCoasters:output_type -> coaster.CoastersMessage
	0, // 5: coaster.CoasterService.getCoaster:output_type -> coaster.CoasterMessage
	3, // 6: coaster.CoasterService.createCoaster:output_type -> coaster.Empty
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_coaster_proto_init() }
func file_coaster_proto_init() {
	if File_coaster_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_coaster_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoasterMessage); i {
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
		file_coaster_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoastersMessage); i {
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
		file_coaster_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoasterIDMessage); i {
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
		file_coaster_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_coaster_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_coaster_proto_goTypes,
		DependencyIndexes: file_coaster_proto_depIdxs,
		MessageInfos:      file_coaster_proto_msgTypes,
	}.Build()
	File_coaster_proto = out.File
	file_coaster_proto_rawDesc = nil
	file_coaster_proto_goTypes = nil
	file_coaster_proto_depIdxs = nil
}
