package mongodb

import (
	"context"

	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DefaultUrl = "mongodb+srv://andrewvb2012:U3aQdGnoYzznB1d2@cluster0.gmjmvad.mongodb.net/"
)

type MongoConnection struct {
	Session   *mongo.Client
	Datastore *mongo.Database
	cfg       *config.Persistence
	log       logg.Logger
}

func NewConnection(cfg *config.Persistence, log logg.Logger) MongoConnection {
	cli := MongoConnection{cfg: cfg, log: log}
	cli.connect()
	return cli
}

func (con *MongoConnection) connect() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(getUrlConnection(con.cfg)))
	if err != nil {
		con.log.Panic("Error connecting to mongodb. Reason: " + err.Error())
	}
	if defaultDB := getDefaultDatabase(con.cfg); len(defaultDB) > 0 {
		con.Datastore = client.Database(defaultDB)
	}
	con.Session = client
}

func getDefaultDatabase(cfg *config.Persistence) string {
	if len(cfg.MongoDB.Database) > 0 {
		return cfg.MongoDB.Database
	} else {
		return ""
	}
}

func getUrlConnection(cfg *config.Persistence) string {
	if len(cfg.MongoDB.URL) > 0 {
		return cfg.MongoDB.URL
	} else {
		return DefaultUrl
	}
}
