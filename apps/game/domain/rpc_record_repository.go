package domain

import (
	"context"
	"errors"
	"log"
	"moj/domain/record"
	"moj/game/etc"
	"moj/game/pkg/app_err"
	red_pb "moj/record/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type RPCRecordRepository struct {
	conf   *etc.Config
	conn   *grpc.ClientConn
	client red_pb.RecordServiceClient
}

// FindBestRecord implements record.RecordRepository.
func (r *RPCRecordRepository) FindBestRecord(ctx context.Context, uid string, qid string, gid string) (*record.Record, error) {
	var gameID *string
	if gid != "" {
		gameID = &gid
	}
	req := &red_pb.GetBestRecordRequest{
		QuestionID: qid,
		AccountID:  uid,
		GameID:     gameID,
	}
	resp, err := r.client.GetBestRecord(ctx, req)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, err)
		return nil, err
	}
	return &record.Record{
		RecordID:         resp.Record.RecordID,
		AccountID:        resp.Record.AccountID,
		GameID:           resp.Record.GetGameID(),
		QuestionID:       resp.Record.QuestionID,
		Language:         resp.Record.Language,
		Code:             resp.Record.Code,
		CodeHash:         resp.Record.CodeHash,
		JudgeStatus:      resp.Record.JudgeStatus,
		FailedReason:     resp.Record.GetFailedReason(),
		NumberFinishedAt: int(resp.Record.NumberFinishedAt),
		TotalQuestion:    int(resp.Record.TotalQuestion),
		CreateTime:       resp.Record.CreateTime,
		FinishTime:       resp.Record.FinishTime,
		MemoryUsed:       int(resp.Record.MemoryUsed),
		TimeUsed:         int(resp.Record.TimeUsed),
		CPUTimeUsed:      int(resp.Record.CpuTimeUsed),
	}, err
}

// FindRecordByID implements record.RecordRepository.
func (r *RPCRecordRepository) FindRecordByID(ctx context.Context, recordID string) (*record.Record, error) {
	panic("unimplemented")
}

// Save implements record.RecordRepository.
func (r *RPCRecordRepository) Save(context.Context, *record.Record) (string, error) {
	panic("unimplemented")
}

func NewRPCRecordRepository(conf *etc.Config) record.RecordRepository {
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

	conn, err := grpc.NewClient(conf.RecordRPCAddr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := red_pb.NewRecordServiceClient(conn)
	return &RPCRecordRepository{
		conf:   conf,
		conn:   conn,
		client: client,
	}
}
