package MongoDB

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var m *Mongo

func init() {
	m = New()
}

type Mongo struct {
	Client *mongo.Client
	Config Config
}

func New() *Mongo {
	return new(Mongo)
}

func GetDB() *mongo.Client { return m.GetDB() }
func (m *Mongo) GetDB() *mongo.Client {
	if m.Client == nil {
		fmt.Print("Mongo: initDB()! \n")
		m.initDB()
	}
	return m.Client
}

// InitDB
//
//	@Description: 初始化mongo
//	@param mongoUrl
//	@return error
func (m *Mongo) initDB() *mongo.Client {
	c := m.Config
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=%s", c.User, c.Password, c.Host, c.Port, c.Auth)
	ctx, cancel := context.WithTimeout(context.Background(), c.TimeOut*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url).SetMaxPoolSize(c.MaxPoolSize))
	if err != nil {
		panic(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return client
}

// CreateMongoCollection
//
//	@Description: 创建mongo集合的服务
//	@param dbName
//	@param colName
//	@return BaseCollection
//	@return error
func (m *Mongo) CreateMongoCollection(dbName, colName string) BaseCollection {
	dataBase := m.Client.Database(dbName)
	return &BaseCollectionImpl{
		DbName:     dbName,
		ColName:    colName,
		DataBase:   dataBase,
		Collection: dataBase.Collection(colName),
	}
}
