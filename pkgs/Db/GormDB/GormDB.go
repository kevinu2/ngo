package GormDB

import (
	"errors"
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var g *Gorm

type Gorm struct {
	GormDB *gorm.DB
	Config *GormConfig
}

func init() {
	g = New()
}

func New() *Gorm {
	return new(Gorm)
}

func AddConfig(dbType, dbUser, dbPass, dbHost string, dbPort int, dbTime, dbName string, dbMaxIdle, dbMaxOpen, dbMaxLifeTime int) {
	g.AddConfig(dbType, dbUser, dbPass, dbHost, dbPort, dbTime, dbName, dbMaxIdle, dbMaxOpen, dbMaxLifeTime)
}
func (g *Gorm) AddConfig(dbType, dbUser, dbPass, dbHost string, dbPort int, dbTime, dbName string, dbMaxIdle, dbMaxOpen, dbMaxLifeTime int) {
	g.Config = &GormConfig{
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

func GetDB() *gorm.DB { return g.GetDB() }
func (g *Gorm) GetDB() *gorm.DB {
	if g.GormDB == nil {
		fmt.Printf("Gorm: initDB()!")
		g.initDB()
	}
	return g.GormDB
}

func (g *Gorm) initDB() {
	var (
		db  *gorm.DB
		err error
	)

	switch g.Config.DbType {
	case DbPostgres.GetType():
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s  sslmode=disable TimeZone=%s dbname=%s", g.Config.DbHost, g.Config.DbPort, g.Config.DbUser, g.Config.DbPass, "UTC", g.Config.DbName)
		db, err = gorm.Open(
			postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			}))
		if err != nil {
			panic("Failed to connect to DB: " + err.Error())
		}
	case DbMySQL.GetType():
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
			panic("Failed to connect to DB: " + err.Error())
		}
	case DbClickHouse.GetType():
		dsn := fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20", g.Config.DbHost, g.Config.DbPort, g.Config.DbName, g.Config.DbUser, g.Config.DbPass)
		db, err = gorm.Open(
			clickhouse.Open(dsn),
			&gorm.Config{},
		)
	default:
		panic(errors.New("Unsupported DB Type: " + g.Config.DbType))
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(g.Config.DbMaxIdle)
	sqlDb.SetMaxOpenConns(g.Config.DbMaxOpen)
	sqlDb.SetConnMaxLifetime(time.Minute * time.Duration(g.Config.DbMaxLifeTime))

	g.GormDB = db
}
