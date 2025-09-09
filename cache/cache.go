package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var ctx = context.Background()

func Connect(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func Set(shortCode string, longURL string, ttl time.Duration) error {
	return RDB.Set(ctx, shortCode, longURL, ttl).Err()
}

func Get(shortCode string) (string, error) {
	return RDB.Get(ctx, shortCode).Result()
}
