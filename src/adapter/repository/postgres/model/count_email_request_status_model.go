package model

type CountEmailRequestStatusModel struct {
	Status string `json:"status"`
	Total  int64  `json:"total"`
}
