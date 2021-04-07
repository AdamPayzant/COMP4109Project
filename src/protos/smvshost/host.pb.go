// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.6
// source: smvshost/host.proto

package smvshost

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// As the name describes
type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[0]
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
	return file_smvshost_host_proto_rawDescGZIP(), []int{0}
}

// The CA of a user
type CA struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ca []byte `protobuf:"bytes,1,opt,name=ca,proto3" json:"ca,omitempty"`
}

func (x *CA) Reset() {
	*x = CA{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CA) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CA) ProtoMessage() {}

func (x *CA) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CA.ProtoReflect.Descriptor instead.
func (*CA) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{1}
}

func (x *CA) GetCa() []byte {
	if x != nil {
		return x.Ca
	}
	return nil
}

// Contains the client device's auth token
type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{2}
}

func (x *Token) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// A generic status response message
type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status int32 `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{3}
}

func (x *Status) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type DeleteReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User      string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	MessageID int64  `protobuf:"varint,2,opt,name=messageID,proto3" json:"messageID,omitempty"`
	Token     string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *DeleteReq) Reset() {
	*x = DeleteReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteReq) ProtoMessage() {}

func (x *DeleteReq) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteReq.ProtoReflect.Descriptor instead.
func (*DeleteReq) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteReq) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *DeleteReq) GetMessageID() int64 {
	if x != nil {
		return x.MessageID
	}
	return 0
}

func (x *DeleteReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// Contains a secret to be shared between two hosts
type InitMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secret string `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
}

func (x *InitMessage) Reset() {
	*x = InitMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitMessage) ProtoMessage() {}

func (x *InitMessage) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitMessage.ProtoReflect.Descriptor instead.
func (*InitMessage) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{5}
}

func (x *InitMessage) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

// The message to be sent between hosts for texts
type H2HText struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *ListofMessages `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	User    string          `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Secret  string          `protobuf:"bytes,3,opt,name=secret,proto3" json:"secret,omitempty"`
}

func (x *H2HText) Reset() {
	*x = H2HText{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *H2HText) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*H2HText) ProtoMessage() {}

func (x *H2HText) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use H2HText.ProtoReflect.Descriptor instead.
func (*H2HText) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{6}
}

func (x *H2HText) GetMessage() *ListofMessages {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *H2HText) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *H2HText) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

// Because the way protobuffers (don't) do lists, this is the best way to do this
type ListofMessages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []string `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *ListofMessages) Reset() {
	*x = ListofMessages{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListofMessages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListofMessages) ProtoMessage() {}

func (x *ListofMessages) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListofMessages.ProtoReflect.Descriptor instead.
func (*ListofMessages) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{7}
}

func (x *ListofMessages) GetMessages() []string {
	if x != nil {
		return x.Messages
	}
	return nil
}

type ClientText struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetUser string          `protobuf:"bytes,1,opt,name=targetUser,proto3" json:"targetUser,omitempty"`
	Message    *ListofMessages `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Token      string          `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"` // The client device's auth token
}

func (x *ClientText) Reset() {
	*x = ClientText{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientText) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientText) ProtoMessage() {}

func (x *ClientText) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientText.ProtoReflect.Descriptor instead.
func (*ClientText) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{8}
}

func (x *ClientText) GetTargetUser() string {
	if x != nil {
		return x.TargetUser
	}
	return ""
}

func (x *ClientText) GetMessage() *ListofMessages {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *ClientText) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type Username struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token    string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *Username) Reset() {
	*x = Username{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Username) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Username) ProtoMessage() {}

func (x *Username) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Username.ProtoReflect.Descriptor instead.
func (*Username) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{9}
}

func (x *Username) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Username) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// This will just send the full conversation
type Conversation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Convo []byte `protobuf:"bytes,1,opt,name=convo,proto3" json:"convo,omitempty"`
}

func (x *Conversation) Reset() {
	*x = Conversation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smvshost_host_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Conversation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Conversation) ProtoMessage() {}

