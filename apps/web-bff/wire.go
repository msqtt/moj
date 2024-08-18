//go:build wireinject
// +build wireinject

package main

import (
	"moj/apps/web-bff/etc"
	"moj/apps/web-bff/graph"
	"moj/apps/web-bff/handler"
	"moj/apps/web-bff/oss"
	"moj/apps/web-bff/rpc"
	"moj/apps/web-bff/token"

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
