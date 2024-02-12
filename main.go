package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/invest-scraping/api"
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

	registerRoutes(cfg, &conn)

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
func registerRoutes(cfg *config.Config, conn *mongodb.MongoConnection) {
	r := gin.Default()
	r.Use(Cors())
	// r.LoadHTMLGlob("./web/templates/*")
	routes := api.NewRoutes(r.RouterGroup, conn, cfg)

	routes.StockPriceRoutes()
	routes.AssetRoutes()
	r.Run(":" + getDefaultPort(cfg))
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Content-Type, X-Auth-Token, Content-Length, Accept-Encoding,X-CSRF-Token, Authorization, Nonce, Nonce-Signature, Address, Access-Control-Allow-Origin, Set-Cookie")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, HEAD, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusOK)
			return
		}

		c.Next()
	}
}

func getDefaultPort(cfg *config.Config) string {
	port := flag.String("port", cfg.Server.Port, "Instance port")
	flag.Parse()
	return *port
}
