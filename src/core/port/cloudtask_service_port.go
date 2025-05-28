package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
)

type ICloudTaskServicePort interface {
	CreateNewTask(ctx context.Context, request *request.CreateTaskDto) error
}
