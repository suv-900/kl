package controllers

import "github.com/redis/go-redis/v9"

func ConnectRedis() {
	rdb := redis.NewClient()
}
