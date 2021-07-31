package enum

import "github.com/kevinu2/ngo/constant"

type DbType uint8

const (
	DbPG DbType = iota + 1
	DbMySQL
)

func (dt DbType) GetType() string {
	switch dt {
	case DbPG:
		return "postgres"
	case DbMySQL:
		return "mysql"
	default:
		return constant.DefaultEmpty

	}
}
