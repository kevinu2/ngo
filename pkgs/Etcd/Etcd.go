package Etcd

import (
	"context"
	"ngo2/pkgs/Log"
	"time"

	"github.com/pkg/errors"
	clientV3 "go.etcd.io/etcd/client/v3"
)

type etcdv3DataSource struct {
	propertyKey         string
	lastUpdatedRevision int64
	client              *Client
	// cancel is the func, call cancel will stop watching on the propertyKey
	cancel context.CancelFunc
	// closed indicate whether continuing to watch on the propertyKey
	// closed util.AtomicBool

	logger *Log.Log

	changed chan struct{}
}

// NewDataSource new a etcdv3DataSource instance.
// client is the etcdv3 client, it must be useful and should be release by User.
func NewDataSource(client *Client, key string) DataSource {
	ds := &etcdv3DataSource{
		client:      client,
		propertyKey: key,
	}
	go ds.watch()
	return ds
}

// ReadConfig ...
func (s *etcdv3DataSource) ReadConfig() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := s.client.Get(ctx, s.propertyKey)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, errors.New("empty response")
	}
	s.lastUpdatedRevision = resp.Header.GetRevision()
	return resp.Kvs[0].Value, nil
}

// IsConfigChanged ...
func (s *etcdv3DataSource) IsConfigChanged() <-chan struct{} {
	return s.changed
}

func (s *etcdv3DataSource) handle(resp *clientV3.WatchResponse) {
	if resp.CompactRevision > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.CompactRevision
	}
	if resp.Header.GetRevision() > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.Header.GetRevision()
	}

	if err := resp.Err(); err != nil {
		return
	}

	for _, ev := range resp.Events {
		if ev.Type == clientV3.EventTypePut || ev.Type == clientV3.EventTypeDelete {
			select {
			case s.changed <- struct{}{}:
			default:
			}
		}
	}
}

func (s *etcdv3DataSource) watch() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	rch := s.client.Watch(ctx, s.propertyKey, clientV3.WithCreatedNotify(), clientV3.WithRev(s.lastUpdatedRevision))
	for {
		for resp := range rch {
			s.handle(&resp)
		}
		time.Sleep(time.Second)

		ctx, cancel = context.WithCancel(context.Background())
		if s.lastUpdatedRevision > 0 {
			rch = s.client.Watch(ctx, s.propertyKey, clientV3.WithCreatedNotify(), clientV3.WithRev(s.lastUpdatedRevision))
		} else {
			rch = s.client.Watch(ctx, s.propertyKey, clientV3.WithCreatedNotify())
		}
		s.cancel = cancel
	}
}

// Close ...
func (s *etcdv3DataSource) Close() error {
	s.cancel()
	return nil
}
