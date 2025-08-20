package Elasticsearch

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch"
)

var e *Elasticsearch

type Elasticsearch struct {
	ElasticDB *elasticsearch.Client
	Config    *Config
}

func init() {
	e = New()
}

func New() *Elasticsearch {
	return new(Elasticsearch)
}

func AddConfig(hosts []string) {
	e.AddConfig(hosts)
}
func (e *Elasticsearch) AddConfig(hosts []string) {
	e.Config = &Config{
		Addresses: hosts,
	}
}

func GetDB() *elasticsearch.Client { return e.GetDB() }

func (e *Elasticsearch) GetDB() *elasticsearch.Client {
	if e.ElasticDB == nil {
		fmt.Printf("ElasticDB: initDB()!")
		e.initDB()
	}
	return e.ElasticDB
}

func (e *Elasticsearch) initDB() {
	var (
		db  *elasticsearch.Client
		err error
	)
	cfg := elasticsearch.Config{
		Addresses: e.Config.Addresses,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 60 * time.Second,
			DisableKeepAlives:     true,
		},
	}
	if db, err = elasticsearch.NewClient(cfg); err != nil {
		panic("Failed to connect to Elasticsearch: " + err.Error())
	}
	e.ElasticDB = db
}
