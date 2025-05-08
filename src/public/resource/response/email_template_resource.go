package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type EmailTemplateResponse struct {
	Id          int64                        `json:"id"`
	Name        string                       `json:"name"`
	WorkspaceID int64                        `json:"workspace_id"`
	CreatedAt   int64                        `json:"created_at"`
	UpdatedAt   int64                        `json:"updated_at"`
	Subject     string                       `json:"subject,omitempty"`
	Body        string                       `json:"body,omitempty"`
	Variables   interface{}                  `json:"variables,omitempty"`
	Metric      *EmailTemplateMetricResponse `json:"metric,omitempty"`
}
type EmailTemplateMetricResponse struct {
	TotalSent   int64 `json:"total_sent"`
	TotalErrors int64 `json:"total_errors"`
}

func ToEmailTemplateResponse(templateEntity *entity.EmailTemplateEntity) *EmailTemplateResponse {

	template := &EmailTemplateResponse{
		Id:          templateEntity.ID,
		Name:        templateEntity.Name,
		WorkspaceID: templateEntity.WorkspaceId,
		CreatedAt:   templateEntity.CreatedAt,
		UpdatedAt:   templateEntity.UpdatedAt,
		Subject:     templateEntity.Subject,
		Body:        templateEntity.Body,
		Variables:   templateEntity.Variables,
	}
	if templateEntity.Metric != nil {
		template.Metric = &EmailTemplateMetricResponse{
			TotalSent:   templateEntity.Metric.TotalSent,
			TotalErrors: templateEntity.Metric.TotalErrors,
		}
	}
	return template
}
func ToListEmailTemplateResponse(templates []*entity.EmailTemplateEntity) []*EmailTemplateResponse {
	res := make([]*EmailTemplateResponse, len(templates))
	for i, template := range templates {
		res[i] = ToEmailTemplateResponse(template)
	}
	return res
}
