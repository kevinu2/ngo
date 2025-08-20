package Etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"strings"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/kevinu2/ngo/v2/pkgs/Log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/grpc"
)

// Client ...
type Client struct {
	*clientv3.Client
	config *Config
}

// NewClient ...
func NewClient(config *Config) (*Client, error) {
	conf := clientv3.Config{}

	if config.Endpoints == nil {
		return nil, errors.New("client etcd endpoints empty, empty endpoints")
	}

	// build dial options (fix previously ignored insecure option)
	var dialOpts []grpc.DialOption
	if !config.Secure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if len(dialOpts) > 0 {
		conf.DialOptions = dialOpts
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
		certBytes, err := os.ReadFile(config.CaCert)
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

	client, err := clientv3.New(conf)

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

// GetKeyValue queries etcd key, returns *clientv3.KeyValue
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

	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
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
	resp, err := client.Delete(ctx, prefix, clientv3.WithPrefix())
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
		txnOps := make([]clientv3.Op, 0, maxTxnOps)

		for _, k := range ops {
			txnOps = append(txnOps, clientv3.OpGet(k,
				clientv3.WithPrefix(),
				clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend),
				clientv3.WithRev(firstRevision)))
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
		start := i * maxTxnOps
		end := (i + 1) * maxTxnOps
		if end > len(keys) {
			end = len(keys)
		}
		if start < end {
			if err := doTxn(keys[start:end]); err != nil {
				return vars, err
			}
		}
	}

	return vars, nil
}

// GetLeaseSession 创建租约会话
func (client *Client) GetLeaseSession(ctx context.Context, opts ...concurrency.SessionOption) (leaseSession *concurrency.Session, err error) {
	return concurrency.NewSession(client.Client, opts...)
}
