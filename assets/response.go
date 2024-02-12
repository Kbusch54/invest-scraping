package assets

type StocksResponse struct {
	Name      string `json:"name"`
	StockType string `json:"stockType"`
	Symbol    string `json:"symbol"`
	StockData []StockPriceResponse
}

type StockPriceResponse struct {
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Time   string  `json:"time"`
}
