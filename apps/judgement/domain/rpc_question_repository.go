package domain

import (
	"context"
	"errors"
	"log"
	"moj/apps/judgement/etc"
	"moj/apps/judgement/pkg/app_err"
	ques_pb "moj/apps/judgement/rpc/question"
	"moj/domain/question"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type RPCQuestionRepository struct {
	conf   *etc.Config
	conn   *grpc.ClientConn
	client ques_pb.QuestionServiceClient
}

// FindQuestionByID implements question.QuestionRepository.
func (r *RPCQuestionRepository) FindQuestionByID(questionID string) (*question.Question, error) {
	resp, err := r.client.GetQuestionInfo(context.TODO(), &ques_pb.GetQuestionInfoRequest{
		QuestionID: questionID,
	})

	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, err)
		return nil, err
	}

	langs := make([]question.QuestionLanguage, len(resp.QuestionInfo.AllowedLanguages))
	for id, l := range resp.QuestionInfo.AllowedLanguages {
		langs[id] = question.QuestionLanguage(l)
	}

	cases := make([]question.Case, len(resp.QuestionInfo.Cases))

	for id, ca := range resp.QuestionInfo.Cases {
		cases[id] = question.Case{
			Number:         int(ca.Number),
			InputFilePath:  ca.InputFilePath,
			OutputFilePath: ca.OutFilePath,
		}
	}

	ret := &question.Question{
		QuestionID:       resp.QuestionInfo.QuestionID,
		AccountID:        resp.QuestionInfo.AccountID,
		Enabled:          resp.QuestionInfo.Enabled,
		Title:            resp.QuestionInfo.Title,
		Content:          resp.QuestionInfo.Content,
		Level:            question.QuestionLevel(resp.QuestionInfo.Level),
		AllowedLanguages: langs,
		Cases:            cases,
		TimeLimit:        int(resp.QuestionInfo.TimeLimit),
		MemoryLimit:      int(resp.QuestionInfo.MemoryLimit),
		Tags:             resp.QuestionInfo.Tags,
		CreateTime:       resp.QuestionInfo.CreateTime,
		ModifyTime:       resp.QuestionInfo.ModifyTime,
	}
	return ret, err
}

// Save implements question.QuestionRepository.
func (r *RPCQuestionRepository) Save(*question.Question) (questionID string, err error) {
	panic("unimplemented")
}

func NewRPCQuestionRepository(conf *etc.Config) question.QuestionRepository {
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

	conn, err := grpc.NewClient(conf.QuestionRPCAddr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := ques_pb.NewQuestionServiceClient(conn)
	return &RPCQuestionRepository{
		conf:   conf,
		conn:   conn,
		client: client,
	}
}
