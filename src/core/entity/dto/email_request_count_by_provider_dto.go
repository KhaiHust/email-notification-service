package dto

type EmailRequestCountByProviderDTO struct {
	ProviderID int64
	Provider   string
	Total      int64
	Date       string
}
