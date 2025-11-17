package goqueue

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"
)

type Engine struct {
	queues     map[string]*QueueManager
	mu         sync.Mutex
	queueEvent chan *QueueManager
}

// NewEngine cria uma engine central que mantém todas as filas registradas.
func NewEngine() *Engine {
	e := &Engine{
		queues:     make(map[string]*QueueManager),
		queueEvent: make(chan *QueueManager, 16),
	}

	go e.cleanupExpiredQueues()

	return e
}

// GetOrCreateQueue retorna a fila existente ou cria uma nova fila com buffer padrão.
func (e *Engine) GetOrCreateQueue(key string) *QueueManager {
	e.mu.Lock()
	defer e.mu.Unlock()

	if q, exists := e.queues[key]; exists {
		return q
	}

	qm := NewQueueManager(key, 100)
	e.queues[key] = qm

	select {
	case e.queueEvent <- qm:
	default:
	}

	return qm
}

// GetQueue retorna a fila se existir.
func (e *Engine) GetQueue(key string) (*QueueManager, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	q, ok := e.queues[key]
	return q, ok
}

// DeleteQueue remove uma fila e fecha seus canais.
func (e *Engine) DeleteQueue(key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	q, ok := e.queues[key]
	if !ok {
		return errors.New("queue not found")
	}

	q.Close()
	delete(e.queues, key)
	return nil
}

func (e *Engine) cleanupExpiredQueues() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		e.mu.Lock()
		for key, q := range e.queues {
			if q.HasExpired() {
				q.Close()
				delete(e.queues, key)
			}
		}
		e.mu.Unlock()
	}
}

// ListQueues retorna todas as chaves registradas.
func (e *Engine) ListQueues() []string {
	e.mu.Lock()
	defer e.mu.Unlock()

	keys := make([]string, 0, len(e.queues))
	for k := range e.queues {
		keys = append(keys, k)
	}
	return keys
}

// WatchQueues expõe um canal para observar criação de filas.
func (e *Engine) WatchQueues() <-chan *QueueManager {
	return e.queueEvent
}

// NextWork permite consumir o próximo work de uma fila pelo canal principal.
func (e *Engine) NextWork(ctx context.Context, queueKey string) (*Work, bool) {
	q, ok := e.GetQueue(queueKey)
	if !ok {
		return nil, false
	}
	return q.NextWork(ctx)
}

// NextWorkFromQueues permite consumir o próximo work disponível entre várias filas.
// Retorna o work, a chave da fila correspondente e um indicador de sucesso.
func (e *Engine) NextWorkFromQueues(ctx context.Context, queueKeys []string) (*Work, string, bool) {
	e.mu.Lock()
	cases := []reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
	}
	keyIndex := make([]string, 0, len(queueKeys))

	for _, key := range queueKeys {
		if q, ok := e.queues[key]; ok {
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(q.WorkQueue)})
			keyIndex = append(keyIndex, key)
		}
	}
	e.mu.Unlock()

	if len(cases) == 1 {
		return nil, "", false
	}

	for len(cases) > 1 {
		chosen, recv, ok := reflect.Select(cases)
		if chosen == 0 {
			return nil, "", false
		}
		if !ok {
			cases = append(cases[:chosen], cases[chosen+1:]...)
			keyIndex = append(keyIndex[:chosen-1], keyIndex[chosen:]...)
			continue
		}

		work := recv.Interface().(*Work)
		return work, keyIndex[chosen-1], true
	}

	return nil, "", false
}
