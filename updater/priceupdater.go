package updater

import (
	"github.com/invest-scraping/assets"
	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
	"github.com/invest-scraping/scheduler"
)

type PriceUpdaterSchedulerImpl struct {
	aSvc     assets.Service
	monitors []*config.Monitor
	log      logg.Logger
}

func NewPriceUpdaterScheduler(conn *mongodb.MongoConnection, cfg *config.Config) scheduler.Scheduler {
	var (
		log      = logg.NewDefaultLog()
		monitors = cfg.Monitors
		aSvc     = assets.NewAssetsService(conn, cfg)
	)
	return &PriceUpdaterSchedulerImpl{
		aSvc:     aSvc,
		monitors: monitors,
		log:      log,
	}
}
func (s *PriceUpdaterSchedulerImpl) Execute() {
	for _, m := range s.monitors {
		err := s.aSvc.RunPricePerAssetUpdate(m)
		if err != nil {
			s.log.Error("Error updating price for asset. Reason: ", err.Error())
		}
	}
}

func (s *PriceUpdaterSchedulerImpl) Expression() string {
	return "@every 5m"
}
