package graph

import (
	"moj/apps/web-bff/etc"
	"moj/apps/web-bff/rpc"
	"moj/apps/web-bff/token"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Conf           *etc.Config
	RpcClients     *rpc.RpcClients
	sessionManager *token.SessionManager
}

func NewResolver(
	conf *etc.Config,
	rpcClients *rpc.RpcClients,
	sessionManager *token.SessionManager,
) *Resolver {
	return &Resolver{
		Conf:           conf,
		RpcClients:     rpcClients,
		sessionManager: sessionManager,
	}
}
