// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: appuser.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// RESOURCES
type AppUserStatus int32

const (
	AppUserStatus_APP_USER_STATUS_UNSPECIFIED AppUserStatus = 0
	AppUserStatus_APP_USER_STATUS_INACTIVE    AppUserStatus = 1
	AppUserStatus_APP_USER_STATUS_ONLINE      AppUserStatus = 2
	AppUserStatus_APP_USER_STATUS_OFFLINE     AppUserStatus = 3
	AppUserStatus_APP_USER_STATUS_AWAY        AppUserStatus = 4
)

// Enum value maps for AppUserStatus.
var (
	AppUserStatus_name = map[int32]string{
		0: "APP_USER_STATUS_UNSPECIFIED",
		1: "APP_USER_STATUS_INACTIVE",
		2: "APP_USER_STATUS_ONLINE",
		3: "APP_USER_STATUS_OFFLINE",
		4: "APP_USER_STATUS_AWAY",
	}
	AppUserStatus_value = map[string]int32{
		"APP_USER_STATUS_UNSPECIFIED": 0,
		"APP_USER_STATUS_INACTIVE":    1,
		"APP_USER_STATUS_ONLINE":      2,
		"APP_USER_STATUS_OFFLINE":     3,
		"APP_USER_STATUS_AWAY":        4,
	}
)

func (x AppUserStatus) Enum() *AppUserStatus {
	p := new(AppUserStatus)
	*p = x
	return p
}

func (x AppUserStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AppUserStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_appuser_proto_enumTypes[0].Descriptor()
}

func (AppUserStatus) Type() protoreflect.EnumType {
	return &file_appuser_proto_enumTypes[0]
}

func (x AppUserStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AppUserStatus.Descriptor instead.
func (AppUserStatus) EnumDescriptor() ([]byte, []int) {
	return file_appuser_proto_rawDescGZIP(), []int{0}
}

type Appuser struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	OnlineStatus  AppUserStatus          `protobuf:"varint,3,opt,name=online_status,json=onlineStatus,proto3,enum=v1.appuser.AppUserStatus" json:"online_status,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Appuser) Reset() {
	*x = Appuser{}
	mi := &file_appuser_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Appuser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Appuser) ProtoMessage() {}

func (x *Appuser) ProtoReflect() protoreflect.Message {
	mi := &file_appuser_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Appuser.ProtoReflect.Descriptor instead.
func (*Appuser) Descriptor() ([]byte, []int) {
	return file_appuser_proto_rawDescGZIP(), []int{0}
}

func (x *Appuser) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Appuser) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Appuser) GetOnlineStatus() AppUserStatus {
	if x != nil {
		return x.OnlineStatus
	}
	return AppUserStatus_APP_USER_STATUS_UNSPECIFIED
}

func (x *Appuser) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Appuser) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// ----- APPUSER -----
type CreateAppuserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateAppuserRequest) Reset() {
	*x = CreateAppuserRequest{}
	mi := &file_appuser_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateAppuserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAppuserRequest) ProtoMessage() {}

func (x *CreateAppuserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_appuser_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAppuserRequest.ProtoReflect.Descriptor instead.
func (*CreateAppuserRequest) Descriptor() ([]byte, []int) {
	return file_appuser_proto_rawDescGZIP(), []int{1}
}

func (x *CreateAppuserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreateAppuserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type CreateAppuserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateAppuserResponse) Reset() {
	*x = CreateAppuserResponse{}
	mi := &file_appuser_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateAppuserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAppuserResponse) ProtoMessage() {}

func (x *CreateAppuserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_appuser_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAppuserResponse.ProtoReflect.Descriptor instead.
func (*CreateAppuserResponse) Descriptor() ([]byte, []int) {
	return file_appuser_proto_rawDescGZIP(), []int{2}
}

var File_appuser_proto protoreflect.FileDescriptor

var file_appuser_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x76, 0x31, 0x2e, 0x61, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xeb, 0x01, 0x0a,
	0x07, 0x41, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x0d, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x76, 0x31,
	0x2e, 0x61, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x41, 0x70, 0x70, 0x55, 0x73, 0x65, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0c, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x42, 0x0a, 0x14, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x17,
	0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0xa1, 0x01, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x1b, 0x41, 0x50, 0x50,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x41, 0x50,
	0x50, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e,
	0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x41, 0x50, 0x50, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x4f, 0x4e, 0x4c, 0x49,
	0x4e, 0x45, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x41, 0x50, 0x50, 0x5f, 0x55, 0x53, 0x45, 0x52,
	0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x4f, 0x46, 0x46, 0x4c, 0x49, 0x4e, 0x45, 0x10,
	0x03, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x50, 0x50, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x41, 0x57, 0x41, 0x59, 0x10, 0x04, 0x32, 0x66, 0x0a, 0x0e, 0x41,
	0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a,
	0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x12, 0x20,
	0x2e, 0x76, 0x31, 0x2e, 0x61, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x41, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x21, 0x2e, 0x76, 0x31, 0x2e, 0x61, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x75, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_appuser_proto_rawDescOnce sync.Once
	file_appuser_proto_rawDescData = file_appuser_proto_rawDesc
)

func file_appuser_proto_rawDescGZIP() []byte {
	file_appuser_proto_rawDescOnce.Do(func() {
		file_appuser_proto_rawDescData = protoimpl.X.CompressGZIP(file_appuser_proto_rawDescData)
	})
	return file_appuser_proto_rawDescData
}

var file_appuser_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_appuser_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_appuser_proto_goTypes = []any{
	(AppUserStatus)(0),            // 0: v1.appuser.AppUserStatus
	(*Appuser)(nil),               // 1: v1.appuser.Appuser
	(*CreateAppuserRequest)(nil),  // 2: v1.appuser.CreateAppuserRequest
	(*CreateAppuserResponse)(nil), // 3: v1.appuser.CreateAppuserResponse
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_appuser_proto_depIdxs = []int32{
	0, // 0: v1.appuser.Appuser.online_status:type_name -> v1.appuser.AppUserStatus
	4, // 1: v1.appuser.Appuser.created_at:type_name -> google.protobuf.Timestamp
	4, // 2: v1.appuser.Appuser.updated_at:type_name -> google.protobuf.Timestamp
	2, // 3: v1.appuser.AppuserService.CreateAppuser:input_type -> v1.appuser.CreateAppuserRequest
	3, // 4: v1.appuser.AppuserService.CreateAppuser:output_type -> v1.appuser.CreateAppuserResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_appuser_proto_init() }
func file_appuser_proto_init() {
	if File_appuser_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_appuser_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_appuser_proto_goTypes,
		DependencyIndexes: file_appuser_proto_depIdxs,
		EnumInfos:         file_appuser_proto_enumTypes,
		MessageInfos:      file_appuser_proto_msgTypes,
	}.Build()
	File_appuser_proto = out.File
	file_appuser_proto_rawDesc = nil
	file_appuser_proto_goTypes = nil
	file_appuser_proto_depIdxs = nil
}
