package GormDB

type GormConfig struct {
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
