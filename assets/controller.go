package assets

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/invest-scraping/api/httputil"
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
	GetPricesByAssetName(c *gin.Context)
}

type StockRequest struct {
	Name  string `json:"name"`
	Since string `json:"since"`
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

func (ctrl *ControllerImpl) GetPricesByAssetName(c *gin.Context) {
	var req StockRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	ctrl.log.Info("Request: ", req)
	assetName := req.Name
	since := req.Since
	if assetName == "" {
		c.JSON(http.StatusBadRequest, httputil.BadRequestError("Invalid asset name"))
		return
	}
	//since string to time
	layout := "2006-01-02T15:04:05Z07:00"
	sinceTime, err := time.Parse(layout, since) // "2006-01-02T15:04:05Z07:00"
	if err != nil {
		ctrl.log.Warnf("Error parsing time. Reason: %v", err.Error())
		sinceTime = time.Now().AddDate(0, -1, 0)
	}
	asset, err := ctrl.svc.GetPricesSince(assetName, sinceTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.InternalServerError[string](err.Error()))
		return
	}
	c.JSON(http.StatusOK, httputil.OK("Asset prices", asset))
	return
}
