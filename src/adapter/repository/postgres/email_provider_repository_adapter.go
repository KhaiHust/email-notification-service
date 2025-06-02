package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailProviderRepositoryAdapter struct {
	base
}

func (e EmailProviderRepositoryAdapter) GetProvidersByIds(ctx context.Context, ids []int64) ([]*entity.EmailProviderEntity, error) {
	var emailProviderModels []*model.EmailProviderModel
	if err := e.db.WithContext(ctx).Model(&model.EmailProviderModel{}).Where("id IN ?", ids).Find(&emailProviderModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailProviderEntity(emailProviderModels), nil
}

func (e EmailProviderRepositoryAdapter) GetAllEmailProviders(ctx context.Context, filter *request.GetEmailProviderRequestFilter) ([]*entity.EmailProviderEntity, error) {
	//build query
	query := e.db.WithContext(ctx).Model(&model.EmailProviderModel{}).
		Select("id, workspace_id, provider, email, from_name, environment")

	if filter != nil {
		conditions := make(map[string]interface{})
		if filter.WorkspaceID != nil {
			conditions["workspace_id"] = *filter.WorkspaceID
		}
		if filter.Provider != nil {
			conditions["provider"] = *filter.Provider
		}
		if filter.Environment != nil {
			conditions["environment"] = *filter.Environment
		}
		if len(conditions) > 0 {
			query = query.Where(conditions)
		}
	}
	var emailProviderModels []*model.EmailProviderModel
	if err := query.Find(&emailProviderModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailProviderEntity(emailProviderModels), nil
}

func (e EmailProviderRepositoryAdapter) GetEmailProviderByIds(ctx context.Context, ids []int64) ([]*entity.EmailProviderEntity, error) {
	var emailProviderModels []*model.EmailProviderModel
	if err := e.db.WithContext(ctx).Select("id,provider").Where("id IN ?", ids).Find(&emailProviderModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailProviderEntity(emailProviderModels), nil
}

func (e EmailProviderRepositoryAdapter) GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) (*entity.EmailProviderEntity, error) {
	//build raw sql
	sql := "SELECT ep.id, ep.workspace_id, ep.provider, ep.email, ep.from_name, ep.environment FROM email_providers ep " +
		"JOIN workspaces w ON ep.workspace_id = w.id " +
		"WHERE w.code = ? AND ep.provider = ?"
	var emailProviderModel model.EmailProviderModel
	if err := e.db.WithContext(ctx).Raw(sql, workspaceCode, provider).Scan(&emailProviderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailProviderEntity(&emailProviderModel), nil
}

func (e EmailProviderRepositoryAdapter) GetEmailProviderByIDAndWorkspaceID(ctx context.Context, providerID, workspaceID int64) (*entity.EmailProviderEntity, error) {
	var emailProviderModels *model.EmailProviderModel
	if err := e.db.WithContext(ctx).Where("id = ? AND workspace_id = ?", providerID, workspaceID).First(&emailProviderModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailProviderEntity(emailProviderModels), nil
}

func (e EmailProviderRepositoryAdapter) UpdateEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error) {
	emailProviderModel := mapper.ToEmailProviderModel(emailProvider)
	if err := tx.WithContext(ctx).Model(&model.EmailProviderModel{}).
		Where("id = ?", emailProvider.ID).
		Updates(emailProviderModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailProviderEntity(emailProviderModel), nil
}

func (e EmailProviderRepositoryAdapter) GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error) {
	var emailProviderModel model.EmailProviderModel
	if err := e.db.WithContext(ctx).Where("id = ?", ID).First(&emailProviderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailProviderEntity(&emailProviderModel), nil
}

func (e EmailProviderRepositoryAdapter) SaveEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error) {
	emailProviderModel := mapper.ToEmailProviderModel(emailProvider)

	if err := tx.WithContext(ctx).Create(emailProviderModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailProviderEntity(emailProviderModel), nil
}

func NewEmailProviderRepositoryAdapter(db *gorm.DB) port.IEmailProviderRepositoryPort {
	return &EmailProviderRepositoryAdapter{
		base: base{db},
	}
}
