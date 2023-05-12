package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"fmt"
)

var (
	ctx     = context.Background()
	redisDb *redis.Client
)

func RedisInit(redisAddr, password string, db int) {
	redisDb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: password,
		DB: db,
	})

	pong, err := redisDb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("connect redis fail: %v", err))
	} else {
		fmt.Printf("connect redis succ %v\n", pong)
	}

}