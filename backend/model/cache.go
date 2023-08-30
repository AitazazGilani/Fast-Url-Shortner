package model

import(
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var Ctx = context.Background()


//create a redis client
func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("CACHE_ADDR"),
		Password: os.Getenv("CACHE_PASS"),
		DB: dbNo,
	})
	return rdb
}