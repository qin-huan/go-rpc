package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	url, err := redis.ParseURL("redis://:@127.0.0.1:6379/1")
	if err != nil {
		panic(err)
	}
	fmt.Println("url is ", url)
	client := redis.NewClient(url)

	keys := client.WithContext(ctx).Keys(ctx, "*")
	fmt.Println(keys)
}