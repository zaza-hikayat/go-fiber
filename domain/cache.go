package domain

import (
	"context"
	"time"

	r "github.com/go-redis/redis/v8"
)

type RedisCache interface {
	Keys(ctx context.Context, pattern string) *r.StringSliceCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) IStatusCmd
	Get(ctx context.Context, key string) IStringCmd
	Del(ctx context.Context, keys ...string) *r.IntCmd
	HSet(ctx context.Context, key, field string, value interface{}) *r.IntCmd
	HGet(ctx context.Context, key, field string) *r.StringCmd
	GetAllKeys(ctx context.Context, key string) ([]string, error)
	GetClientConnection() *r.Client
}

type IStringCmd interface {
	Bool() (bool, error)
	Bytes() ([]byte, error)
	Float32() (float32, error)
	Float64() (float64, error)
	Int() (int, error)
	Int64() (int64, error)
	Result() (string, error)
	Scan(val interface{}) error
	String() string
	Time() (time.Time, error)
	Uint64() (uint64, error)
	Val() string
}

type IStatusCmd interface {
	Result() (string, error)
	String() string
	Val() string
	Err() error
}
