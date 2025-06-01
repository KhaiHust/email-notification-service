package dto

type SendVolumeDTO struct {
	TotalSend           int64
	TotalSendByProvider map[string]int64
	Date                int64
}
type SendVolumeByProviderDto struct {
	ProviderID int64
	Provider   string
	Total      int64
	TotalError int64
	TotalSent  int64
}
