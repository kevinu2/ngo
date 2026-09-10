package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Config struct {
	Endpoints   []string
	Prefix      string
	DialTimeout time.Duration
	LeaseTTL    int64
}

type Instance struct {
	Name       string   `json:"name"`
	InstanceID string   `json:"instanceId"`
	Address    string   `json:"address"`
	Routes     []string `json:"routes"`
}

type Registrar struct {
	Address   string
	client    *clientv3.Client
	leaseID   clientv3.LeaseID
	errors    chan error
	closeErr  error
	close     sync.Once
	leaseLost atomic.Bool
}

func NewInstance(name, address string, routes []string) Instance {
	hostname, _ := os.Hostname()
	return Instance{
		Name:       name,
		InstanceID: fmt.Sprintf("%s-%d", hostname, os.Getpid()),
		Address:    address,
		Routes:     routes,
	}
}

func Register(ctx context.Context, config Config, instance Instance) (*Registrar, error) {
	if err := validate(config, instance); err != nil {
		return nil, err
	}
	client, err := clientv3.New(clientv3.Config{Endpoints: config.Endpoints, DialTimeout: config.DialTimeout})
	if err != nil {
		return nil, fmt.Errorf("connect etcd: %w", err)
	}
	lease, err := client.Grant(ctx, config.LeaseTTL)
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("grant etcd lease: %w", err)
	}
	payload, err := json.Marshal(instance)
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	key := path.Join(config.Prefix, instance.Name, instance.InstanceID)
	if _, err := client.Put(ctx, key, string(payload), clientv3.WithLease(lease.ID)); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("register etcd service: %w", err)
	}
	keepAlive, err := client.KeepAlive(ctx, lease.ID)
	if err != nil {
		_, _ = client.Revoke(context.Background(), lease.ID)
		_ = client.Close()
		return nil, fmt.Errorf("keep etcd lease alive: %w", err)
	}

	registrar := &Registrar{Address: instance.Address, client: client, leaseID: lease.ID, errors: make(chan error, 1)}
	go registrar.monitor(ctx, keepAlive)
	return registrar, nil
}

func (r *Registrar) Errors() <-chan error {
	return r.errors
}

func (r *Registrar) Close() error {
	r.close.Do(r.closeResources)
	return r.closeErr
}

func (r *Registrar) closeResources() {
	var revokeErr error
	if !r.leaseLost.Load() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, revokeErr = r.client.Revoke(ctx, r.leaseID)
		if errors.Is(revokeErr, rpctypes.ErrLeaseNotFound) {
			revokeErr = nil
		}
	}
	r.closeErr = errors.Join(revokeErr, r.client.Close())
}

func (r *Registrar) monitor(ctx context.Context, keepAlive <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-keepAlive:
			if ok {
				continue
			}
			if ctx.Err() == nil {
				r.leaseLost.Store(true)
				select {
				case r.errors <- errors.New("etcd lease keepalive stopped"):
				default:
				}
			}
			return
		}
	}
}

func ResolveAdvertiseAddress(listenAddress, endpoint string) (string, error) {
	host, port, err := net.SplitHostPort(listenAddress)
	if err != nil || port == "" {
		return "", fmt.Errorf("invalid listen address %q", listenAddress)
	}
	if host != "" && host != "0.0.0.0" && host != "::" {
		return net.JoinHostPort(host, port), nil
	}
	target := endpoint
	if parsed, parseErr := url.Parse(endpoint); parseErr == nil && parsed.Host != "" {
		target = parsed.Host
	}
	connection, err := net.DialTimeout("udp", target, 3*time.Second)
	if err != nil {
		return "", fmt.Errorf("resolve advertise address via %s: %w", endpoint, err)
	}
	defer connection.Close()
	udpAddress, ok := connection.LocalAddr().(*net.UDPAddr)
	if !ok || udpAddress.IP == nil || udpAddress.IP.IsUnspecified() {
		return "", errors.New("cannot determine advertise address")
	}
	return net.JoinHostPort(udpAddress.IP.String(), port), nil
}

func validate(config Config, instance Instance) error {
	if len(config.Endpoints) == 0 || strings.TrimSpace(config.Prefix) == "" || config.DialTimeout <= 0 || config.LeaseTTL <= 0 {
		return errors.New("etcd endpoints, prefix, dial timeout and lease TTL are required")
	}
	if instance.Name == "" || instance.InstanceID == "" || instance.Address == "" || len(instance.Routes) == 0 {
		return errors.New("service name, instance ID, address and routes are required")
	}
	if _, _, err := net.SplitHostPort(instance.Address); err != nil {
		return fmt.Errorf("invalid service address %q: %w", instance.Address, err)
	}
	for _, route := range instance.Routes {
		if !strings.HasPrefix(route, "/") {
			return fmt.Errorf("invalid service route %q", route)
		}
	}
	return nil
}

