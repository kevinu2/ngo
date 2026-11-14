package Cmd

import (
	"context"
	"encoding/json"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// OnMessage processes a command and returns its response.
type OnMessage func(Message) (json.RawMessage, error)

// Config is initialized by the caller using struct fields or DefaultConfig.
type Config struct {
	OnEvent          func(context.Context, Event) error
	Role             RoleType
	Service          CmdType
	InstanceID       string
	ListenAddress    string
	AdvertiseAddress string
	Etcd             clientv3.Config
	Prefix           string
	LeaseTTL         int64
	Timeout          time.Duration
	RouteTimeouts    map[CmdType]time.Duration
	Concurrency      int
	OnMessage        OnMessage
	// Routes registers additional command identities on the same listener and
	// etcd lease. A nil onMessage registers a send-only identity.
	Routes map[CmdType]OnMessage
}

type Message struct {
	SourceRole RoleType        `json:"source_role"`
	TargetRole RoleType        `json:"target_role"`
	TaskID     string          `json:"task_id"`
	Action     ActionType      `json:"action"`
	Source     CmdType         `json:"source"`
	Target     CmdType         `json:"target"`
	Data       json.RawMessage `json:"data,omitempty"`
}

type Reply struct {
	TaskID     string          `json:"task_id"`
	InstanceID string          `json:"instance_id"`
	Data       json.RawMessage `json:"data,omitempty"`
	Error      string          `json:"error,omitempty"`
}

type Instance struct {
	Role    RoleType `json:"role"`
	Service CmdType  `json:"service"`
	ID      string   `json:"id"`
	Address string   `json:"address"`
}

type Result struct {
	Instance Instance
	Reply    Reply
	Err      error
}
