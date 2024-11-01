package SqlXDB

import "github.com/kevinu2/ngo2/pkgs/Default"

type DbType uint8

const (
	DbPostgres DbType = iota + 1
	DbMySQL
	DbMariaDB
	DbPercona
	DbDoris
	DbClickHouse
	DbMSSQL
	DbOracle
	DbSQLite
)

func (dt DbType) GetType() string {
	switch dt {
	case DbPostgres:
		return "postgres"
	case DbMySQL:
		return "mysql"
	case DbMariaDB:
		return "mysql"
	case DbPercona:
		return "mysql"
	case DbDoris:
		return "mysql"
	case DbClickHouse:
		return "clickhouse"
	case DbMSSQL:
		return "mssql"
	case DbOracle:
		return "oracle"
	case DbSQLite:
		return "sqlite"
	default:
		return Default.DefaultEmpty
	}
}