func (x *Conversation) ProtoReflect() protoreflect.Message {
	mi := &file_smvshost_host_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Conversation.ProtoReflect.Descriptor instead.
func (*Conversation) Descriptor() ([]byte, []int) {
	return file_smvshost_host_proto_rawDescGZIP(), []int{10}
}

func (x *Conversation) GetConvo() []byte {
	if x != nil {
		return x.Convo
	}
	return nil
}

var File_smvshost_host_proto protoreflect.FileDescriptor

var file_smvshost_host_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x6d, 0x76, 0x73, 0x68, 0x6f, 0x73, 0x74, 0x2f, 0x68, 0x6f, 0x73, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x73, 0x6d, 0x76, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x14, 0x0a, 0x02, 0x43, 0x41, 0x12, 0x0e, 0x0a, 0x02, 0x63, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x63, 0x61, 0x22, 0x1d, 0x0a, 0x05, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x20, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x53, 0x0a, 0x09, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0x25, 0x0a, 0x0b, 0x49, 0x6e, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x65, 0x0a, 0x07, 0x48, 0x32, 0x48, 0x54, 0x65,
	0x78, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x6d, 0x76, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x6f,
	0x66, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x2c,
	0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x6f, 0x66, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0x72, 0x0a, 0x0a,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x6d,
	0x76, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x6f, 0x66, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0x3c, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x24,
	0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14,
	0x0a, 0x05, 0x63, 0x6f, 0x6e, 0x76, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x63,
	0x6f, 0x6e, 0x76, 0x6f, 0x32, 0x2e, 0x0a, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x41,
	0x48, 0x6f, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x05, 0x47, 0x65, 0x74, 0x43, 0x41, 0x12, 0x0b, 0x2e,
	0x73, 0x6d, 0x76, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x08, 0x2e, 0x73, 0x6d, 0x76,
	0x73, 0x2e, 0x43, 0x41, 0x32, 0xd4, 0x02, 0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x48,
	0x6f, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x05, 0x52, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x0b, 0x2e, 0x73,
	0x6d, 0x76, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x1a, 0x0c, 0x2e, 0x73, 0x6d, 0x76, 0x73,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2e, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0f, 0x2e, 0x73, 0x6d, 0x76, 0x73, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x73, 0x6d, 0x76, 0x73,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x32, 0x0a, 0x0f, 0x49, 0x6e, 0x69, 0x74, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x43, 0x6f, 0x6e, 0x76, 0x6f, 0x12, 0x11, 0x2e, 0x73, 0x6d, 0x76,
	0x73, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0c, 0x2e,
	0x73, 0x6d, 0x76, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x0c, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x43, 0x6f, 0x6e, 0x76, 0x6f, 0x12, 0x11, 0x2e, 0x73, 0x6d,
	0x76, 0x73, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0c,
	0x2e, 0x73, 0x6d, 0x76, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2a, 0x0a, 0x08,
	0x53, 0x65, 0x6e, 0x64, 0x54, 0x65, 0x78, 0x74, 0x12, 0x10, 0x2e, 0x73, 0x6d, 0x76, 0x73, 0x2e,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x78, 0x74, 0x1a, 0x0c, 0x2e, 0x73, 0x6d, 0x76,
	0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2a, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x69,
	0x65, 0x76, 0x65, 0x54, 0x65, 0x78, 0x74, 0x12, 0x0d, 0x2e, 0x73, 0x6d, 0x76, 0x73, 0x2e, 0x48,
	0x32, 0x48, 0x54, 0x65, 0x78, 0x74, 0x1a, 0x0c, 0x2e, 0x73, 0x6d, 0x76, 0x73, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x35, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x76, 0x65,
	0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x2e, 0x73, 0x6d, 0x76, 0x73, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x1a, 0x12, 0x2e, 0x73, 0x6d, 0x76, 0x73, 0x2e, 0x43,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x3c, 0x5a, 0x3a, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x64, 0x61, 0x6d, 0x50, 0x61,
	0x79, 0x7a, 0x61, 0x6e, 0x74, 0x2f, 0x43, 0x4f, 0x4d, 0x50, 0x34, 0x31, 0x30, 0x39, 0x50, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2f, 0x73, 0x6d, 0x76, 0x73, 0x68, 0x6f, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_smvshost_host_proto_rawDescOnce sync.Once
	file_smvshost_host_proto_rawDescData = file_smvshost_host_proto_rawDesc
)

