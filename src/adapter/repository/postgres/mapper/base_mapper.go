package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

func ToBaseModelMapper(en *entity.BaseEntity) model.BaseModel {
	return model.BaseModel{
		ID: en.ID,
	}
}
func ToBaseEntityMapper(baseModel *model.BaseModel) entity.BaseEntity {
	return entity.BaseEntity{
		ID:        baseModel.ID,
		CreatedAt: baseModel.CreatedAt.Unix(),
		UpdatedAt: baseModel.UpdatedAt.Unix(),
	}
}
