package consumer

import (
	"encoding/json"
	"log"
	"moj/apps/judgement/domain"
	"moj/apps/judgement/etc"
	"moj/domain/judgement"
	"moj/domain/policy"
	"moj/domain/question"

	"github.com/nsqio/go-nsq"
)

type NsqExecuteJudgementConsumer struct {
	consumer     *nsq.Consumer
	conf         *etc.Config
	handler      *judgement.ExecutionCmdHandler
	service      policy.CaseFileService
	questionRepo question.QuestionRepository
	dispatcher   domain.EventDispatcher
}

func (t *NsqExecuteJudgementConsumer) Close() {
	t.consumer.Stop()
}

func (t *NsqExecuteJudgementConsumer) Start() {
	err := t.consumer.ConnectToNSQLookupd(t.conf.NsqLookUpAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *NsqExecuteJudgementConsumer) RegisterListener() {
	fn := func(msg *nsq.Message) error {
		if len(msg.Body) == 0 {
			return nil
		}
		evt := judgement.ExecutionCmd{}
		err := json.Unmarshal(msg.Body, &evt)
		if err != nil {
			return err
		}
		queue := domain.NewSimpleEventQueue()
		err = policy.NewJudgeOnSubmitPolicy(t.service, t.handler, t.questionRepo, queue).OnEvent(evt)
		if err != nil {
			return err
		}
		t.dispatcher.Dispatch(queue)
		return nil
	}
	t.consumer.AddHandler(nsq.HandlerFunc(fn))
}

func NewNsqExecuteJudgementConsumer(
	conf *etc.Config,
	handler *judgement.ExecutionCmdHandler,
	service policy.CaseFileService,
	questionRepo question.QuestionRepository,
	dispatcher domain.EventDispatcher,
) *NsqExecuteJudgementConsumer {
	consumer, err := nsq.NewConsumer(conf.ExecuteJudgeTopic,
		conf.ExecuteJudgeTopicChannel, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	return &NsqExecuteJudgementConsumer{
		consumer:     consumer,
		conf:         conf,
		handler:      handler,
		service:      service,
		questionRepo: questionRepo,
		dispatcher:   dispatcher,
	}
}
