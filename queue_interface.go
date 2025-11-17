package goqueue

import "time"

// QueueExpiration define operações para controlar tempo de expiração das filas.
type QueueExpiration interface {
	SetExpiration(d time.Duration)  // define tempo de expiração da fila
	ResetExpiration()               // reseta contador de expiração
	GetExpiration() time.Duration   // retorna tempo configurado
	GetTimeToExpire() time.Duration // quanto falta para expirar
	HasExpired() bool               // indica se já expirou
}

// QueueInspector expõe informações sobre trabalhos existentes na fila.
type QueueInspector interface {
	CountAllWorks() int     // todos os works
	CountPendingWorks() int // apenas pendentes
	CountRunningWorks() int // em processamento
	CountDoneWorks() int    // finalizados
	CountFailedWorks() int  // falhados
	ListWorksByState(state WorkState) []*Work
}

// QueueOperations contém ações de escrita e identificação da fila.
type QueueOperations interface {
	AddWork(w *Work) error // adiciona um work na fila
	GetKey() string        // retorna a chave da fila
}
