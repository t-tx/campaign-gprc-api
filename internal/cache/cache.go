package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func NewClient(host string) *redis.Client {
	ctx := context.Background()

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: host, // Redis server address
	})

	// Ensure Redis is reachable
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	return rdb
}

func main() {
	// Create a context
	ctx := context.Background()

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	// Ensure Redis is reachable
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	// Increment the key "counter"
	key := "counter"
	newValue, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		log.Fatalf("could not increment key: %v", err)
	}

	fmt.Printf("The new value of '%s' is: %d\n", key, newValue)
}
