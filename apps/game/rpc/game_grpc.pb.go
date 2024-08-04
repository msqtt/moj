// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: rpc/game.proto

package game_pb

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

// GameServiceClient is the client API for GameService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameServiceClient interface {
	CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameResponse, error)
	GetGameInfo(ctx context.Context, in *GetGameInfoRequest, opts ...grpc.CallOption) (*GetGameInfoResponse, error)
	GetGamePage(ctx context.Context, in *GetGamePageRequest, opts ...grpc.CallOption) (*GetGamePageResponse, error)
	GetScore(ctx context.Context, in *GetScoreRequest, opts ...grpc.CallOption) (*GetScoreResponse, error)
	GetScorePage(ctx context.Context, in *GetScorePageRequest, opts ...grpc.CallOption) (*GetScorePageResponse, error)
	UpdateGame(ctx context.Context, in *UpdateGameRequest, opts ...grpc.CallOption) (*UpdateGameResponse, error)
	SignUpGame(ctx context.Context, in *SignUpGameRequest, opts ...grpc.CallOption) (*SignUpGameResponse, error)
	CancelSignUpGame(ctx context.Context, in *CancelSignUpGameRequest, opts ...grpc.CallOption) (*CancelSignUpGameResponse, error)
}

type gameServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGameServiceClient(cc grpc.ClientConnInterface) GameServiceClient {
	return &gameServiceClient{cc}
}

