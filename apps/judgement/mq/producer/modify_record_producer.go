package producer

import (
	"encoding/json"
	"errors"
	"log"
	"moj/judgement/etc"

	"github.com/nsqio/go-nsq"
)

var ErrFailedSendModifyRecordMessage = errors.New("failed to send modify record message")

type NsqModifyRecordProducer struct {
	producer *nsq.Producer
	conf     *etc.Config
}

// Close implements Producer.
func (n *NsqModifyRecordProducer) Close() {
	n.producer.Stop()
}

// Send implements Producer.
func (n *NsqModifyRecordProducer) Send(message any) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return errors.Join(ErrFailedSendModifyRecordMessage, err)
	}
	err = n.producer.Publish(n.conf.FinishJudgementTopic, msg)
	if err != nil {
		err = errors.Join(ErrFailedSendModifyRecordMessage, err)
	}
	return err
}

func NewNsqModifyRecordProducer(conf *etc.Config) Producer {
	producer, err := nsq.NewProducer(conf.NsqdAddr, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	return &NsqModifyRecordProducer{
		producer: producer,
		conf:     conf,
	}
}
