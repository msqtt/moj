package mq

import (
	"encoding/json"
	"log"
	"moj/apps/game/domain"
	"moj/apps/game/etc"
	"moj/domain/game"
	"moj/domain/policy"
	"moj/domain/record"

	"github.com/nsqio/go-nsq"
)

type NsqCalculateScoreConsumer struct {
	consumer   *nsq.Consumer
	conf       *etc.Config
	handler    *game.CalculateScoreCmdHandler
	dispatcher domain.EventDispatcher
}

// RegisterListener implements Consumer.
func (n *NsqCalculateScoreConsumer) RegisterListener() {
	fn := func(msg *nsq.Message) error {
		if len(msg.Body) == 0 {
			return nil
		}
		evt := record.ModifyRecordEvent{}
		err := json.Unmarshal(msg.Body, &evt)
		if err != nil {
			return err
		}
		synQueue := domain.NewSimpleEventQueue()
		err = policy.NewCalculateScorePolicy(n.handler, synQueue).OnEvent(evt)
		if err != nil {
			return err
		}
		n.dispatcher.Dispatch(synQueue)
		return nil
	}
	n.consumer.AddHandler(nsq.HandlerFunc(fn))
}

func NewNsqCalculateScoreConsumer(
	conf *etc.Config,
	handler *game.CalculateScoreCmdHandler,
	dispatcher domain.EventDispatcher,
) *NsqCalculateScoreConsumer {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(conf.CalculateScoreTopic,
		conf.CalculateScoreChannel, config)
	if err != nil {
		log.Fatal(err)
	}
	return &NsqCalculateScoreConsumer{
		consumer:   consumer,
		conf:       conf,
		handler:    handler,
		dispatcher: dispatcher,
	}
}

func (n *NsqCalculateScoreConsumer) Start() {
	err := n.consumer.ConnectToNSQLookupd(n.conf.NSQLookUpAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func (n *NsqCalculateScoreConsumer) Stop() {
	n.consumer.Stop()
}
