package entity

type Conversion struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from" example:"usd"`
	To     string  `json:"to" example:"bitcoin"`
}
