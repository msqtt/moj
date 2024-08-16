package main

import (
	"log/slog"
	"moj/apps/web-bff/etc"
	"moj/apps/web-bff/graph"
	"moj/apps/web-bff/rpc"
	"os"
	"time"
)

type App struct {
	config     *etc.Config
	resolver   *graph.Resolver
	rpcClients *rpc.RpcClients
}

func NewApp(
	config *etc.Config,
	resolver *graph.Resolver,
	rpcClients *rpc.RpcClients,
) *App {
	return &App{
		config:     config,
		resolver:   resolver,
		rpcClients: rpcClients,
	}
}

// setup default logger
func setupLogger(debug bool) {
	opts := &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(time.Now().Format(time.DateTime))
		}
		return a
	}}
	if debug {
		opts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
}

// start rpc server
func (a *App) Start() {
	setupLogger(a.config.Debug)
	startHttpServer(a.resolver)
}

func (a *App) Stop() {
	for _, c := range a.rpcClients.Connects {
		c.Close()
	}
}
