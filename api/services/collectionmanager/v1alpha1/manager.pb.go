// Events RPC protocol version v1alpha1
//
// This file defines version v1alpha of the RPC protocol. To implement a plugin
// against this protocol, copy this definition into your own codebase and
// use protoc to generate stubs for your target language.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: api/services/collectionmanager/v1alpha1/manager.proto

package v1alpha1

import (
	_struct "github.com/golang/protobuf/ptypes/struct"
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

type Diagnostic_Severity int32

const (
	Diagnostic_INVALID Diagnostic_Severity = 0
	Diagnostic_ERROR   Diagnostic_Severity = 1
	Diagnostic_WARNING Diagnostic_Severity = 2
)

// Enum value maps for Diagnostic_Severity.
var (
	Diagnostic_Severity_name = map[int32]string{
		0: "INVALID",
		1: "ERROR",
		2: "WARNING",
	}
	Diagnostic_Severity_value = map[string]int32{
		"INVALID": 0,
		"ERROR":   1,
		"WARNING": 2,
	}
)

func (x Diagnostic_Severity) Enum() *Diagnostic_Severity {
	p := new(Diagnostic_Severity)
	*p = x
	return p
}

func (x Diagnostic_Severity) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Diagnostic_Severity) Descriptor() protoreflect.EnumDescriptor {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_enumTypes[0].Descriptor()
}

func (Diagnostic_Severity) Type() protoreflect.EnumType {
	return &file_api_services_collectionmanager_v1alpha1_manager_proto_enumTypes[0]
}

func (x Diagnostic_Severity) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Diagnostic_Severity.Descriptor instead.
func (Diagnostic_Severity) EnumDescriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{0, 0}
}

type Diagnostic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Severity Diagnostic_Severity `protobuf:"varint,1,opt,name=severity,proto3,enum=manager.Diagnostic_Severity" json:"severity,omitempty"`
	Summary  string              `protobuf:"bytes,2,opt,name=summary,proto3" json:"summary,omitempty"`
	Detail   string              `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
}

func (x *Diagnostic) Reset() {
	*x = Diagnostic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Diagnostic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Diagnostic) ProtoMessage() {}

func (x *Diagnostic) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Diagnostic.ProtoReflect.Descriptor instead.
func (*Diagnostic) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{0}
}

func (x *Diagnostic) GetSeverity() Diagnostic_Severity {
	if x != nil {
		return x.Severity
	}
	return Diagnostic_INVALID
}

func (x *Diagnostic) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *Diagnostic) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

type Retrieve struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Retrieve) Reset() {
	*x = Retrieve{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Retrieve) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Retrieve) ProtoMessage() {}

func (x *Retrieve) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Retrieve.ProtoReflect.Descriptor instead.
func (*Retrieve) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{1}
}

type Publish struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Publish) Reset() {
	*x = Publish{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Publish) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Publish) ProtoMessage() {}

func (x *Publish) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Publish.ProtoReflect.Descriptor instead.
func (*Publish) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{2}
}

// Collection contains configuration information for a collection.
type Collection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SchemaAddress     string   `protobuf:"bytes,1,opt,name=schemaAddress,proto3" json:"schemaAddress,omitempty"`
	LinkedCollections []string `protobuf:"bytes,2,rep,name=linkedCollections,proto3" json:"linkedCollections,omitempty"`
	Files             []*File  `protobuf:"bytes,3,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *Collection) Reset() {
	*x = Collection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Collection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Collection) ProtoMessage() {}

func (x *Collection) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Collection.ProtoReflect.Descriptor instead.
func (*Collection) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{3}
}

func (x *Collection) GetSchemaAddress() string {
	if x != nil {
		return x.SchemaAddress
	}
	return ""
}

func (x *Collection) GetLinkedCollections() []string {
	if x != nil {
		return x.LinkedCollections
	}
	return nil
}

func (x *Collection) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

// File contains a regular expression for file name matching and associated
// attributes to apply the the descriptor for matching file.
type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File       string          `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	Attributes *_struct.Struct `protobuf:"bytes,2,opt,name=attributes,proto3" json:"attributes,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{4}
}

func (x *File) GetFile() string {
	if x != nil {
		return x.File
	}
	return ""
}

func (x *File) GetAttributes() *_struct.Struct {
	if x != nil {
		return x.Attributes
	}
	return nil
}

// AuthConfig contains authorization information for connecting to a registry.
type AuthConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username      string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password      string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Auth          string `protobuf:"bytes,3,opt,name=auth,proto3" json:"auth,omitempty"`
	ServerAddress string `protobuf:"bytes,4,opt,name=server_address,json=serverAddress,proto3" json:"server_address,omitempty"`
	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	AccessToken string `protobuf:"bytes,5,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	// RegistryToken is a bearer token to be sent to a registry
	RefreshToken string `protobuf:"bytes,6,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
}

