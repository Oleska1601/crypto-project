package entity

type Response struct {
	ID     int     `json:"id,omitempty"`
	UserID int     `json:"user_id,omitempty"`
	BTC    float64 `json:"btc"`
	ETH    float64 `json:"eth"`
	LTC    float64 `json:"ltc"`
	USDT   float64 `json:"usdt"`
}
