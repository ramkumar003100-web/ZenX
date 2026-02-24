package events

import "context"

type NATSConn interface {
	Publish(subject string, data []byte) error
	RequestWithContext(context.Context, string, []byte) ([]byte, error)
}
