package mapper

import (
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/event/message"
)

func ToEmailSendingDto(msg *message.EmailRequestSendingMessage) *request.EmailSendingRequestDto {
	return &request.EmailSendingRequestDto{
		TemplateId: msg.TemplateId,
		Data:       ToListSendingDataDto(msg.SendData),
	}
}
func ToListSendingDataDto(data []*message.EmailSendData) []*request.EmailSendingData {
	list := make([]*request.EmailSendingData, 0)
	for _, item := range data {
		list = append(list, ToSendingDataDto(item))
	}
	return list
}
func ToSendingDataDto(data *message.EmailSendData) *request.EmailSendingData {
	return &request.EmailSendingData{
		To:      data.To,
		Subject: data.Subject,
		Body:    data.Body,
	}
}
