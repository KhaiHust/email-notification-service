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
