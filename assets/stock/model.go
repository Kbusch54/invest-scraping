package stock

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	StockType string             `bson:"stockType"`
	Endpoint  string             `bson:"endpoint"`
	Symbol    string             `bson:"symbol"`
	LastPrice float64            `bson:"last_price"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (s *Stock) GetID() any {
	return s.ID
}

func (s *Stock) toResponse() *StockResponse {
	return &StockResponse{
		Name:      s.Name,
		StockType: s.StockType,
		Symbol:    s.Symbol,
		LastPrice: s.LastPrice,
		UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *Stock) UpdateLastPrice(price float64) {
	s.LastPrice = price
	s.UpdatedAt = time.Now()
}

func (s *Stock) NewStock(name, symbol, stockType, endpoint string, last_price float64) *Stock {
	return &Stock{
		Name:      name,
		StockType: stockType,
		Endpoint:  endpoint,
		Symbol:    symbol,
		LastPrice: last_price,
		UpdatedAt: time.Now(),
	}
}

func ToBatchResponse(stocks []Stock) []StockResponse {
	var responses []StockResponse
	for _, stock := range stocks {
		responses = append(responses, *stock.toResponse())
	}
	return responses
}
