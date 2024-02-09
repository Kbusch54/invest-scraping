package stock

import (
	"time"

	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	CreateOrUpdateStock(name, symbol, stockType, endpoint string, last_price float64) error
	FindAll() ([]StockResponse, error)
	FindByName(name string) (StockResponse, error)
}

type ServiceDefaultImpl struct {
	repo Repository
	log  logg.Logger
}

func NewStockService(conn *mongodb.MongoConnection) Service {
	log := logg.NewDefaultLog()
	repos := NewMongoRepository(conn)
	return &ServiceDefaultImpl{
		repo: repos,
		log:  log,
	}
}

func (s *ServiceDefaultImpl) CreateOrUpdateStock(name, symbol, stockType, endpoint string, last_price float64) error {
	stock := &Stock{}
	stock, err := s.repo.FindByName(name)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			stock.NewStock(name, symbol, stockType, endpoint, last_price)
			return s.repo.UpdateOrInsert(stock)
		}
		s.log.Error("Error finding stock. Reason: ", err.Error())
		return err
	}
	stock.LastPrice = last_price
	stock.UpdatedAt = time.Now()
	return s.repo.UpdateOrInsert(stock)

}

func (s *ServiceDefaultImpl) FindAll() ([]StockResponse, error) {
	stocks, err := s.repo.FindAll()
	if err != nil {
		s.log.Error("Error finding stocks. Reason: ", err.Error())
		return nil, err
	}

	return ToBatchResponse(*stocks), nil
}

func (s *ServiceDefaultImpl) FindByName(name string) (StockResponse, error) {
	stock, err := s.repo.FindByName(name)
	if err != nil {
		s.log.Error("Error finding stock. Reason: ", err.Error())
		return StockResponse{}, err
	}
	return *stock.toResponse(), nil
}
