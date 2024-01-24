package redis

import (
	"context"
	"log"
	"time"

	r "github.com/go-redis/redis/v8"
	"github.com/zaza-hikayat/go-fiber/configs"
	"github.com/zaza-hikayat/go-fiber/domain"
)

type redisCache struct {
	client *r.Client
}

func NewRedisCahce(conf *configs.Config) domain.RedisCache {
	redisClient := r.NewClient(&r.Options{
		Network:  "tcp",
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Database,
	})

	if res := redisClient.Ping(context.TODO()); res.Err() != nil {
		log.Fatalf("err connect redis >>> %v", res.Err())
	}

	return &redisCache{
		client: redisClient,
	}
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) domain.IStatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

func (r *redisCache) Get(ctx context.Context, key string) domain.IStringCmd {
	return r.client.Get(ctx, key)
}

func (r *redisCache) Del(ctx context.Context, keys ...string) *r.IntCmd {
	return r.client.Del(ctx, keys...)
}

func (r *redisCache) HSet(ctx context.Context, key, field string, value interface{}) *r.IntCmd {
	return r.client.HSet(ctx, key, field, value)
}

func (r *redisCache) HGet(ctx context.Context, key, field string) *r.StringCmd {
	return r.client.HGet(ctx, key, field)
}

func (r *redisCache) GetAllKeys(ctx context.Context, pattern string) ([]string, error) {
	keys := []string{}

	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return keys, err
	}

	return keys, nil
}

func (r *redisCache) Keys(ctx context.Context, pattern string) *r.StringSliceCmd {
	return r.client.Keys(ctx, pattern)
}

func (r *redisCache) GetClientConnection() *r.Client {
	return r.client
}
