// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: oidc/registration/v1/client.proto

package registrationv1

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Client struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RedirectUris    []string              `protobuf:"bytes,1,rep,name=redirect_uris,json=redirectUris,proto3" json:"redirect_uris,omitempty"`
	ResponseTypes   []string              `protobuf:"bytes,2,rep,name=response_types,json=responseTypes,proto3" json:"response_types,omitempty"`
	GrantTypes      []string              `protobuf:"bytes,3,rep,name=grant_types,json=grantTypes,proto3" json:"grant_types,omitempty"`
	ApplicationType *wrappers.StringValue `protobuf:"bytes,4,opt,name=application_type,json=applicationType,proto3" json:"application_type,omitempty"`
	Contacts        []string              `protobuf:"bytes,5,rep,name=contacts,proto3" json:"contacts,omitempty"`
	ClientName      *wrappers.StringValue `protobuf:"bytes,6,opt,name=client_name,json=clientName,proto3" json:"client_name,omitempty"`
	LogoUri         *wrappers.StringValue `protobuf:"bytes,7,opt,name=logo_uri,json=logoUri,proto3" json:"logo_uri,omitempty"`
	ClientUri       *wrappers.StringValue `protobuf:"bytes,8,opt,name=client_uri,json=clientUri,proto3" json:"client_uri,omitempty"`
	PolicyUri       *wrappers.StringValue `protobuf:"bytes,9,opt,name=policy_uri,json=policyUri,proto3" json:"policy_uri,omitempty"`
	TosUri          *wrappers.StringValue `protobuf:"bytes,10,opt,name=tos_uri,json=tosUri,proto3" json:"tos_uri,omitempty"`
	JwksUri         *wrappers.StringValue `protobuf:"bytes,11,opt,name=jwks_uri,json=jwksUri,proto3" json:"jwks_uri,omitempty"`
	Jwks            *wrappers.BytesValue  `protobuf:"bytes,12,opt,name=jwks,proto3" json:"jwks,omitempty"`
}

func (x *Client) Reset() {
	*x = Client{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oidc_registration_v1_client_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Client) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Client) ProtoMessage() {}

func (x *Client) ProtoReflect() protoreflect.Message {
	mi := &file_oidc_registration_v1_client_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Client.ProtoReflect.Descriptor instead.
func (*Client) Descriptor() ([]byte, []int) {
	return file_oidc_registration_v1_client_proto_rawDescGZIP(), []int{0}
}

func (x *Client) GetRedirectUris() []string {
	if x != nil {
		return x.RedirectUris
	}
	return nil
}

func (x *Client) GetResponseTypes() []string {
	if x != nil {
		return x.ResponseTypes
	}
	return nil
}

func (x *Client) GetGrantTypes() []string {
	if x != nil {
		return x.GrantTypes
	}
	return nil
}

func (x *Client) GetApplicationType() *wrappers.StringValue {
	if x != nil {
		return x.ApplicationType
	}
	return nil
}

func (x *Client) GetContacts() []string {
	if x != nil {
		return x.Contacts
	}
	return nil
}

func (x *Client) GetClientName() *wrappers.StringValue {
	if x != nil {
		return x.ClientName
	}
	return nil
}

func (x *Client) GetLogoUri() *wrappers.StringValue {
	if x != nil {
		return x.LogoUri
	}
	return nil
}

func (x *Client) GetClientUri() *wrappers.StringValue {
	if x != nil {
		return x.ClientUri
	}
	return nil
}

func (x *Client) GetPolicyUri() *wrappers.StringValue {
	if x != nil {
		return x.PolicyUri
	}
	return nil
}

func (x *Client) GetTosUri() *wrappers.StringValue {
	if x != nil {
		return x.TosUri
	}
	return nil
}

func (x *Client) GetJwksUri() *wrappers.StringValue {
	if x != nil {
		return x.JwksUri
	}
	return nil
}

func (x *Client) GetJwks() *wrappers.BytesValue {
	if x != nil {
		return x.Jwks
	}
	return nil
}

var File_oidc_registration_v1_client_proto protoreflect.FileDescriptor

var file_oidc_registration_v1_client_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6f, 0x69, 0x64, 0x63, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x14, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xed, 0x04, 0x0a, 0x06, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x5f, 0x75, 0x72, 0x69, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x55, 0x72, 0x69, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73,
	0x12, 0x1f, 0x0a, 0x0b, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x12, 0x47, 0x0a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f,
	0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f,
	0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x12, 0x3d, 0x0a, 0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x6f, 0x5f, 0x75, 0x72,
	0x69, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x6f, 0x55, 0x72, 0x69, 0x12, 0x3b,
	0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x55, 0x72, 0x69, 0x12, 0x3b, 0x0a, 0x0a, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x55, 0x72, 0x69, 0x12, 0x35, 0x0a, 0x07, 0x74, 0x6f, 0x73, 0x5f,
	0x75, 0x72, 0x69, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x74, 0x6f, 0x73, 0x55, 0x72, 0x69, 0x12,
	0x37, 0x0a, 0x08, 0x6a, 0x77, 0x6b, 0x73, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x07, 0x6a, 0x77, 0x6b, 0x73, 0x55, 0x72, 0x69, 0x12, 0x2f, 0x0a, 0x04, 0x6a, 0x77, 0x6b, 0x73,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x04, 0x6a, 0x77, 0x6b, 0x73, 0x42, 0x25, 0x5a, 0x23, 0x6f, 0x69, 0x64,
	0x63, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76,
	0x31, 0x3b, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x76, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_oidc_registration_v1_client_proto_rawDescOnce sync.Once
	file_oidc_registration_v1_client_proto_rawDescData = file_oidc_registration_v1_client_proto_rawDesc
)

