package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"gorm.io/gorm"
)

type IEmailTemplateRepositoryPort interface {
	SaveNewTemplate(ctx context.Context, tx *gorm.DB, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error)
	GetTemplateByID(ctx context.Context, ID int64) (*entity.EmailTemplateEntity, error)
	GetAllTemplates(ctx context.Context, filter *request.GetListEmailTemplateFilter) ([]*entity.EmailTemplateEntity, error)
	CountAllTemplates(ctx context.Context, filter *request.GetListEmailTemplateFilter) (int64, error)
	GetTemplateByIDAndWorkspaceID(ctx context.Context, ID int64, workspaceID int64) (*entity.EmailTemplateEntity, error)
}
