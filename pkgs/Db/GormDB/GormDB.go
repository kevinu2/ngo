package GormDB

import (
	"errors"
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
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

func GetDB() *gorm.DB { return g.GetDB() }
func (g *Gorm) GetDB() *gorm.DB {
	if g.GormDB == nil {
		fmt.Print("Gorm: initDB()! \n")
		g.initDB()
	}
	return g.GormDB
}

func (g *Gorm) initDB() {
	var (
		db  *gorm.DB
		err error
	)

	if g.Config == nil {
		panic(errors.New("Gorm.Config is nil"))
	}
	_, err = time.LoadLocation(g.Config.DbTimeZone)
	if err != nil {
		g.Config.DbTimeZone = "UTC"
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,                        // Ignore ErrRecordNotFound error for logger
			LogLevel:                  LogLevel(g.Config.LogLevel), // Log level
		},
	)
	switch g.Config.DbType {
	case DbPostgres.GetType():
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s  sslmode=disable TimeZone=%s dbname=%s", g.Config.DbHost, g.Config.DbPort, g.Config.DbUser, g.Config.DbPass, g.Config.DbTimeZone, g.Config.DbName)
		db, err = gorm.Open(
			postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			}),
			&gorm.Config{
				Logger: newLogger,
			},
		)
		if err != nil {
			panic("Failed to connect to DB: " + err.Error())
		}
	case DbMySQL.GetType():
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s", g.Config.DbUser, g.Config.DbPass, g.Config.DbHost, g.Config.DbPort, g.Config.DbName, g.Config.DbTimeZone)
		db, err = gorm.Open(

			mysql.New(mysql.Config{
				DSN:                       dsn,
				DefaultStringSize:         256,
				DisableDatetimePrecision:  true,
				DontSupportRenameIndex:    true,
				DontSupportRenameColumn:   true,
				SkipInitializeWithVersion: false,
			}),
			&gorm.Config{
				Logger: newLogger,
			},
		)
		if err != nil {
			panic("Failed to connect to DB: " + err.Error())
		}
	case DbClickHouse.GetType():
		dsn := fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20", g.Config.DbHost, g.Config.DbPort, g.Config.DbName, g.Config.DbUser, g.Config.DbPass)
		db, err = gorm.Open(
			clickhouse.Open(dsn),
			&gorm.Config{
				Logger: newLogger,
			},
		)
	default:
		panic(errors.New("Unsupported DB Type: " + g.Config.DbType))
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic("Failed to get DB instance: " + err.Error())
	}
	sqlDb.SetMaxIdleConns(g.Config.DbMaxIdle)
	sqlDb.SetMaxOpenConns(g.Config.DbMaxOpen)
	sqlDb.SetConnMaxLifetime(time.Minute * time.Duration(g.Config.DbMaxLifeTime))

	g.GormDB = db
}
