package models

type CommonResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
