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

// NewQueueManager cria uma fila com canal e buffer interno.
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

// SetExpiration define o tempo de expiração da fila e reinicia o contador.
func (q *QueueManager) SetExpiration(d time.Duration) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.expiration = d
	q.lastResetTime = time.Now()
}

// ResetExpiration reinicia o contador de expiração sem alterar o valor configurado.
func (q *QueueManager) ResetExpiration() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.lastResetTime = time.Now()
}

// GetExpiration retorna a duração configurada para expiração da fila.
func (q *QueueManager) GetExpiration() time.Duration {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.expiration
}

// GetTimeToExpire informa quanto tempo falta para a fila expirar.
func (q *QueueManager) GetTimeToExpire() time.Duration {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.expiration == 0 {
		return 0
	}

	elapsed := time.Since(q.lastResetTime)
	return q.expiration - elapsed
}

// HasExpired indica se a fila já atingiu o tempo de expiração configurado.
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

// AddWork insere um novo trabalho na fila, preenchendo campos padrão quando necessário.
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

// GetKey retorna a chave associada à fila.
func (q *QueueManager) GetKey() string {
	return q.Key
}

// NextWork consome o próximo trabalho do canal, respeitando o contexto informado.
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

// Close encerra a fila e fecha o canal principal.
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

// CountAllWorks retorna a quantidade total de trabalhos já registrados.
func (q *QueueManager) CountAllWorks() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.WorkBuffer)
}

// CountPendingWorks retorna quantos trabalhos estão pendentes.
func (q *QueueManager) CountPendingWorks() int {
	return len(q.ListWorksByState(WorkPending))
}

// CountRunningWorks retorna quantos trabalhos estão em execução.
func (q *QueueManager) CountRunningWorks() int {
	return len(q.ListWorksByState(WorkRunning))
}

// CountDoneWorks retorna quantos trabalhos foram finalizados.
func (q *QueueManager) CountDoneWorks() int {
	return len(q.ListWorksByState(WorkDone))
}

// CountFailedWorks retorna quantos trabalhos falharam.
func (q *QueueManager) CountFailedWorks() int {
	return len(q.ListWorksByState(WorkFailed))
}

// ListWorksByState lista os trabalhos que possuem o estado informado.
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