type DiscoveryConfig struct {
	Endpoints   []string
	Prefix      string
	DialTimeout time.Duration
}

type registeredInstance struct {
	Name       string   `json:"name"`
	InstanceID string   `json:"instanceId"`
	Address    string   `json:"address"`
	Routes     []string `json:"routes"`
}

type DiscoveredInstance struct {
	Service string
	Address string
	Routes  []string
}

type Watcher struct {
	client      *clientv3.Client
	prefix      string
	dialTimeout time.Duration
	apply       func([]DiscoveredInstance) error
	errors      chan error
	close       sync.Once
	closeErr    error
}

func Start(ctx context.Context, config DiscoveryConfig, apply func([]DiscoveredInstance) error) (*Watcher, error) {
	if len(config.Endpoints) == 0 || strings.TrimSpace(config.Prefix) == "" || config.DialTimeout <= 0 || apply == nil {
		return nil, errors.New("etcd endpoints, prefix, dial timeout and snapshot callback are required")
	}
	client, err := clientv3.New(clientv3.Config{Endpoints: config.Endpoints, DialTimeout: config.DialTimeout})
	if err != nil {
		return nil, fmt.Errorf("connect etcd: %w", err)
	}
	watcher := &Watcher{client: client, prefix: strings.TrimSuffix(path.Clean(config.Prefix), "/") + "/", dialTimeout: config.DialTimeout, apply: apply, errors: make(chan error, 1)}
	revision, err := watcher.reload(ctx)
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	go watcher.watch(ctx, revision+1)
	return watcher, nil
}

func (w *Watcher) Errors() <-chan error { return w.errors }

func (w *Watcher) Close() error {
	w.close.Do(func() { w.closeErr = w.client.Close() })
	return w.closeErr
}

func (w *Watcher) reload(ctx context.Context) (int64, error) {
	requestCtx, cancel := context.WithTimeout(ctx, w.dialTimeout)
	defer cancel()
	response, err := w.client.Get(requestCtx, w.prefix, clientv3.WithPrefix())
	if err != nil {
		return 0, fmt.Errorf("load etcd services: %w", err)
	}
	instances := make([]DiscoveredInstance, 0, len(response.Kvs))
	for _, item := range response.Kvs {
		instance, err := decodeInstance(item.Value)
		if err != nil {
			return 0, fmt.Errorf("decode etcd service %q: %w", string(item.Key), err)
		}
		instances = append(instances, instance)
	}
	if err := w.apply(instances); err != nil {
		return 0, fmt.Errorf("apply etcd services: %w", err)
	}
	return response.Header.Revision, nil
}

func (w *Watcher) watch(ctx context.Context, revision int64) {
	responses := w.client.Watch(ctx, w.prefix, clientv3.WithPrefix(), clientv3.WithRev(revision))
	for {
		select {
		case <-ctx.Done():
			return
		case response, ok := <-responses:
			if !ok {
				if ctx.Err() == nil {
					w.report(errors.New("etcd watch stopped"))
				}
				return
			}
			if err := response.Err(); err != nil {
				w.report(fmt.Errorf("etcd watch failed: %w", err))
				return
			}
			if _, err := w.reload(ctx); err != nil {
				w.report(err)
				return
			}
		}
	}
}

func (w *Watcher) report(err error) {
	select {
	case w.errors <- err:
	default:
	}
}

func decodeInstance(data []byte) (DiscoveredInstance, error) {
	registered := registeredInstance{}
	if err := json.Unmarshal(data, &registered); err != nil {
		return DiscoveredInstance{}, err
	}
	if registered.Name == "" || registered.InstanceID == "" || registered.Address == "" || len(registered.Routes) == 0 {
		return DiscoveredInstance{}, errors.New("registration requires name, instanceId, address and routes")
	}
	if _, _, err := net.SplitHostPort(registered.Address); err != nil {
		return DiscoveredInstance{}, fmt.Errorf("invalid address %q: %w", registered.Address, err)
	}
	for _, pattern := range registered.Routes {
		if !strings.HasPrefix(pattern, "/") {
			return DiscoveredInstance{}, fmt.Errorf("invalid route %q", pattern)
		}
	}
	return DiscoveredInstance{Service: registered.Name, Address: registered.Address, Routes: registered.Routes}, nil
}