func (c *gameServiceClient) CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameResponse, error) {
	out := new(CreateGameResponse)
	err := c.cc.Invoke(ctx, "/game.GameService/CreateGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) GetGameInfo(ctx context.Context, in *GetGameInfoRequest, opts ...grpc.CallOption) (*GetGameInfoResponse, error) {
	out := new(GetGameInfoResponse)
	err := c.cc.Invoke(ctx, "/game.GameService/GetGameInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) GetGamePage(ctx context.Context, in *GetGamePageRequest, opts ...grpc.CallOption) (*GetGamePageResponse, error) {
	out := new(GetGamePageResponse)
	err := c.cc.Invoke(ctx, "/game.GameService/GetGamePage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) GetScore(ctx context.Context, in *GetScoreRequest, opts ...grpc.CallOption) (*GetScoreResponse, error) {
	out := new(GetScoreResponse)
	err := c.cc.Invoke(ctx, "/game.GameService/GetScore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) GetScorePage(ctx context.Context, in *GetScorePageRequest, opts ...grpc.CallOption) (*GetScorePageResponse, error) {
	out := new(GetScorePageResponse)
	err := c.cc.Invoke(ctx, "/game.GameService/GetScorePage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) UpdateGame(ctx context.Context, in *UpdateGameRequest, opts ...grpc.CallOption) (*UpdateGameResponse, error) {
	out := new(UpdateGameResponse)
	err := c.cc.Invoke(ctx, "/game.GameService/UpdateGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) SignUpGame(ctx context.Context, in *SignUpGameRequest, opts ...grpc.CallOption) (*SignUpGameResponse, error) {
	out := new(SignUpGameResponse)
	err := c.cc.Invoke(ctx, "/game.GameService/SignUpGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) CancelSignUpGame(ctx context.Context, in *CancelSignUpGameRequest, opts ...grpc.CallOption) (*CancelSignUpGameResponse, error) {
	out := new(CancelSignUpGameResponse)
	err := c.cc.Invoke(ctx, "/game.GameService/CancelSignUpGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GameServiceServer is the server API for GameService service.
// All implementations must embed UnimplementedGameServiceServer
// for forward compatibility
type GameServiceServer interface {
	CreateGame(context.Context, *CreateGameRequest) (*CreateGameResponse, error)
	GetGameInfo(context.Context, *GetGameInfoRequest) (*GetGameInfoResponse, error)
	GetGamePage(context.Context, *GetGamePageRequest) (*GetGamePageResponse, error)
	GetScore(context.Context, *GetScoreRequest) (*GetScoreResponse, error)
	GetScorePage(context.Context, *GetScorePageRequest) (*GetScorePageResponse, error)
	UpdateGame(context.Context, *UpdateGameRequest) (*UpdateGameResponse, error)
	SignUpGame(context.Context, *SignUpGameRequest) (*SignUpGameResponse, error)
	CancelSignUpGame(context.Context, *CancelSignUpGameRequest) (*CancelSignUpGameResponse, error)
	mustEmbedUnimplementedGameServiceServer()
}

// UnimplementedGameServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGameServiceServer struct {
}

func (UnimplementedGameServiceServer) CreateGame(context.Context, *CreateGameRequest) (*CreateGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGame not implemented")
}
func (UnimplementedGameServiceServer) GetGameInfo(context.Context, *GetGameInfoRequest) (*GetGameInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGameInfo not implemented")
}
func (UnimplementedGameServiceServer) GetGamePage(context.Context, *GetGamePageRequest) (*GetGamePageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGamePage not implemented")
}
func (UnimplementedGameServiceServer) GetScore(context.Context, *GetScoreRequest) (*GetScoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScore not implemented")
}
func (UnimplementedGameServiceServer) GetScorePage(context.Context, *GetScorePageRequest) (*GetScorePageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScorePage not implemented")
}
func (UnimplementedGameServiceServer) UpdateGame(context.Context, *UpdateGameRequest) (*UpdateGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGame not implemented")
}
func (UnimplementedGameServiceServer) SignUpGame(context.Context, *SignUpGameRequest) (*SignUpGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUpGame not implemented")
}
func (UnimplementedGameServiceServer) CancelSignUpGame(context.Context, *CancelSignUpGameRequest) (*CancelSignUpGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelSignUpGame not implemented")
}
func (UnimplementedGameServiceServer) mustEmbedUnimplementedGameServiceServer() {}

// UnsafeGameServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameServiceServer will
// result in compilation errors.
type UnsafeGameServiceServer interface {
	mustEmbedUnimplementedGameServiceServer()
}

func RegisterGameServiceServer(s grpc.ServiceRegistrar, srv GameServiceServer) {
	s.RegisterService(&GameService_ServiceDesc, srv)
}

func _GameService_CreateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).CreateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GameService/CreateGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).CreateGame(ctx, req.(*CreateGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_GetGameInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGameInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).GetGameInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GameService/GetGameInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).GetGameInfo(ctx, req.(*GetGameInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_GetGamePage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGamePageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).GetGamePage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GameService/GetGamePage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).GetGamePage(ctx, req.(*GetGamePageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_GetScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).GetScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GameService/GetScore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).GetScore(ctx, req.(*GetScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_GetScorePage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScorePageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).GetScorePage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GameService/GetScorePage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).GetScorePage(ctx, req.(*GetScorePageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_UpdateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).UpdateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GameService/UpdateGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).UpdateGame(ctx, req.(*UpdateGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_SignUpGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).SignUpGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GameService/SignUpGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).SignUpGame(ctx, req.(*SignUpGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_CancelSignUpGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelSignUpGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).CancelSignUpGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GameService/CancelSignUpGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).CancelSignUpGame(ctx, req.(*CancelSignUpGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GameService_ServiceDesc is the grpc.ServiceDesc for GameService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GameService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "game.GameService",
	HandlerType: (*GameServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGame",
			Handler:    _GameService_CreateGame_Handler,
		},
		{
			MethodName: "GetGameInfo",
			Handler:    _GameService_GetGameInfo_Handler,
		},
		{
			MethodName: "GetGamePage",
			Handler:    _GameService_GetGamePage_Handler,
		},
		{
			MethodName: "GetScore",
			Handler:    _GameService_GetScore_Handler,
		},
		{
			MethodName: "GetScorePage",
			Handler:    _GameService_GetScorePage_Handler,
		},
		{
			MethodName: "UpdateGame",
			Handler:    _GameService_UpdateGame_Handler,
		},
		{
			MethodName: "SignUpGame",
			Handler:    _GameService_SignUpGame_Handler,
		},
		{
			MethodName: "CancelSignUpGame",
			Handler:    _GameService_CancelSignUpGame_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/game.proto",
}
