package request

type GetListEmailTemplateFilter struct {
	WorkspaceID        *int64
	Name               *string
	Limit              *int64
	Since              *int64
	Until              *int64
	SortOrder          string
	CreatedAtFrom      *int64
	CreatedAtTo        *int64
	UpdatedAtFrom      *int64
	UpdatedAtTo        *int64
	EmailRequestFilter *EmailRequestFilter
}
