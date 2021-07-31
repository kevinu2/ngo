package XormDB

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/xormplus/xorm"
	"time"
)

var Xorm *xorm.Engine

func GetXormDB() *xorm.Engine {
	if Xorm == nil {
		initXormDB()
	}
	return Xorm
}

func initXormDB() {
	username := viper.Get("datasource.username")
	password := viper.Get("datasource.password")
	url := viper.Get("datasource.url")
	var err error
	db, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp%s", username, password, url))
	if err != nil {
		fmt.Printf("xorm.NewEngine() fails, err: %v", err.Error())
		panic(err)
	}
	db.SetTZLocation(time.UTC)
	db.SetTZDatabase(time.UTC)
	db.SetMaxIdleConns(viper.GetInt("datasource.max-idle"))
	db.SetMaxOpenConns(viper.GetInt("datasource.max-open"))
	db.ShowSQL(viper.GetBool("datasource.sql-show"))
	db.SetSchema("anime")
	err = db.Ping()
	if err != nil {
		fmt.Printf("Xorm.Ping() fails, err: %v", err.Error())
		panic(err)
	}
	Xorm = db
}
