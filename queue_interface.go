package goqueue

import "time"

type QueueExpiration interface {
	SetExpiration(d time.Duration)  // define tempo de expiração da fila
	ResetExpiration()               // reseta contador de expiração
	GetExpiration() time.Duration   // retorna tempo configurado
	GetTimeToExpire() time.Duration // quanto falta para expirar
}

type QueueInspector interface {
	CountAllWorks() int     // todos os works
	CountPendingWorks() int // apenas pendentes
	CountRunningWorks() int // em processamento
	CountDoneWorks() int    // finalizados
	CountFailedWorks() int  // falhados
	ListWorksByState(state WorkState) []*Work
}

type QueueOperations interface {
	AddWork(w *Work) error // adiciona um work na fila
	GetKey() string        // retorna a chave da fila
}
