package api

import (
	"github.com/gin-gonic/gin"
	"github.com/invest-scraping/assets"
	"github.com/invest-scraping/assets/stock"
	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
)

type Routes struct {
	routerGroup       gin.RouterGroup
	cfg               *config.Config
	log               *logg.Logger
	conn              *mongodb.MongoConnection
	publicRouterGroup *gin.RouterGroup
}

func NewRoutes(routerGroup gin.RouterGroup, conn *mongodb.MongoConnection, cfg *config.Config) Routes {

	publicRouterGroup := routerGroup.Group("/public/investments")
	return Routes{
		routerGroup:       routerGroup,
		conn:              conn,
		cfg:               cfg,
		publicRouterGroup: publicRouterGroup,
	}
}

// func (r *Routes) OnBoardRoutes() error {
// 	onboarding := onboarding.NewOnboardingController(r.conn, *r.cfg)
// 	{
// 		r.publicRouterGroup.POST("/onboardUser", onboarding.OnboardUser)
// 	}
// 	return nil
// }

func (r *Routes) StockPriceRoutes() error {
	stock := stock.NewController(r.conn)
	{
		r.publicRouterGroup.GET("/stock/:name", stock.GetCurrentStockPrice)
		// r.publicRouterGroup.GET("/stocks", stock.FindAllStocks)
		// r.publicRouterGroup.GET("/stock/:name", stock.FindStockByName)
	}
	return nil

}

func (r *Routes) AssetRoutes() error {
	asset := assets.NewController(r.conn, r.cfg)
	// {
	r.publicRouterGroup.GET("/asset", asset.GetPricesByAssetName)
	// 	r.publicRouterGroup.GET("/assets", asset.FindAllAssets)
	// 	r.publicRouterGroup.GET("/asset/:name", asset.FindAssetByName)
	// }
	return nil
}
