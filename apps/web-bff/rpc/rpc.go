package rpc

import (
	"log"
	game_pb "moj/game/rpc"
	ques_pb "moj/question/rpc"
	red_pb "moj/record/rpc"
	user_pb "moj/user/rpc"
	"moj/web-bff/etc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type RpcClients struct {
	UserClient     user_pb.UserServiceClient
	CaptchaClient  user_pb.CaptchaServiceClient
	QuestionClient ques_pb.QuestionServiceClient
	GameClient     game_pb.GameServiceClient
	RecordClient   red_pb.RecordServiceClient
	Connects       []*grpc.ClientConn
}

func NewConn(conf *etc.Config, addr string) *grpc.ClientConn {
	var creds credentials.TransportCredentials
	if conf.TLS {
		creds1, err := credentials.NewClientTLSFromFile(conf.CertFile, conf.KeyFile)
		if err != nil {
			log.Fatalln("failed to load credentials", "error", err)
		}
		creds = creds1
	} else {
		creds = insecure.NewCredentials()
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func NewUserAndCaptchaClient(conf *etc.Config) (user_pb.UserServiceClient, user_pb.CaptchaServiceClient, *grpc.ClientConn) {
	conn := NewConn(conf, conf.UserRPCAddr)
	return user_pb.NewUserServiceClient(conn), user_pb.NewCaptchaServiceClient(conn), conn
}

func NewQuestionClient(conf *etc.Config) (ques_pb.QuestionServiceClient, *grpc.ClientConn) {
	conn := NewConn(conf, conf.QuestionRPCAddr)
	return ques_pb.NewQuestionServiceClient(conn), conn
}

func NewGameClient(conf *etc.Config) (game_pb.GameServiceClient, *grpc.ClientConn) {
	conn := NewConn(conf, conf.GameRPCAddr)
	return game_pb.NewGameServiceClient(conn), conn
}

func NewRecordClient(conf *etc.Config) (red_pb.RecordServiceClient, *grpc.ClientConn) {
	conn := NewConn(conf, conf.RecordRPCAddr)
	return red_pb.NewRecordServiceClient(conn), conn
}
