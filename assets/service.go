package assets

import (
	"strconv"

	"github.com/invest-scraping/assets/scraping"
	"github.com/invest-scraping/assets/stock"
	"github.com/invest-scraping/assets/stockprice"
	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
)

type Service interface {
	NewPriceInput(name, symbol string, price float64) error
	InitAssets(mon *[]config.Monitor)
	RunPricePerAssetUpdate(m *config.Monitor) error
}

type ServiceDefaultImpl struct {
	log      logg.Logger
	spSvc    stockprice.Service
	scrapers scraping.Service
	sSvc     stock.Service
}

func NewAssetsService(conn *mongodb.MongoConnection, cfg *config.Config) Service {
	log := logg.NewDefaultLog()
	spSvc := stockprice.NewStockPriceService(conn)
	sSvc := stock.NewStockService(conn)
	scrapers := scraping.NewScraperService(cfg.Monitors)
	return &ServiceDefaultImpl{
		log:      log,
		spSvc:    spSvc,
		sSvc:     sSvc,
		scrapers: scrapers,
	}
}

func (s *ServiceDefaultImpl) NewPriceInput(name, symbol string, price float64) error {
	stock, err := s.sSvc.GetStockByName(name)
	if err != nil {
		s.log.Error("Error finding stock. Reason: ", err.Error())
		return err
	}
	stock.UpdateLastPrice(price)
	err = s.sSvc.UpdateStock(stock)
	if err != nil {
		s.log.Error("Error creating stock. Reason: ", err.Error())
		return err
	}
	return s.spSvc.CreateStockPrice(name, symbol, price)
}

func (s *ServiceDefaultImpl) InitAssets(mon *[]config.Monitor) {
	for _, m := range *mon {
		err := s.sSvc.CreateOrUpdateStock(m.Name, m.Symbol, m.Type, m.Endpoint, 0)
		if err != nil {
			s.log.Error("Error creating stock. Reason: ", err.Error())
		}
		err = s.RunPricePerAssetUpdate(&m)
		if err != nil {
			s.log.Error("Error running price update. Reason: ", err.Error())
		}
	}
}

func (s *ServiceDefaultImpl) RunPricePerAssetUpdate(m *config.Monitor) error {
	s.log.Info("Running full price update for: ", m.Name)
	priceStr, err := s.scrapers.RunPriceUpdate(m.Symbol)
	if err != nil {
		s.log.Error("Error running price update. Reason: ", err.Error())
		return err
	}
	parsedPrice, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		s.log.Error("Error parsing price. Reason: ", err.Error())
		return err
	}
	err = s.NewPriceInput(m.Name, m.Symbol, parsedPrice)
	if err != nil {
		s.log.Error("Error creating stock price. Reason: ", err.Error())
		return err
	}

	return nil
}
