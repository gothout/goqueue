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

// GetState retorna o estado atual do trabalho.
func (w *Work) GetState() WorkState {
	return w.State
}

// GetInboundPayload retorna a mensagem recebida pelo trabalho.
func (w *Work) GetInboundPayload() string {
	return w.InboundPayload
}

// GetOutboundPayload retorna a resposta enviada pelo trabalho.
func (w *Work) GetOutboundPayload() string {
	return w.OutboundPayload
}

// GetErrorMessage retorna a mensagem de erro do trabalho, quando houver.
func (w *Work) GetErrorMessage() string {
	return w.ErrorMsg
}

// =========================
// Work Writer helpers
// =========================

func (w *Work) SetState(state WorkState) {
	w.State = state
}

// SetOutboundPayload registra a resposta produzida pelo trabalho.
func (w *Work) SetOutboundPayload(payload string) {
	w.OutboundPayload = payload
}

// SetErrorMessage salva a mensagem de erro vinculada ao trabalho.
func (w *Work) SetErrorMessage(err string) {
	w.ErrorMsg = err
}

// =========================
// Work Controller helpers
// =========================

// Start marca o trabalho como em execução e registra o horário de início.
func (w *Work) Start() {
	w.State = WorkRunning
	w.StartedAt = time.Now()
}

// Finish marca o trabalho como concluído e salva o horário de finalização.
func (w *Work) Finish() {
	w.State = WorkDone
	w.FinishedAt = time.Now()
}

// Fail marca o trabalho como falho, armazenando a razão informada.
func (w *Work) Fail(reason string) {
	w.State = WorkFailed
	w.ErrorMsg = reason
	w.FinishedAt = time.Now()
}

// NewWork é um helper para criar um work pronto para enfileirar.
func NewWork(key, payload string) *Work {
	return &Work{
		Key:            key,
		State:          WorkPending,
		InboundPayload: payload,
		CreatedAt:      time.Now(),
	}
}
