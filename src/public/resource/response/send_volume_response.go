package response

import "github.com/KhaiHust/email-notification-service/core/entity/dto"

type SendVolumeResponse struct {
	MapSendVolume map[string]SendVolumeResponseDto `json:"map_send_volume"`
}
type SendVolumeResponseDto struct {
	TotalSend           int64            `json:"total_send"`
	TotalSendByProvider map[string]int64 `json:"total_send_by_provider"`
}

func ToSendVolumeResponse(send map[string]*dto.SendVolumeDTO) *SendVolumeResponse {
	response := &SendVolumeResponse{
		MapSendVolume: make(map[string]SendVolumeResponseDto),
	}
	for date, sendVolume := range send {
		response.MapSendVolume[date] = SendVolumeResponseDto{
			TotalSend:           sendVolume.TotalSend,
			TotalSendByProvider: sendVolume.TotalSendByProvider,
		}
	}
	return response
}

type SendVolumeByProviderResponse struct {
	ProviderID int64  `json:"provider_id"`
	Provider   string `json:"provider"`
	Total      int64  `json:"total"`
	TotalError int64  `json:"total_error"`
	TotalSent  int64  `json:"total_sent"`
}

func ToSendVolumeByProviderResponse(send []*dto.SendVolumeByProviderDto) []*SendVolumeByProviderResponse {
	var response []*SendVolumeByProviderResponse
	for _, v := range send {
		response = append(response, &SendVolumeByProviderResponse{
			ProviderID: v.ProviderID,
			Provider:   v.Provider,
			Total:      v.Total,
			TotalError: v.TotalError,
			TotalSent:  v.TotalSent,
		})
	}
	return response
}
