package dataprotection

import "context"

type Deleter interface {
	DeleteByUser(context.Context, string) error
}

func GDPRDelete(ctx context.Context, userID string, stores ...Deleter) error {
	for _, s := range stores {
		if err := s.DeleteByUser(ctx, userID); err != nil {
			return err
		}
	}
	return nil
}
