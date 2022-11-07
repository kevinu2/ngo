package RedisDB

import (
	"encoding/json"
	"fmt"
	"github.com/kevinu2/redisgo"
)

var r *DB

type DB struct {
	RedisDB *redisgo.Cacher
	Config  *RedisConfig
}

func init() {
	r = New()
}

func New() *DB {
	return new(DB)
}

func AddConfig(dbUser, dbPass, dbHost string, dbPort int, dbName, dbMaxIdle, dbMaxActive, dbIdleTimeout int) {
	r.AddConfig(dbUser, dbPass, dbHost, dbPort, dbName, dbMaxIdle, dbMaxActive, dbIdleTimeout)
}
func (r *DB) AddConfig(dbUser, dbPass, dbHost string, dbPort int, dbName, dbMaxIdle, dbMaxActive, dbIdleTimeout int) {
	r.Config = &RedisConfig{
		DbUser:        dbUser,
		DbPass:        dbPass,
		DbHost:        dbHost,
		DbPort:        dbPort,
		DbName:        dbName,
		DbMaxIdle:     dbMaxIdle,
		DbMaxActive:   dbMaxActive,
		DbIdleTimeout: dbIdleTimeout,
	}
}

func GetDB() *redisgo.Cacher { return r.GetDB() }
func (r *DB) GetDB() *redisgo.Cacher {
	if r.RedisDB == nil {
		fmt.Print("Redis: initDB()! \n")
		r.initDB()
	}
	return r.RedisDB
}

func (r *DB) initDB() {
	c, err := redisgo.New(
		redisgo.Options{
			Network:     "tcp",
			Addr:        fmt.Sprintf("%s:%d", r.Config.DbHost, r.Config.DbPort),
			Password:    r.Config.DbPass,
			Db:          r.Config.DbName,
			MaxActive:   r.Config.DbMaxActive,
			MaxIdle:     r.Config.DbMaxIdle,
			IdleTimeout: r.Config.DbIdleTimeout,
			Marshal:     json.Marshal,
			Unmarshal:   json.Unmarshal,
		})
	if err != nil {
		fmt.Printf("redisgo.New() fails, err: %v", err.Error())
		panic(err)
	}
	r.RedisDB = c
}
