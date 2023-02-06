package MongoDB

import "time"

type Config struct {
	User        string
	Password    string
	Host        string
	Port        int
	Auth        string
	MaxPoolSize uint64
	TimeOut     time.Duration
}
