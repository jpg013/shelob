package db

import (
	"github.com/go-redis/redis/v7"
)

var (
	redisClient *redis.Client
)

func ExampleNewClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetCache(key string) (string, error) {
	return redisClient.Get(key).Result()
}
