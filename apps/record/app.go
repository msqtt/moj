package main

import (
	"fmt"
	"log"
	"log/slog"
	"moj/record/db"
	"moj/record/etc"
	"moj/record/mq/consumer"
	red_pb "moj/record/rpc"
	"moj/record/schedule"
	"moj/record/svc"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type App struct {
	server                  *svc.Server
	mongodb                 *db.MongoDB
	config                  *etc.Config
	taskers                 []*schedule.TikerTasker
	nsqFinishRecordConsumer *consumer.NsqFinishRecordConsumer
}

func NewApp(server *svc.Server, taskers []*schedule.TikerTasker,
	db *db.MongoDB, config *etc.Config,
	nsqFinishRecordConsumer *consumer.NsqFinishRecordConsumer,
) *App {
	return &App{
		server:                  server,
		mongodb:                 db,
		config:                  config,
		taskers:                 taskers,
		nsqFinishRecordConsumer: nsqFinishRecordConsumer,
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
	red_pb.RegisterRecordServiceServer(grpcServer, a.server)

	// step3: launch schedule task
	for _, task := range a.taskers {
		cancel := task.Launch()
		defer cancel()
	}

	// step4: launch consumer listen
	a.nsqFinishRecordConsumer.RegisterListener()
	a.nsqFinishRecordConsumer.Start()

	// step4: launch grpc server
	log.Println("starting server at", addr)
	log.Fatalln(grpcServer.Serve(lis))
}

func (a *App) Stop() {
	a.nsqFinishRecordConsumer.Close()
	a.mongodb.Close()
}
