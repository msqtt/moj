package main

import (
	"fmt"
	"log"
	"log/slog"
	"moj/user/db"
	"moj/user/etc"
	user_pb "moj/user/rpc"
	service "moj/user/svc"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type App struct {
	server  *service.Server
	mongodb *db.MongoDB
	config  *etc.Config
}

func NewApp(server *service.Server, db *db.MongoDB, config *etc.Config) *App {
	return &App{
		server:  server,
		mongodb: db,
		config:  config,
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

	addr := fmt.Sprintf("0.0.0.0:%d", a.config.AppPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("failed to listen", "error", err)
	}
	var opts []grpc.ServerOption
	if a.config.TLS {
		creds, err := credentials.NewServerTLSFromFile(a.config.CertFile, a.config.KeyFile)
		if err != nil {
			log.Fatalln("failed to load credentials", "error", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	user_pb.RegisterUserServiceServer(grpcServer, a.server)
	user_pb.RegisterCaptchaServiceServer(grpcServer, a.server)

	log.Println("starting server at", addr)
	log.Fatalln(grpcServer.Serve(lis))
}

func (a *App) Stop() {
	a.mongodb.Close()
}
