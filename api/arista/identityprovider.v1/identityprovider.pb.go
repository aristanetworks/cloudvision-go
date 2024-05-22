// Copyright (c) 2023 Arista Networks, Inc.  All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: arista/identityprovider.v1/identityprovider.proto

package identityprovider

import (
	fmp "github.com/aristanetworks/cloudvision-go/api/fmp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ProtocolBinding indicates SAML protocol binding to be used.
type ProtocolBinding int32

const (
	// PROTOCOL_BINDING_UNSPECIFIED indicates that a protocol binding is unspecified.
	ProtocolBinding_PROTOCOL_BINDING_UNSPECIFIED ProtocolBinding = 0
	// PROTOCOL_BINDING_HTTP_POST indicates HTTP-POST SAML protocol binding.
	ProtocolBinding_PROTOCOL_BINDING_HTTP_POST ProtocolBinding = 1
	// PROTOCOL_BINDING_HTTP_REDIRECT indicates HTTP-Redirect SAML protocol binding.
	ProtocolBinding_PROTOCOL_BINDING_HTTP_REDIRECT ProtocolBinding = 2
)

// Enum value maps for ProtocolBinding.
var (
	ProtocolBinding_name = map[int32]string{
		0: "PROTOCOL_BINDING_UNSPECIFIED",
		1: "PROTOCOL_BINDING_HTTP_POST",
		2: "PROTOCOL_BINDING_HTTP_REDIRECT",
	}
	ProtocolBinding_value = map[string]int32{
		"PROTOCOL_BINDING_UNSPECIFIED":   0,
		"PROTOCOL_BINDING_HTTP_POST":     1,
		"PROTOCOL_BINDING_HTTP_REDIRECT": 2,
	}
)

func (x ProtocolBinding) Enum() *ProtocolBinding {
	p := new(ProtocolBinding)
	*p = x
	return p
}

func (x ProtocolBinding) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProtocolBinding) Descriptor() protoreflect.EnumDescriptor {
	return file_arista_identityprovider_v1_identityprovider_proto_enumTypes[0].Descriptor()
}

func (ProtocolBinding) Type() protoreflect.EnumType {
	return &file_arista_identityprovider_v1_identityprovider_proto_enumTypes[0]
}

func (x ProtocolBinding) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProtocolBinding.Descriptor instead.
func (ProtocolBinding) EnumDescriptor() ([]byte, []int) {
	return file_arista_identityprovider_v1_identityprovider_proto_rawDescGZIP(), []int{0}
}

// OAuthKey contains OAuth provider ID.
type OAuthKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// provider_id is the ID of the OAuth provider.
	ProviderId *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=provider_id,json=providerId,proto3" json:"provider_id,omitempty"`
}

func (x *OAuthKey) Reset() {
	*x = OAuthKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_identityprovider_v1_identityprovider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OAuthKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OAuthKey) ProtoMessage() {}

