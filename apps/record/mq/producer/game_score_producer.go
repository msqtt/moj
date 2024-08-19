package producer

import (
	"encoding/json"
	"errors"
	"log"
	"moj/record/etc"

	"github.com/nsqio/go-nsq"
)

var ErrFailedSendGameScoreMessage = errors.New("failed to send game score message")

type Producer interface {
	Send(message any) error
	Close()
}

type RecordGameScoreProducer struct {
	producer *nsq.Producer
	conf     *etc.Config
}

// Close implements Producer.
func (r *RecordGameScoreProducer) Close() {
	r.producer.Stop()
}

// Send implements Producer.
func (r *RecordGameScoreProducer) Send(message any) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return errors.Join(ErrFailedSendGameScoreMessage, err)
	}
	err = r.producer.Publish(r.conf.CalculateScoreTopic, msg)
	if err != nil {
		return errors.Join(ErrFailedSendGameScoreMessage, err)
	}
	return err
}

func NewRecordGameScoreProducer(conf *etc.Config) Producer {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(conf.NsqdAddr, config)
	if err != nil {
		log.Fatal(err)
	}
	return &RecordGameScoreProducer{
		producer: producer,
		conf:     conf,
	}
}
