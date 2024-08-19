package producer

import (
	"encoding/json"
	"errors"
	"log"
	"moj/record/etc"

	"github.com/nsqio/go-nsq"
)

var ErrFailedSendJudgeMessage = errors.New("failed to send judge message")

type ExecuteJudgeProducer struct {
	producer *nsq.Producer
	conf     *etc.Config
}

// Close implements Producer.
func (e *ExecuteJudgeProducer) Close() {
	e.producer.Stop()
}

// Send implements Producer.
func (e *ExecuteJudgeProducer) Send(message any) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return errors.Join(ErrFailedSendJudgeMessage, err)
	}
	err = e.producer.Publish(e.conf.ExecuteJudgeTopic, msg)
	if err != nil {
		err = errors.Join(ErrFailedSendJudgeMessage, err)
	}
	return err
}

func NewExecuteJudgeProducer(conf *etc.Config) Producer {
	producer, err := nsq.NewProducer(conf.NsqdAddr, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	return &ExecuteJudgeProducer{
		producer: producer,
		conf:     conf,
	}
}
