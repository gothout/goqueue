package goqueue

import "context"

type EngineManager interface {
	GetOrCreateQueue(key string) *QueueManager
	GetQueue(key string) (*QueueManager, bool)
	DeleteQueue(key string) error
	ListQueues() []string
	WatchQueues() <-chan *QueueManager
	NextWork(ctx context.Context, queueKey string) (*Work, bool)
}

type QueueManagerInterface interface {
	QueueExpiration
	QueueOperations
	QueueInspector
}
