package request

type GetListEmailTemplateFilter struct {
	WorkspaceID *int64
	Name        *string
	*BaseFilter
	EmailRequestFilter *EmailRequestFilter
}
