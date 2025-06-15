package cronjob

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	golibcron "github.com/golibs-starter/golib-cron"
	"github.com/golibs-starter/golib/log"
)

const (
	EmailSendRetryCronJobLockKey        = "email_send_retry_cron_job_lock"
	ExpiredTimeForEmailSendRetryCronJob = 4 // minutes
)

type EmailSendRetryCronJob struct {
	emailSendRetryUsecase usecase.IEmailSendRetryUsecase
	redisPort             port.IRedisPort
}

func (e EmailSendRetryCronJob) Run(ctx context.Context) {
	//check if the job is already running
	ok, err := e.redisPort.SetLock(ctx, EmailSendRetryCronJobLockKey, "1", ExpiredTimeForEmailSendRetryCronJob)
	if err != nil {
		log.Error(ctx, "Error when set lock for email send retry cron job", err)
		return
	}
	if !ok {
		log.Warn(ctx, "Email send retry cron job is already running")
		return
	}
	defer func() {
		if err := e.redisPort.DeleteKey(ctx, EmailSendRetryCronJobLockKey); err != nil {
			log.Error(ctx, "Error when delete lock for email send retry cron job", err)
		} else {
			log.Info(ctx, "Deleted lock for email send retry cron job")
		}
	}()
	if err := e.emailSendRetryUsecase.ProcessBatches(ctx); err != nil {
		log.Error(ctx, "Error when process email send retry batches", err)
		return
	}
	log.Info(ctx, "Email send retry cron job completed successfully")
}

func NewEmailSendRetryCronJob(
	emailSendRetryUsecase usecase.IEmailSendRetryUsecase,
	redisPort port.IRedisPort,
) golibcron.Job {
	return &EmailSendRetryCronJob{
		emailSendRetryUsecase: emailSendRetryUsecase,
		redisPort:             redisPort,
	}
}
