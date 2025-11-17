package goqueue

import (
	"context"
	"errors"
	"sync"
)

type Engine struct {
	queues     map[string]*QueueManager
	mu         sync.Mutex
	queueEvent chan *QueueManager
}

// NewEngine cria uma engine central que mantém todas as filas registradas
func NewEngine() *Engine {
	return &Engine{
		queues:     make(map[string]*QueueManager),
		queueEvent: make(chan *QueueManager, 16),
	}
}

// GetOrCreateQueue retorna a fila existente ou cria uma nova fila com buffer padrão
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

// GetQueue retorna a fila se existir
func (e *Engine) GetQueue(key string) (*QueueManager, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	q, ok := e.queues[key]
	return q, ok
}

// DeleteQueue remove uma fila e fecha seus canais
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

// ListQueues retorna todas as chaves registradas
func (e *Engine) ListQueues() []string {
	e.mu.Lock()
	defer e.mu.Unlock()

	keys := make([]string, 0, len(e.queues))
	for k := range e.queues {
		keys = append(keys, k)
	}
	return keys
}

// WatchQueues expõe um canal para observar criação de filas
func (e *Engine) WatchQueues() <-chan *QueueManager {
	return e.queueEvent
}

// NextWork permite consumir o próximo work de uma fila pelo canal principal
func (e *Engine) NextWork(ctx context.Context, queueKey string) (*Work, bool) {
	q, ok := e.GetQueue(queueKey)
	if !ok {
		return nil, false
	}
	return q.NextWork(ctx)
}
