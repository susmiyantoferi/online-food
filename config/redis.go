package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func RedisCLient() *redis.Client {
	rdsAddrs := os.Getenv("RDS_ADDRS")
	rdsPwd := os.Getenv("RDS_PWD")
	rdsDb := os.Getenv("RDS_DB")
	rdb, _ := strconv.Atoi(rdsDb)

	Redis := redis.NewClient(&redis.Options{
		Addr:     rdsAddrs,
		Password: rdsPwd,
		DB:       rdb,
	})

	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("redis: %v", err)
	}

	return Redis
}
