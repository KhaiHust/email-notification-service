package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IUserRepositoryPort interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	SaveUser(ctx context.Context, tx *gorm.DB, user *entity.UserEntity) (*entity.UserEntity, error)
}
