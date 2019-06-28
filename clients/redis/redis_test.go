package redis_test

import (
	cache "ark-common/clients/redis"
	"os"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/go-redis/redis"
)

func TestRedisLuaRateLimit(t *testing.T) {
	opt := redis.Options{
		Addr:     os.Getenv("ARK_REDIS_ADDR"),
		Password: os.Getenv("ARK_REDIS_PASSWD"),
		DB:       1,
	}
	client := cache.NewClient(&opt)
	if client == nil {
		t.Fatalf("redis connect failed")
	}
	key := "tencent-api-ratelimit"
	limiter := cache.NewLimiter(client, key, 1*time.Minute, 10, 1)
	time.Sleep(20 * time.Second)
	if limiter == nil {
		t.Fatalf("new limiter failed")
	}
	allowd := limiter.Take()
	t.Logf("Current ratelimit: %v", allowd)
}

func TestRedisRateLimit(t *testing.T) {
	opt := &redis.RingOptions{
		Addrs: map[string]string{
			"s1": os.Getenv("ARK_REDIS_ADDR"),
		},
		Password: os.Getenv("ARK_REDIS_PASSWD"),
		DB:       1,
	}
	client := cache.NewRingClient(opt)
	key := "aliyun-api-ratelimit"
	limiter := cache.NewRateLimiter(client)
	delay, allowd := limiter.AllowRate(key, 2*rate.Every(time.Minute))
	time.Sleep(30 * time.Second)
	t.Logf("Current request Delay: %v, Allowed: %v", delay, allowd)
	delay, allowd = limiter.AllowRate(key, 2*rate.Every(time.Minute))
	t.Logf("Current request Delay: %v, Allowed: %v", delay, allowd)
}
