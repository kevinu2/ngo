package RedisDB

import (
	"encoding/json"
	"fmt"
	"github.com/aiscrm/redisgo"
	"github.com/spf13/viper"
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

func AddConfig(dbType, dbUser, dbPass, dbHost string, dbPort int, dbTime string, dbName, dbMaxIdle, dbMaxActive, dbIdleTimeout int) {
	r.AddConfig(dbType, dbUser, dbPass, dbHost, dbPort, dbTime, dbName, dbMaxIdle, dbMaxActive, dbIdleTimeout)
}
func (r *DB) AddConfig(dbType, dbUser, dbPass, dbHost string, dbPort int, dbTime string, dbName, dbMaxIdle, dbMaxActive, dbIdleTimeout int) {
	r.Config = &RedisConfig{
		DbUser:        dbUser,
		DbPass:        dbPass,
		DbHost:        dbHost,
		DbPort:        dbPort,
		DbType:        dbType,
		DbName:        dbName,
		DbTimeZone:    dbTime,
		DbMaxIdle:     dbMaxIdle,
		DbMaxActive:   dbMaxActive,
		DbIdleTimeout: dbIdleTimeout,
	}
}

func GetDB() *redisgo.Cacher { return r.GetDB() }
func (r *DB) GetDB() *redisgo.Cacher {
	if r.RedisDB == nil {
		fmt.Printf("Redis: initDB()!")
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
			IdleTimeout: viper.GetInt("redis.pool.idle-timeout"),
			Prefix:      viper.GetString("redis.prefix"),
			Marshal:     json.Marshal,
			Unmarshal:   json.Unmarshal,
		})
	if err != nil {
		fmt.Printf("redisgo.New() fails, err: %v", err.Error())
		panic(err)
	}
	r.RedisDB = c
}
