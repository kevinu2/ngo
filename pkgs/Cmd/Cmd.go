// Package Cmd provides independent rpcx command transport and etcd discovery.
package Cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kevinu2/ngo/v2/pkgs/Registry/Etcd"
	uuid "github.com/satori/go.uuid"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const rpcService = "CMDx"

type serviceState uint8

const (
	stateNew serviceState = iota
	stateStarted
	stateClosed
)

var cmdMu sync.RWMutex
var cmd *Service

// Service owns one listener and one lease for all command routes. No Kafka dependencies.
type Service struct {
	events        *eventQueue
	cfg           Config
	instance      Instance
	etcd          *clientv3.Client
	rpc           *server.Server
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	lifecycle     sync.Mutex
	state         serviceState
	errors        chan error
	registrations chan Instance
}

// Do returns the service initialized by New, or nil before initialization.
func Do() *Service {
	cmdMu.RLock()
	defer cmdMu.RUnlock()
	return cmd
}

// Start initializes defaults if needed and panics if initialization or startup fails.
func Start() {
	v := Do()
	if v == nil {
		var err error
		v, err = New(nil)
		if err != nil {
			panic(err)
		}
	}
	if err := v.Start(); err != nil {
		panic(err)
	}
}

// New validates and copies configuration, using defaults when config is nil.
func New(config *Config) (*Service, error) {
	if config == nil {
		config = DefaultConfig()
	}
	cfg := *config
	cfg.Etcd.Endpoints = append([]string(nil), config.Etcd.Endpoints...)
	cfg.Routes = make(map[CmdType]OnMessage, len(config.Routes)+1)
	cfg.Routes[cfg.Service] = cfg.OnMessage
	for target, onMessage := range config.Routes {
		if target.Get().Name() == "" {
			return nil, errors.New("invalid CMD route")
		}
		cfg.Routes[target] = onMessage
	}

	if cfg.Service.Get().Name() == "" {
		return nil, errors.New("CMD service engine is required")
	}
	if cfg.Role != RoleCtl && cfg.Role != RoleExec {
		return nil, errors.New("CMD role must be Ctl or Exec")
	}
	if cfg.InstanceID == "" || cfg.InstanceID == "." || cfg.InstanceID == ".." || strings.ContainsAny(cfg.InstanceID, "/\\\x00") {
		return nil, errors.New("invalid CMD instance ID")
	}
	cfg.Prefix = strings.TrimRight(cfg.Prefix, "/")
	if cfg.LeaseTTL < 5 {
		return nil, errors.New("CMD lease TTL must be at least 5 seconds")
	}
	if cfg.Timeout <= 0 {
		return nil, errors.New("CMD timeout must be positive")
	}
	timeouts := make(map[CmdType]time.Duration, len(cfg.RouteTimeouts))
	for route, timeout := range cfg.RouteTimeouts {
		if route.Get().Name() == "" || timeout <= 0 {
			return nil, errors.New("invalid CMD route timeout")
		}
		timeouts[route] = timeout
	}
	cfg.RouteTimeouts = timeouts
	if cfg.Concurrency < 1 {
		return nil, errors.New("CMD concurrency must be positive")
	}
	if len(cfg.Etcd.Endpoints) == 0 {
		return nil, errors.New("CMD etcd endpoints are required")
	}

	var err error
	address := cfg.AdvertiseAddress
	if address == "" {
		address, err = Etcd.ResolveAdvertiseAddress(cfg.ListenAddress, cfg.Etcd.Endpoints[0])
		if err != nil {
			return nil, err
		}
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil || host == "" || port == "0" || net.ParseIP(host) != nil && net.ParseIP(host).IsUnspecified() {
		return nil, errors.New("CMD requires a reachable advertise address with a nonzero port")
	}
	cfg.AdvertiseAddress = address

	s := &Service{
		cfg:           cfg,
		instance:      Instance{Role: cfg.Role, Service: cfg.Service, ID: cfg.InstanceID, Address: address},
		errors:        make(chan error, 8),
		registrations: make(chan Instance, 1),
	}
	s.startEvents()
	cmdMu.Lock()
	cmd = s
	cmdMu.Unlock()
	return s, nil
}

// Start opens the receiver and maintains its leased registration.
func (s *Service) Start() error {
	s.lifecycle.Lock()
	defer s.lifecycle.Unlock()
	if s.state == stateClosed {
		return errors.New("CMD service is closed")
	}
	if s.state == stateStarted {
		return errors.New("CMD service is already started")
	}
	cfg := s.cfg
	ln, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		return err
	}
	cli, err := clientv3.New(cfg.Etcd)
	if err != nil {
		ln.Close()
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.etcd = cli
	s.cancel = cancel
	serverTimeout := cfg.Timeout
	for _, timeout := range cfg.RouteTimeouts {
		if timeout > serverTimeout {
			serverTimeout = timeout
		}
	}
	s.rpc = server.NewServer(server.WithReadTimeout(serverTimeout), server.WithWriteTimeout(serverTimeout))
	s.rpc.DisableHTTPGateway = true
	s.rpc.DisableJSONRPC = true
	if err = s.rpc.RegisterName(rpcService, s, ""); err != nil {
		cancel()
		cli.Close()
		ln.Close()
		return err
	}
	lease, keep, err := s.register(ctx)
	if err != nil {
		cancel()
		cli.Close()
		ln.Close()
		return err
	}
	s.wg.Add(2)
	go func() {
		defer s.wg.Done()
		if err := s.rpc.ServeListener("tcp", ln); err != nil && ctx.Err() == nil {
			select {
			case s.errors <- err:
			default:
			}
			cancel()
		}
	}()
	go func() { defer s.wg.Done(); defer ln.Close(); defer s.rpc.Close(); s.maintain(ctx, lease, keep) }()
	<-s.rpc.Started
	s.state = stateStarted
	return nil
}

func (s *Service) register(ctx context.Context) (clientv3.LeaseID, <-chan *clientv3.LeaseKeepAliveResponse, error) {
	op, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()
	lease, err := s.etcd.Grant(op, s.cfg.LeaseTTL)
	if err != nil {
		return 0, nil, err
	}
	var checks []clientv3.Cmp
	var writes []clientv3.Op
	for service := range s.cfg.Routes {
		instance := s.instance
		instance.Service = service
		value, _ := json.Marshal(instance)
		key := s.cfg.Prefix + "/" + service.Get().Role().name(s.cfg.Role) + "/" + s.cfg.InstanceID
		checks = append(checks, clientv3.Compare(clientv3.Version(key), "=", 0))
		writes = append(writes, clientv3.OpPut(key, string(value), clientv3.WithLease(lease.ID)))
	}
	tx, err := s.etcd.Txn(op).If(checks...).Then(writes...).Commit()
	if err != nil || !tx.Succeeded {
		revokeCtx, revokeCancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
		_, revokeErr := s.etcd.Revoke(revokeCtx, lease.ID)
		revokeCancel()
		if revokeErr != nil {
			select {
			case s.errors <- fmt.Errorf("CMD revoke: %w", revokeErr):
			default:
			}
		}
		if err == nil {
			err = errors.New("CMD instance ID is already registered")
		}
		return 0, nil, err
	}
	keep, err := s.etcd.KeepAlive(ctx, lease.ID)
	if err != nil {
		revokeCtx, revokeCancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
		_, revokeErr := s.etcd.Revoke(revokeCtx, lease.ID)
		revokeCancel()
		if revokeErr != nil {
			select {
			case s.errors <- fmt.Errorf("CMD revoke: %w", revokeErr):
			default:
			}
		}
		return 0, nil, err
	}
	select {
	case <-s.registrations:
	default:
	}
	select {
	case s.registrations <- s.instance:
	default:
	}
	return lease.ID, keep, nil
}

func (s *Service) maintain(ctx context.Context, lease clientv3.LeaseID, keep <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case <-ctx.Done():
			revokeCtx, revokeCancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
			_, revokeErr := s.etcd.Revoke(revokeCtx, lease)
			revokeCancel()
			if revokeErr != nil {
				select {
				case s.errors <- fmt.Errorf("CMD revoke: %w", revokeErr):
				default:
				}
			}
			return
		case response, ok := <-keep:
			if ok && response != nil {
				continue
			}
		}
		select {
		case s.errors <- errors.New("CMD lease interrupted; re-registering"):
		default:
		}
		revokeCtx, revokeCancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
		_, revokeErr := s.etcd.Revoke(revokeCtx, lease)
		revokeCancel()
		if revokeErr != nil {
			select {
			case s.errors <- fmt.Errorf("CMD revoke: %w", revokeErr):
			default:
			}
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			id, ch, err := s.register(ctx)
			if err != nil {
				select {
				case s.errors <- err:
				default:
				}
				continue
			}
			lease, keep = id, ch
			break
		}
	}
}

