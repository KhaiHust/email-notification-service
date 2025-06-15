package request

type EmailRequestFilter struct {
	WorkspaceIDs     []int64
	EmailTemplateIDs []int64
	Statuses         []string
	RequestID        *string
	*BaseFilter
	SentAtFrom *int64
	SentAtTo   *int64
	Email      *string
}
