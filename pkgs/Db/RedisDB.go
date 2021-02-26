package Db

import (
	"encoding/json"
	"fmt"
	"github.com/aiscrm/redisgo"
	"github.com/spf13/viper"
)

var Redis *redisgo.Cacher

func GetRedisDB() *redisgo.Cacher {
	if Redis == nil {
		initRedisDB()
	}
	return Redis
}

func initRedisDB() {
	c, err := redisgo.New(
		redisgo.Options{
			Addr:        viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
			Password:    viper.GetString("redis.password"),
			Db:          viper.GetInt("redis.database"),
			MaxActive:   viper.GetInt("redis.pool.max-active"),
			MaxIdle:     viper.GetInt("redis.pool.max-idle"),
			IdleTimeout: viper.GetInt("redis.pool.idle-timeout"),
			Prefix:      viper.GetString("redis.prefix"),
			Marshal:     json.Marshal,
			Unmarshal:   json.Unmarshal,
		})
	if err != nil {
		fmt.Printf("redisgo.New() fails, err: %v", err.Error())
		panic(err)
	}
	Redis = c
}
