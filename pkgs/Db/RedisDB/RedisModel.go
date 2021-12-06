package RedisDB

type RedisConfig struct {
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
