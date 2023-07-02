package account

type OrderAccrualResponse struct {
	Order   string  `json:"account"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
