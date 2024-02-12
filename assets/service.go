package assets

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/invest-scraping/assets/scraping"
	"github.com/invest-scraping/assets/stock"
	"github.com/invest-scraping/assets/stockprice"
	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
	"github.com/invest-scraping/stream"
	"github.com/invest-scraping/stream/kafka"
)

type Service interface {
	NewPriceInput(name, symbol string, price float64, time time.Time) error
	InitAssets(mon []*config.Monitor) error
	RunPricePerAssetUpdate(m *config.Monitor) error
	StreamPriceUpdate(name, symbol string, price float64, time time.Time) error
	ReplayPriceUpdateSince(name string, time time.Time) error
	GetPricesSince(name string, since time.Time) (StocksResponse, error)
}

type ServiceDefaultImpl struct {
	log      logg.Logger
	spSvc    stockprice.Service
	scrapers scraping.Service
	producer stream.Producer
	sSvc     stock.Service
}

func NewAssetsService(conn *mongodb.MongoConnection, cfg *config.Config) Service {
	log := logg.NewDefaultLog()
	spSvc := stockprice.NewStockPriceService(conn)
	sSvc := stock.NewStockService(conn)
	scrapers := scraping.NewScraperService(cfg.Monitors)
	producer := kafka.NewProducer(&cfg.Stream, EventTopic, log)
	return &ServiceDefaultImpl{
		log:      log,
		spSvc:    spSvc,
		sSvc:     sSvc,
		producer: producer,
		scrapers: scrapers,
	}
}

func (s *ServiceDefaultImpl) NewPriceInput(name, symbol string, price float64, time time.Time) error {
	stock, err := s.sSvc.GetStockByName(name)
	if err != nil {
		s.log.Error("Assets.Service Error finding stock. Reason: ", err.Error())
		return err
	}
	stock.UpdateLastPrice(price, time)
	err = s.sSvc.UpdateStock(stock)
	if err != nil {
		s.log.Error("Error creating stock. Reason: ", err.Error())
		return err
	}
	return s.spSvc.CreateStockPrice(name, symbol, price, time)
}

func (s *ServiceDefaultImpl) InitAssets(mon []*config.Monitor) error {
	for _, m := range mon {
		err := s.sSvc.CreateOrUpdateStock(m.Name, m.Symbol, m.Type, m.Endpoint, 0, time.Now())
		if err != nil {
			s.log.Error("Error creating stock. Reason: ", err.Error())
			return err
		}
		err = s.RunPricePerAssetUpdate(m)
		if err != nil {
			s.log.Error("Error running price update. Reason: ", err.Error())
			return err

		}
	}
	return nil
}

func (s *ServiceDefaultImpl) RunPricePerAssetUpdate(m *config.Monitor) error {
	s.log.Info("Running full price update for: ", m.Name)
	priceStr, err := s.scrapers.RunPriceUpdate(m.Symbol)
	if err != nil {
		s.log.Error("Error running price update. Reason: ", err.Error())
		return err
	}
	parsedPrice, err := strconv.ParseFloat(strings.ReplaceAll(priceStr, ",", ""), 64)
	if err != nil {
		s.log.Error("Error parsing price. Reason: ", err.Error())
		return err
	}
	lastUpdated, err := s.scrapers.GetLastUpdatedBySymbol(m.Symbol)
	if err != nil {
		s.log.Error("Error getting last updated. Reason: ", err.Error())
		return err
	}
	err = s.NewPriceInput(m.Name, m.Symbol, parsedPrice, lastUpdated)
	if err != nil {
		s.log.Error("Error creating stock price. Reason: ", err.Error())
		return err
	}

	return nil
}

func (s *ServiceDefaultImpl) StreamPriceUpdate(name, symbol string, price float64, time time.Time) error {
	evt := NewEvent(name, symbol, price, EventName, time)
	msg, err := json.Marshal(evt)
	if err != nil {
		s.log.Errorf("Error marshalling Price update json. Reason: ", err.Error())
		return err
	}
	err = s.producer.Produce(msg)
	if err != nil {
		s.log.Error("Error producing event. Reason: ", err.Error())
		return err
	}
	return nil
}

func (s *ServiceDefaultImpl) ReplayPriceUpdateSince(name string, since time.Time) error {
	stockPrices, err := s.spSvc.GetPricesSince(name, since)
	if err != nil {
		s.log.Error("Error finding stock prices. Reason: ", err.Error())
		return err
	}
	for _, sp := range stockPrices {
		newTime, err := time.Parse(time.RFC3339, sp.Time)
		if err != nil {
			s.log.Error("Error parsing time. Reason: ", err.Error())
			return err
		}
		err = s.StreamPriceUpdate(sp.Name, sp.Symbol, sp.Price, newTime)
	}
	return nil
}

func (s *ServiceDefaultImpl) GetPricesSince(name string, since time.Time) (StocksResponse, error) {
	stocksResp := StocksResponse{}
	stockData, err := s.sSvc.GetStockByName(name)
	if err != nil {
		s.log.Error("Error finding stock. Reason: ", err.Error())
		return stocksResp, err
	}
	stData, err := s.spSvc.GetPricesSince(name, since)
	if err != nil {
		s.log.Error("Error finding stock prices. Reason: ", err.Error())
		return stocksResp, err
	}
	stocksResp.Name = stockData.Name
	stocksResp.StockType = stockData.StockType
	stocksResp.Symbol = stockData.Symbol
	stResp := []StockPriceResponse{}
	for _, sp := range stData {
		stResp = append(stResp,
			StockPriceResponse{
				Name:   sp.Name,
				Symbol: sp.Symbol,
				Price:  sp.Price,
				Time:   sp.Time,
			})
	}
	stocksResp.StockData = stResp
	return stocksResp, nil

}
