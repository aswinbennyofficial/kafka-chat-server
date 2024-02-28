package main

import (
	"context"
	"log"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func ConnectRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisClient = client
}

func PublishToRedis(msg string) {
	ctx := context.Background()
	err := redisClient.Publish(ctx, "broadcast", msg).Err()
	if err != nil {
		panic(err)
	}

}

func SubscribeToRedis() {
	ctx := context.Background()
	// There is no error because go-redis automatically reconnects on error.
	pubsub := redisClient.Subscribe(ctx, "broadcast")
	// Close the subscription when we are done.
	defer pubsub.Close()

	// or using go routines
	ch := pubsub.Channel()

	for msg := range ch {
		log.Println(msg.Channel, msg.Payload)
	}

}
