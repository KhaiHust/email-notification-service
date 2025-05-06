package request

import "github.com/KhaiHust/email-notification-service/core/entity/dto/request"

type GetListApiKeyRequest struct {
	WorkspaceID int64
	Environment []string
}

func NewGetListApiKeyFilter(params *GetListApiKeyRequest) *request.GetApiKeyRequestFilter {
	return &request.GetApiKeyRequestFilter{
		WorkspaceIDs: []int64{params.WorkspaceID},
		Environments: params.Environment,
	}
}
