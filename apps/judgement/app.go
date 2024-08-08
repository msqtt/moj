package main

import (
	"fmt"
	"log"
	"log/slog"
	"moj/apps/judgement/db"
	"moj/apps/judgement/etc"
	"moj/apps/judgement/mq/consumer"
	jud_pb "moj/apps/judgement/rpc"
	"moj/apps/judgement/svc"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type App struct {
	server                      *svc.Server
	mongodb                     *db.MongoDB
	config                      *etc.Config
	nsqExecuteJudgementConsumer *consumer.NsqExecuteJudgementConsumer
}

func NewApp(server *svc.Server,
	db *db.MongoDB, config *etc.Config,
	nsqExecuteJudgementConsumer *consumer.NsqExecuteJudgementConsumer,
) *App {
	return &App{
		server:                      server,
		mongodb:                     db,
		config:                      config,
		nsqExecuteJudgementConsumer: nsqExecuteJudgementConsumer,
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
	// step1: setup log
	setupLogger(a.config.Debug)

	// step2: setup grpc service
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
	jud_pb.RegisterJudgeServiceServer(grpcServer, a.server)

	// step3: launch consumer listen
	a.nsqExecuteJudgementConsumer.RegisterListener()
	a.nsqExecuteJudgementConsumer.Start()

	// step4: launch grpc server
	log.Println("starting server at", addr)
	log.Fatalln(grpcServer.Serve(lis))
}

func (a *App) Stop() {
	a.nsqExecuteJudgementConsumer.Close()
	a.mongodb.Close()
}
