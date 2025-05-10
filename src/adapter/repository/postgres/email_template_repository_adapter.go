package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/specification"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmailTemplateRepositoryAdapter struct {
	base
}

func (e EmailTemplateRepositoryAdapter) UpdateTemplate(ctx context.Context, tx *gorm.DB, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error) {
	emailTemplateModel := mapper.ToEmailTemplateModel(template)
	if err := tx.WithContext(ctx).Model(&model.EmailTemplateModel{}).Where("id = ?", template.ID).Updates(emailTemplateModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailTemplateEntity(emailTemplateModel), nil

}

func (e EmailTemplateRepositoryAdapter) GetTemplateForUpdateByIDAndWorkspaceID(ctx context.Context, tx *gorm.DB, ID int64, workspaceID int64) (*entity.EmailTemplateEntity, error) {
	var templateModel model.EmailTemplateModel
	if err := tx.WithContext(ctx).Clauses(
		clause.Locking{Strength: clause.LockingStrengthUpdate},
	).Where("id = ? AND workspace_id = ?", ID, workspaceID).First(&templateModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailTemplateEntity(&templateModel), nil
}

func (e EmailTemplateRepositoryAdapter) GetTemplateByIDAndWorkspaceID(ctx context.Context, ID int64, workspaceID int64) (*entity.EmailTemplateEntity, error) {
	var emailTemplateModel model.EmailTemplateModel
	if err := e.db.WithContext(ctx).Where("id = ? AND workspace_id = ?", ID, workspaceID).First(&emailTemplateModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailTemplateEntity(&emailTemplateModel), nil
}

func (e EmailTemplateRepositoryAdapter) CountAllTemplates(ctx context.Context, filter *request.GetListEmailTemplateFilter) (int64, error) {
	emailTemplateSpecification := specification.ToEmailTemplateSpecification(filter)
	query, args, err := specification.NewEmailTemplateSpecificationQueryWithCount(emailTemplateSpecification)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := e.db.WithContext(ctx).Raw(query, args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e EmailTemplateRepositoryAdapter) GetAllTemplates(ctx context.Context, filter *request.GetListEmailTemplateFilter) ([]*entity.EmailTemplateEntity, error) {
	emailTemplateSpecification := specification.ToEmailTemplateSpecification(filter)
	query, args, err := specification.NewEmailTemplateSpecificationQuery(emailTemplateSpecification)
	if err != nil {
		return nil, err
	}
	var emailTemplateModels []*model.EmailTemplateModel
	if err := e.db.WithContext(ctx).Raw(query, args...).Scan(&emailTemplateModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailTemplateEntities(emailTemplateModels), nil
}

func (e EmailTemplateRepositoryAdapter) GetTemplateByID(ctx context.Context, ID int64) (*entity.EmailTemplateEntity, error) {
	var emailTemplateModel model.EmailTemplateModel
	if err := e.db.WithContext(ctx).Where("id = ?", ID).First(&emailTemplateModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailTemplateEntity(&emailTemplateModel), nil
}

func (e EmailTemplateRepositoryAdapter) SaveTemplate(ctx context.Context, tx *gorm.DB, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error) {
	emailTemplateModel := mapper.ToEmailTemplateModel(template)
	if err := tx.WithContext(ctx).Model(&model.EmailTemplateModel{}).Create(emailTemplateModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailTemplateEntity(emailTemplateModel), nil
}

func NewEmailTemplateRepositoryAdapter(db *gorm.DB) port.IEmailTemplateRepositoryPort {
	return &EmailTemplateRepositoryAdapter{
		base: base{db: db},
	}
}
