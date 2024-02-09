package scraping

import (
	"errors"
	"time"

	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/scrape/goquery"
)

var (
	ErrScraperNotFound = errors.New("Scraper not found")
)

type Service interface {
	RunPriceUpdate(key string) (string, error)
	GetLastUpdatedBySymbol(symbol string) (time.Time, error)
}

type ServiceDefaultImpl struct {
	log      logg.Logger
	scrapers map[string]*goquery.Scraper
}

func NewScraperService(mon []*config.Monitor) Service {
	log := logg.NewDefaultLog()
	scrapers := make(map[string]*goquery.Scraper)
	for _, m := range mon {
		key := m.Symbol // or m.Name
		scrapers[key] = goquery.NewScraper(m.Endpoint+m.EndpointExt, m.PriceXpath)
	}
	return &ServiceDefaultImpl{
		log:      log,
		scrapers: scrapers,
	}
}

func (s *ServiceDefaultImpl) RunPriceUpdate(key string) (string, error) {
	scraper, ok := s.scrapers[key]
	if !ok {
		return "", ErrScraperNotFound
	}
	return scraper.RunQuery()
}

func (s *ServiceDefaultImpl) GetLastUpdatedBySymbol(symbol string) (time.Time, error) {
	scraper, ok := s.scrapers[symbol]
	if !ok {
		return time.Time{}, ErrScraperNotFound
	}
	return scraper.LastUpdated(), nil
}
