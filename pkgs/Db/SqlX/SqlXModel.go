package SqlX

type Config struct {
	User        string
	Pass        string
	Host        string
	Port        int
	Type        string
	Name        string
	Charset     string
	ParseTime   bool
	TimeZone    string
	MaxSize     int
	MaxIdle     int
	MaxOpen     int
	MaxLifeTime int
	LogLevel    string
}
