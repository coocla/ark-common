package redis

import (
	"bk-cmdb/src/framework/core/log"
	"time"

	"github.com/go-redis/redis_rate"

	"github.com/go-redis/redis"
	"github.com/vearne/ratelimit"
)

// NewClient 生成新的redis客户端
func NewClient(opt *redis.Options) *redis.Client {
	r := redis.NewClient(opt)
	_, err := r.Ping().Result()
	if err != nil {
		return nil
	}
	return r
}

// NewLimiter 初始化一个配额限速
func NewLimiter(r *redis.Client, key string, duration time.Duration, capacity int, size int) *ratelimit.RedisRateLimiter {
	limiter, err := ratelimit.NewRedisRateLimiter(r, key, duration, capacity, size, ratelimit.TokenBucketAlg)
	if err != nil {
		log.Errorf("initializate limiter failed: %v", err)
		return nil
	}
	return limiter
}

func NewRingClient(opt *redis.RingOptions) *redis.Ring {
	r := redis.NewRing(opt)
	if err := r.FlushDb().Err(); err != nil {
		return nil
	}
	return r
}

func NewRateLimiter(r *redis.Ring) *redis_rate.Limiter {
	limiter := redis_rate.NewLimiter(r)
	return limiter
}
