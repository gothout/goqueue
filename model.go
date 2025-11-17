package goqueue

import "time"

type Queue struct {
	UUID      string // ID único da fila
	Key       string // chave que identifica a fila (ex: numero_protocolo)
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WorkState string

const (
	WorkPending WorkState = "pending"
	WorkRunning WorkState = "running"
	WorkDone    WorkState = "done"
	WorkFailed  WorkState = "failed"
)

type Work struct {
	UUID  string    // ID único do trabalho
	Key   string    // mesmo identificador da fila
	State WorkState // pending, running, done, failed

	InboundPayload  string // mensagem que chegou
	OutboundPayload string // resposta

	ErrorMsg string // caso falhe

	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
}
