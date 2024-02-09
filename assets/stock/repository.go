package stock

import (
	"context"
	"errors"

	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/persistence"
	"github.com/invest-scraping/persistence/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

const COLLECTION = "stocks"

type MongoRepository struct {
	conn    *mongodb.MongoConnection
	absrepo *persistence.AbstractMongoRepository[*Stock]
	log     logg.Logger
}
type Repository interface {
	FindAll() (*[]Stock, error)
	FindByID(id string) (*Stock, error)
	UpdateOrInsert(stfxFollow *Stock) error
	FindByName(name string) (*Stock, error)
}

var (
	ErrRange = errors.New("requested page is out of range")
)

func NewMongoRepository(conn *mongodb.MongoConnection) Repository {
	log := logg.NewDefaultLog()
	absrepo := persistence.NewAbstractRepository[*Stock](conn, COLLECTION)
	return &MongoRepository{
		conn:    conn,
		log:     log,
		absrepo: absrepo,
	}
}

func (r *MongoRepository) FindAll() (*[]Stock, error) {
	var stocks []Stock
	res, err := r.conn.Datastore.Collection(r.absrepo.Collection).Find(context.Background(), bson.M{})
	if err != nil {
		r.log.Error("FindAll Error finding stocks. Reason: ", err.Error())
		return nil, err
	}
	if err = res.All(context.Background(), &stocks); err != nil {
		r.log.Error("FindAll Error finding stocks. Reason: ", err.Error())
		return nil, err
	}
	return &stocks, nil
}

func (r *MongoRepository) UpdateOrInsert(stock *Stock) error {
	err := r.absrepo.InsertOrUpdate(stock)
	if err != nil {
		r.log.Error("UpdateOrInsert Error updating or inserting stock. Reason: ", err.Error())
		return err
	}
	return nil
}

func (r *MongoRepository) FindByID(id string) (*Stock, error) {
	var stk Stock
	filter := bson.M{"_id": id}
	err := r.conn.Datastore.Collection(r.absrepo.Collection).FindOne(context.Background(), filter).Decode(&stk)
	if err != nil {
		r.log.Error("FindByID Error finding Stock. Reason: ", err.Error())
		return nil, err
	}
	return &stk, nil
}

func (r *MongoRepository) FindByName(name string) (*Stock, error) {
	var stk Stock
	filter := bson.M{"name": name}
	err := r.conn.Datastore.Collection(r.absrepo.Collection).FindOne(context.Background(), filter).Decode(&stk)
	if err != nil {
		r.log.Error("FindByName Error finding Stock. Reason: ", err.Error())
		return nil, err
	}
	return &stk, nil
}
