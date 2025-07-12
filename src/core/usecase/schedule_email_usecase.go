package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/golibs-starter/golib/log"
	"sync"
	"time"
)

const (
	SendEmailPath      = "v1/tasks/email-request/schedule"
	X_Signature        = "X-Signature"
	HTTPMethodPostCode = 1
)

type IScheduleEmailUsecase interface {
	ScheduleEmail(ctx context.Context, emailRequest *entity.EmailRequestEntity) error
	ScheduleEmails(ctx context.Context, emailRequests []*entity.EmailRequestEntity) error
}
type ScheduleEmailUsecase struct {
	cloudTaskServicePort      port.ICloudTaskServicePort
	taskProps                 *properties.TaskProperties
	eventPublisher            port.IEventPublisher
	updateEmailRequestUsecase IUpdateEmailRequestUsecase
}

func (s ScheduleEmailUsecase) ScheduleEmails(ctx context.Context, emailRequests []*entity.EmailRequestEntity) error {
	if len(emailRequests) == 0 {
		log.Error(ctx, "Email requests data is empty")
		return fmt.Errorf("email requests data cannot be empty")
	}

	var wg sync.WaitGroup

	for _, emailRequest := range emailRequests {
		wg.Add(1)
		go func(req *entity.EmailRequestEntity) {
			defer wg.Done()
			if err := s.ScheduleEmail(ctx, req); err != nil {
				log.Error(ctx, "Error when scheduling email", err)
				req.Status = constant.EmailSendingStatusFailed
				req.ErrorMessage = fmt.Sprintf("Failed to schedule email: %v", err)
			} else {
				req.Status = constant.EmailSendingStatusScheduled
			}
		}(emailRequest)
	}

	wg.Wait()
	if _, err := s.updateEmailRequestUsecase.UpdateStatusByBatches(ctx, emailRequests); err != nil {
		log.Error(ctx, "Error when updating email requests status", err)
		return fmt.Errorf("failed to update email requests status: %v", err)
	}

	log.Info(ctx, "All email tasks scheduled successfully")
	return nil
}

func (s ScheduleEmailUsecase) ScheduleEmail(ctx context.Context, emailRequest *entity.EmailRequestEntity) error {
	if emailRequest == nil {
		log.Error(ctx, "Email data is nil")
		return fmt.Errorf("email data cannot be nil")
	}
	payloadData := &request.ScheduleEmailPayload{
		EmailRequestID: emailRequest.ID,
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
		ScheduleTime: emailRequest.SendAt,
		Headers: map[string]string{
			X_Signature: signature,
		},
		MethodCode: HTTPMethodPostCode,
		Payload:    payloadData,
	}
	err = s.cloudTaskServicePort.CreateNewTask(ctx, taskRequestBody)
	if err != nil {
		log.Error(ctx, "Error when create new task", err)
		return err
	}
	//todo: fire event to update email request status to scheduled
	//emailRequest.Status = constant.EmailSendingStatusScheduled
	//ev := event.NewEventEmailRequestSync(ctx, emailRequest)
	//s.eventPublisher.Publish(ev)
	//log.Info(ctx,
	//	fmt.Sprintf("Schedule email task created successfully for EmailRequestID: %s",
	//		emailRequest.RequestID))
	return nil
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
func NewScheduleEmailUsecase(
	cloudTaskServicePort port.ICloudTaskServicePort,
	taskProps *properties.TaskProperties,
	eventPublisher port.IEventPublisher,
	updateEmailRequestUsecase IUpdateEmailRequestUsecase,
) IScheduleEmailUsecase {
	return &ScheduleEmailUsecase{
		cloudTaskServicePort:      cloudTaskServicePort,
		taskProps:                 taskProps,
		eventPublisher:            eventPublisher,
		updateEmailRequestUsecase: updateEmailRequestUsecase,
	}
}
