package accrual

type AccrualResponseDTO struct {
	Order   int64   `json:"account"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
