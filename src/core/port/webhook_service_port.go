package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

type IWebhookServicePort interface {
	Send(ctx context.Context, webhook *entity.WebhookEntity, message string) error
}
