// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: proto/rndz.proto

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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Types that are assignable to Cmd:
	//	*Request_Ping
	//	*Request_Isync
	//	*Request_Fsync
	//	*Request_Rsync
	//	*Request_Bye
	Cmd isRequest_Cmd `protobuf_oneof:"cmd"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (m *Request) GetCmd() isRequest_Cmd {
	if m != nil {
		return m.Cmd
	}
	return nil
}

func (x *Request) GetPing() *Ping {
	if x, ok := x.GetCmd().(*Request_Ping); ok {
		return x.Ping
	}
	return nil
}

func (x *Request) GetIsync() *Isync {
	if x, ok := x.GetCmd().(*Request_Isync); ok {
		return x.Isync
	}
	return nil
}

func (x *Request) GetFsync() *Fsync {
	if x, ok := x.GetCmd().(*Request_Fsync); ok {
		return x.Fsync
	}
	return nil
}

func (x *Request) GetRsync() *Rsync {
	if x, ok := x.GetCmd().(*Request_Rsync); ok {
		return x.Rsync
	}
	return nil
}

func (x *Request) GetBye() *Bye {
	if x, ok := x.GetCmd().(*Request_Bye); ok {
		return x.Bye
	}
	return nil
}

type isRequest_Cmd interface {
	isRequest_Cmd()
}

type Request_Ping struct {
	Ping *Ping `protobuf:"bytes,2,opt,name=Ping,proto3,oneof"`
}

type Request_Isync struct {
	Isync *Isync `protobuf:"bytes,3,opt,name=Isync,proto3,oneof"`
}

type Request_Fsync struct {
	Fsync *Fsync `protobuf:"bytes,4,opt,name=Fsync,proto3,oneof"`
}

type Request_Rsync struct {
	Rsync *Rsync `protobuf:"bytes,5,opt,name=Rsync,proto3,oneof"`
}

type Request_Bye struct {
	Bye *Bye `protobuf:"bytes,6,opt,name=Bye,proto3,oneof"`
}

func (*Request_Ping) isRequest_Cmd() {}

func (*Request_Isync) isRequest_Cmd() {}

func (*Request_Fsync) isRequest_Cmd() {}

func (*Request_Rsync) isRequest_Cmd() {}

func (*Request_Bye) isRequest_Cmd() {}

type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{1}
}

type Isync struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Isync) Reset() {
	*x = Isync{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Isync) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Isync) ProtoMessage() {}

func (x *Isync) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Isync.ProtoReflect.Descriptor instead.
func (*Isync) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{2}
}

func (x *Isync) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Types that are assignable to Cmd:
	//	*Response_Pong
	//	*Response_Redirect
	//	*Response_Fsync
	Cmd isResponse_Cmd `protobuf_oneof:"cmd"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{3}
}

func (x *Response) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (m *Response) GetCmd() isResponse_Cmd {
	if m != nil {
		return m.Cmd
	}
	return nil
}

func (x *Response) GetPong() *Pong {
	if x, ok := x.GetCmd().(*Response_Pong); ok {
		return x.Pong
	}
	return nil
}

func (x *Response) GetRedirect() *Redirect {
	if x, ok := x.GetCmd().(*Response_Redirect); ok {
		return x.Redirect
	}
	return nil
}

func (x *Response) GetFsync() *Fsync {
	if x, ok := x.GetCmd().(*Response_Fsync); ok {
		return x.Fsync
	}
	return nil
}

type isResponse_Cmd interface {
	isResponse_Cmd()
}

type Response_Pong struct {
	Pong *Pong `protobuf:"bytes,2,opt,name=Pong,proto3,oneof"`
}

type Response_Redirect struct {
	Redirect *Redirect `protobuf:"bytes,3,opt,name=Redirect,proto3,oneof"`
}

type Response_Fsync struct {
	Fsync *Fsync `protobuf:"bytes,4,opt,name=Fsync,proto3,oneof"`
}

func (*Response_Pong) isResponse_Cmd() {}

func (*Response_Redirect) isResponse_Cmd() {}

func (*Response_Fsync) isResponse_Cmd() {}

type Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Pong) Reset() {
	*x = Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pong) ProtoMessage() {}

func (x *Pong) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pong.ProtoReflect.Descriptor instead.
func (*Pong) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{4}
}

type Redirect struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
}

func (x *Redirect) Reset() {
	*x = Redirect{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Redirect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Redirect) ProtoMessage() {}

func (x *Redirect) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Redirect.ProtoReflect.Descriptor instead.
func (*Redirect) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{5}
}

func (x *Redirect) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Redirect) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

type Fsync struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
}

func (x *Fsync) Reset() {
	*x = Fsync{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Fsync) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Fsync) ProtoMessage() {}

func (x *Fsync) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Fsync.ProtoReflect.Descriptor instead.
func (*Fsync) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{6}
}

func (x *Fsync) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Fsync) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

type Rsync struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Rsync) Reset() {
	*x = Rsync{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rsync) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rsync) ProtoMessage() {}

func (x *Rsync) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rsync.ProtoReflect.Descriptor instead.
func (*Rsync) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{7}
}

func (x *Rsync) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Bye struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Bye) Reset() {
	*x = Bye{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rndz_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bye) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bye) ProtoMessage() {}

func (x *Bye) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rndz_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bye.ProtoReflect.Descriptor instead.
func (*Bye) Descriptor() ([]byte, []int) {
	return file_proto_rndz_proto_rawDescGZIP(), []int{8}
}

