// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: sendMsg.proto

//指定包名

package sendmsg

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

//定义结构体
type SendMsgRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//定义用户名
	RequestType string            `protobuf:"bytes,1,opt,name=requestType,proto3" json:"requestType,omitempty"`
	Receiver    string            `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Data        map[string]string `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	UpdateType  string            `protobuf:"bytes,4,opt,name=updateType,proto3" json:"updateType,omitempty"`
	Msg         string            `protobuf:"bytes,5,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *SendMsgRequest) Reset() {
	*x = SendMsgRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sendMsg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMsgRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMsgRequest) ProtoMessage() {}

func (x *SendMsgRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sendMsg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMsgRequest.ProtoReflect.Descriptor instead.
func (*SendMsgRequest) Descriptor() ([]byte, []int) {
	return file_sendMsg_proto_rawDescGZIP(), []int{0}
}

func (x *SendMsgRequest) GetRequestType() string {
	if x != nil {
		return x.RequestType
	}
	return ""
}

func (x *SendMsgRequest) GetReceiver() string {
	if x != nil {
		return x.Receiver
	}
	return ""
}

func (x *SendMsgRequest) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *SendMsgRequest) GetUpdateType() string {
	if x != nil {
		return x.UpdateType
	}
	return ""
}

func (x *SendMsgRequest) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

//响应结构体
type SendMsgResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *SendMsgResponse) Reset() {
	*x = SendMsgResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sendMsg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMsgResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMsgResponse) ProtoMessage() {}

func (x *SendMsgResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sendMsg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMsgResponse.ProtoReflect.Descriptor instead.
func (*SendMsgResponse) Descriptor() ([]byte, []int) {
	return file_sendMsg_proto_rawDescGZIP(), []int{1}
}

func (x *SendMsgResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *SendMsgResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_sendMsg_proto protoreflect.FileDescriptor

var file_sendMsg_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x65, 0x6e, 0x64, 0x6d, 0x73, 0x67, 0x22, 0xf0, 0x01, 0x0a, 0x0e, 0x53, 0x65, 0x6e,
	0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x6d, 0x73,
	0x67, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d,
	0x73, 0x67, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x37, 0x0a, 0x0f, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x32, 0x50, 0x0a, 0x0e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73,
	0x67, 0x12, 0x17, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x6d, 0x73, 0x67, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x65, 0x6e,
	0x64, 0x6d, 0x73, 0x67, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x14, 0x5a, 0x12, 0x2e, 0x2e, 0x2f, 0x73, 0x65, 0x6e,
	0x64, 0x6d, 0x73, 0x67, 0x3b, 0x73, 0x65, 0x6e, 0x64, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sendMsg_proto_rawDescOnce sync.Once
	file_sendMsg_proto_rawDescData = file_sendMsg_proto_rawDesc
)

func file_sendMsg_proto_rawDescGZIP() []byte {
	file_sendMsg_proto_rawDescOnce.Do(func() {
		file_sendMsg_proto_rawDescData = protoimpl.X.CompressGZIP(file_sendMsg_proto_rawDescData)
	})
	return file_sendMsg_proto_rawDescData
}

var file_sendMsg_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_sendMsg_proto_goTypes = []interface{}{
	(*SendMsgRequest)(nil),  // 0: sendmsg.SendMsgRequest
	(*SendMsgResponse)(nil), // 1: sendmsg.SendMsgResponse
	nil,                     // 2: sendmsg.SendMsgRequest.DataEntry
}
var file_sendMsg_proto_depIdxs = []int32{
	2, // 0: sendmsg.SendMsgRequest.data:type_name -> sendmsg.SendMsgRequest.DataEntry
	0, // 1: sendmsg.SendMsgService.SendMsg:input_type -> sendmsg.SendMsgRequest
	1, // 2: sendmsg.SendMsgService.SendMsg:output_type -> sendmsg.SendMsgResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_sendMsg_proto_init() }
func file_sendMsg_proto_init() {
	if File_sendMsg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sendMsg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMsgRequest); i {
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
		file_sendMsg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMsgResponse); i {
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
			RawDescriptor: file_sendMsg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sendMsg_proto_goTypes,
		DependencyIndexes: file_sendMsg_proto_depIdxs,
		MessageInfos:      file_sendMsg_proto_msgTypes,
	}.Build()
	File_sendMsg_proto = out.File
	file_sendMsg_proto_rawDesc = nil
	file_sendMsg_proto_goTypes = nil
	file_sendMsg_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SendMsgServiceClient is the client API for SendMsgService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SendMsgServiceClient interface {
	SendMsg(ctx context.Context, in *SendMsgRequest, opts ...grpc.CallOption) (*SendMsgResponse, error)
}

type sendMsgServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSendMsgServiceClient(cc grpc.ClientConnInterface) SendMsgServiceClient {
	return &sendMsgServiceClient{cc}
}

func (c *sendMsgServiceClient) SendMsg(ctx context.Context, in *SendMsgRequest, opts ...grpc.CallOption) (*SendMsgResponse, error) {
	out := new(SendMsgResponse)
	err := c.cc.Invoke(ctx, "/sendmsg.SendMsgService/SendMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SendMsgServiceServer is the server API for SendMsgService service.
type SendMsgServiceServer interface {
	SendMsg(context.Context, *SendMsgRequest) (*SendMsgResponse, error)
}

// UnimplementedSendMsgServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSendMsgServiceServer struct {
}

func (*UnimplementedSendMsgServiceServer) SendMsg(context.Context, *SendMsgRequest) (*SendMsgResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg not implemented")
}

func RegisterSendMsgServiceServer(s *grpc.Server, srv SendMsgServiceServer) {
	s.RegisterService(&_SendMsgService_serviceDesc, srv)
}

func _SendMsgService_SendMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendMsgServiceServer).SendMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sendmsg.SendMsgService/SendMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendMsgServiceServer).SendMsg(ctx, req.(*SendMsgRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SendMsgService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sendmsg.SendMsgService",
	HandlerType: (*SendMsgServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMsg",
			Handler:    _SendMsgService_SendMsg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sendMsg.proto",
}
