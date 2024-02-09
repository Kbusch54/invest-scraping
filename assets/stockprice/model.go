package stockprice

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockPrice struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
	Symbol string             `bson:"symbol"`
	Price  float64            `bson:"price"`
	Time   time.Time          `bson:"time"`
}

func (s *StockPrice) GetID() any {
	return s.ID
}

func (s *StockPrice) toResponse() *StockPriceResponse {
	return &StockPriceResponse{
		Name:   s.Name,
		Symbol: s.Symbol,
		Price:  s.Price,
		Time:   s.Time.Format(time.RFC3339),
	}
}

func (s *StockPrice) NewStockPrice(name, symbol string, price float64, time time.Time) *StockPrice {
	return &StockPrice{
		ID:     primitive.NewObjectID(),
		Name:   name,
		Symbol: symbol,
		Price:  price,
		Time:   time,
	}
}

func ToBatchResponse(stockPrices []StockPrice) []StockPriceResponse {
	var responses []StockPriceResponse
	for _, stockPrice := range stockPrices {
		responses = append(responses, *stockPrice.toResponse())
	}
	return responses
}
