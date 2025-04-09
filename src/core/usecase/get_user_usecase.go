package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

const (
	// UserCacheKeyByEmail Cache key for user
	UserCacheKeyByEmail = "user:email:%s"
)

type IGetUserUseCase interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
}
type GetUserUseCase struct {
	userRepo  port.IUserRepositoryPort
	redisPort port.IRedisPort
}

func (g GetUserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	cacheKey := fmt.Sprint(UserCacheKeyByEmail, email)
	userData, err := g.redisPort.GetFromRedis(ctx, cacheKey)
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when get data from redis: ", err))
		return nil, err
	}
	if userData != nil {
		user := &entity.UserEntity{}
		err = user.UnmarshalBinary(userData)
		if err != nil {
			log.Error(ctx, fmt.Sprint("Error when unmarshal data from redis: ", err))
			return nil, err
		}
		return user, nil
	}
	userEntity, err := g.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when get user from db: ", err))
		return nil, err
	}

	// Set user to redis
	err = g.redisPort.SetToRedis(ctx, cacheKey, userEntity, constant.DefaultTTL)
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when set data to redis: ", err))
		return nil, err
	}
	return userEntity, nil
}

func NewGetUserUseCase(userRepo port.IUserRepositoryPort, redisPort port.IRedisPort) IGetUserUseCase {
	return &GetUserUseCase{
		userRepo:  userRepo,
		redisPort: redisPort,
	}
}
