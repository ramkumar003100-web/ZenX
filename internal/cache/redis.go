package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"zenx/pkg/metrics"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db, PoolSize: 20})
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	observe("set", err)
	return err
}
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	v, err := r.client.Get(ctx, key).Result()
	observe("get", err)
	return v, err
}
func (r *RedisCache) HSet(ctx context.Context, key string, values ...any) error {
	err := r.client.HSet(ctx, key, values...).Err()
	observe("hset", err)
	return err
}
func (r *RedisCache) HGet(ctx context.Context, key, field string) (string, error) {
	v, err := r.client.HGet(ctx, key, field).Result()
	observe("hget", err)
	return v, err
}
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	observe("exists", err)
	return n > 0, err
}
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	observe("delete", err)
	return err
}
func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	observe("ttl", err)
	return ttl, err
}
func (r *RedisCache) Close() error { return r.client.Close() }

func (r *RedisCache) Client() *redis.Client { return r.client }

func observe(operation string, err error) {
	result := "ok"
	if err != nil {
		result = "error"
	}
	metrics.CacheOps.WithLabelValues(operation, result).Inc()
}