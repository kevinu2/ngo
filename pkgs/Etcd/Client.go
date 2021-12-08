package Etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/kevinu2/ngo2/pkgs/Log"
	"io/ioutil"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3/concurrency"

	grpcProm "github.com/grpc-ecosystem/go-grpc-prometheus"
	clientV3 "go.etcd.io/etcd/client/v3"

	//"go.etcd.io/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc"
)

// Client ...
type Client struct {
	*clientV3.Client
	config *Config
}

// New ...
func newClient(config *Config) (*Client, error) {
	conf := clientV3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(grpcProm.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(grpcProm.StreamClientInterceptor),
		},
		AutoSyncInterval: config.AutoSyncInterval,
	}

	if config.Endpoints == nil {
		return nil, errors.New("client etcd endpoints empty, empty endpoints")
	}

	if !config.Secure {
		conf.DialOptions = append(conf.DialOptions, grpc.WithInsecure())
	}

	if config.BasicAuth {
		conf.Username = config.UserName
		conf.Password = config.Password
	}

	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	if config.CaCert != "" {
		certBytes, err := ioutil.ReadFile(config.CaCert)
		if err != nil {
			Log.Logger().Panic("parse CaCert failed, err: " + err.Error())
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}

	if config.CertFile != "" && config.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			Log.Logger().Panic("load CertFile or KeyFile failed" + err.Error())
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}

	if tlsEnabled {
		conf.TLS = tlsConfig
	}

	client, err := clientV3.New(conf)

	if err != nil {
		// Log.Logger().Panic("client etcd start panic, err: " + err.Error())
		return nil, errors.New("client etcd start failed: " + err.Error())
	}

	cc := &Client{
		Client: client,
		config: config,
	}

	Log.Logger().Info("dial etcd server")
	return cc, nil
}

// GetKeyValue queries etcd key, returns mvccpb.KeyValue
func (client *Client) GetKeyValue(ctx context.Context, key string) (kv *mvccpb.KeyValue, err error) {
	rp, err := client.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(rp.Kvs) > 0 {
		return rp.Kvs[0], nil
	}

	return
}

// GetPrefix get prefix
func (client *Client) GetPrefix(ctx context.Context, prefix string) (map[string]string, error) {
	var (
		vars = make(map[string]string)
	)

	resp, err := client.Get(ctx, prefix, clientV3.WithPrefix())
	if err != nil {
		return vars, err
	}

	for _, kv := range resp.Kvs {
		vars[string(kv.Key)] = string(kv.Value)
	}

	return vars, nil
}

// DelPrefix 按前缀删除
func (client *Client) DelPrefix(ctx context.Context, prefix string) (deleted int64, err error) {
	resp, err := client.Delete(ctx, prefix, clientV3.WithPrefix())
	if err != nil {
		return 0, err
	}
	return resp.Deleted, err
}

// GetValues queries etcd for keys prefixed by prefix.
func (client *Client) GetValues(ctx context.Context, keys ...string) (map[string]string, error) {
	var (
		firstRevision = int64(0)
		vars          = make(map[string]string)
		maxTxnOps     = 128
		// getOps        = make([]string, 0, maxTxnOps)
	)

	doTxn := func(ops []string) error {
		txnOps := make([]clientV3.Op, 0, maxTxnOps)

		for _, k := range ops {
			txnOps = append(txnOps, clientV3.OpGet(k,
				clientV3.WithPrefix(),
				clientV3.WithSort(clientV3.SortByKey, clientV3.SortDescend),
				clientV3.WithRev(firstRevision)))
		}

		result, err := client.Txn(ctx).Then(txnOps...).Commit()
		if err != nil {
			return err
		}
		for i, r := range result.Responses {
			originKey := ops[i]
			originKeyFixed := originKey
			if !strings.HasSuffix(originKeyFixed, "/") {
				originKeyFixed = originKey + "/"
			}
			for _, ev := range r.GetResponseRange().Kvs {
				k := string(ev.Key)
				if k == originKey || strings.HasPrefix(k, originKeyFixed) {
					vars[string(ev.Key)] = string(ev.Value)
				}
			}
		}
		if firstRevision == 0 {
			firstRevision = result.Header.GetRevision()
		}
		return nil
	}
	cnt := len(keys) / maxTxnOps
	for i := 0; i <= cnt; i++ {
		switch temp := i == cnt; temp {
		case false:
			if err := doTxn(keys[i*maxTxnOps : (i+1)*maxTxnOps]); err != nil {
				return vars, err
			}
		case true:
			if err := doTxn(keys[i*maxTxnOps:]); err != nil {
				return vars, err
			}
		}
	}
	return vars, nil
}

//GetLeaseSession 创建租约会话
func (client *Client) GetLeaseSession(ctx context.Context, opts ...concurrency.SessionOption) (leaseSession *concurrency.Session, err error) {
	return concurrency.NewSession(client.Client, opts...)
}
