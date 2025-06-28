package thirdparty

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	nrredis "github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisService struct {
	redisClient *redis.Client
}

func (r RedisService) SetLock(ctx context.Context, key string, value string, expired int64) (bool, error) {
	ok, err := r.redisClient.SetNX(ctx, key, value, time.Duration(expired)*time.Second).Result()
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when set lock to redis: ", err))
		return false, err
	}
	if !ok {
		log.Warn(ctx, fmt.Sprintf("Key %s already exists", key))
		return false, nil
	}
	log.Info(ctx, fmt.Sprintf("Set lock for key %s with value %s", key, value))
	return true, nil
}

func (r RedisService) DeleteKey(ctx context.Context, key string) error {
	err := r.redisClient.Del(ctx, key).Err()
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when delete key from redis: ", err))
		if errors.Is(err, redis.Nil) {
			log.Warn(ctx, fmt.Sprintf("Key %s does not exist", key))
			return nil
		}
		return err
	}
	log.Info(ctx, fmt.Sprintf("Deleted key %s from redis", key))
	return nil
}

func (r RedisService) GetHSetFromRedis(ctx context.Context, key string) (map[string]string, error) {
	data, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when get data from redis: ", err))
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func (r RedisService) SetHSetToRedis(ctx context.Context, key string, mapFieldValue map[string]interface{}, expired int64) error {
	err := r.redisClient.HSet(ctx, key, mapFieldValue, time.Duration(expired)*time.Second).Err()
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when set data to redis: ", err))
		return err
	}
	return nil
}

func (r RedisService) SetToRedis(ctx context.Context, key string, value interface{}, expired int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when marshal data to json: ", err))
		return err
	}
	result, err := r.redisClient.Set(ctx, key, data, time.Duration(expired)*time.Second).Result()
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when set data to redis: ", err))
		return err
	}
	log.Info(ctx, fmt.Sprintf("Set data to redis with key: %s, result: %s", key, result))
	return nil
}

func (r RedisService) GetFromRedis(ctx context.Context, key string) ([]byte, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		log.Info(ctx, fmt.Sprintf("Key %s does not exist", key))
		return nil, nil
	case err != nil:
		log.Error(ctx, fmt.Sprint("Error when get data from redis: ", err))
		return nil, err
	case val == "":
		log.Warn(ctx, fmt.Sprintf("Key %s is empty", key))
		return nil, nil
	}
	return []byte(val), nil
}

func NewRedisService(redisClient *redis.Client) port.IRedisPort {

	// Add the New Relic hook
	redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	return &RedisService{
		redisClient: redisClient,
	}
}
