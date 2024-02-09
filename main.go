package main

import (
	"fmt"

	"github.com/invest-scraping/assets"
	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
	"github.com/invest-scraping/scheduler"
	"github.com/invest-scraping/updater"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.Load("./config/env")

	logger := logg.NewDefaultLog()

	// create a database connection
	conn := mongodb.NewConnection(&cfg.Persistence, logger)

	// register monitors
	config.RegisterMonitors(cfg)

	err := registerAssets(&conn, cfg)
	if err != nil {
		logger.Error("Error registering assets. Reason: ", err.Error())
		panic(err)
	}
	go registerSchedulers(cfg, &conn)

}

func registerAssets(conn *mongodb.MongoConnection, cfg *config.Config) error {
	fmt.Println("Starting application")
	aSvc := assets.NewAssetsService(conn, cfg)
	return aSvc.InitAssets(cfg.Monitors)
}

func registerSchedulers(cfg *config.Config, conn *mongodb.MongoConnection) {
	c := cron.New()
	jobs := Schedulers(cfg, conn)
	for _, job := range jobs {
		_, _ = c.AddFunc(job.Expression(), job.Execute)
		// job.Execute()
	}
	c.Start()
}

func Schedulers(cfg *config.Config, conn *mongodb.MongoConnection) []scheduler.Scheduler {
	var (
		updatePrice = updater.NewPriceUpdaterScheduler(conn, cfg)
	)
	return []scheduler.Scheduler{
		updatePrice,
	}
}