func (x *AuthConfig) Reset() {
	*x = AuthConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthConfig) ProtoMessage() {}

func (x *AuthConfig) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthConfig.ProtoReflect.Descriptor instead.
func (*AuthConfig) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{5}
}

func (x *AuthConfig) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *AuthConfig) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *AuthConfig) GetAuth() string {
	if x != nil {
		return x.Auth
	}
	return ""
}

func (x *AuthConfig) GetServerAddress() string {
	if x != nil {
		return x.ServerAddress
	}
	return ""
}

func (x *AuthConfig) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *AuthConfig) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

type Retrieve_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source      string          `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Destination string          `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	Filter      *_struct.Struct `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
	Auth        *AuthConfig     `protobuf:"bytes,4,opt,name=auth,proto3" json:"auth,omitempty"`
}

func (x *Retrieve_Request) Reset() {
	*x = Retrieve_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Retrieve_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Retrieve_Request) ProtoMessage() {}

func (x *Retrieve_Request) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Retrieve_Request.ProtoReflect.Descriptor instead.
func (*Retrieve_Request) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Retrieve_Request) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Retrieve_Request) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

func (x *Retrieve_Request) GetFilter() *_struct.Struct {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *Retrieve_Request) GetAuth() *AuthConfig {
	if x != nil {
		return x.Auth
	}
	return nil
}

type Retrieve_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Digests     []string      `protobuf:"bytes,1,rep,name=digests,proto3" json:"digests,omitempty"`
	Diagnostics []*Diagnostic `protobuf:"bytes,2,rep,name=diagnostics,proto3" json:"diagnostics,omitempty"`
}

func (x *Retrieve_Response) Reset() {
	*x = Retrieve_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Retrieve_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Retrieve_Response) ProtoMessage() {}

func (x *Retrieve_Response) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Retrieve_Response.ProtoReflect.Descriptor instead.
func (*Retrieve_Response) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Retrieve_Response) GetDigests() []string {
	if x != nil {
		return x.Digests
	}
	return nil
}

func (x *Retrieve_Response) GetDiagnostics() []*Diagnostic {
	if x != nil {
		return x.Diagnostics
	}
	return nil
}

type Publish_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source      string      `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Destination string      `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	Collection  *Collection `protobuf:"bytes,3,opt,name=collection,proto3" json:"collection,omitempty"`
	Auth        *AuthConfig `protobuf:"bytes,4,opt,name=auth,proto3" json:"auth,omitempty"`
}

func (x *Publish_Request) Reset() {
	*x = Publish_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Publish_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Publish_Request) ProtoMessage() {}

func (x *Publish_Request) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Publish_Request.ProtoReflect.Descriptor instead.
func (*Publish_Request) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Publish_Request) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Publish_Request) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

func (x *Publish_Request) GetCollection() *Collection {
	if x != nil {
		return x.Collection
	}
	return nil
}

func (x *Publish_Request) GetAuth() *AuthConfig {
	if x != nil {
		return x.Auth
	}
	return nil
}

type Publish_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Digest      string        `protobuf:"bytes,1,opt,name=digest,proto3" json:"digest,omitempty"`
	Diagnostics []*Diagnostic `protobuf:"bytes,2,rep,name=diagnostics,proto3" json:"diagnostics,omitempty"`
}

func (x *Publish_Response) Reset() {
	*x = Publish_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Publish_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Publish_Response) ProtoMessage() {}

func (x *Publish_Response) ProtoReflect() protoreflect.Message {
	mi := &file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Publish_Response.ProtoReflect.Descriptor instead.
func (*Publish_Response) Descriptor() ([]byte, []int) {
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP(), []int{2, 1}
}

func (x *Publish_Response) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *Publish_Response) GetDiagnostics() []*Diagnostic {
	if x != nil {
		return x.Diagnostics
	}
	return nil
}

var File_api_services_collectionmanager_v1alpha1_manager_proto protoreflect.FileDescriptor

var file_api_services_collectionmanager_v1alpha1_manager_proto_rawDesc = []byte{
	0x0a, 0x35, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa9,
	0x01, 0x0a, 0x0a, 0x44, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x12, 0x38, 0x0a,
	0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1c, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x61, 0x67, 0x6e, 0x6f,
	0x73, 0x74, 0x69, 0x63, 0x2e, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x52, 0x08, 0x73,
	0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x22, 0x2f, 0x0a, 0x08, 0x53, 0x65, 0x76,
	0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44,
	0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x01, 0x12, 0x0b, 0x0a,
	0x07, 0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x22, 0x87, 0x02, 0x0a, 0x08, 0x52,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x1a, 0x9d, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a,
	0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x27,
	0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x04, 0x61, 0x75, 0x74, 0x68, 0x1a, 0x5b, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x73, 0x12, 0x35, 0x0a,
	0x0b, 0x64, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x61,
	0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x52, 0x0b, 0x64, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73,
	0x74, 0x69, 0x63, 0x73, 0x22, 0x88, 0x02, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x1a, 0xa1, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x33, 0x0a, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x04, 0x61,
	0x75, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x04,
	0x61, 0x75, 0x74, 0x68, 0x1a, 0x59, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x0b, 0x64, 0x69, 0x61, 0x67,
	0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x74,
	0x69, 0x63, 0x52, 0x0b, 0x64, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x73, 0x22,
	0x85, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24,
	0x0a, 0x0d, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x11, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x23, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x53, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x22, 0xc7, 0x01, 0x0a,
	0x0a, 0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x75, 0x74, 0x68, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x21,
	0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73,
	0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xa8, 0x01, 0x0a, 0x11, 0x43, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x47, 0x0a, 0x0e,
	0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x18,
	0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0f, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76,
	0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x52, 0x65,
	0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x50, 0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x75, 0x6f, 0x72, 0x2d, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x75, 0x6f,
	0x72, 0x2d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescOnce sync.Once
	file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescData = file_api_services_collectionmanager_v1alpha1_manager_proto_rawDesc
)

func file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescGZIP() []byte {
	file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescOnce.Do(func() {
		file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescData)
	})
	return file_api_services_collectionmanager_v1alpha1_manager_proto_rawDescData
}

var file_api_services_collectionmanager_v1alpha1_manager_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_api_services_collectionmanager_v1alpha1_manager_proto_goTypes = []interface{}{
	(Diagnostic_Severity)(0),  // 0: manager.Diagnostic.Severity
	(*Diagnostic)(nil),        // 1: manager.Diagnostic
	(*Retrieve)(nil),          // 2: manager.Retrieve
	(*Publish)(nil),           // 3: manager.Publish
	(*Collection)(nil),        // 4: manager.Collection
	(*File)(nil),              // 5: manager.File
	(*AuthConfig)(nil),        // 6: manager.AuthConfig
	(*Retrieve_Request)(nil),  // 7: manager.Retrieve.Request
	(*Retrieve_Response)(nil), // 8: manager.Retrieve.Response
	(*Publish_Request)(nil),   // 9: manager.Publish.Request
	(*Publish_Response)(nil),  // 10: manager.Publish.Response
	(*_struct.Struct)(nil),    // 11: google.protobuf.Struct
}
var file_api_services_collectionmanager_v1alpha1_manager_proto_depIdxs = []int32{
	0,  // 0: manager.Diagnostic.severity:type_name -> manager.Diagnostic.Severity
	5,  // 1: manager.Collection.files:type_name -> manager.File
	11, // 2: manager.File.attributes:type_name -> google.protobuf.Struct
	11, // 3: manager.Retrieve.Request.filter:type_name -> google.protobuf.Struct
	6,  // 4: manager.Retrieve.Request.auth:type_name -> manager.AuthConfig
	1,  // 5: manager.Retrieve.Response.diagnostics:type_name -> manager.Diagnostic
	4,  // 6: manager.Publish.Request.collection:type_name -> manager.Collection
	6,  // 7: manager.Publish.Request.auth:type_name -> manager.AuthConfig
	1,  // 8: manager.Publish.Response.diagnostics:type_name -> manager.Diagnostic
	9,  // 9: manager.CollectionManager.PublishContent:input_type -> manager.Publish.Request
	7,  // 10: manager.CollectionManager.RetrieveContent:input_type -> manager.Retrieve.Request
	10, // 11: manager.CollectionManager.PublishContent:output_type -> manager.Publish.Response
	8,  // 12: manager.CollectionManager.RetrieveContent:output_type -> manager.Retrieve.Response
	11, // [11:13] is the sub-list for method output_type
	9,  // [9:11] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_api_services_collectionmanager_v1alpha1_manager_proto_init() }
func file_api_services_collectionmanager_v1alpha1_manager_proto_init() {
	if File_api_services_collectionmanager_v1alpha1_manager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Diagnostic); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Retrieve); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Publish); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Collection); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthConfig); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Retrieve_Request); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Retrieve_Response); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Publish_Request); i {
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
		file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Publish_Response); i {
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
			RawDescriptor: file_api_services_collectionmanager_v1alpha1_manager_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_services_collectionmanager_v1alpha1_manager_proto_goTypes,
		DependencyIndexes: file_api_services_collectionmanager_v1alpha1_manager_proto_depIdxs,
		EnumInfos:         file_api_services_collectionmanager_v1alpha1_manager_proto_enumTypes,
		MessageInfos:      file_api_services_collectionmanager_v1alpha1_manager_proto_msgTypes,
	}.Build()
	File_api_services_collectionmanager_v1alpha1_manager_proto = out.File
	file_api_services_collectionmanager_v1alpha1_manager_proto_rawDesc = nil
	file_api_services_collectionmanager_v1alpha1_manager_proto_goTypes = nil
	file_api_services_collectionmanager_v1alpha1_manager_proto_depIdxs = nil
}
