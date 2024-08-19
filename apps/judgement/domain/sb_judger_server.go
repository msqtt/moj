package domain

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"moj/game/pkg/app_err"
	"moj/judgement/etc"
	sb_pb "moj/judgement/rpc/sb-judger"
	"moj/domain/judgement"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type SbJudger struct {
	conf   *etc.Config
	conn   *grpc.ClientConn
	client sb_pb.CodeClient
}

func toLang(lang string) sb_pb.Language {
	switch lang {
	case "c":
		return sb_pb.Language_c
	case "cpp":
		return sb_pb.Language_cpp
	case "java":
		return sb_pb.Language_java
	case "python":
		return sb_pb.Language_python
	case "go":
		return sb_pb.Language_golang
	case "rust":
		return sb_pb.Language_rust
	default:
		return sb_pb.Language_c
	}
}

// execute implements judgement.ExecutionService.
func (s *SbJudger) Execute(cmd judgement.ExecutionCmd) (res judgement.ExecuteResult, err error) {
	req := &sb_pb.JudgeCodeRequest{
		Lang:        toLang(cmd.Language),
		Code:        cmd.Code,
		Time:        uint32(cmd.TimeLimit),
		Memory:      uint32(cmd.MemoryLimit),
		OutMsgLimit: uint32(s.conf.OutPutMsgLimit),
		Case:        []*sb_pb.Case{},
	}
	resp, err := s.client.JudgeCode(context.TODO(), req)
	if err != nil {
		slog.Error("failed to call JudgeCode", "error", err)
		err = errors.Join(app_err.ErrServerInternal, err)
	}

	maxCpuTimeUse := 0
	for _, cr := range resp.CodeResults {
		if int(cr.CpuTimeUsage) > maxCpuTimeUse {
			maxCpuTimeUse = int(cr.CpuTimeUsage)
		}
	}
	res = judgement.ExecuteResult{
		JudgeStatus:      judgement.JudgeStatusType(resp.State.String()),
		NumberFinishedAt: len(resp.CodeResults),
		MemoryUsed:       int(resp.MaxMemoryUsage),
		TimeUsed:         int(resp.MaxTimeUsage),
		CPUTimeUsed:      int(maxCpuTimeUse),
		FailedReason:     resp.OutPut,
	}
	return res, err
}

func NewSbJudger(conf *etc.Config) judgement.ExecutionService {
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

	conn, err := grpc.NewClient(conf.SbJudgerRPCAddr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := sb_pb.NewCodeClient(conn)
	return &SbJudger{
		conf:   conf,
		conn:   conn,
		client: client,
	}
}

func (s *SbJudger) Close() {
	s.conn.Close()
}
