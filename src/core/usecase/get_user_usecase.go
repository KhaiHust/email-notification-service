package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IGetUserUseCase interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	GetListUsers(ctx context.Context, workspaceID int64) ([]*entity.UserEntity, error)
}
type GetUserUseCase struct {
	userRepo                    port.IUserRepositoryPort
	workspaceUserRepositoryPort port.IWorkspaceUserRepositoryPort
	redisPort                   port.IRedisPort
}

func (g GetUserUseCase) GetListUsers(ctx context.Context, workspaceID int64) ([]*entity.UserEntity, error) {
	workspaceUsers, err := g.workspaceUserRepositoryPort.GetWorkspaceUserByWorkspaceID(ctx, workspaceID)
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when get workspace users: ", err))
		return nil, err
	}
	userIDs := make([]int64, 0, len(workspaceUsers))
	for _, workspaceUser := range workspaceUsers {
		userIDs = append(userIDs, workspaceUser.UserID)
	}
	users, err := g.userRepo.GetUsersByIds(ctx, userIDs)
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when get users by IDs: ", err))
		return nil, err
	}
	mapWorkspaceUser := make(map[int64]*entity.WorkspaceUserEntity)
	for _, workspaceUser := range workspaceUsers {
		mapWorkspaceUser[workspaceUser.UserID] = workspaceUser
	}
	for _, user := range users {
		if workspaceUser, ok := mapWorkspaceUser[user.ID]; ok {
			user.WorkspaceUserEntity = workspaceUser
		}
	}
	return users, nil
}

func (g GetUserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	cacheKey := fmt.Sprintf(common.UserCacheKeyByEmail, email)
	userData, err := g.redisPort.GetFromRedis(ctx, cacheKey)
	if err != nil {
		log.Error(ctx, fmt.Sprint("Error when get data from redis: ", err))
		return nil, err
	}
	if userData != nil {
		user := &entity.UserEntity{}
		err = json.Unmarshal(userData, user)
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
	go func() {
		errSaveRedis := g.redisPort.SetToRedis(ctx, cacheKey, userEntity, constant.DefaultTTL)
		if errSaveRedis != nil {
			log.Error(ctx, fmt.Sprint("Error when set data to redis: ", errSaveRedis))
		}
	}()
	return userEntity, nil
}

func NewGetUserUseCase(
	userRepo port.IUserRepositoryPort,
	workspaceUserRepositoryPort port.IWorkspaceUserRepositoryPort,
	redisPort port.IRedisPort,
) IGetUserUseCase {
	return &GetUserUseCase{
		userRepo:                    userRepo,
		workspaceUserRepositoryPort: workspaceUserRepositoryPort,
		redisPort:                   redisPort,
	}
}
