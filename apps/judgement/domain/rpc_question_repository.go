package domain

import (
	"context"
	"errors"
	"log"
	"moj/apps/judgement/etc"
	"moj/apps/judgement/pkg/app_err"
	ques_pb "moj/apps/question/rpc"
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
func (r *RPCQuestionRepository) FindQuestionByID(ctx context.Context, questionID string) (*question.Question, error) {
	resp, err := r.client.GetQuestion(ctx, &ques_pb.GetQuestionRequest{
		QuestionID: questionID,
	})

	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, err)
		return nil, err
	}

	langs := make([]question.QuestionLanguage, len(resp.Question.AllowedLanguages))
	for id, l := range resp.Question.AllowedLanguages {
		langs[id] = question.QuestionLanguage(l)
	}

	cases := make([]question.Case, len(resp.Question.Cases))

	for id, ca := range resp.Question.Cases {
		cases[id] = question.Case{
			Number:         int(ca.Number),
			InputFilePath:  ca.InputFilePath,
			OutputFilePath: ca.OutputFilePath,
		}
	}

	ret := &question.Question{
		QuestionID:       resp.Question.QuestionID,
		AccountID:        resp.Question.AccountID,
		Enabled:          resp.Question.Enabled,
		Title:            resp.Question.Title,
		Content:          resp.Question.Content,
		Level:            question.QuestionLevel(resp.Question.Level),
		AllowedLanguages: langs,
		Cases:            cases,
		TimeLimit:        int(resp.Question.TimeLimit),
		MemoryLimit:      int(resp.Question.MemoryLimit),
		Tags:             resp.Question.Tags,
		CreateTime:       resp.Question.CreateTime,
		ModifyTime:       resp.Question.ModifyTime,
	}
	return ret, err
}

// Save implements question.QuestionRepository.
func (r *RPCQuestionRepository) Save(context.Context, *question.Question) (questionID string, err error) {
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
