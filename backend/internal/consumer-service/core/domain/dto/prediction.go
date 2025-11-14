package dto

type PredictionRequest struct {
	Comment string `json:"comment"`
}

type PredictionResponse struct {
	IsToxic bool `json:"toxic"`
}
