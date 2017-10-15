// Code generated by protoc-gen-go. DO NOT EDIT.
// source: redis.proto

/*
Package redis is a generated protocol buffer package.

It is generated from these files:
	redis.proto

It has these top-level messages:
	RedisRequest
	RedisReply
*/
package redis

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RedisRequest struct {
	Action string `protobuf:"bytes,1,opt,name=action" json:"action,omitempty"`
	Param  string `protobuf:"bytes,2,opt,name=param" json:"param,omitempty"`
}

func (m *RedisRequest) Reset()                    { *m = RedisRequest{} }
func (m *RedisRequest) String() string            { return proto.CompactTextString(m) }
func (*RedisRequest) ProtoMessage()               {}
func (*RedisRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RedisRequest) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *RedisRequest) GetParam() string {
	if m != nil {
		return m.Param
	}
	return ""
}

// The response message containing the greetings
type RedisReply struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *RedisReply) Reset()                    { *m = RedisReply{} }
func (m *RedisReply) String() string            { return proto.CompactTextString(m) }
func (*RedisReply) ProtoMessage()               {}
func (*RedisReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RedisReply) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func init() {
	proto.RegisterType((*RedisRequest)(nil), "redis.RedisRequest")
	proto.RegisterType((*RedisReply)(nil), "redis.RedisReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Redis service

type RedisClient interface {
	Command(ctx context.Context, in *RedisRequest, opts ...grpc.CallOption) (*RedisReply, error)
}

type redisClient struct {
	cc *grpc.ClientConn
}

func NewRedisClient(cc *grpc.ClientConn) RedisClient {
	return &redisClient{cc}
}

func (c *redisClient) Command(ctx context.Context, in *RedisRequest, opts ...grpc.CallOption) (*RedisReply, error) {
	out := new(RedisReply)
	err := grpc.Invoke(ctx, "/redis.Redis/Command", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Redis service

type RedisServer interface {
	Command(context.Context, *RedisRequest) (*RedisReply, error)
}

func RegisterRedisServer(s *grpc.Server, srv RedisServer) {
	s.RegisterService(&_Redis_serviceDesc, srv)
}

func _Redis_Command_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RedisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RedisServer).Command(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redis.Redis/Command",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RedisServer).Command(ctx, req.(*RedisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Redis_serviceDesc = grpc.ServiceDesc{
	ServiceName: "redis.Redis",
	HandlerType: (*RedisServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Command",
			Handler:    _Redis_Command_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "redis.proto",
}

func init() { proto.RegisterFile("redis.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 147 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x4a, 0x4d, 0xc9,
	0x2c, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x6c, 0xb8, 0x78, 0x82,
	0x40, 0x8c, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x31, 0x2e, 0xb6, 0xc4, 0xe4, 0x92,
	0xcc, 0xfc, 0x3c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x28, 0x4f, 0x48, 0x84, 0x8b, 0xb5,
	0x20, 0xb1, 0x28, 0x31, 0x57, 0x82, 0x09, 0x2c, 0x0c, 0xe1, 0x28, 0xa9, 0x70, 0x71, 0x41, 0x75,
	0x17, 0xe4, 0x54, 0x82, 0xf4, 0x16, 0xa5, 0x16, 0x97, 0xe6, 0x94, 0xc0, 0xf4, 0x42, 0x78, 0x46,
	0x36, 0x5c, 0xac, 0x60, 0x55, 0x42, 0xc6, 0x5c, 0xec, 0xce, 0xf9, 0xb9, 0xb9, 0x89, 0x79, 0x29,
	0x42, 0xc2, 0x7a, 0x10, 0xc7, 0x20, 0x5b, 0x2e, 0x25, 0x88, 0x2a, 0x58, 0x90, 0x53, 0xa9, 0xc4,
	0x90, 0xc4, 0x06, 0x76, 0xaf, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xfc, 0xd2, 0xdd, 0xf7, 0xbe,
	0x00, 0x00, 0x00,
}
