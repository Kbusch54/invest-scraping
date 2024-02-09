package stock

type StockResponse struct {
	Name      string  `json:"name"`
	StockType string  `json:"stockType"`
	Symbol    string  `json:"symbol"`
	LastPrice float64 `json:"last_price"`
	UpdatedAt string  `json:"updated_at"`
}