func (x *OAuthKey) ProtoReflect() protoreflect.Message {
	mi := &file_arista_identityprovider_v1_identityprovider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OAuthKey.ProtoReflect.Descriptor instead.
func (*OAuthKey) Descriptor() ([]byte, []int) {
	return file_arista_identityprovider_v1_identityprovider_proto_rawDescGZIP(), []int{0}
}

func (x *OAuthKey) GetProviderId() *wrapperspb.StringValue {
	if x != nil {
		return x.ProviderId
	}
	return nil
}

// OAuthConfig holds the configuration for an OAuth provider.
type OAuthConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// key is the ID of the OAuth provider.
	Key *OAuthKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// endpoint is the URL that identifies an OAuth authorization server.
	// This endpoint is used to interact with the provider. It must be a
	// URI [RFC3986] with a scheme component that must be https, a host component,
	// and optionally, port and path components, but no query or fragment components.
	Endpoint *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// client_id is the ID that the OAuth authorization server issues to the
	// registered client.
	ClientId *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	// client_secret is the secret that the OAuth authorization server issues
	// to the registered client.
	ClientSecret *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
	// algorithms is the set of signing algorithms. This is an optional field.
	// If specified, only this set of algorithms may be used to sign the JWT.
	// Otherwise, this defaults to the set of algorithms that the provider supports.
	Algorithms *fmp.RepeatedString `protobuf:"bytes,5,opt,name=algorithms,proto3" json:"algorithms,omitempty"`
	// link_to_shared_provider indicates whether or not use the provider as a shared
	// provider. This is an optional field and set to false by default.
	LinkToSharedProvider *wrapperspb.BoolValue `protobuf:"bytes,6,opt,name=link_to_shared_provider,json=linkToSharedProvider,proto3" json:"link_to_shared_provider,omitempty"`
	// jwks_uri is where signing keys are downloaded. This is an optional field.
	// Only needed if the default construction from endpoint would be incorrect.
	JwksUri *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=jwks_uri,json=jwksUri,proto3" json:"jwks_uri,omitempty"`
	// permitted_email_domains are domains of emails that users are allowed to use.
	// This is an optional field. If not set, all domains are accepted by default.
	PermittedEmailDomains *fmp.RepeatedString `protobuf:"bytes,8,opt,name=permitted_email_domains,json=permittedEmailDomains,proto3" json:"permitted_email_domains,omitempty"`
	// roles_scope_name is the name for a scope tied to a claim that holds
	// CloudVision roles in ID Token. CloudVision uses scope values to specify
	// what access privileges are being requested for id token. CloudVision
	// appends this value to `scope` query parameter in the authorization request URL.
	// This is an optional field. If not set, CloudVision determines that
	// mapping roles from the provider is disabled. If it's set, roles_claim_name
	// also needs to be set.
	RolesScopeName *wrapperspb.StringValue `protobuf:"bytes,9,opt,name=roles_scope_name,json=rolesScopeName,proto3" json:"roles_scope_name,omitempty"`
	// bearer_token_introspection_endpoint is the provider instrospection endpoint used
	// in Bearer Token based login support for CloudVision. This is an optional field.
	// If specified, this endpoint will be used to verify bearer tokens generated via
	// the provider to log in automated user accounts.
	BearerTokenIntrospectionEndpoint *wrapperspb.StringValue `protobuf:"bytes,10,opt,name=bearer_token_introspection_endpoint,json=bearerTokenIntrospectionEndpoint,proto3" json:"bearer_token_introspection_endpoint,omitempty"`
	// roles_claim_name is the name for a claim that holds CloudVision roles in ID Token.
	// CloudVision uses this value to look up roles in the ID Token.
	// This is an optional field. If not set, CloudVision determines that
	// mapping roles from the provider is disabled. If it's set, roles_scope_name
	// also needs to be set.
	RolesClaimName *wrapperspb.StringValue `protobuf:"bytes,11,opt,name=roles_claim_name,json=rolesClaimName,proto3" json:"roles_claim_name,omitempty"`
}

func (x *OAuthConfig) Reset() {
	*x = OAuthConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_identityprovider_v1_identityprovider_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OAuthConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OAuthConfig) ProtoMessage() {}

func (x *OAuthConfig) ProtoReflect() protoreflect.Message {
	mi := &file_arista_identityprovider_v1_identityprovider_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OAuthConfig.ProtoReflect.Descriptor instead.
func (*OAuthConfig) Descriptor() ([]byte, []int) {
	return file_arista_identityprovider_v1_identityprovider_proto_rawDescGZIP(), []int{1}
}

func (x *OAuthConfig) GetKey() *OAuthKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *OAuthConfig) GetEndpoint() *wrapperspb.StringValue {
	if x != nil {
		return x.Endpoint
	}
	return nil
}

func (x *OAuthConfig) GetClientId() *wrapperspb.StringValue {
	if x != nil {
		return x.ClientId
	}
	return nil
}

func (x *OAuthConfig) GetClientSecret() *wrapperspb.StringValue {
	if x != nil {
		return x.ClientSecret
	}
	return nil
}

func (x *OAuthConfig) GetAlgorithms() *fmp.RepeatedString {
	if x != nil {
		return x.Algorithms
	}
	return nil
}

