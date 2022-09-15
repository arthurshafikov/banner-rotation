package queue

import (
	"context"
	"encoding/json"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

type element struct {
	Key   []byte
	Value []byte
	Topic string
}

type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type Queue struct {
	ctx         context.Context
	elements    chan element
	KafkaWriter KafkaWriter
}

func NewQueue(ctx context.Context, kafkaWriter KafkaWriter) *Queue {
	return &Queue{
		ctx:         ctx,
		elements:    make(chan element, 100),
		KafkaWriter: kafkaWriter,
	}
}

func (q *Queue) Dispatch() error {
OUTER:
	for {
		fmt.Println("For...")
		select {
		case <-q.ctx.Done():
			fmt.Println("Ctx done...")
			break OUTER
		case el := <-q.elements:
			fmt.Println("new elem!")
			if err := q.writeMessageToKafka(el); err != nil {
				return err
			}
		}
	}

	return nil
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
