package redis

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func InitClient() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	if redisPassword == "" {
		redisPassword = ""
	}
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Printf("Invalid Redis DB number, defaulting to 0: %v", err)
		redisDB = 0
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	_, err = Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}
