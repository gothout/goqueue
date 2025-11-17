package goqueue

import "context"

// EngineManager expõe operações disponíveis para coordenar filas e consumir trabalhos.
type EngineManager interface {
	GetOrCreateQueue(key string) *QueueManager
	GetQueue(key string) (*QueueManager, bool)
	DeleteQueue(key string) error
	ListQueues() []string
	WatchQueues() <-chan *QueueManager
	NextWork(ctx context.Context, queueKey string) (*Work, bool)
	NextWorkFromQueues(ctx context.Context, queueKeys []string) (*Work, string, bool)
}

// QueueManagerInterface agrupa as interfaces públicas de gerenciamento de fila.
type QueueManagerInterface interface {
	QueueExpiration
	QueueOperations
	QueueInspector
}
