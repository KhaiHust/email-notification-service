package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailRequestRepositoryAdapter struct {
	base
}

func (e EmailRequestRepositoryAdapter) UpdateStatusByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmailRequestRepositoryAdapter) SaveEmailRequestByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error) {
	emailRequestModels := mapper.ToListEmailRequestModel(emailRequests)
	if err := tx.WithContext(ctx).Model(&model.EmailRequestModel{}).Create(emailRequestModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailRequestEntity(emailRequestModels), nil
}

func NewEmailRequestRepositoryAdapter(db *gorm.DB) port.IEmailRequestRepositoryPort {
	return &EmailRequestRepositoryAdapter{
		base: base{
			db: db,
		}}
}
