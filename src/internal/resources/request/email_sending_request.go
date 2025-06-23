package request

import "github.com/KhaiHust/email-notification-service/core/entity/dto/request"

type EmailSendingRequest struct {
	TemplateId  int64               `json:"template_id" validate:"required"`
	Datas       []*EmailSendingData `json:"datas" validate:"required,dive"`
	ProviderID  int64               `json:"provider_id"`
	Provider    string              `json:"provider" validate:"required,oneof=gmail outlook"`
	Environment string              `json:"-"`
}
type EmailSendingData struct {
	To      string            `json:"to" validate:"required,email"`
	SendAt  *int64            `json:"send_at,omitempty"`
	Subject map[string]string `json:"subject"`
	Body    map[string]string `json:"body"`
}

func ToEmailSendingRequestDto(req *EmailSendingRequest) *request.EmailSendingRequestDto {
	if req == nil {
		return nil
	}
	return &request.EmailSendingRequestDto{
		TemplateID: req.TemplateId,
		ProviderID: req.ProviderID,
		Datas:      ToListEmailSendingDataDto(req.Datas),
	}
}
func ToEmailSendingDataDto(req *EmailSendingData) *request.EmailSendingData {
	if req == nil {
		return nil
	}
	return &request.EmailSendingData{
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
		SendAt:  req.SendAt,
	}
}
func ToListEmailSendingDataDto(req []*EmailSendingData) []*request.EmailSendingData {
	if req == nil {
		return nil
	}
	var data []*request.EmailSendingData
	for _, item := range req {
		data = append(data, ToEmailSendingDataDto(item))
	}
	return data
}
