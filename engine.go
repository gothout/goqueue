package goqueue

import (
	"sync"
)

type Engine struct {
	queues map[string]*QueueManager
	mu     sync.Mutex
}

type QueueManager struct {
	Key        string
	WorkQueue  chan *Work // canal de polling da fila
	WorkBuffer []*Work    // fila interna de works pendentes
	mu         sync.Mutex
}

// NewEngine cria uma engine central
func NewEngine() *Engine {
	return &Engine{
		queues: make(map[string]*QueueManager),
	}
}

// GetOrCreateQueue retorna a fila ou cria se não existir
func (e *Engine) GetOrCreateQueue(key string) *QueueManager {
	e.mu.Lock()
	defer e.mu.Unlock()

	q, exists := e.queues[key]
	if exists {
		return q
	}

	// cria canal para polling dessa fila
	qm := &QueueManager{
		Key:        key,
		WorkQueue:  make(chan *Work, 100), // capacidade default, pode ajustar depois
		WorkBuffer: make([]*Work, 0),
	}

	e.queues[key] = qm
	return qm
}