func file_oidc_registration_v1_client_proto_rawDescGZIP() []byte {
	file_oidc_registration_v1_client_proto_rawDescOnce.Do(func() {
		file_oidc_registration_v1_client_proto_rawDescData = protoimpl.X.CompressGZIP(file_oidc_registration_v1_client_proto_rawDescData)
	})
	return file_oidc_registration_v1_client_proto_rawDescData
}

var (
	file_oidc_registration_v1_client_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
	file_oidc_registration_v1_client_proto_goTypes  = []interface{}{
		(*Client)(nil),               // 0: oidc.registration.v1.Client
		(*wrappers.StringValue)(nil), // 1: google.protobuf.StringValue
		(*wrappers.BytesValue)(nil),  // 2: google.protobuf.BytesValue
	}
)

var file_oidc_registration_v1_client_proto_depIdxs = []int32{
	1, // 0: oidc.registration.v1.Client.application_type:type_name -> google.protobuf.StringValue
	1, // 1: oidc.registration.v1.Client.client_name:type_name -> google.protobuf.StringValue
	1, // 2: oidc.registration.v1.Client.logo_uri:type_name -> google.protobuf.StringValue
	1, // 3: oidc.registration.v1.Client.client_uri:type_name -> google.protobuf.StringValue
	1, // 4: oidc.registration.v1.Client.policy_uri:type_name -> google.protobuf.StringValue
	1, // 5: oidc.registration.v1.Client.tos_uri:type_name -> google.protobuf.StringValue
	1, // 6: oidc.registration.v1.Client.jwks_uri:type_name -> google.protobuf.StringValue
	2, // 7: oidc.registration.v1.Client.jwks:type_name -> google.protobuf.BytesValue
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_oidc_registration_v1_client_proto_init() }
func file_oidc_registration_v1_client_proto_init() {
	if File_oidc_registration_v1_client_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_oidc_registration_v1_client_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Client); i {
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
			RawDescriptor: file_oidc_registration_v1_client_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_oidc_registration_v1_client_proto_goTypes,
		DependencyIndexes: file_oidc_registration_v1_client_proto_depIdxs,
		MessageInfos:      file_oidc_registration_v1_client_proto_msgTypes,
	}.Build()
	File_oidc_registration_v1_client_proto = out.File
	file_oidc_registration_v1_client_proto_rawDesc = nil
	file_oidc_registration_v1_client_proto_goTypes = nil
	file_oidc_registration_v1_client_proto_depIdxs = nil
}
