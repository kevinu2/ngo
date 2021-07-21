package Db

import (
	"github.com/aiscrm/redisgo"
	"github.com/jmoiron/sqlx"
	"github.com/xormplus/xorm"
	"gorm.io/gorm"
)

type Pool struct {
	Redis *redisgo.Cacher
	SqlX  *sqlx.DB
	Gorm  *gorm.DB
	Xorm  *xorm.Engine
}

func InitPool() *Pool {
	return &Pool{
		Redis: GetRedisDB(),
		SqlX:  GetSqlX(),
		Gorm:  GetGormDB(),
		Xorm:  GetXormDB(),
	}
}
