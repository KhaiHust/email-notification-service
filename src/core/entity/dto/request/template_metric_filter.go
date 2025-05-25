package request

type TemplateMetricFilter struct {
	TemplateID  int64
	WorkspaceID int64
	StartDate   *int64
	EndDate     *int64
	Interval    string
	Duration    int
	IsChart     bool
}
