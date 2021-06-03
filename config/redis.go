package config

import (
	// "github.com/go-redis/redis/v7"

	"os"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func CreateRedisClient() {

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))

	if err != nil {
		panic(err)
	}

	Redis = redis.NewClient(opt)
}
