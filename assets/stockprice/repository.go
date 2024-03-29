package stockprice

import (
	"context"
	"time"

	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence"
	"github.com/invest-scraping/persistence/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COLLECTION = "stock.price"

type MongoRepository struct {
	conn    *mongodb.MongoConnection
	absrepo *persistence.AbstractMongoRepository[*StockPrice]
	log     logg.Logger
}
type Repository interface {
	FindLatestStockPrice(symbol string) (*StockPrice, error)
	FindStockPriceByName(name string, since time.Time) (*[]StockPrice, error)
	CreateStockPrice(stockPrice *StockPrice) error
}

func NewMongoRepository(conn *mongodb.MongoConnection) Repository {
	log := logg.NewDefaultLog()
	absrepo := persistence.NewAbstractRepository[*StockPrice](conn, COLLECTION)
	return &MongoRepository{
		conn:    conn,
		log:     log,
		absrepo: absrepo,
	}
}

func (r *MongoRepository) FindLatestStockPrice(symbol string) (*StockPrice, error) {
	stockPrice := &StockPrice{}
	filter := bson.M{"symbol": symbol}

	// Sort by 'time' in descending order (-1) and limit to 1 document
	opts := options.FindOne().SetSort(bson.D{{"time", -1}})

	err := r.conn.Datastore.Collection(r.absrepo.Collection).FindOne(context.Background(), filter, opts).Decode(stockPrice)
	if err != nil {
		r.log.Error("FindLatestStockPrice Error finding stock price. Reason: ", err.Error())
		return nil, err
	}

	return stockPrice, nil
}

func (r *MongoRepository) FindStockPriceByName(name string, since time.Time) (*[]StockPrice, error) {
	var stockPrices []StockPrice
	// filter := bson.M{"name": name, "time": bson.M{"$gte": since}}
	filter := bson.M{
		"name": name,
		"time": bson.M{
			"$gte": since,
		},
	}
	sort := bson.M{"time": -1}
	opt := options.Find().SetSort(sort)
	r.log.Info("Filter: ", filter)
	res, err := r.conn.Datastore.Collection(r.absrepo.Collection).Find(context.Background(), filter, opt)
	if err != nil {
		r.log.Error("FindStockPriceByName Error finding stock prices. Reason: ", err.Error())
		return nil, err
	}
	if err = res.All(context.Background(), &stockPrices); err != nil {
		r.log.Error("FindStockPriceByName Error finding stock prices. Reason: ", err.Error())
		return nil, err
	}
	return &stockPrices, nil
}

func (r *MongoRepository) CreateStockPrice(stockPrice *StockPrice) error {
	return r.absrepo.InsertOrUpdate(stockPrice)
}
