//go:build wireinject
// +build wireinject

package main

import (
	"moj/web-bff/etc"
	"moj/web-bff/graph"
	"moj/web-bff/handler"
	"moj/web-bff/oss"
	"moj/web-bff/rpc"
	"moj/web-bff/token"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

func provideRpcClients(conf *etc.Config) *rpc.RpcClients {
	ret := &rpc.RpcClients{}
	var conn *grpc.ClientConn

	ret.UserClient, ret.CaptchaClient, conn = rpc.NewUserAndCaptchaClient(conf)
	ret.Connects = append(ret.Connects, conn)

	ret.QuestionClient, conn = rpc.NewQuestionClient(conf)
	ret.Connects = append(ret.Connects, conn)

	ret.GameClient, conn = rpc.NewGameClient(conf)
	ret.Connects = append(ret.Connects, conn)

	ret.RecordClient, conn = rpc.NewRecordClient(conf)
	ret.Connects = append(ret.Connects, conn)

	return ret
}

var providers = wire.NewSet(
	etc.NewAppConfig,
	graph.NewResolver,
	token.NewPasetoTokener,
	token.NewSessionManager,

	handler.NewAvatarHandler,
	handler.NewCaseFileHandler,

	oss.NewMinioOssUploader,

	provideRpcClients,
)

func InitializeApplication() *App {
	wire.Build(NewApp, providers)
	return nil
}
