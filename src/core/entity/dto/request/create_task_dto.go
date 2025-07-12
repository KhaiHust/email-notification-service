package request

type CreateTaskDto struct {
	TaskName     string
	ScheduleTime *int64
	TargetUrl    string
	Headers      map[string]string
	Payload      interface{}
	MethodCode   int32
}
