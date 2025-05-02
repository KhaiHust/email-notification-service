package request

type EmailRequestFilter struct {
	EmailTemplateIDs []int64
	Statuses         []string
	CreatedAtFrom    *int64
	CreatedAtTo      *int64
	SentAtFrom       *int64
	SentAtTo         *int64
}
