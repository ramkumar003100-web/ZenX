package performance

import "context"

type Loader func(context.Context, string) (string, error)
type Writer func(context.Context, string, string) error

type ReadThroughCache struct {
	Get  func(context.Context, string) (string, error)
	Set  func(context.Context, string, string) error
	Load Loader
}

func (c ReadThroughCache) Read(ctx context.Context, key string) (string, error) {
	v, err := c.Get(ctx, key)
	if err == nil {
		return v, nil
	}
	v, err = c.Load(ctx, key)
	if err != nil {
		return "", err
	}
	_ = c.Set(ctx, key, v)
	return v, nil
}

type WriteThroughCache struct {
	Set     func(context.Context, string, string) error
	Persist Writer
}

func (c WriteThroughCache) Write(ctx context.Context, key, val string) error {
	if err := c.Persist(ctx, key, val); err != nil {
		return err
	}
	return c.Set(ctx, key, val)
}
