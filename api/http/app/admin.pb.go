// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: http/app/admin.proto

package app

import (
	_ "github.com/tenz-io/gokit/genproto/go/custom/common"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type AdminLoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @inject_tag: bind:"form,name=username" validate:"required,non_blank,min_len=2,pattern=#abc123"
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty" bind:"form,name=username" validate:"required,non_blank,min_len=2,pattern=#abc123"`
	// @inject_tag: bind:"form,name=password" validate:"required,non_blank,min_len=4"
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty" bind:"form,name=password" validate:"required,non_blank,min_len=4"`
}

func (x *AdminLoginRequest) Reset() {
	*x = AdminLoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_http_app_admin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AdminLoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminLoginRequest) ProtoMessage() {}

func (x *AdminLoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_http_app_admin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminLoginRequest.ProtoReflect.Descriptor instead.
func (*AdminLoginRequest) Descriptor() ([]byte, []int) {
	return file_http_app_admin_proto_rawDescGZIP(), []int{0}
}

func (x *AdminLoginRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *AdminLoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type AdminLoginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AdminLoginResponse) Reset() {
	*x = AdminLoginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_http_app_admin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AdminLoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminLoginResponse) ProtoMessage() {}

func (x *AdminLoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_http_app_admin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminLoginResponse.ProtoReflect.Descriptor instead.
func (*AdminLoginResponse) Descriptor() ([]byte, []int) {
	return file_http_app_admin_proto_rawDescGZIP(), []int{1}
}

type AdminAddTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @inject_tag: bind:"form,name=userid" validate="required,gt=0"
	Userid int64 `protobuf:"varint,1,opt,name=userid,proto3" json:"userid,omitempty" bind:"form,name=userid"`
	// @inject_tag: bind:"form,name=expire" validate="required,gt=0" default:"60"
	Expire int32 `protobuf:"varint,2,opt,name=expire,proto3" json:"expire,omitempty" bind:"form,name=expire" default:"60"` // expire time in days
}

func (x *AdminAddTokenRequest) Reset() {
	*x = AdminAddTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_http_app_admin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AdminAddTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminAddTokenRequest) ProtoMessage() {}

func (x *AdminAddTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_http_app_admin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminAddTokenRequest.ProtoReflect.Descriptor instead.
func (*AdminAddTokenRequest) Descriptor() ([]byte, []int) {
	return file_http_app_admin_proto_rawDescGZIP(), []int{2}
}

func (x *AdminAddTokenRequest) GetUserid() int64 {
	if x != nil {
		return x.Userid
	}
	return 0
}

func (x *AdminAddTokenRequest) GetExpire() int32 {
	if x != nil {
		return x.Expire
	}
	return 0
}

type AdminAddTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *AdminAddTokenResponse) Reset() {
	*x = AdminAddTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_http_app_admin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AdminAddTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminAddTokenResponse) ProtoMessage() {}

func (x *AdminAddTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_http_app_admin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminAddTokenResponse.ProtoReflect.Descriptor instead.
func (*AdminAddTokenResponse) Descriptor() ([]byte, []int) {
	return file_http_app_admin_proto_rawDescGZIP(), []int{3}
}

func (x *AdminAddTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

var File_http_app_admin_proto protoreflect.FileDescriptor

var file_http_app_admin_proto_rawDesc = []byte{
	0x0a, 0x14, 0x68, 0x74, 0x74, 0x70, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x61, 0x70, 0x70,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19,
	0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x61, 0x75,
	0x74, 0x68, 0x7a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4b, 0x0a, 0x11, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x14, 0x0a, 0x12, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x46, 0x0a, 0x14,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x65, 0x22, 0x3a, 0x0a, 0x15, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x64, 0x64,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x32, 0xd8, 0x01, 0x0a, 0x0b, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x12, 0x5c, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1b, 0x2e, 0x68, 0x74, 0x74, 0x70,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x61, 0x70,
	0x70, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x8a, 0xb5, 0x18, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e,
	0x22, 0x0c, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x6b,
	0x0a, 0x08, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1e, 0x2e, 0x68, 0x74, 0x74,
	0x70, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x68, 0x74, 0x74,
	0x70, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x8a, 0xb5, 0x18,
	0x02, 0x08, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x22, 0x10, 0x2f, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2f, 0x61, 0x64, 0x64, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x22, 0x5a, 0x20, 0x67,
	0x6f, 0x2d, 0x77, 0x65, 0x62, 0x2d, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2f, 0x61, 0x70, 0x70, 0x3b, 0x61, 0x70, 0x70, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_http_app_admin_proto_rawDescOnce sync.Once
	file_http_app_admin_proto_rawDescData = file_http_app_admin_proto_rawDesc
)

func file_http_app_admin_proto_rawDescGZIP() []byte {
	file_http_app_admin_proto_rawDescOnce.Do(func() {
		file_http_app_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_http_app_admin_proto_rawDescData)
	})
	return file_http_app_admin_proto_rawDescData
}

var file_http_app_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_http_app_admin_proto_goTypes = []any{
	(*AdminLoginRequest)(nil),     // 0: http.app.AdminLoginRequest
	(*AdminLoginResponse)(nil),    // 1: http.app.AdminLoginResponse
	(*AdminAddTokenRequest)(nil),  // 2: http.app.AdminAddTokenRequest
	(*AdminAddTokenResponse)(nil), // 3: http.app.AdminAddTokenResponse
}
var file_http_app_admin_proto_depIdxs = []int32{
	0, // 0: http.app.AdminServer.Login:input_type -> http.app.AdminLoginRequest
	2, // 1: http.app.AdminServer.AddToken:input_type -> http.app.AdminAddTokenRequest
	1, // 2: http.app.AdminServer.Login:output_type -> http.app.AdminLoginResponse
	3, // 3: http.app.AdminServer.AddToken:output_type -> http.app.AdminAddTokenResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_http_app_admin_proto_init() }
func file_http_app_admin_proto_init() {
	if File_http_app_admin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_http_app_admin_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*AdminLoginRequest); i {
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
		file_http_app_admin_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AdminLoginResponse); i {
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
		file_http_app_admin_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*AdminAddTokenRequest); i {
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
		file_http_app_admin_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*AdminAddTokenResponse); i {
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
			RawDescriptor: file_http_app_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_http_app_admin_proto_goTypes,
		DependencyIndexes: file_http_app_admin_proto_depIdxs,
		MessageInfos:      file_http_app_admin_proto_msgTypes,
	}.Build()
	File_http_app_admin_proto = out.File
	file_http_app_admin_proto_rawDesc = nil
	file_http_app_admin_proto_goTypes = nil
	file_http_app_admin_proto_depIdxs = nil
}
