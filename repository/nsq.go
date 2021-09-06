package repository

import (
	"github.com/nsqio/go-nsq"
)

var mq *nsq.Producer

func NewNsqConn() (*nsq.Producer, error) {
	config := nsq.NewConfig()
	mq, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		return nil, err
	}
	return mq, nil
}
