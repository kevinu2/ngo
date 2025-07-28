package XormDB

import "ngo2/pkgs/Default"

type DbType uint8

const (
	DbPostgres DbType = iota + 1
	DbMySQL
)

func (dt DbType) GetType() string {
	switch dt {
	case DbPostgres:
		return "postgres"
	case DbMySQL:
		return "mysql"
	default:
		return Default.StringEmpty
	}
}
