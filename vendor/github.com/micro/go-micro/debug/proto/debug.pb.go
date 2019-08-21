// Code generated by protoc-gen-go. DO NOT EDIT.
// source: micro/go-micro/debug/proto/debug.proto

package debug

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type HealthRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthRequest) Reset()         { *m = HealthRequest{} }
func (m *HealthRequest) String() string { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()    {}
func (*HealthRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f25415e61bccfa1f, []int{0}
}

func (m *HealthRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthRequest.Unmarshal(m, b)
}
func (m *HealthRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthRequest.Marshal(b, m, deterministic)
}
func (m *HealthRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthRequest.Merge(m, src)
}
func (m *HealthRequest) XXX_Size() int {
	return xxx_messageInfo_HealthRequest.Size(m)
}
func (m *HealthRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HealthRequest proto.InternalMessageInfo

type HealthResponse struct {
	// default: ok
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthResponse) Reset()         { *m = HealthResponse{} }
func (m *HealthResponse) String() string { return proto.CompactTextString(m) }
func (*HealthResponse) ProtoMessage()    {}
func (*HealthResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f25415e61bccfa1f, []int{1}
}

func (m *HealthResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthResponse.Unmarshal(m, b)
}
func (m *HealthResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthResponse.Marshal(b, m, deterministic)
}
func (m *HealthResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthResponse.Merge(m, src)
}
func (m *HealthResponse) XXX_Size() int {
	return xxx_messageInfo_HealthResponse.Size(m)
}
func (m *HealthResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HealthResponse proto.InternalMessageInfo

func (m *HealthResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type StatsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatsRequest) Reset()         { *m = StatsRequest{} }
func (m *StatsRequest) String() string { return proto.CompactTextString(m) }
func (*StatsRequest) ProtoMessage()    {}
func (*StatsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f25415e61bccfa1f, []int{2}
}

func (m *StatsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatsRequest.Unmarshal(m, b)
}
func (m *StatsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatsRequest.Marshal(b, m, deterministic)
}
func (m *StatsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatsRequest.Merge(m, src)
}
func (m *StatsRequest) XXX_Size() int {
	return xxx_messageInfo_StatsRequest.Size(m)
}
func (m *StatsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatsRequest proto.InternalMessageInfo

type StatsResponse struct {
	// unix timestamp
	Started uint64 `protobuf:"varint,1,opt,name=started,proto3" json:"started,omitempty"`
	// in seconds
	Uptime uint64 `protobuf:"varint,2,opt,name=uptime,proto3" json:"uptime,omitempty"`
	// in bytes
	Memory uint64 `protobuf:"varint,3,opt,name=memory,proto3" json:"memory,omitempty"`
	// num threads
	Threads uint64 `protobuf:"varint,4,opt,name=threads,proto3" json:"threads,omitempty"`
	// total gc in nanoseconds
	Gc                   uint64   `protobuf:"varint,5,opt,name=gc,proto3" json:"gc,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatsResponse) Reset()         { *m = StatsResponse{} }
func (m *StatsResponse) String() string { return proto.CompactTextString(m) }
func (*StatsResponse) ProtoMessage()    {}
func (*StatsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f25415e61bccfa1f, []int{3}
}

func (m *StatsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatsResponse.Unmarshal(m, b)
}
func (m *StatsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatsResponse.Marshal(b, m, deterministic)
}
func (m *StatsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatsResponse.Merge(m, src)
}
func (m *StatsResponse) XXX_Size() int {
	return xxx_messageInfo_StatsResponse.Size(m)
}
func (m *StatsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatsResponse proto.InternalMessageInfo

func (m *StatsResponse) GetStarted() uint64 {
	if m != nil {
		return m.Started
	}
	return 0
}

func (m *StatsResponse) GetUptime() uint64 {
	if m != nil {
		return m.Uptime
	}
	return 0
}

func (m *StatsResponse) GetMemory() uint64 {
	if m != nil {
		return m.Memory
	}
	return 0
}

func (m *StatsResponse) GetThreads() uint64 {
	if m != nil {
		return m.Threads
	}
	return 0
}

func (m *StatsResponse) GetGc() uint64 {
	if m != nil {
		return m.Gc
	}
	return 0
}

func init() {
	proto.RegisterType((*HealthRequest)(nil), "HealthRequest")
	proto.RegisterType((*HealthResponse)(nil), "HealthResponse")
	proto.RegisterType((*StatsRequest)(nil), "StatsRequest")
	proto.RegisterType((*StatsResponse)(nil), "StatsResponse")
}

func init() {
	proto.RegisterFile("micro/go-micro/debug/proto/debug.proto", fileDescriptor_f25415e61bccfa1f)
}

var fileDescriptor_f25415e61bccfa1f = []byte{
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0xb7, 0x75, 0x5b, 0x71, 0xb0, 0x59, 0xc8, 0x41, 0xc2, 0x9e, 0x24, 0x07, 0x29, 0x88,
	0x59, 0xd0, 0xbf, 0xe0, 0xc1, 0x73, 0xbd, 0x0b, 0xd9, 0x76, 0xe8, 0x16, 0xac, 0xa9, 0x99, 0xe9,
	0xc1, 0xb3, 0x7f, 0x5c, 0x9a, 0xa4, 0x60, 0x6f, 0xef, 0xbd, 0xf0, 0x1e, 0xf9, 0x06, 0x1e, 0xc6,
	0xa1, 0xf5, 0xee, 0xd4, 0xbb, 0xa7, 0x28, 0x3a, 0x3c, 0xcf, 0xfd, 0x69, 0xf2, 0x8e, 0x93, 0x36,
	0x41, 0xeb, 0x03, 0x54, 0x6f, 0x68, 0x3f, 0xf9, 0xd2, 0xe0, 0xf7, 0x8c, 0xc4, 0xba, 0x06, 0xb1,
	0x06, 0x34, 0xb9, 0x2f, 0x42, 0x79, 0x07, 0x25, 0xb1, 0xe5, 0x99, 0x54, 0x76, 0x9f, 0xd5, 0x37,
	0x4d, 0x72, 0x5a, 0xc0, 0xed, 0x3b, 0x5b, 0xa6, 0xb5, 0xf9, 0x9b, 0x41, 0x95, 0x82, 0xd4, 0x54,
	0x70, 0x4d, 0x6c, 0x3d, 0x63, 0x17, 0xaa, 0xfb, 0x66, 0xb5, 0xcb, 0xe6, 0x3c, 0xf1, 0x30, 0xa2,
	0xca, 0xc3, 0x43, 0x72, 0x4b, 0x3e, 0xe2, 0xe8, 0xfc, 0x8f, 0xba, 0x8a, 0x79, 0x74, 0xcb, 0x12,
	0x5f, 0x3c, 0xda, 0x8e, 0xd4, 0x3e, 0x2e, 0x25, 0x2b, 0x05, 0xe4, 0x7d, 0xab, 0x8a, 0x10, 0xe6,
	0x7d, 0xfb, 0xfc, 0x01, 0xc5, 0xeb, 0xc2, 0x27, 0x1f, 0xa1, 0x8c, 0x20, 0x52, 0x98, 0x0d, 0xe2,
	0xf1, 0x60, 0xb6, 0x84, 0x7a, 0x27, 0x6b, 0x28, 0xc2, 0xd7, 0x65, 0x65, 0xfe, 0x33, 0x1d, 0x85,
	0xd9, 0x10, 0xe9, 0xdd, 0xb9, 0x0c, 0x77, 0x7b, 0xf9, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xfe, 0xb9,
	0x5f, 0xf7, 0x61, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DebugClient is the client API for Debug service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DebugClient interface {
	Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error)
	Stats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
}

type debugClient struct {
	cc *grpc.ClientConn
}

func NewDebugClient(cc *grpc.ClientConn) DebugClient {
	return &debugClient{cc}
}

func (c *debugClient) Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error) {
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, "/Debug/Health", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debugClient) Stats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, "/Debug/Stats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DebugServer is the server API for Debug service.
type DebugServer interface {
	Health(context.Context, *HealthRequest) (*HealthResponse, error)
	Stats(context.Context, *StatsRequest) (*StatsResponse, error)
}

func RegisterDebugServer(s *grpc.Server, srv DebugServer) {
	s.RegisterService(&_Debug_serviceDesc, srv)
}

func _Debug_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Debug/Health",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugServer).Health(ctx, req.(*HealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Debug_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Debug/Stats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugServer).Stats(ctx, req.(*StatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Debug_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Debug",
	HandlerType: (*DebugServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Health",
			Handler:    _Debug_Health_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _Debug_Stats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "micro/go-micro/debug/proto/debug.proto",
}
