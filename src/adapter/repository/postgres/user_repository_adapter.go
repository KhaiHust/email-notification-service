package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type UserRepositoryAdapter struct {
	base
}

func (u UserRepositoryAdapter) SaveUser(ctx context.Context, tx *gorm.DB, user *entity.UserEntity) (*entity.UserEntity, error) {
	userModel := mapper.ToUserModel(user)
	if err := tx.WithContext(ctx).Model(&model.UserModel{}).Create(userModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToUserEntity(userModel), nil
}

func (u UserRepositoryAdapter) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	var userModel model.UserModel
	if err := u.db.WithContext(ctx).Model(&model.UserModel{}).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToUserEntity(&userModel), nil

}

func NewUserRepositoryAdapter(base base) port.IUserRepositoryPort {
	return &UserRepositoryAdapter{
		base: base,
	}
}
