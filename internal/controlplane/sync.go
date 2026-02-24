package controlplane

import (
	"context"
	"encoding/json"
)

type PubSub interface {
	Publish(context.Context, string, []byte) error
	Subscribe(context.Context, string) (<-chan []byte, error)
}

type Syncer struct {
	Channel string
	Bus     PubSub
}

func (s Syncer) Broadcast(ctx context.Context, v ConfigVersion) error {
	b, _ := json.Marshal(v)
	return s.Bus.Publish(ctx, s.Channel, b)
}

func (s Syncer) Consume(ctx context.Context, onChange func(ConfigVersion)) error {
	ch, err := s.Bus.Subscribe(ctx, s.Channel)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case b, ok := <-ch:
				if !ok {
					return
				}
				var v ConfigVersion
				if json.Unmarshal(b, &v) == nil {
					onChange(v)
				}
			}
		}
	}()
	return nil
}
