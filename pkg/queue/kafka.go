package queue

import (
	"context"
	"encoding/json"
	"log"
	"os"

	kafka "github.com/segmentio/kafka-go"
)

type element struct {
	Key   []byte
	Value []byte
	Topic string
}

type Queue struct {
	ctx         context.Context
	elements    chan element
	KafkaWriter *kafka.Writer
}

func NewQueue(ctx context.Context, brokerAddress string) *Queue {
	l := log.New(os.Stdout, "kafka writer: ", 0)

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Logger:  l,
	})

	return &Queue{
		ctx:         ctx,
		elements:    make(chan element, 100),
		KafkaWriter: w,
	}
}

func (q *Queue) Dispatch() {
OUTER:
	for {
		select {
		case <-q.ctx.Done():
			break OUTER
		case el := <-q.elements:
			err := q.writeMessageToKafka(el)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (q *Queue) AddToQueue(topic string, value interface{}) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	q.elements <- element{
		Value: valueJSON,
		Topic: topic,
	}

	return nil
}

func (q *Queue) writeMessageToKafka(el element) error {
	return q.KafkaWriter.WriteMessages(q.ctx, kafka.Message{
		Key:   el.Key,
		Value: el.Value,
		Topic: el.Topic,
	})
}
