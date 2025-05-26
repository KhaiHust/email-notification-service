package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/golibs-starter/golib/log"
	"time"
)

const (
	SendEmailPath      = "/v1/email-request/schedule"
	X_Signature        = "X-Signature"
	HTTPMethodPostCode = 1
)

type IScheduleEmailUsecase interface {
	ScheduleEmail(ctx context.Context, emailData *request.EmailSendingData) error
}
type ScheduleEmailUsecase struct {
	cloudTaskServicePort port.ICloudTaskServicePort
	taskProps            *properties.TaskProperties
}

func (s ScheduleEmailUsecase) ScheduleEmail(ctx context.Context, emailData *request.EmailSendingData) error {
	if emailData == nil {
		log.Error(ctx, "Email data is nil")
		return fmt.Errorf("email data cannot be nil")
	}
	payloadData := &request.ScheduleEmailPayload{
		EmailRequestID: emailData.EmailRequestID,
	}
	payloadJson, err := json.Marshal(payloadData)
	if err != nil {
		log.Error(ctx, "Error when marshal payload data", err)
		return err
	}
	signature := s.GenerateChecksum(payloadJson)
	taskRequestBody := &request.CreateTaskDto{
		TaskName:     fmt.Sprintf("Email-Task-%d", time.Now().UnixNano()),
		TargetUrl:    s.taskProps.BaseUrl + SendEmailPath,
		ScheduleTime: emailData.SendAt,
		Headers: map[string]string{
			X_Signature: signature,
		},
		MethodCode: HTTPMethodPostCode,
		Payload:    payloadJson,
	}
	err = s.cloudTaskServicePort.CreateNewTask(ctx, taskRequestBody)
	if err != nil {
		log.Error(ctx, "Error when create new task", err)
		return err
	}
	log.Info(ctx,
		fmt.Sprintf("Schedule email task created successfully for EmailRequestID: %s",
			emailData.EmailRequestID))
	return nil
}

func NewScheduleEmailUsecase(cloudTaskServicePort port.ICloudTaskServicePort) IScheduleEmailUsecase {
	return &ScheduleEmailUsecase{
		cloudTaskServicePort: cloudTaskServicePort,
	}
}
func (s ScheduleEmailUsecase) GenerateChecksum(body []byte) string {
	hc := hmac.New(sha256.New, []byte(s.taskProps.SecretKey))
	hc.Write(body)
	return hex.EncodeToString(hc.Sum(nil))
}
func (s ScheduleEmailUsecase) ValidateChecksum(body []byte, signature string) bool {
	expectedSignature := s.GenerateChecksum(body)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
