package deployment

import "context"

type Hook func(context.Context) error

type RollingPlanner struct{ Before, After Hook }

func (p RollingPlanner) Execute(ctx context.Context) error {
	if p.Before != nil {
		if err := p.Before(ctx); err != nil {
			return err
		}
	}
	if p.After != nil {
		return p.After(ctx)
	}
	return nil
}

func CheckSchemaCompatibility(current, target int) bool { return target-current <= 1 }
