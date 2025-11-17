package goqueue

// WorkReader define métodos de leitura dos dados de um trabalho.
type WorkReader interface {
	GetState() WorkState
	GetInboundPayload() string
	GetOutboundPayload() string
	GetErrorMessage() string
}

// WorkWriter define métodos de escrita para atualizar um trabalho.
type WorkWriter interface {
	SetState(state WorkState)
	SetOutboundPayload(payload string)
	SetErrorMessage(err string)
}

// WorkController define operações de ciclo de vida de um trabalho.
type WorkController interface {
	Start()  // WorkRunning + timestamp
	Finish() // WorkDone + timestamp
	Fail(reason string)
}
