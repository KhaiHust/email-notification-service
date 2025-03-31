package mapper

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/utils"
)

func ToEmailProviderModel(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity) (*model.EmailProviderModel, error) {
	var emailProviderModel model.EmailProviderModel
	if err := utils.CopyStruct(&emailProviderModel, emailProviderEntity); err != nil {
		return nil, err
	}
	return &emailProviderModel, nil
}