func file_smvshost_host_proto_rawDescGZIP() []byte {
	file_smvshost_host_proto_rawDescOnce.Do(func() {
		file_smvshost_host_proto_rawDescData = protoimpl.X.CompressGZIP(file_smvshost_host_proto_rawDescData)
	})
	return file_smvshost_host_proto_rawDescData
}

var file_smvshost_host_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_smvshost_host_proto_goTypes = []interface{}{
	(*Empty)(nil),          // 0: smvs.Empty
	(*CA)(nil),             // 1: smvs.CA
	(*Token)(nil),          // 2: smvs.Token
	(*Status)(nil),         // 3: smvs.Status
	(*DeleteReq)(nil),      // 4: smvs.DeleteReq
	(*InitMessage)(nil),    // 5: smvs.InitMessage
	(*H2HText)(nil),        // 6: smvs.H2HText
	(*ListofMessages)(nil), // 7: smvs.ListofMessages
	(*ClientText)(nil),     // 8: smvs.ClientText
	(*Username)(nil),       // 9: smvs.Username
	(*Conversation)(nil),   // 10: smvs.Conversation
}
var file_smvshost_host_proto_depIdxs = []int32{
	7,  // 0: smvs.H2HText.message:type_name -> smvs.ListofMessages
	7,  // 1: smvs.ClientText.message:type_name -> smvs.ListofMessages
	0,  // 2: smvs.clientCAHost.GetCA:input_type -> smvs.Empty
	2,  // 3: smvs.clientHost.ReKey:input_type -> smvs.Token
	4,  // 4: smvs.clientHost.DeleteMessage:input_type -> smvs.DeleteReq
	5,  // 5: smvs.clientHost.InitializeConvo:input_type -> smvs.InitMessage
	5,  // 6: smvs.clientHost.ConfirmConvo:input_type -> smvs.InitMessage
	8,  // 7: smvs.clientHost.SendText:input_type -> smvs.ClientText
	6,  // 8: smvs.clientHost.RecieveText:input_type -> smvs.H2HText
	9,  // 9: smvs.clientHost.GetConversation:input_type -> smvs.Username
	1,  // 10: smvs.clientCAHost.GetCA:output_type -> smvs.CA
	3,  // 11: smvs.clientHost.ReKey:output_type -> smvs.Status
	3,  // 12: smvs.clientHost.DeleteMessage:output_type -> smvs.Status
	3,  // 13: smvs.clientHost.InitializeConvo:output_type -> smvs.Status
	3,  // 14: smvs.clientHost.ConfirmConvo:output_type -> smvs.Status
	3,  // 15: smvs.clientHost.SendText:output_type -> smvs.Status
	3,  // 16: smvs.clientHost.RecieveText:output_type -> smvs.Status
	10, // 17: smvs.clientHost.GetConversation:output_type -> smvs.Conversation
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_smvshost_host_proto_init() }
func file_smvshost_host_proto_init() {
	if File_smvshost_host_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_smvshost_host_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_smvshost_host_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CA); i {
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
		file_smvshost_host_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Token); i {
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
		file_smvshost_host_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_smvshost_host_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteReq); i {
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
		file_smvshost_host_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitMessage); i {
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
		file_smvshost_host_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*H2HText); i {
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
		file_smvshost_host_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListofMessages); i {
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
		file_smvshost_host_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientText); i {
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
		file_smvshost_host_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Username); i {
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
		file_smvshost_host_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Conversation); i {
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
			RawDescriptor: file_smvshost_host_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_smvshost_host_proto_goTypes,
		DependencyIndexes: file_smvshost_host_proto_depIdxs,
		MessageInfos:      file_smvshost_host_proto_msgTypes,
	}.Build()
	File_smvshost_host_proto = out.File
	file_smvshost_host_proto_rawDesc = nil
	file_smvshost_host_proto_goTypes = nil
	file_smvshost_host_proto_depIdxs = nil
}
