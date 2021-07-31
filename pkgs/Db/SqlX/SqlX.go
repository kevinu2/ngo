package SqlX

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var SqlX *sqlx.DB

func GetSqlX() *sqlx.DB {
	if SqlX == nil {
		initSqlX()
	}
	return SqlX
}

func initSqlX() {
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	database := viper.GetString("datasource.database")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	parameters := viper.GetString("datasource.parameters")
	maxIdle := viper.GetInt("datasource.max-idle")
	maxOpen := viper.GetInt("datasource.max-open")

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", username, password, host, port, database, parameters)

	db, err := sqlx.Open("mysql", url)
	if err != nil {
		panic("连接数据库失败:" + err.Error())
	}

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpen)

	SqlX = db
}
