// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.25.3
// source: messages.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Profile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login *string `protobuf:"bytes,1,opt,name=login,proto3,oneof" json:"login,omitempty"`
	Email *string `protobuf:"bytes,2,opt,name=email,proto3,oneof" json:"email,omitempty"`
	Phone *string `protobuf:"bytes,3,opt,name=phone,proto3,oneof" json:"phone,omitempty"`
}

func (x *Profile) Reset() {
	*x = Profile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Profile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profile) ProtoMessage() {}

func (x *Profile) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profile.ProtoReflect.Descriptor instead.
func (*Profile) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0}
}

func (x *Profile) GetLogin() string {
	if x != nil && x.Login != nil {
		return *x.Login
	}
	return ""
}

func (x *Profile) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *Profile) GetPhone() string {
	if x != nil && x.Phone != nil {
		return *x.Phone
	}
	return ""
}

type AccessToken struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *AccessToken) Reset() {
	*x = AccessToken{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessToken) ProtoMessage() {}

func (x *AccessToken) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessToken.ProtoReflect.Descriptor instead.
func (*AccessToken) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{1}
}

func (x *AccessToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RefreshToken struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *RefreshToken) Reset() {
	*x = RefreshToken{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshToken) ProtoMessage() {}

func (x *RefreshToken) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshToken.ProtoReflect.Descriptor instead.
func (*RefreshToken) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{2}
}

func (x *RefreshToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type SuperAccessToken struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *SuperAccessToken) Reset() {
	*x = SuperAccessToken{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SuperAccessToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SuperAccessToken) ProtoMessage() {}

func (x *SuperAccessToken) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SuperAccessToken.ProtoReflect.Descriptor instead.
func (*SuperAccessToken) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{3}
}

func (x *SuperAccessToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type SessionInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Browser *string `protobuf:"bytes,1,opt,name=browser,proto3,oneof" json:"browser,omitempty"`
	Os      *string `protobuf:"bytes,2,opt,name=os,proto3,oneof" json:"os,omitempty"`
	Ip      *string `protobuf:"bytes,3,opt,name=ip,proto3,oneof" json:"ip,omitempty"`
}

func (x *SessionInfo) Reset() {
	*x = SessionInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionInfo) ProtoMessage() {}

func (x *SessionInfo) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionInfo.ProtoReflect.Descriptor instead.
func (*SessionInfo) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{4}
}

func (x *SessionInfo) GetBrowser() string {
	if x != nil && x.Browser != nil {
		return *x.Browser
	}
	return ""
}

func (x *SessionInfo) GetOs() string {
	if x != nil && x.Os != nil {
		return *x.Os
	}
	return ""
}

func (x *SessionInfo) GetIp() string {
	if x != nil && x.Ip != nil {
		return *x.Ip
	}
	return ""
}

type Sessions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sessions []*Sessions_Session `protobuf:"bytes,1,rep,name=sessions,proto3" json:"sessions,omitempty"`
}

func (x *Sessions) Reset() {
	*x = Sessions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sessions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sessions) ProtoMessage() {}

func (x *Sessions) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sessions.ProtoReflect.Descriptor instead.
func (*Sessions) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{5}
}

func (x *Sessions) GetSessions() []*Sessions_Session {
	if x != nil {
		return x.Sessions
	}
	return nil
}

type ProfileChangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AnyAccessToken string   `protobuf:"bytes,1,opt,name=any_access_token,json=anyAccessToken,proto3" json:"any_access_token,omitempty"`
	NewProfile     *Profile `protobuf:"bytes,2,opt,name=new_profile,json=newProfile,proto3" json:"new_profile,omitempty"`
}

func (x *ProfileChangeRequest) Reset() {
	*x = ProfileChangeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileChangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileChangeRequest) ProtoMessage() {}

func (x *ProfileChangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileChangeRequest.ProtoReflect.Descriptor instead.
func (*ProfileChangeRequest) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{6}
}

func (x *ProfileChangeRequest) GetAnyAccessToken() string {
	if x != nil {
		return x.AnyAccessToken
	}
	return ""
}

func (x *ProfileChangeRequest) GetNewProfile() *Profile {
	if x != nil {
		return x.NewProfile
	}
	return nil
}

type Sessions_Session struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info         *SessionInfo           `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
	TerminatedAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=terminated_at,json=terminatedAt,proto3,oneof" json:"terminated_at,omitempty"`
}

