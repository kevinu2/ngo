package SqlX

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

var sqlX *SqlX

type SqlX struct {
	SqlXDb *sqlx.DB
	Config Config
}

func GetSqlX() *SqlX { return sqlX.GetSqlX() }
func (s *SqlX) GetSqlX() *SqlX {
	if sqlX == nil {
		s.initSqlX()
	}
	return sqlX
}

func (s *SqlX) initSqlX() {
	var (
		db  *sqlx.DB
		err error
	)

	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&maxAllowedPacket=%d",
		s.Config.User, s.Config.Pass, s.Config.Host, s.Config.Port, s.Config.Name, s.Config.Charset, s.Config.ParseTime, s.Config.TimeZone, s.Config.MaxSize)

	switch s.Config.Type {
	case DbMySQL.GetType(), DbMariaDB.GetType(), DbPercona.GetType(), DbDoris.GetType():
		db, err = sqlx.Open("mysql", url)
		if err != nil {
			panic("DB open failed: " + err.Error())
		}
	case DbPostgres.GetType():
		db, err = sqlx.Open("postgres", url)
		if err != nil {
			panic("DB open failed: " + err.Error())
		}
	default:
		panic("Unsupported DB Type: " + s.Config.Type)
	}
	db.SetMaxIdleConns(s.Config.MaxIdle)
	db.SetMaxOpenConns(s.Config.MaxOpen)
	db.SetConnMaxLifetime(time.Duration(s.Config.MaxLifeTime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(s.Config.MaxIdle) * time.Second)

	s.SqlXDb = db
}