func (x *OAuthConfig) GetLinkToSharedProvider() *wrapperspb.BoolValue {
	if x != nil {
		return x.LinkToSharedProvider
	}
	return nil
}

func (x *OAuthConfig) GetJwksUri() *wrapperspb.StringValue {
	if x != nil {
		return x.JwksUri
	}
	return nil
}

func (x *OAuthConfig) GetPermittedEmailDomains() *fmp.RepeatedString {
	if x != nil {
		return x.PermittedEmailDomains
	}
	return nil
}

func (x *OAuthConfig) GetRolesScopeName() *wrapperspb.StringValue {
	if x != nil {
		return x.RolesScopeName
	}
	return nil
}

func (x *OAuthConfig) GetBearerTokenIntrospectionEndpoint() *wrapperspb.StringValue {
	if x != nil {
		return x.BearerTokenIntrospectionEndpoint
	}
	return nil
}

func (x *OAuthConfig) GetRolesClaimName() *wrapperspb.StringValue {
	if x != nil {
		return x.RolesClaimName
	}
	return nil
}

// SAMLKey contains SAML Provider ID.
type SAMLKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// provider_id is the ID of the SAML provider.
	ProviderId *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=provider_id,json=providerId,proto3" json:"provider_id,omitempty"`
}

func (x *SAMLKey) Reset() {
	*x = SAMLKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_identityprovider_v1_identityprovider_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SAMLKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SAMLKey) ProtoMessage() {}

func (x *SAMLKey) ProtoReflect() protoreflect.Message {
	mi := &file_arista_identityprovider_v1_identityprovider_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SAMLKey.ProtoReflect.Descriptor instead.
func (*SAMLKey) Descriptor() ([]byte, []int) {
	return file_arista_identityprovider_v1_identityprovider_proto_rawDescGZIP(), []int{2}
}

func (x *SAMLKey) GetProviderId() *wrapperspb.StringValue {
	if x != nil {
		return x.ProviderId
	}
	return nil
}

// SAMLConfig holds the configuration for a SAML provider.
type SAMLConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// key is the ID of the SAML provider.
	Key *SAMLKey `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// idp_issuer identifies the SAML provider. There is no restriction on its format
	// other than a string to carry the issuer's name.
	IdpIssuer *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=idp_issuer,json=idpIssuer,proto3" json:"idp_issuer,omitempty"`
	// idp_metadata_url is the URL that CloudVision uses to fetch the
	// SAML provider metadata.
	IdpMetadataUrl *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=idp_metadata_url,json=idpMetadataUrl,proto3" json:"idp_metadata_url,omitempty"`
	// authreq_binding specifies the ProtocolBinding used to send SAML authentication
	// request to the SAML provider.
	AuthreqBinding ProtocolBinding `protobuf:"varint,4,opt,name=authreq_binding,json=authreqBinding,proto3,enum=arista.identityprovider.v1.ProtocolBinding" json:"authreq_binding,omitempty"`
	// email_attrname specifies the Attribute name for email ID in Assertion of SAMLResponse
	// from the SAML provider.
	EmailAttrname *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=email_attrname,json=emailAttrname,proto3" json:"email_attrname,omitempty"`
	// link_to_shared_provider indicates whether or not use the provider as a shared
	// provider. This is an optional field and set to false by default.
	LinkToSharedProvider *wrapperspb.BoolValue `protobuf:"bytes,6,opt,name=link_to_shared_provider,json=linkToSharedProvider,proto3" json:"link_to_shared_provider,omitempty"`
	// permitted_email_domains are domains of emails that users are allowed to use.
	// This is an optional field. If not set, all domains are accepted by default.
	PermittedEmailDomains *fmp.RepeatedString `protobuf:"bytes,7,opt,name=permitted_email_domains,json=permittedEmailDomains,proto3" json:"permitted_email_domains,omitempty"`
	// force_saml_authn indicates wether or not enable force authentication in SAML login.
	// This is an optional field. If not set, it defaults to false.
	ForceSamlAuthn *wrapperspb.BoolValue `protobuf:"bytes,8,opt,name=force_saml_authn,json=forceSamlAuthn,proto3" json:"force_saml_authn,omitempty"`
	// roles_attrname specifies the Attribute name for CloudVision roles in the Assertion
	// of SAMLResponse. This is an optional field. If not set, CloudVision determines that
	// mapping roles from the provider is disabled.
	RolesAttrname *wrapperspb.StringValue `protobuf:"bytes,9,opt,name=roles_attrname,json=rolesAttrname,proto3" json:"roles_attrname,omitempty"`
	// org_attrname specifies the Attribute name for CloudVision organization/tenant in
	// the Assertion of SAMLResponse. This is an optional field. CloudVision supports use
	// of certain shared SAML Identity Providers for authenticating users across multiple
	// CloudVision organizations/tenants. In case a given organization uses a shared provider,
	// then, CloudVision needs this attribute to determine if the organization that
	// the shared SAML Identity Provider is sending the assertion for is the same as the
	// one the user requested to be logged into. For an existing user on CloudVision,
	// the user's email is used to determine which organization the user belongs to do
	// the same verification but in case a dynamic user creation is needed and the given
	// user doesn't exist on CloudVision currently then the matching organization attribute
	// from the shared Identity Privder becomes necessary. Dynamic user creation is
	// disabled for a given organization using shared Identity Provider if this attribute
	// is not specified.
	OrgAttrname *wrapperspb.StringValue `protobuf:"bytes,10,opt,name=org_attrname,json=orgAttrname,proto3" json:"org_attrname,omitempty"`
	// username_attrname specifies Attribute name for CloudVision users' username in the
	// Assertion of SAMLResponse. This is an optional field as long as mapping roles from
	// provider is not enabled. Once enabled, this field becomes mandatory.
	UsernameAttrname *wrapperspb.StringValue `protobuf:"bytes,11,opt,name=username_attrname,json=usernameAttrname,proto3" json:"username_attrname,omitempty"`
}