func (s *Service) Close() error {
	s.lifecycle.Lock()
	defer s.lifecycle.Unlock()
	if s.state == stateClosed {
		return nil
	}
	if s.state == stateStarted {
		s.cancel()
		s.wg.Wait()
		s.etcd.Close()
	}
	s.closeEvents()
	s.state = stateClosed
	close(s.errors)
	close(s.registrations)
	return nil
}

// Receive handles an RPC command and returns its execution result.
func (s *Service) Receive(_ context.Context, cmd *Message, reply *Reply) (receiveErr error) {
	started := time.Now()
	s.emit("receive", *cmd, "", nil, nil, time.Time{})
	defer func() {
		err := receiveErr
		if err == nil && reply.Error != "" {
			err = errors.New(reply.Error)
		}
		s.emit("complete", *cmd, "", reply.Data, err, started)
	}()
	reply.TaskID = cmd.TaskID
	reply.InstanceID = s.instance.ID
	if cmd.TargetRole != s.cfg.Role || (cmd.SourceRole != RoleCtl && cmd.SourceRole != RoleExec) {
		return errors.New("CMD role mismatch")
	}
	if cmd.TaskID == "" || cmd.Action == 0 {
		return errors.New("CMD task ID and action are required")
	}
	onMessage, registered := s.cfg.Routes[cmd.Target]
	if !registered {
		return errors.New("CMD target mismatch")
	}
	if onMessage == nil {
		return fmt.Errorf("CMD unsupported action: %d", cmd.Action)
	}
	data, err := onMessage(*cmd)
	reply.Data = data
	if err != nil {
		reply.Error = err.Error()
	}
	return nil
}

