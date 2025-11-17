package goqueue

import (
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
	mu            sync.Mutex
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

// ============================
//  QueueOperations
// ============================

func (q *QueueManager) AddWork(w *Work) error {
	if w == nil {
		return errors.New("work cannot be nil")
	}

	q.mu.Lock()
	q.WorkBuffer = append(q.WorkBuffer, w)
	q.mu.Unlock()

	q.WorkQueue <- w
	return nil
}

func (q *QueueManager) GetKey() string {
	return q.Key
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
