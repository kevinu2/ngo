package Cmd

import (
	"time"

	uuid "github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// DefaultConfig returns an independent configuration for a new instance.
func DefaultConfig() *Config {
	return &Config{
		Service:          0,
		Role:             RoleExec,
		InstanceID:       uuid.NewV4().String(),
		ListenAddress:    "0.0.0.0:9102",
		AdvertiseAddress: "",
		Etcd:             clientv3.Config{DialTimeout: 5 * time.Second},
		Prefix:           "/cmd",
		LeaseTTL:         15,
		Timeout:          5 * time.Second,
		Concurrency:      8,
	}
}
