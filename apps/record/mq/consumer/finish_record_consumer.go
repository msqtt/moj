package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"moj/apps/record/domain"
	"moj/apps/record/etc"
	"moj/domain/judgement"
	"moj/domain/policy"
	"moj/domain/record"

	"github.com/nsqio/go-nsq"
)

type NsqFinishRecordConsumer struct {
	consumer   *nsq.Consumer
	conf       *etc.Config
	handler    *record.ModifyRecordCmdHandler
	dispatcher domain.EventDispatcher
}

func (c *NsqFinishRecordConsumer) Close() {
	c.consumer.Stop()
}

func (n *NsqFinishRecordConsumer) Start() {
	fmt.Println(n.conf.NsqLookUpAddr)
	err := n.consumer.ConnectToNSQLookupd(n.conf.NsqLookUpAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func (n *NsqFinishRecordConsumer) RegisterListener() {
	fn := func(msg *nsq.Message) error {
		if len(msg.Body) == 0 {
			return nil
		}
		evt := judgement.ExecutionFinishEvent{}
		err := json.Unmarshal(msg.Body, &evt)
		if err != nil {
			return err
		}
		queue := domain.NewSimpleEventQueue()
		err = policy.NewModifyRecordAfterExecutionPolicy(n.handler, queue).OnEvent(evt)
		if err != nil {
			return err
		}
		n.dispatcher.Dispatch(queue)
		return nil
	}
	n.consumer.AddHandler(nsq.HandlerFunc(fn))
}

func NewNsqFinishRecordConsumer(
	conf *etc.Config,
	handler *record.ModifyRecordCmdHandler,
	dispatcher domain.EventDispatcher,

) *NsqFinishRecordConsumer {
	consumer, err := nsq.NewConsumer(conf.FinishRecordTopic,
		conf.FinishRecordTopicChannel, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	return &NsqFinishRecordConsumer{
		consumer:   consumer,
		conf:       conf,
		handler:    handler,
		dispatcher: dispatcher,
	}
}