// Discover returns every live registration of the exact service name.
func (s *Service) Discover(ctx context.Context, service string) ([]Instance, error) {
	s.lifecycle.Lock()
	ready := s.state == stateStarted
	s.lifecycle.Unlock()
	if !ready {
		return nil, errors.New("CMD service is not started")
	}
	if service == "" || strings.ContainsAny(service, "/\\\x00") {
		return nil, errors.New("invalid CMD service engine")
	}
	ctx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()
	prefix := s.cfg.Prefix + "/" + service + "/"
	resp, err := s.etcd.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	instances := make([]Instance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var i Instance
		if err = json.Unmarshal(kv.Value, &i); err != nil {
			return nil, fmt.Errorf("invalid CMD registration %s: %w", kv.Key, err)
		}
		if i.Service.Get().Role().name(i.Role) != service || i.Service.Get().Name() == "" || i.ID == "" || i.ID == "." || i.ID == ".." || strings.ContainsAny(i.ID, "/\\\x00") || string(kv.Key) != prefix+i.ID || i.Address == "" {
			return nil, fmt.Errorf("invalid CMD registration %s", kv.Key)
		}
		instances = append(instances, i)
	}
	sort.Slice(instances, func(i, j int) bool { return instances[i].ID < instances[j].ID })
	return instances, nil
}

