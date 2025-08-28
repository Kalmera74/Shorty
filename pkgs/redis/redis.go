package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func InitClient() {
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
}