func (x *Sessions_Session) Reset() {
	*x = Sessions_Session{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sessions_Session) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sessions_Session) ProtoMessage() {}

func (x *Sessions_Session) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sessions_Session.ProtoReflect.Descriptor instead.
func (*Sessions_Session) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{5, 0}
}

func (x *Sessions_Session) GetInfo() *SessionInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

func (x *Sessions_Session) GetTerminatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.TerminatedAt
	}
	return nil
}

var File_messages_proto protoreflect.FileDescriptor

var file_messages_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x03, 0x73, 0x73, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x78, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x12, 0x19, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x88,
	0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x42, 0x08, 0x0a, 0x06,
	0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x22, 0x23, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x24, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x28, 0x0a, 0x10, 0x53,
	0x75, 0x70, 0x65, 0x72, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x70, 0x0a, 0x0b, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1d, 0x0a, 0x07, 0x62, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x62, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x72,
	0x88, 0x01, 0x01, 0x12, 0x13, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x02, 0x6f, 0x73, 0x88, 0x01, 0x01, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x02, 0x69, 0x70, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x62, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x72, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x6f, 0x73,
	0x42, 0x05, 0x0a, 0x03, 0x5f, 0x69, 0x70, 0x22, 0xc7, 0x01, 0x0a, 0x08, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x31, 0x0a, 0x08, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x73, 0x6f, 0x2e, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x87, 0x01, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x73, 0x73, 0x6f, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x44, 0x0a, 0x0d, 0x74, 0x65, 0x72,
	0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x0c,
	0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x42,
	0x10, 0x0a, 0x0e, 0x5f, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x22, 0x6f, 0x0a, 0x14, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x6e, 0x79,
	0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x6e, 0x79, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x2d, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x73, 0x73, 0x6f, 0x2e, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x0a, 0x6e, 0x65, 0x77, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x42, 0x17, 0x5a, 0x05, 0x2e, 0x3b, 0x67, 0x65, 0x6e, 0xaa, 0x02, 0x0d, 0x4f, 0x6e,
	0x6e, 0x79, 0x77, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x53, 0x53, 0x4f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_messages_proto_rawDescOnce sync.Once
	file_messages_proto_rawDescData = file_messages_proto_rawDesc
)

func file_messages_proto_rawDescGZIP() []byte {
	file_messages_proto_rawDescOnce.Do(func() {
		file_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_proto_rawDescData)
	})
	return file_messages_proto_rawDescData
}

var file_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_messages_proto_goTypes = []interface{}{
	(*Profile)(nil),               // 0: sso.Profile
	(*AccessToken)(nil),           // 1: sso.AccessToken
	(*RefreshToken)(nil),          // 2: sso.RefreshToken
	(*SuperAccessToken)(nil),      // 3: sso.SuperAccessToken
	(*SessionInfo)(nil),           // 4: sso.SessionInfo
	(*Sessions)(nil),              // 5: sso.Sessions
	(*ProfileChangeRequest)(nil),  // 6: sso.ProfileChangeRequest
	(*Sessions_Session)(nil),      // 7: sso.Sessions.Session
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_messages_proto_depIdxs = []int32{
	7, // 0: sso.Sessions.sessions:type_name -> sso.Sessions.Session
	0, // 1: sso.ProfileChangeRequest.new_profile:type_name -> sso.Profile
	4, // 2: sso.Sessions.Session.info:type_name -> sso.SessionInfo
	8, // 3: sso.Sessions.Session.terminated_at:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_messages_proto_init() }
func file_messages_proto_init() {
	if File_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Profile); i {
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
		file_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessToken); i {
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
		file_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshToken); i {
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
		file_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SuperAccessToken); i {
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
		file_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionInfo); i {
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
		file_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sessions); i {
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
		file_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileChangeRequest); i {
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
		file_messages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sessions_Session); i {
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
	file_messages_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_messages_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_messages_proto_msgTypes[7].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_proto_goTypes,
		DependencyIndexes: file_messages_proto_depIdxs,
		MessageInfos:      file_messages_proto_msgTypes,
	}.Build()
	File_messages_proto = out.File
	file_messages_proto_rawDesc = nil
	file_messages_proto_goTypes = nil
	file_messages_proto_depIdxs = nil
}