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

func (s *Stock) GetID() primitive.ObjectID {
	return s.ID
}