func (s *Service) call(ctx context.Context, target Instance, cmd Message) (result Result) {
	result = Result{Instance: target}
	cmd.Target = target.Service
	cmd.TargetRole = target.Role
	started := time.Now()
	s.emit("send", cmd, target.ID, nil, nil, time.Time{})
	defer func() {
		event := "reply"
		if result.Err != nil {
			event = "error"
		}
		s.emit(event, cmd, target.ID, result.Reply.Data, result.Err, started)
	}()
	timeout := s.cfg.Timeout
	if routeTimeout := s.cfg.RouteTimeouts[target.Service]; routeTimeout > 0 {
		timeout = routeTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := ctx.Err(); err != nil {
		result.Err = err
		return result
	}
	option := client.DefaultOption
	option.SerializeType = protocol.JSON
	option.Heartbeat = false
	option.Retries = 0
	option.ConnectTimeout = s.cfg.Timeout
	if d, ok := ctx.Deadline(); ok {
		option.ConnectTimeout = time.Until(d)
	}
	rpc := client.NewClient(option)
	if err := rpc.Connect("tcp", target.Address); err != nil {
		result.Err = err
		return result
	}
	defer rpc.Close()
	cmd.Target = target.Service
	cmd.TargetRole = target.Role
	metadata := make(map[string]string)
	if deadline, ok := ctx.Deadline(); ok {
		metadata[share.ServerTimeout] = strconv.FormatInt(max(1, time.Until(deadline).Milliseconds()), 10)
	}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, metadata)
	stop := context.AfterFunc(ctx, func() { _ = rpc.Conn.Close() })
	defer stop()
	call := <-rpc.Go(context.WithoutCancel(ctx), rpcService, "Receive", &cmd, &result.Reply, make(chan *client.Call, 1)).Done
	result.Err = call.Error
	if ctx.Err() != nil {
		result.Err = ctx.Err()
	}
	if result.Err == nil && result.Reply.Error != "" {
		result.Err = errors.New(result.Reply.Error)
	}
	return result
}

func (s *Service) prepare(cmd Message) (Message, error) {
	if cmd.Action == 0 {
		return cmd, errors.New("CMD action is required")
	}
	if cmd.TaskID == "" {
		cmd.TaskID = uuid.NewV4().String()
	}
	cmd.SourceRole = s.cfg.Role
	if cmd.Source == 0 {
		cmd.Source = s.instance.Service
	}
	if _, registered := s.cfg.Routes[cmd.Source]; !registered {
		return Message{}, errors.New("CMD source is not registered on this instance")
	}
	return cmd, nil
}

// Send notifies every instance registered under the service name.
func (s *Service) Send(ctx context.Context, service string, cmd Message) ([]Result, error) {
	cmd, err := s.prepare(cmd)
	if err != nil {
		return nil, err
	}
	instances, err := s.Discover(ctx, service)
	if err != nil {
		s.emit("error", cmd, "", nil, err, time.Time{})
		return nil, err
	}
	if len(instances) == 0 {
		err = fmt.Errorf("no CMD instances for %s", service)
		s.emit("error", cmd, "", nil, err, time.Time{})
		return nil, err
	}
	results := make([]Result, len(instances))
	sem := make(chan struct{}, s.cfg.Concurrency)
	var wg sync.WaitGroup
	for index, i := range instances {
		wg.Add(1)
		go func(index int, i Instance) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				results[index] = Result{Instance: i, Err: ctx.Err()}
				s.emit("error", cmd, i.ID, nil, ctx.Err(), time.Time{})
				return
			}
			results[index] = s.call(ctx, i, cmd)
		}(index, i)
	}
	wg.Wait()
	var failures []error
	for _, r := range results {
		if r.Err != nil {
			failures = append(failures, fmt.Errorf("%s/%s: %w", service, r.Instance.ID, r.Err))
		}
	}
	return results, errors.Join(failures...)
}

// Registrations emits a snapshot after initial registration and lease recovery.
func (s *Service) Registrations() <-chan Instance { return s.registrations }

// SendCommand serializes a typed payload and waits for all route instances.
func SendCommand(ctx context.Context, source, target CmdType, action ActionType, taskID string, data any) error {
	s := Do()
	if s == nil {
		return fmt.Errorf("CMD is not initialized")
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	role := RoleExec
	if s.cfg.Role == RoleExec {
		role = RoleCtl
	}
	_, err = s.Send(ctx, target.Get().Role().name(role), Message{Source: source, Action: action, TaskID: taskID, Data: payload})
	return err
}
