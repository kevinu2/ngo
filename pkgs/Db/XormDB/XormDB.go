package XormDB

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"github.com/xormplus/xorm/log"
	"github.com/xormplus/xorm/names"
	"time"
)

var x *Xorm

type Xorm struct {
	XormDB *xorm.Engine
	Config *Config
}

func init() {
	x = New()
}

func New() *Xorm {
	return new(Xorm)
}

func GetDB() *xorm.Engine { return x.GetDB() }
func (x *Xorm) GetDB() *xorm.Engine {
	if x.XormDB == nil {
		fmt.Print("Xorm: initDB()! \n")
		x.initDB()
	}
	return x.XormDB
}

func (x *Xorm) initDB() {
	var (
		db  *xorm.Engine
		err error
	)
	if x.Config == nil {
		panic(errors.New("Xorm.Config is nil"))
	}
	switch x.Config.Type {
	case DbPostgres.GetType():
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			x.Config.User, x.Config.Pass, x.Config.Host, x.Config.Port, x.Config.Db)
		x.XormDB, err = xorm.NewEngine(DbPostgres.GetType(), dsn)
		if err != nil {
			fmt.Printf("xorm.NewEngine(%s, %s) fails, err: %s", DbPostgres.GetType(), dsn, err.Error())
			panic(err)
		}
	case DbMySQL.GetType():
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
			x.Config.User, x.Config.Pass, x.Config.Host, x.Config.Port, x.Config.Db, x.Config.CharSet, x.Config.ParseTime, x.Config.TimeZone)
		db, err = xorm.NewEngine(DbMySQL.GetType(), dsn)
		if err != nil {
			fmt.Printf("xorm.NewEngine(%s, %s) fails, err: %s", DbMySQL.GetType(), dsn, err.Error())
			panic(err)
		}
	default:
		fmt.Println(x.Config)
		panic(errors.New("Unsupported DB Type: " + x.Config.Type))
	}
	loc, err := time.LoadLocation(x.Config.TimeZone)
	if err != nil {
		loc = time.UTC
	}
	db.SetTZLocation(loc)
	db.SetTZDatabase(loc)
	db.SetMaxIdleConns(x.Config.MaxIdle)
	db.SetMaxOpenConns(x.Config.MaxOpen)
	db.SetConnMaxLifetime(time.Duration(x.Config.MaxLifeTime))
	db.ShowSQL(x.Config.ShowSQL)
	db.Logger().SetLevel(log.LOG_DEBUG)
	db.SetMapper(names.GonicMapper{})
	db.SetTableMapper(names.NewPrefixMapper(names.SnakeMapper{}, x.Config.Prefix))
	db.SetSchema("anime")
	err = db.Ping()
	if err != nil {
		fmt.Printf("Xorm.Ping() fails, err: %v", err.Error())
		panic(err)
	}
	x.XormDB = db
}
