package stock

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/invest-scraping/api/httputil"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type ControllerImpl struct {
	svc Service
	log logg.Logger
}

type Controller interface {
	GetCurrentStockPrice(c *gin.Context)
}

const (
	DefaultPaginationSize = 20
	DefaultDsc            = true
)

func NewController(conn *mongodb.MongoConnection) Controller {
	log := logg.NewDefaultLog()
	servSvc := NewStockService(conn)
	return &ControllerImpl{
		svc: servSvc,
		log: log,
	}
}

func (ctrl *ControllerImpl) GetCurrentStockPrice(c *gin.Context) {
	stockName := c.Param("name")
	if stockName == "" {
		c.JSON(http.StatusBadRequest, httputil.BadRequestError("Invalid stock name"))
		return
	}
	stock, err := ctrl.svc.GetStockByName(stockName)
	if err != nil {
		ctrl.log.Error("Error getting stock data: ", err)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, httputil.BadRequestError("Stock not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.InternalServerError("Error getting stock data"))
		return
	}
	c.JSON(http.StatusOK, httputil.OK("Stock Data", stock))
	return
}