func (x *SAMLConfig) Reset() {
	*x = SAMLConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arista_identityprovider_v1_identityprovider_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SAMLConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SAMLConfig) ProtoMessage() {}

func (x *SAMLConfig) ProtoReflect() protoreflect.Message {
	mi := &file_arista_identityprovider_v1_identityprovider_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SAMLConfig.ProtoReflect.Descriptor instead.
func (*SAMLConfig) Descriptor() ([]byte, []int) {
	return file_arista_identityprovider_v1_identityprovider_proto_rawDescGZIP(), []int{3}
}

func (x *SAMLConfig) GetKey() *SAMLKey {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *SAMLConfig) GetIdpIssuer() *wrapperspb.StringValue {
	if x != nil {
		return x.IdpIssuer
	}
	return nil
}

func (x *SAMLConfig) GetIdpMetadataUrl() *wrapperspb.StringValue {
	if x != nil {
		return x.IdpMetadataUrl
	}
	return nil
}

func (x *SAMLConfig) GetAuthreqBinding() ProtocolBinding {
	if x != nil {
		return x.AuthreqBinding
	}
	return ProtocolBinding_PROTOCOL_BINDING_UNSPECIFIED
}

func (x *SAMLConfig) GetEmailAttrname() *wrapperspb.StringValue {
	if x != nil {
		return x.EmailAttrname
	}
	return nil
}

func (x *SAMLConfig) GetLinkToSharedProvider() *wrapperspb.BoolValue {
	if x != nil {
		return x.LinkToSharedProvider
	}
	return nil
}

func (x *SAMLConfig) GetPermittedEmailDomains() *fmp.RepeatedString {
	if x != nil {
		return x.PermittedEmailDomains
	}
	return nil
}

func (x *SAMLConfig) GetForceSamlAuthn() *wrapperspb.BoolValue {
	if x != nil {
		return x.ForceSamlAuthn
	}
	return nil
}

func (x *SAMLConfig) GetRolesAttrname() *wrapperspb.StringValue {
	if x != nil {
		return x.RolesAttrname
	}
	return nil
}

func (x *SAMLConfig) GetOrgAttrname() *wrapperspb.StringValue {
	if x != nil {
		return x.OrgAttrname
	}
	return nil
}

func (x *SAMLConfig) GetUsernameAttrname() *wrapperspb.StringValue {
	if x != nil {
		return x.UsernameAttrname
	}
	return nil
}

var File_arista_identityprovider_v1_identityprovider_proto protoreflect.FileDescriptor

var file_arista_identityprovider_v1_identityprovider_proto_rawDesc = []byte{
	0x0a, 0x31, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2f, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x14, 0x66, 0x6d, 0x70, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x66, 0x6d, 0x70, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4f, 0x0a, 0x08, 0x4f, 0x41, 0x75,
	0x74, 0x68, 0x4b, 0x65, 0x79, 0x12, 0x3d, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x3a, 0x04, 0x80, 0x8e, 0x19, 0x01, 0x22, 0x90, 0x06, 0x0a, 0x0b, 0x4f,
	0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x36, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x38, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x09,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x41, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0c, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x33, 0x0a, 0x0a, 0x61, 0x6c,
	0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x66, 0x6d, 0x70, 0x2e, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x52, 0x0a, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x73, 0x12,
	0x51, 0x0a, 0x17, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x74, 0x6f, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x14, 0x6c, 0x69,
	0x6e, 0x6b, 0x54, 0x6f, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x12, 0x37, 0x0a, 0x08, 0x6a, 0x77, 0x6b, 0x73, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x07, 0x6a, 0x77, 0x6b, 0x73, 0x55, 0x72, 0x69, 0x12, 0x4b, 0x0a, 0x17, 0x70,
	0x65, 0x72, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x66,
	0x6d, 0x70, 0x2e, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x52, 0x15, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x12, 0x46, 0x0a, 0x10, 0x72, 0x6f, 0x6c, 0x65,
	0x73, 0x5f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x0e, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x6b, 0x0a, 0x23, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x5f, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65,
	0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x20, 0x62, 0x65, 0x61,
	0x72, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x73, 0x70, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x46, 0x0a,
	0x10, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x5f, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0e, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x43, 0x6c, 0x61, 0x69,
	0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02, 0x72, 0x77, 0x22, 0x4e, 0x0a,
	0x07, 0x53, 0x41, 0x4d, 0x4c, 0x4b, 0x65, 0x79, 0x12, 0x3d, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x49, 0x64, 0x3a, 0x04, 0x80, 0x8e, 0x19, 0x01, 0x22, 0xa2, 0x06,
	0x0a, 0x0a, 0x53, 0x41, 0x4d, 0x4c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x35, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x41, 0x4d, 0x4c, 0x4b, 0x65, 0x79, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x3b, 0x0a, 0x0a, 0x69, 0x64, 0x70, 0x5f, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x69, 0x64, 0x70, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72,
	0x12, 0x46, 0x0a, 0x10, 0x69, 0x64, 0x70, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0e, 0x69, 0x64, 0x70, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x55, 0x72, 0x6c, 0x12, 0x54, 0x0a, 0x0f, 0x61, 0x75, 0x74, 0x68,
	0x72, 0x65, 0x71, 0x5f, 0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x2b, 0x2e, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x0e,
	0x61, 0x75, 0x74, 0x68, 0x72, 0x65, 0x71, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x43,
	0x0a, 0x0e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x0d, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x41, 0x74, 0x74, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x51, 0x0a, 0x17, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x74, 0x6f, 0x5f, 0x73,
	0x68, 0x61, 0x72, 0x65, 0x64, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x14, 0x6c, 0x69, 0x6e, 0x6b, 0x54, 0x6f, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x50, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x4b, 0x0a, 0x17, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x74,
	0x74, 0x65, 0x64, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x66, 0x6d, 0x70, 0x2e, 0x52, 0x65,
	0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x15, 0x70, 0x65,
	0x72, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x44, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x73, 0x12, 0x44, 0x0a, 0x10, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x5f, 0x73, 0x61, 0x6d,
	0x6c, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0e, 0x66, 0x6f, 0x72, 0x63, 0x65,
	0x53, 0x61, 0x6d, 0x6c, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x12, 0x43, 0x0a, 0x0e, 0x72, 0x6f, 0x6c,
	0x65, 0x73, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x0d, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x41, 0x74, 0x74, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3f,
	0x0a, 0x0c, 0x6f, 0x72, 0x67, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x0b, 0x6f, 0x72, 0x67, 0x41, 0x74, 0x74, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x49, 0x0a, 0x11, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x61, 0x74, 0x74, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x10, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x41, 0x74, 0x74, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x3a, 0x06, 0xfa, 0x8d, 0x19, 0x02,
	0x72, 0x77, 0x2a, 0x77, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x42, 0x69,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x20, 0x0a, 0x1c, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f,
	0x4c, 0x5f, 0x42, 0x49, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4f, 0x4c, 0x5f, 0x42, 0x49, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x5f, 0x48, 0x54, 0x54, 0x50,
	0x5f, 0x50, 0x4f, 0x53, 0x54, 0x10, 0x01, 0x12, 0x22, 0x0a, 0x1e, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4f, 0x4c, 0x5f, 0x42, 0x49, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x5f, 0x48, 0x54, 0x54, 0x50,
	0x5f, 0x52, 0x45, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x10, 0x02, 0x42, 0x5a, 0x5a, 0x58, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x76, 0x69,
	0x73, 0x69, 0x6f, 0x6e, 0x2d, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x3b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arista_identityprovider_v1_identityprovider_proto_rawDescOnce sync.Once
	file_arista_identityprovider_v1_identityprovider_proto_rawDescData = file_arista_identityprovider_v1_identityprovider_proto_rawDesc
)

func file_arista_identityprovider_v1_identityprovider_proto_rawDescGZIP() []byte {
	file_arista_identityprovider_v1_identityprovider_proto_rawDescOnce.Do(func() {
		file_arista_identityprovider_v1_identityprovider_proto_rawDescData = protoimpl.X.CompressGZIP(file_arista_identityprovider_v1_identityprovider_proto_rawDescData)
	})
	return file_arista_identityprovider_v1_identityprovider_proto_rawDescData
}

var file_arista_identityprovider_v1_identityprovider_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_arista_identityprovider_v1_identityprovider_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_arista_identityprovider_v1_identityprovider_proto_goTypes = []interface{}{
	(ProtocolBinding)(0),           // 0: arista.identityprovider.v1.ProtocolBinding
	(*OAuthKey)(nil),               // 1: arista.identityprovider.v1.OAuthKey
	(*OAuthConfig)(nil),            // 2: arista.identityprovider.v1.OAuthConfig
	(*SAMLKey)(nil),                // 3: arista.identityprovider.v1.SAMLKey
	(*SAMLConfig)(nil),             // 4: arista.identityprovider.v1.SAMLConfig
	(*wrapperspb.StringValue)(nil), // 5: google.protobuf.StringValue
	(*fmp.RepeatedString)(nil),     // 6: fmp.RepeatedString
	(*wrapperspb.BoolValue)(nil),   // 7: google.protobuf.BoolValue
}
var file_arista_identityprovider_v1_identityprovider_proto_depIdxs = []int32{
	5,  // 0: arista.identityprovider.v1.OAuthKey.provider_id:type_name -> google.protobuf.StringValue
	1,  // 1: arista.identityprovider.v1.OAuthConfig.key:type_name -> arista.identityprovider.v1.OAuthKey
	5,  // 2: arista.identityprovider.v1.OAuthConfig.endpoint:type_name -> google.protobuf.StringValue
	5,  // 3: arista.identityprovider.v1.OAuthConfig.client_id:type_name -> google.protobuf.StringValue
	5,  // 4: arista.identityprovider.v1.OAuthConfig.client_secret:type_name -> google.protobuf.StringValue
	6,  // 5: arista.identityprovider.v1.OAuthConfig.algorithms:type_name -> fmp.RepeatedString
	7,  // 6: arista.identityprovider.v1.OAuthConfig.link_to_shared_provider:type_name -> google.protobuf.BoolValue
	5,  // 7: arista.identityprovider.v1.OAuthConfig.jwks_uri:type_name -> google.protobuf.StringValue
	6,  // 8: arista.identityprovider.v1.OAuthConfig.permitted_email_domains:type_name -> fmp.RepeatedString
	5,  // 9: arista.identityprovider.v1.OAuthConfig.roles_scope_name:type_name -> google.protobuf.StringValue
	5,  // 10: arista.identityprovider.v1.OAuthConfig.bearer_token_introspection_endpoint:type_name -> google.protobuf.StringValue
	5,  // 11: arista.identityprovider.v1.OAuthConfig.roles_claim_name:type_name -> google.protobuf.StringValue
	5,  // 12: arista.identityprovider.v1.SAMLKey.provider_id:type_name -> google.protobuf.StringValue
	3,  // 13: arista.identityprovider.v1.SAMLConfig.key:type_name -> arista.identityprovider.v1.SAMLKey
	5,  // 14: arista.identityprovider.v1.SAMLConfig.idp_issuer:type_name -> google.protobuf.StringValue
	5,  // 15: arista.identityprovider.v1.SAMLConfig.idp_metadata_url:type_name -> google.protobuf.StringValue
	0,  // 16: arista.identityprovider.v1.SAMLConfig.authreq_binding:type_name -> arista.identityprovider.v1.ProtocolBinding
	5,  // 17: arista.identityprovider.v1.SAMLConfig.email_attrname:type_name -> google.protobuf.StringValue
	7,  // 18: arista.identityprovider.v1.SAMLConfig.link_to_shared_provider:type_name -> google.protobuf.BoolValue
	6,  // 19: arista.identityprovider.v1.SAMLConfig.permitted_email_domains:type_name -> fmp.RepeatedString
	7,  // 20: arista.identityprovider.v1.SAMLConfig.force_saml_authn:type_name -> google.protobuf.BoolValue
	5,  // 21: arista.identityprovider.v1.SAMLConfig.roles_attrname:type_name -> google.protobuf.StringValue
	5,  // 22: arista.identityprovider.v1.SAMLConfig.org_attrname:type_name -> google.protobuf.StringValue
	5,  // 23: arista.identityprovider.v1.SAMLConfig.username_attrname:type_name -> google.protobuf.StringValue
	24, // [24:24] is the sub-list for method output_type
	24, // [24:24] is the sub-list for method input_type
	24, // [24:24] is the sub-list for extension type_name
	24, // [24:24] is the sub-list for extension extendee
	0,  // [0:24] is the sub-list for field type_name
}

func init() { file_arista_identityprovider_v1_identityprovider_proto_init() }
func file_arista_identityprovider_v1_identityprovider_proto_init() {
	if File_arista_identityprovider_v1_identityprovider_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arista_identityprovider_v1_identityprovider_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OAuthKey); i {
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
		file_arista_identityprovider_v1_identityprovider_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OAuthConfig); i {
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
		file_arista_identityprovider_v1_identityprovider_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SAMLKey); i {
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
		file_arista_identityprovider_v1_identityprovider_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SAMLConfig); i {
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
			RawDescriptor: file_arista_identityprovider_v1_identityprovider_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arista_identityprovider_v1_identityprovider_proto_goTypes,
		DependencyIndexes: file_arista_identityprovider_v1_identityprovider_proto_depIdxs,
		EnumInfos:         file_arista_identityprovider_v1_identityprovider_proto_enumTypes,
		MessageInfos:      file_arista_identityprovider_v1_identityprovider_proto_msgTypes,
	}.Build()
	File_arista_identityprovider_v1_identityprovider_proto = out.File
	file_arista_identityprovider_v1_identityprovider_proto_rawDesc = nil
	file_arista_identityprovider_v1_identityprovider_proto_goTypes = nil
	file_arista_identityprovider_v1_identityprovider_proto_depIdxs = nil
}
