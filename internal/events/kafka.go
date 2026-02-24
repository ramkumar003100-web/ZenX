package events

import (
	"context"
	"encoding/json"
	"errors"
)

type KafkaProducer interface {
	Send(context.Context, string, []byte) error
}
type KafkaConsumer interface {
	Receive(context.Context, string) ([]byte, error)
}

type KafkaAdapter struct {
	Producer KafkaProducer
	Consumer KafkaConsumer
}

func (k KafkaAdapter) Publish(ctx context.Context, topic string, e Event) error {
	if k.Producer == nil {
		return errors.New("producer nil")
	}
	b, _ := json.Marshal(e)
	return k.Producer.Send(ctx, topic, b)
}
func (k KafkaAdapter) Consume(ctx context.Context, topic string) (Event, error) {
	if k.Consumer == nil {
		return Event{}, errors.New("consumer nil")
	}
	b, err := k.Consumer.Receive(ctx, topic)
	if err != nil {
		return Event{}, err
	}
	var e Event
	return e, json.Unmarshal(b, &e)
}
