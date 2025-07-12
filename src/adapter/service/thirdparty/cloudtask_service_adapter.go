package thirdparty

import (
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"context"
	"encoding/json"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	coreProperties "github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golibs-starter/golib/log"
	"math/rand"
	"strings"
)

type CloudTaskServiceAdapter struct {
	taskProps         *properties.GoogleCloudTaskProperties
	appConfig         *coreProperties.AppProperties
	googleCloudClient *cloudtasks.Client
}

func (c CloudTaskServiceAdapter) CreateNewTask(ctx context.Context, request *request.CreateTaskDto) error {
	queueName := c.getRandomQueueName()
	body, err := json.Marshal(request.Payload)
	if err != nil {
		log.Error(ctx, "json marshal error: ", err)
		return err
	}
	parent := "projects/" + c.appConfig.GoogleCloudProject + "/locations/" + c.taskProps.GCPTaskLocation + "/queues/" + queueName
	taskNamePath := parent + "/tasks/" + request.TaskName
	req := &taskspb.CreateTaskRequest{
		Parent: parent,
		Task: &taskspb.Task{
			Name:         taskNamePath,
			ScheduleTime: c.convertTimestamp(request.ScheduleTime),
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod(request.MethodCode),
					Url:        request.TargetUrl,
					Headers:    request.Headers,
					Body:       body,
				},
			},
		},
	}
	_, err = c.googleCloudClient.CreateTask(ctx, req)
	if err != nil {
		log.Error(ctx, "create task error: ", err)
		return err
	}
	log.Info(ctx, "create task success: ", request.TaskName)
	return nil
}
func (c CloudTaskServiceAdapter) convertTimestamp(scheduleTime *int64) *timestamp.Timestamp {
	if scheduleTime == nil {
		return nil
	}
	return &timestamp.Timestamp{
		Seconds: *scheduleTime,
	}
}
func (c CloudTaskServiceAdapter) getRandomQueueName() string {
	queues := c.taskProps.CloudTaskQueues
	queueList := strings.Split(queues, ",")
	if len(queueList) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(queueList))
	return queueList[randomIndex]
}
func NewCloudTaskServiceAdapter(appConfig *coreProperties.AppProperties,
	taskProps *properties.GoogleCloudTaskProperties) port.ICloudTaskServicePort {
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Error("Cloudtask initialize app backend with service account error", err)
		panic(err)
	}
	log.Info("cloudtask initialized")
	return &CloudTaskServiceAdapter{
		appConfig:         appConfig,
		taskProps:         taskProps,
		googleCloudClient: client,
	}
}
