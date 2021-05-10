// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/consuming/consuming.proto

package consuming

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Message send to service when event occurs
type ConsumeRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Topic                string   `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Body                 []byte   `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsumeRequest) Reset()         { *m = ConsumeRequest{} }
func (m *ConsumeRequest) String() string { return proto.CompactTextString(m) }
func (*ConsumeRequest) ProtoMessage()    {}
func (*ConsumeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_61f24985cffdbec6, []int{0}
}

func (m *ConsumeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumeRequest.Unmarshal(m, b)
}
func (m *ConsumeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumeRequest.Marshal(b, m, deterministic)
}
func (m *ConsumeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumeRequest.Merge(m, src)
}
func (m *ConsumeRequest) XXX_Size() int {
	return xxx_messageInfo_ConsumeRequest.Size(m)
}
func (m *ConsumeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumeRequest proto.InternalMessageInfo

func (m *ConsumeRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ConsumeRequest) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *ConsumeRequest) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type ConsumeResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsumeResponse) Reset()         { *m = ConsumeResponse{} }
func (m *ConsumeResponse) String() string { return proto.CompactTextString(m) }
func (*ConsumeResponse) ProtoMessage()    {}
func (*ConsumeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_61f24985cffdbec6, []int{1}
}

func (m *ConsumeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumeResponse.Unmarshal(m, b)
}
func (m *ConsumeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumeResponse.Marshal(b, m, deterministic)
}
func (m *ConsumeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumeResponse.Merge(m, src)
}
func (m *ConsumeResponse) XXX_Size() int {
	return xxx_messageInfo_ConsumeResponse.Size(m)
}
func (m *ConsumeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumeResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ConsumeRequest)(nil), "consuming.ConsumeRequest")
	proto.RegisterType((*ConsumeResponse)(nil), "consuming.ConsumeResponse")
}

func init() { proto.RegisterFile("proto/consuming/consuming.proto", fileDescriptor_61f24985cffdbec6) }

var fileDescriptor_61f24985cffdbec6 = []byte{
	// 170 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2f, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xce, 0xcf, 0x2b, 0x2e, 0xcd, 0xcd, 0xcc, 0x4b, 0x47, 0xb0, 0xf4, 0xc0, 0x32,
	0x42, 0x9c, 0x70, 0x01, 0x25, 0x2f, 0x2e, 0x3e, 0x67, 0x30, 0x27, 0x35, 0x28, 0xb5, 0xb0, 0x34,
	0xb5, 0xb8, 0x44, 0x88, 0x8f, 0x8b, 0x29, 0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x88,
	0x29, 0x33, 0x45, 0x48, 0x84, 0x8b, 0xb5, 0x24, 0xbf, 0x20, 0x33, 0x59, 0x82, 0x09, 0x2c, 0x04,
	0xe1, 0x08, 0x09, 0x71, 0xb1, 0x24, 0xe5, 0xa7, 0x54, 0x4a, 0x30, 0x2b, 0x30, 0x6a, 0xf0, 0x04,
	0x81, 0xd9, 0x4a, 0x82, 0x5c, 0xfc, 0x70, 0xb3, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x8d, 0xc2,
	0xb8, 0x04, 0x9c, 0x61, 0x76, 0x05, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0x0a, 0x39, 0x71, 0xb1,
	0x43, 0x95, 0x09, 0x49, 0xea, 0x21, 0x9c, 0x86, 0xea, 0x0c, 0x29, 0x29, 0x6c, 0x52, 0x10, 0x53,
	0x95, 0x18, 0x9c, 0xb8, 0xa3, 0x10, 0x7e, 0x48, 0x62, 0x03, 0xfb, 0xca, 0x18, 0x10, 0x00, 0x00,
	0xff, 0xff, 0xb9, 0x8c, 0xb7, 0x80, 0xf8, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConsumingServiceClient is the client API for ConsumingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConsumingServiceClient interface {
	Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (*ConsumeResponse, error)
}

type consumingServiceClient struct {
	cc *grpc.ClientConn
}

func NewConsumingServiceClient(cc *grpc.ClientConn) ConsumingServiceClient {
	return &consumingServiceClient{cc}
}

func (c *consumingServiceClient) Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (*ConsumeResponse, error) {
	out := new(ConsumeResponse)
	err := c.cc.Invoke(ctx, "/consuming.ConsumingService/Consume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsumingServiceServer is the server API for ConsumingService service.
type ConsumingServiceServer interface {
	Consume(context.Context, *ConsumeRequest) (*ConsumeResponse, error)
}

// UnimplementedConsumingServiceServer can be embedded to have forward compatible implementations.
type UnimplementedConsumingServiceServer struct {
}

func (*UnimplementedConsumingServiceServer) Consume(ctx context.Context, req *ConsumeRequest) (*ConsumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Consume not implemented")
}

func RegisterConsumingServiceServer(s *grpc.Server, srv ConsumingServiceServer) {
	s.RegisterService(&_ConsumingService_serviceDesc, srv)
}

func _ConsumingService_Consume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConsumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsumingServiceServer).Consume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/consuming.ConsumingService/Consume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsumingServiceServer).Consume(ctx, req.(*ConsumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConsumingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "consuming.ConsumingService",
	HandlerType: (*ConsumingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Consume",
			Handler:    _ConsumingService_Consume_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/consuming/consuming.proto",
}