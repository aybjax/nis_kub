package adapter

import (
	"fmt"
	"github.com/aybjax/nis_lib/cmntypes"
	"github.com/nats-io/nats.go"
)

func NewQueueEngine() (cmntypes.AppQueue, error) {
	// Connect to NATS
	nc, err := nats.Connect(fmt.Sprintf("%s:%s", ENV.NATS_URL, ENV.NATS_PORT))
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	return &jetStreamQueue{js}, nil
}

type jetStreamQueue struct {
	js nats.JetStreamContext
}

func (j *jetStreamQueue) Publish(data []byte, queue string, topic string) error {
	_, err := j.js.Publish(queue, data)

	return err
}

func (j *jetStreamQueue) Subscribe(queue string, topic string, cb func(data []byte) error) {
	j.js.QueueSubscribe(queue, topic, func(msg *nats.Msg) {
		if err := cb(msg.Data); err != nil {
			// TODO
		}
	})
}
