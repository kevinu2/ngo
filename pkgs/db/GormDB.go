package db

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Gorm *gorm.DB

func GetGormDB() *gorm.DB {
	if Gorm == nil {
		initGorm()
	}
	return Gorm
}

func initGorm() {
	username := viper.Get("datasource.username")
	password := viper.Get("datasource.password")
	//maxIdle := viper.GetInt("datasource.max-idle")
	//maxOpen := viper.GetInt("datasource.max-open")

	url := viper.Get("datasource.url")
	dsn := fmt.Sprintf("%s:%s@tcp%s", username, password, url)
	db, err := gorm.Open(
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
	Gorm = db
}
