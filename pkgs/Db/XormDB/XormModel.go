package XormDB

type XormConfig struct {
	User        string
	Pass        string
	Host        string
	Port        int
	Type        string
	Db          string
	CharSet     string // utf8
	ParseTime   string // True | False
	ShowSQL     bool
	TimeZone    string // UTC | Asia/Shanghai
	MaxIdle     int
	MaxOpen     int
	MaxLifeTime int
	Prefix      string
}
