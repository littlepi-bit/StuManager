package Model

import (
	"context"
	"fmt"
	"time"

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

var RedisErr = redis.Nil

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

func SetHash(key, value string, expiration time.Duration) error {
	err := client.Set(context.Background(), key, value, expiration).Err()
	return err
}

func GetHash(key string) (string, error) {
	value, err := client.Get(context.Background(), key).Result()
	return value, err
}
