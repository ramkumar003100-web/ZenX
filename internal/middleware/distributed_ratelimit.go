package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"zenx/internal/router"
)

// DistributedRateLimitTokenBucket uses Redis to coordinate limits across nodes.
func DistributedRateLimitTokenBucket(client *redis.Client, namespace string, burst int, refill time.Duration) router.Middleware {
	if namespace == "" {
		namespace = "zenx:rl"
	}
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			key := fmt.Sprintf("%s:%s", namespace, clientIP(c.Request))
			allowed, err := redisTokenBucket(c.Request.Context(), client, key, burst, refill)
			if err != nil {
				http.Error(c.Writer, "rate limiter unavailable", http.StatusServiceUnavailable)
				return
			}
			if !allowed {
				http.Error(c.Writer, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next(c)
		}
	}
}

func redisTokenBucket(ctx context.Context, client *redis.Client, key string, burst int, refill time.Duration) (bool, error) {
	// Sliding window approximation with sorted set timestamps.
	now := time.Now().UnixMilli()
	window := refill.Milliseconds()
	pipe := client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", now-window))
	count := pipe.ZCard(ctx, key)
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
	pipe.Expire(ctx, key, refill)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	return count.Val() < int64(burst), nil
}
