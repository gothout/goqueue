package goqueue

type EngineManager interface {
	GetOrCreateQueue(key string) *QueueManager
	GetQueue(key string) (*QueueManager, bool)
	DeleteQueue(key string) error
	ListQueues() []string
}

type QueueManagerInterface interface {
	QueueExpiration
	QueueOperations
	QueueInspector
}
