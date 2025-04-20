package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/samber/lo"
	"strconv"
)

type IGetWorkspaceUseCase interface {
	GetWorkspaceByCode(ctx context.Context, userId int64, code string) (*entity.WorkspaceEntity, error)
	GetWorkspaceByUserId(ctx context.Context, userId int64) ([]*entity.WorkspaceEntity, error)
}
type GetWorkspaceUseCase struct {
	workspaceRepositoryPort port.IWorkspaceRepositoryPort
	redisPort               port.IRedisPort
}

func (g GetWorkspaceUseCase) GetWorkspaceByUserId(ctx context.Context, userId int64) ([]*entity.WorkspaceEntity, error) {
	workspaces, err := g.workspaceRepositoryPort.GetWorkspaceByUserId(ctx, userId)
	if err != nil {
		log.Error(ctx, "[GetWorkspaceUseCase] Error getting workspaces by userId: %v", err)
		return nil, err
	}
	//build cache access for validate access workspace
	go func(workspaces []*entity.WorkspaceEntity) {
		if len(workspaces) == 0 {
			return
		}
		mapWSRoles := make(map[string]interface{})
		for _, workspace := range workspaces {
			workspaceUser, isFind := lo.Find(workspace.WorkspaceUserEntity, func(wu entity.WorkspaceUserEntity) bool {
				return wu.UserID == userId
			})
			if !isFind {
				continue
			}
			mapWSRoles[workspace.Code] = workspaceUser.Role
		}
		//save to HSet cache
		err = g.redisPort.SetHSetToRedis(ctx, fmt.Sprintf(common.WorkspaceUserAccessKey, strconv.FormatInt(userId, 10)), mapWSRoles, constant.DefaultTTL)
		if err != nil {
			log.Error(ctx, "[GetWorkspaceUseCase] Error setting workspace access to redis: %v", err)
		}
	}(workspaces)
	return workspaces, nil
}

func (g GetWorkspaceUseCase) GetWorkspaceByCode(ctx context.Context, userId int64, code string) (*entity.WorkspaceEntity, error) {
	return g.workspaceRepositoryPort.GetWorkspaceByCodeAndUserId(ctx, userId, code)
}

func NewGetWorkspaceUseCase(
	workspaceRepositoryPort port.IWorkspaceRepositoryPort,
	redisPort port.IRedisPort,
) IGetWorkspaceUseCase {
	return &GetWorkspaceUseCase{
		workspaceRepositoryPort: workspaceRepositoryPort,
		redisPort:               redisPort,
	}
}
