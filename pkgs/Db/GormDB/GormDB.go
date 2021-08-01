package GormDB

import (
	"fmt"
	"github.com/kevinu2/ngo/enum"
	"github.com/kevinu2/ngo/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var g *Gorm

type Gorm struct {
	GormDB *gorm.DB
	Config *model.GormDbConfig
}

func init() {
	g = New()
}

func New() *Gorm {
	return new(Gorm)
}

func GetDB() *gorm.DB {
	if g.GormDB == nil {
		g.initDB()
	}
	return g.GormDB
}

func AddConfig(dbType, dbUser, dbPass, dbHost string, dbPort int, dbTime, dbName string, dbMaxIdle, dbMaxOpen, dbMaxLifeTime int) {
	g.AddConfig(dbType, dbUser, dbPass, dbHost, dbPort, dbTime, dbName, dbMaxIdle, dbMaxOpen, dbMaxLifeTime)
}
func (g *Gorm) AddConfig(dbType, dbUser, dbPass, dbHost string, dbPort int, dbTime, dbName string, dbMaxIdle, dbMaxOpen, dbMaxLifeTime int) {
	g.Config = &model.GormDbConfig{
		DbUser:        dbUser,
		DbPass:        dbPass,
		DbHost:        dbHost,
		DbPort:        dbPort,
		DbType:        dbType,
		DbName:        dbName,
		DbTimeZone:    dbTime,
		DbMaxIdle:     dbMaxIdle,
		DbMaxOpen:     dbMaxOpen,
		DbMaxLifeTime: dbMaxLifeTime,
	}
}

func (g *Gorm) initDB() {
	var (
		db  *gorm.DB
		err error
	)

	switch g.Config.DbType {
	case enum.DbPG.GetType():
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s  sslmode=disable TimeZone=%s dbname=%s", g.Config.DbHost, g.Config.DbPort, g.Config.DbUser, g.Config.DbPass, "UTC", g.Config.DbName)
		db, err = gorm.Open(
			postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			}))
		if err != nil {
			panic("连接数据库失败:" + err.Error())
		}
	case enum.DbMySQL.GetType():
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s", g.Config.DbUser, g.Config.DbPass, g.Config.DbHost, g.Config.DbPort, g.Config.DbName, "UTC")
		db, err = gorm.Open(

			mysql.New(mysql.Config{
				DSN:                       dsn,
				DefaultStringSize:         256,
				DisableDatetimePrecision:  true,
				DontSupportRenameIndex:    true,
				DontSupportRenameColumn:   true,
				SkipInitializeWithVersion: false,
			}),
			&gorm.Config{},
		)
		if err != nil {
			panic("连接数据库失败:" + err.Error())
		}
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(g.Config.DbMaxIdle)
	sqlDb.SetMaxOpenConns(g.Config.DbMaxOpen)
	sqlDb.SetConnMaxLifetime(time.Hour * time.Duration(g.Config.DbMaxLifeTime))

	g.GormDB = db
}
