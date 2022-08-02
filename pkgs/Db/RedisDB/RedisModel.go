package RedisDB

type RedisConfig struct {
	DbUser        string
	DbPass        string
	DbHost        string
	DbPort        int
	DbName        int
	DbMaxIdle     int
	DbMaxActive   int
	DbIdleTimeout int
}
