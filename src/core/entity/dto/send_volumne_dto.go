package dto

type SendVolumeDTO struct {
	TotalSend           int64
	TotalSendByProvider map[string]int64
	Date                int64
}
