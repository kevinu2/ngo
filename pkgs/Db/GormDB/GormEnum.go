package GormDB

import (
	"github.com/kevinu2/ngo/v2/pkgs/Default"
	"gorm.io/gorm/logger"
)

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
		return Default.StringEmpty
	}
}

func LogLevel(level string) logger.LogLevel {
	switch level {
	case "Silent":
		return logger.Silent
	case "Error":
		return logger.Error
	case "Warn":
		return logger.Warn
	case "Info":
		return logger.Info
	default:
		return logger.Silent
	}
}
