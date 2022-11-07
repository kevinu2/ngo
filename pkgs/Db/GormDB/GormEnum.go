package GormDB

import "github.com/kevinu2/ngo2/pkgs/Default"

type DbType uint8

const (
	DbPostgres DbType = iota + 1
	DbMySQL
	DbClickHouse
)

func (dt DbType) GetType() string {
	switch dt {
	case DbPostgres:
		return "postgres"
	case DbMySQL:
		return "mysql"
	case DbClickHouse:
		return "clickhouse"

	default:
		return Default.DefaultEmpty
	}
}
