package port

import (
	"context"
)

type IRedisPort interface {
	SetToRedis(ctx context.Context, key string, value interface{}, expired int64) error
	GetFromRedis(ctx context.Context, key string) ([]byte, error)
	SetHSetToRedis(ctx context.Context, key string, mapFieldValue map[string]interface{}, expired int64) error
	GetHSetFromRedis(ctx context.Context, key string) (map[string]string, error)
}