var File_proto_rndz_proto protoreflect.FileDescriptor

var file_proto_rndz_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x6e, 0x64, 0x7a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xb7, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b,
	0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x1e, 0x0a, 0x05, 0x49,
	0x73, 0x79, 0x6e, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x49, 0x73, 0x79,
	0x6e, 0x63, 0x48, 0x00, 0x52, 0x05, 0x49, 0x73, 0x79, 0x6e, 0x63, 0x12, 0x1e, 0x0a, 0x05, 0x46,
	0x73, 0x79, 0x6e, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x46, 0x73, 0x79,
	0x6e, 0x63, 0x48, 0x00, 0x52, 0x05, 0x46, 0x73, 0x79, 0x6e, 0x63, 0x12, 0x1e, 0x0a, 0x05, 0x52,
	0x73, 0x79, 0x6e, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x52, 0x73, 0x79,
	0x6e, 0x63, 0x48, 0x00, 0x52, 0x05, 0x52, 0x73, 0x79, 0x6e, 0x63, 0x12, 0x18, 0x0a, 0x03, 0x42,
	0x79, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x42, 0x79, 0x65, 0x48, 0x00,
	0x52, 0x03, 0x42, 0x79, 0x65, 0x42, 0x05, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x22, 0x06, 0x0a, 0x04,
	0x50, 0x69, 0x6e, 0x67, 0x22, 0x17, 0x0a, 0x05, 0x49, 0x73, 0x79, 0x6e, 0x63, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x87, 0x01,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x04, 0x50, 0x6f,
	0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x48,
	0x00, 0x52, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x12, 0x27, 0x0a, 0x08, 0x52, 0x65, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x52, 0x65, 0x64, 0x69,
	0x72, 0x65, 0x63, 0x74, 0x48, 0x00, 0x52, 0x08, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x12, 0x1e, 0x0a, 0x05, 0x46, 0x73, 0x79, 0x6e, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x06, 0x2e, 0x46, 0x73, 0x79, 0x6e, 0x63, 0x48, 0x00, 0x52, 0x05, 0x46, 0x73, 0x79, 0x6e, 0x63,
	0x42, 0x05, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x22, 0x06, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x22,
	0x2e, 0x0a, 0x08, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61,
	0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x22,
	0x2b, 0x0a, 0x05, 0x46, 0x73, 0x79, 0x6e, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x22, 0x17, 0x0a, 0x05,
	0x52, 0x73, 0x79, 0x6e, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x05, 0x0a, 0x03, 0x42, 0x79, 0x65, 0x42, 0x0a, 0x5a, 0x08,
	0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_rndz_proto_rawDescOnce sync.Once
	file_proto_rndz_proto_rawDescData = file_proto_rndz_proto_rawDesc
)

func file_proto_rndz_proto_rawDescGZIP() []byte {
	file_proto_rndz_proto_rawDescOnce.Do(func() {
		file_proto_rndz_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_rndz_proto_rawDescData)
	})
	return file_proto_rndz_proto_rawDescData
}

var file_proto_rndz_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_rndz_proto_goTypes = []interface{}{
	(*Request)(nil),  // 0: Request
	(*Ping)(nil),     // 1: Ping
	(*Isync)(nil),    // 2: Isync
	(*Response)(nil), // 3: Response
	(*Pong)(nil),     // 4: Pong
	(*Redirect)(nil), // 5: Redirect
	(*Fsync)(nil),    // 6: Fsync
	(*Rsync)(nil),    // 7: Rsync
	(*Bye)(nil),      // 8: Bye
}
var file_proto_rndz_proto_depIdxs = []int32{
	1, // 0: Request.Ping:type_name -> Ping
	2, // 1: Request.Isync:type_name -> Isync
	6, // 2: Request.Fsync:type_name -> Fsync
	7, // 3: Request.Rsync:type_name -> Rsync
	8, // 4: Request.Bye:type_name -> Bye
	4, // 5: Response.Pong:type_name -> Pong
	5, // 6: Response.Redirect:type_name -> Redirect
	6, // 7: Response.Fsync:type_name -> Fsync
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_proto_rndz_proto_init() }
func file_proto_rndz_proto_init() {
	if File_proto_rndz_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_rndz_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_proto_rndz_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ping); i {
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
		file_proto_rndz_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Isync); i {
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
		file_proto_rndz_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_proto_rndz_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pong); i {
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
		file_proto_rndz_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Redirect); i {
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
		file_proto_rndz_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Fsync); i {
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
		file_proto_rndz_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rsync); i {
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
		file_proto_rndz_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bye); i {
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
	file_proto_rndz_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Request_Ping)(nil),
		(*Request_Isync)(nil),
		(*Request_Fsync)(nil),
		(*Request_Rsync)(nil),
		(*Request_Bye)(nil),
	}
	file_proto_rndz_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Response_Pong)(nil),
		(*Response_Redirect)(nil),
		(*Response_Fsync)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_rndz_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_rndz_proto_goTypes,
		DependencyIndexes: file_proto_rndz_proto_depIdxs,
		MessageInfos:      file_proto_rndz_proto_msgTypes,
	}.Build()
	File_proto_rndz_proto = out.File
	file_proto_rndz_proto_rawDesc = nil
	file_proto_rndz_proto_goTypes = nil
	file_proto_rndz_proto_depIdxs = nil
}
