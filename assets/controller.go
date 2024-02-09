package assets

import (
	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
)

type ControllerImpl struct {
	svc Service
	log logg.Logger
}

type Controller interface {
	// GetCurrentStockPrice(c *gin.Context)
}

const (
	DefaultPaginationSize = 20
	DefaultDsc            = true
)

func NewController(conn *mongodb.MongoConnection, cfg *config.Config) Controller {
	log := logg.NewDefaultLog()
	servSvc := NewAssetsService(conn, cfg)
	return &ControllerImpl{
		svc: servSvc,
		log: log,
	}
}

// func (ctrl *ControllerImpl) GetCurrentStockPrice(c *gin.Context) {
// 	stockName := c.Param("name")
// 	if stockName == "" {
// 		c.JSON(http.StatusBadRequest, httputil.BadRequestError("Invalid stock name"))
// 		return
// 	}
// 	stock, err := ctrl.svc.(stockName)
// 	c.JSON(http.StatusOK, httputil.OK("FollowUser Success", true))
// 	return
// }
