package stockprice

import (
	"time"

	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
)

type Service interface {
	CreateStockPrice(name, symbol string, last_price float64) error
	GetPricesSince(name string, since time.Time) ([]StockPriceResponse, error)
	FindByName(name string) (StockPriceResponse, error)
}

type ServiceDefaultImpl struct {
	repo Repository
	log  logg.Logger
}

func NewStockPriceService(conn *mongodb.MongoConnection) Service {
	log := logg.NewDefaultLog()
	repos := NewMongoRepository(conn)
	return &ServiceDefaultImpl{
		repo: repos,
		log:  log,
	}
}

func (s *ServiceDefaultImpl) CreateStockPrice(name, symbol string, last_price float64) error {
	stockPrice := &StockPrice{}
	stockPrice.NewStockPrice(name, symbol, last_price)
	return s.repo.CreateStockPrice(stockPrice)

}

func (s *ServiceDefaultImpl) GetPricesSince(name string, since time.Time) ([]StockPriceResponse, error) {
	stockPrices, err := s.repo.FindStockPriceBySymbol(name, since)
	if err != nil {
		s.log.Error("Error finding stock prices. Reason: ", err.Error())
		return nil, err
	}

	return ToBatchResponse(*stockPrices), nil
}

func (s *ServiceDefaultImpl) FindByName(name string) (StockPriceResponse, error) {
	stockPrice, err := s.repo.FindLatestStockPrice(name)
	if err != nil {
		s.log.Error("Error finding stock price. Reason: ", err.Error())
		return StockPriceResponse{}, err
	}

	return *stockPrice.toResponse(), nil
}
