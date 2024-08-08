// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: rpc/judgement.proto

package jud_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// JudgeServiceClient is the client API for JudgeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JudgeServiceClient interface {
	ExecuteJudge(ctx context.Context, in *ExecuteJudgeRequest, opts ...grpc.CallOption) (*ExecuteJudgeResponse, error)
}

type judgeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewJudgeServiceClient(cc grpc.ClientConnInterface) JudgeServiceClient {
	return &judgeServiceClient{cc}
}

func (c *judgeServiceClient) ExecuteJudge(ctx context.Context, in *ExecuteJudgeRequest, opts ...grpc.CallOption) (*ExecuteJudgeResponse, error) {
	out := new(ExecuteJudgeResponse)
	err := c.cc.Invoke(ctx, "/judgement.JudgeService/ExecuteJudge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JudgeServiceServer is the server API for JudgeService service.
// All implementations must embed UnimplementedJudgeServiceServer
// for forward compatibility
type JudgeServiceServer interface {
	ExecuteJudge(context.Context, *ExecuteJudgeRequest) (*ExecuteJudgeResponse, error)
	mustEmbedUnimplementedJudgeServiceServer()
}

// UnimplementedJudgeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedJudgeServiceServer struct {
}

func (UnimplementedJudgeServiceServer) ExecuteJudge(context.Context, *ExecuteJudgeRequest) (*ExecuteJudgeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteJudge not implemented")
}
func (UnimplementedJudgeServiceServer) mustEmbedUnimplementedJudgeServiceServer() {}

// UnsafeJudgeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JudgeServiceServer will
// result in compilation errors.
type UnsafeJudgeServiceServer interface {
	mustEmbedUnimplementedJudgeServiceServer()
}

func RegisterJudgeServiceServer(s grpc.ServiceRegistrar, srv JudgeServiceServer) {
	s.RegisterService(&JudgeService_ServiceDesc, srv)
}

func _JudgeService_ExecuteJudge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteJudgeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JudgeServiceServer).ExecuteJudge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/judgement.JudgeService/ExecuteJudge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JudgeServiceServer).ExecuteJudge(ctx, req.(*ExecuteJudgeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// JudgeService_ServiceDesc is the grpc.ServiceDesc for JudgeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var JudgeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "judgement.JudgeService",
	HandlerType: (*JudgeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteJudge",
			Handler:    _JudgeService_ExecuteJudge_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/judgement.proto",
}
