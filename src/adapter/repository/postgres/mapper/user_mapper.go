package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

func ToUserModel(userEntity *entity.UserEntity) *model.UserModel {
	return &model.UserModel{
		FullName: userEntity.FullName,
		Email:    userEntity.Email,
	}
}
func ToUserEntity(userModel *model.UserModel) *entity.UserEntity {
	return &entity.UserEntity{
		BaseEntity: entity.BaseEntity{
			ID:        userModel.ID,
			CreatedAt: userModel.CreatedAt.Unix(),
			UpdatedAt: userModel.UpdatedAt.Unix(),
		},
		FullName: userModel.FullName,
		Email:    userModel.Email,
	}
}
