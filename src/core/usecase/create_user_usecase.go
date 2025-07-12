package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type ICreateUserUseCase interface {
	CreateNewUser(ctx context.Context, userEntity *entity.UserEntity) (*entity.UserEntity, error)
}
type CreateUserUseCase struct {
	getUserUseCase             IGetUserUseCase
	databaseTransactionUseCase IDatabaseTransactionUseCase
	userRepo                   port.IUserRepositoryPort
	hashPasswordUseCase        IHashPasswordUseCase
	redisPort                  port.IRedisPort
}

func (c CreateUserUseCase) CreateNewUser(ctx context.Context, userEntity *entity.UserEntity) (*entity.UserEntity, error) {
	//check existed user
	existedUser, err := c.getUserUseCase.GetUserByEmail(ctx, userEntity.Email)
	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		log.Error(ctx, "GetUserByEmail error: %v", err)
		return nil, common.ErrInternalServer
	}
	if existedUser != nil {
		log.Error(ctx, "User already existed")
		return nil, common.ErrEmailIsExisted
	}
	// hash password
	hashedPassword, err := c.hashPasswordUseCase.Hash(userEntity.Password)
	if err != nil {
		log.Error(ctx, "Hash error: %v", err)
		return nil, err
	}
	tx := c.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Error(ctx, "Rollback error: %v", errRollback)
		} else {
			log.Info(ctx, "Rollback successfully")
		}
	}()
	userEntity.Password = hashedPassword
	userEntity, err = c.userRepo.SaveUser(ctx, tx, userEntity)
	if err != nil {
		log.Error(ctx, "SaveUser error: %v", err)
		return nil, err
	}
	if errCommit := c.databaseTransactionUseCase.CommitTx(tx); errCommit != nil {
		log.Error(ctx, "Commit error: %v", errCommit)
		return nil, errCommit
	}
	// Set user to redis
	go func() {
		cacheKey := fmt.Sprintf(common.UserCacheKeyByEmail, userEntity.Email)
		errSaveRedis := c.redisPort.SetToRedis(ctx, cacheKey, userEntity, constant.DefaultTTL)
		if errSaveRedis != nil {
			log.Error(ctx, "SetToRedis error: %v", errSaveRedis)
		}
	}()
	//todo: send email verify
	return userEntity, nil

}

func NewCreateUserUseCase(
	getUserUseCase IGetUserUseCase,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	userRepo port.IUserRepositoryPort,
	hashPasswordUseCase IHashPasswordUseCase,
	redisPort port.IRedisPort,
) ICreateUserUseCase {
	return &CreateUserUseCase{
		getUserUseCase:             getUserUseCase,
		databaseTransactionUseCase: databaseTransactionUseCase,
		userRepo:                   userRepo,
		hashPasswordUseCase:        hashPasswordUseCase,
		redisPort:                  redisPort,
	}
}
