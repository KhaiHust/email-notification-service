package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/samber/lo"
	"strconv"
)

type IValidateAccessWorkspaceUsecase interface {
	ValidateAccessWorkspaceByUserIdAndCode(ctx context.Context, userId int64, code string) (map[string]string, int64, error)
}
type ValidateAccessWorkspaceUsecase struct {
	redisPort           port.IRedisPort
	getWorkspaceUseCase IGetWorkspaceUseCase
}

func (v ValidateAccessWorkspaceUsecase) ValidateAccessWorkspaceByUserIdAndCode(ctx context.Context, userId int64, code string) (map[string]string, int64, error) {
	//check in redis
	mapWSUser, err := v.redisPort.GetHSetFromRedis(ctx, common.WorkspaceUserAccessKey+strconv.FormatInt(userId, 10))
	if err != nil {
		log.Error(ctx, "[ValidateAccessWorkspaceUsecase] Error getting workspace access from redis: %v", err)
		return nil, 0, err
	}
	if mapWSUser != nil {
		if role, ok := mapWSUser[code]; ok {
			return map[string]string{code: role}, 0, nil
		}
	}
	//check in db
	workspaces, err := v.getWorkspaceUseCase.GetWorkspaceByUserId(ctx, userId)
	if err != nil {
		log.Error(ctx, "[ValidateAccessWorkspaceUsecase] Error getting workspace by code: %v", err)
		return nil, 0, err
	}
	if len(workspaces) == 0 {
		return nil, 0, fmt.Errorf("workspace %s not found", code)
	}
	workspace, isFind := lo.Find(workspaces, func(w *entity.WorkspaceEntity) bool {
		return w.Code == code
	})
	if !isFind {
		return nil, 0, fmt.Errorf("workspace %s not found", code)
	}
	workspaceUser, isFind := lo.Find(workspace.WorkspaceUserEntity, func(wu entity.WorkspaceUserEntity) bool {
		return wu.UserID == userId
	})
	if !isFind {
		return nil, 0, fmt.Errorf("workspace %s not found", code)
	}
	return map[string]string{code: workspaceUser.Role}, workspace.ID, nil

}

func NewValidateAccessWorkspaceUsecase(
	redisPort port.IRedisPort,
	getWorkspaceUseCase IGetWorkspaceUseCase,
) IValidateAccessWorkspaceUsecase {
	return &ValidateAccessWorkspaceUsecase{
		redisPort:           redisPort,
		getWorkspaceUseCase: getWorkspaceUseCase,
	}
}
