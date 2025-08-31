package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// TODO: Would be better to be DI instead of global
var Client *redis.Client

type Cacher interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type RedisCacher struct {
	Client redis.Cmdable
}

func NewCacher(redisClient redis.Cmdable) Cacher {
	return &RedisCacher{
		Client: redisClient,
	}
}

func (r *RedisCacher) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return r.Client.Set(ctx, key, value, exp).Err()
}

func (r *RedisCacher) Get(ctx context.Context, key string) (string, error) {
	get := r.Client.Get(ctx, key)
	return get.Result()
}
func InitRedisClient() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisDBStr := os.Getenv("REDIS_DB")

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	if redisHost == "" || redisPort == "" {
		redisAddr = "localhost:6379"
		log.Printf("REDIS_HOST or REDIS_PORT not set, defaulting to %s", redisAddr)
	}

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Printf("Invalid Redis DB number, defaulting to 0: %v", err)
		redisDB = 0
	}

	Client = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})

	_, err = Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis at %s: %v", redisAddr, err)
	}

	return Client
}
