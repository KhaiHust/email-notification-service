package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IEmailTemplateRepositoryPort interface {
	SaveNewTemplate(ctx context.Context, tx *gorm.DB, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error)
}
