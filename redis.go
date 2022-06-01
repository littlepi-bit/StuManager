package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func InitRedis(url string) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: url,
		DB:   0,
	})

	return
}

func testRedis() {
	_, err := client.Get(context.Background(), "zngw").Result()
	if err == redis.Nil {
		fmt.Println("key不存在")
	} else {
		fmt.Println("key存在")
	}

	// 设置string
	err = client.Set(context.Background(), "zngw", "hello", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	// 获取string
	val2, err := client.Get(context.Background(), "zngw").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("读取 zngw :", val2)

}
