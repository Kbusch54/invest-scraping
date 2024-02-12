package updater

import (
	"sync"
	"time"

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
	startTime := time.Now()
	s.log.Info("Running scheduled update prices")
	var wg sync.WaitGroup
	for _, m := range s.monitors {
		wg.Add(1)
		go func(mon config.Monitor) {
			defer wg.Done()
			err := s.aSvc.RunPricePerAssetUpdate(&mon)
			if err != nil {
				s.log.Error("Error updating price for asset. Reason: ", err.Error())
			}
		}(*m)
	}
	wg.Wait()
	qualifier := "seconds"
	duration := time.Since(startTime).Seconds()
	if time.Since(startTime).Seconds() < 60 {
		duration = time.Since(startTime).Seconds()
	} else if time.Since(startTime).Seconds() > 60 {
		duration = time.Since(startTime).Minutes()
		qualifier = "minutes"
	} else if time.Since(startTime).Minutes() > 60 {
		duration = time.Since(startTime).Hours()
		qualifier = "hours"
	}
	processed := len(s.monitors)
	s.log.Infof("Finished running scheduled update asset prices. Processed: %v assets. Duration: %v %s", processed, duration, qualifier)

}

func (s *PriceUpdaterSchedulerImpl) Expression() string {
	return "@every 5m"
}
