package goqueue

import (
	"context"
	"errors"
	"sync"
	"time"
)

type QueueManager struct {
	Key           string
	WorkQueue     chan *Work
	WorkBuffer    []*Work
	expiration    time.Duration
	lastResetTime time.Time
	closed        bool
	mu            sync.Mutex
}

func NewQueueManager(key string, bufferSize int) *QueueManager {
	return &QueueManager{
		Key:           key,
		WorkQueue:     make(chan *Work, bufferSize),
		WorkBuffer:    make([]*Work, 0),
		lastResetTime: time.Now(),
	}
}

// ===============================
//  QueueExpiration Implementation
// ===============================

func (q *QueueManager) SetExpiration(d time.Duration) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.expiration = d
	q.lastResetTime = time.Now()
}

func (q *QueueManager) ResetExpiration() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.lastResetTime = time.Now()
}

func (q *QueueManager) GetExpiration() time.Duration {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.expiration
}

func (q *QueueManager) GetTimeToExpire() time.Duration {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.expiration == 0 {
		return 0
	}

	elapsed := time.Since(q.lastResetTime)
	return q.expiration - elapsed
}

func (q *QueueManager) HasExpired() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.expiration == 0 {
		return false
	}

	return time.Since(q.lastResetTime) >= q.expiration
}

// ============================
//  QueueOperations
// ============================

func (q *QueueManager) AddWork(w *Work) error {
	if w == nil {
		return errors.New("work cannot be nil")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return errors.New("queue is closed")
	}

	if w.State == "" {
		w.State = WorkPending
	}
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now()
	}

	q.WorkBuffer = append(q.WorkBuffer, w)

	select {
	case q.WorkQueue <- w:
	default:
		go func(w *Work) {
			q.WorkQueue <- w
		}(w)
	}

	return nil
}

func (q *QueueManager) GetKey() string {
	return q.Key
}

func (q *QueueManager) NextWork(ctx context.Context) (*Work, bool) {
	select {
	case w, ok := <-q.WorkQueue:
		if !ok {
			return nil, false
		}
		return w, true
	case <-ctx.Done():
		return nil, false
	}
}

func (q *QueueManager) Close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	close(q.WorkQueue)
	q.mu.Unlock()
}

// ============================
//  QueueInspector
// ============================

func (q *QueueManager) CountAllWorks() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.WorkBuffer)
}

func (q *QueueManager) CountPendingWorks() int {
	return len(q.ListWorksByState(WorkPending))
}

func (q *QueueManager) CountRunningWorks() int {
	return len(q.ListWorksByState(WorkRunning))
}

func (q *QueueManager) CountDoneWorks() int {
	return len(q.ListWorksByState(WorkDone))
}

func (q *QueueManager) CountFailedWorks() int {
	return len(q.ListWorksByState(WorkFailed))
}

func (q *QueueManager) ListWorksByState(state WorkState) []*Work {
	q.mu.Lock()
	defer q.mu.Unlock()

	list := make([]*Work, 0)
	for _, w := range q.WorkBuffer {
		if w.State == state {
			list = append(list, w)
		}
	}

	return list
}
