package model

type GormDbConfig struct {
	DbUser        string
	DbPass        string
	DbHost        string
	DbPort        int
	DbType        string
	DbName        string
	DbTimeZone    string
	DbMaxIdle     int
	DbMaxOpen     int
	DbMaxLifeTime int
}

type RedisDbConfig struct {
	DbUser        string
	DbPass        string
	DbHost        string
	DbPort        int
	DbType        string
	DbName        int
	DbTimeZone    string
	DbMaxIdle     int
	DbMaxActive   int
	DbIdleTimeout int
}
