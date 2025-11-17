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

// =========================
// Work Reader helpers
// =========================

func (w *Work) GetState() WorkState {
	return w.State
}

func (w *Work) GetInboundPayload() string {
	return w.InboundPayload
}

func (w *Work) GetOutboundPayload() string {
	return w.OutboundPayload
}

func (w *Work) GetErrorMessage() string {
	return w.ErrorMsg
}

// =========================
// Work Writer helpers
// =========================

func (w *Work) SetState(state WorkState) {
	w.State = state
}

func (w *Work) SetOutboundPayload(payload string) {
	w.OutboundPayload = payload
}

func (w *Work) SetErrorMessage(err string) {
	w.ErrorMsg = err
}

// =========================
// Work Controller helpers
// =========================

func (w *Work) Start() {
	w.State = WorkRunning
	w.StartedAt = time.Now()
}

func (w *Work) Finish() {
	w.State = WorkDone
	w.FinishedAt = time.Now()
}

func (w *Work) Fail(reason string) {
	w.State = WorkFailed
	w.ErrorMsg = reason
	w.FinishedAt = time.Now()
}

// NewWork é um helper para criar um work pronto para enfileirar
func NewWork(key, payload string) *Work {
	return &Work{
		Key:            key,
		State:          WorkPending,
		InboundPayload: payload,
		CreatedAt:      time.Now(),
	}
}
