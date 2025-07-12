package client

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/adapter/http/client/dto"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/web/client"
)

type WebhookServiceAdapter struct {
	httpClient client.ContextualHttpClient
}

func (w WebhookServiceAdapter) Send(ctx context.Context, webhook *entity.WebhookEntity, message string) error {
	requestBody := dto.WebhookRequest{Text: message}
	resp, err := w.httpClient.Post(ctx, webhook.URL, requestBody, nil)
	if err != nil {
		log.Error(ctx, "Error when sending webhook", err)
	}
	if resp.StatusCode != 200 {
		log.Error(ctx, "Failed to send webhook", "status_code", resp.StatusCode, "url", webhook.URL)
		return fmt.Errorf("failed to send webhook to %s, status code: %d", webhook.URL, resp.StatusCode)
	}
	return nil
}

func NewWebhookServiceAdapter(httpClient client.ContextualHttpClient) port.IWebhookServicePort {
	return &WebhookServiceAdapter{
		httpClient: httpClient,
	}
}
