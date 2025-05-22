package request

type TemplateMetricFilter struct {
	TemplateID  int64  `json:"template_id"`
	WorkspaceID int64  `json:"workspace_id"`
	StartDate   *int64 `json:"start_date"`
	EndDate     *int64 `json:"end_date"`
}
