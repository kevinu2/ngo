package Db

import (
	"github.com/aiscrm/redisgo"
	"github.com/jmoiron/sqlx"
)

type Pool struct {
	Redis *redisgo.Cacher
	SqlX  *sqlx.DB
}

func InitPool() *Pool {
	return &Pool{
		Redis: GetRedisDB(),
		SqlX:  GetSqlX(),
	}
}
