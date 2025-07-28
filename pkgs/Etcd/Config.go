package Etcd

import (
	"flag"
	"ngo2/pkgs/Log"
	"time"
)

var ConfigPrefix = "etcdv3"

// Config ...
type (
	Config struct {
		Endpoints []string `json:"endpoints"`
		CertFile  string   `json:"certFile"`
		KeyFile   string   `json:"keyFile"`
		CaCert    string   `json:"caCert"`
		BasicAuth bool     `json:"basicAuth"`
		UserName  string   `json:"userName"`
		Password  string   `json:"-"`
		// 连接超时时间
		ConnectTimeout time.Duration `json:"connectTimeout"`
		Secure         bool          `json:"secure"`
		// 自动同步member list的间隔
		AutoSyncInterval time.Duration `json:"autoAsyncInterval"`
		TTL              int           // 单位：s
	}
)

func (config *Config) BindFlags(fs *flag.FlagSet) {
	fs.BoolVar(&config.Secure, "insecure-etcd", true, "--insecure-etcd=true")
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BasicAuth:      false,
		ConnectTimeout: time.Second * 5,
		Secure:         false,
	}
}

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig(ConfigPrefix + "." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := UnmarshalKey(key, config); err != nil {
		Log.Logger().Panic("client etcd parse config panic, err: " + err.Error())
	}
	return config
}

// Build ...
func (config *Config) Build() (*Client, error) {
	return NewClient(config)
}

func (config *Config) MustBuild() *Client {
	client, err := config.Build()
	if err != nil {
		Log.Logger().Panic("build etcd client failed, err: " + err.Error())
	}
	return client
}
