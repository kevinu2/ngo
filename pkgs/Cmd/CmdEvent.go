package Cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Event describes transport activity, independently of any business operation.
type Event struct {
	EventID        string          `json:"event_id"`
	TaskID         string          `json:"task_id"`
	Source         CmdType         `json:"source"`
	Target         CmdType         `json:"target"`
	Action         ActionType      `json:"action"`
	InstanceID     string          `json:"instance_id"`
	PeerInstanceID string          `json:"peer_instance_id"`
	Event          string          `json:"event"`
	RequestData    json.RawMessage `json:"request_data"`
	ResponseData   json.RawMessage `json:"response_data"`
	DurationMS     int64           `json:"duration_ms"`
	Error          string          `json:"error"`
	EventTime      time.Time       `json:"event_time"`
}

type eventQueue struct {
	mu      sync.RWMutex
	closed  bool
	queue   chan Event
	done    chan struct{}
	cancel  context.CancelFunc
	dropped atomic.Uint64
}

func (s *Service) startEvents() {
	ctx, cancel := context.WithCancel(context.Background())
	q := &eventQueue{queue: make(chan Event, 1024), done: make(chan struct{}), cancel: cancel}
	s.events = q
	go func() {
		defer close(q.done)
		for e := range q.queue {
			raw, _ := json.Marshal(e)
			log.Printf("[CMD] %s", raw)
			if s.cfg.OnEvent != nil {
				if err := s.cfg.OnEvent(ctx, e); err != nil {
					log.Printf("[CMD] audit_delivery_failed event_id=%s task_id=%s error=%v", e.EventID, e.TaskID, err)
				}
			}
		}
	}()
}

func (s *Service) emit(event string, cmd Message, peer string, data json.RawMessage, err error, start time.Time) {
	if s.events == nil {
		return
	}
	e := Event{EventID: uuid.NewV4().String(), TaskID: cmd.TaskID, Source: cmd.Source, Target: cmd.Target, Action: cmd.Action, InstanceID: s.instance.ID, PeerInstanceID: peer, Event: event, RequestData: SafeData(cmd.Data), ResponseData: SafeData(data), EventTime: time.Now().UTC()}
	if !start.IsZero() {
		e.DurationMS = time.Since(start).Milliseconds()
	}
	if err != nil {
		e.Error = err.Error()
		if len(e.Error) > 2048 {
			e.Error = e.Error[:2048]
		}
	}
	q := s.events
	q.mu.RLock()
	defer q.mu.RUnlock()
	if q.closed {
		return
	}
	select {
	case q.queue <- e:
	default:
		n := q.dropped.Add(1)
		if n == 1 || n%100 == 0 {
			log.Printf("[CMD] audit_queue_full dropped=%d", n)
		}
	}
}

func (s *Service) closeEvents() {
	if s.events == nil {
		return
	}
	q := s.events
	q.mu.Lock()
	if !q.closed {
		q.closed = true
		close(q.queue)
	}
	q.mu.Unlock()
	select {
	case <-q.done:
	case <-time.After(3 * time.Second):
		q.cancel()
	}
	q.cancel()
}

// DroppedEvents reports local queue overflow, independently of delivery errors.
func (s *Service) DroppedEvents() uint64 {
	if s.events == nil {
		return 0
	}
	return s.events.dropped.Load()
}

// SafeData redacts credentials recursively and bounds audit payload size.
func SafeData(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return json.RawMessage(`null`)
	}
	if len(raw) > 16384 {
		return json.RawMessage(`{"truncated":true}`)
	}
	var value any
	if json.Unmarshal(raw, &value) != nil {
		return json.RawMessage(`{"invalid_json":true}`)
	}
	var clean func(any)
	clean = func(v any) {
		switch x := v.(type) {
		case map[string]any:
			for k, v := range x {
				key := strings.ToLower(k)
				if strings.Contains(key, "password") || strings.Contains(key, "secret") || strings.Contains(key, "token") || strings.Contains(key, "authorization") || strings.Contains(key, "credential") {
					x[k] = "[REDACTED]"
				} else {
					clean(v)
				}
			}
		case []any:
			for _, v := range x {
				clean(v)
			}
		}
	}
	clean(value)
	result, _ := json.Marshal(value)
	return result
}

// HTTPEvents creates a transport observer; destination and credentials belong to the caller.
func HTTPEvents(endpoint, token string) func(context.Context, Event) error {
	if endpoint == "" {
		return nil
	}
	client := &http.Client{Timeout: 2 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return func(ctx context.Context, e Event) error {
		if token == "" {
			return fmt.Errorf("event endpoint requires a token")
		}
		body, err := json.Marshal(e)
		if err != nil {
			return err
		}
		for attempt := 0; attempt < 3; attempt++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Service-Token", token)
			resp, callErr := client.Do(req)
			if callErr == nil {
				io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
				resp.Body.Close()
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					return nil
				}
				callErr = fmt.Errorf("event HTTP status %d", resp.StatusCode)
				if resp.StatusCode < 500 && resp.StatusCode != 429 {
					return callErr
				}
			}
			err = callErr
			if attempt < 2 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Duration(attempt+1) * 100 * time.Millisecond):
				}
			}
		}
		return err
	}
}
