package goqueue

type WorkReader interface {
	GetState() WorkState
	GetInboundPayload() string
	GetOutboundPayload() string
	GetErrorMessage() string
}

type WorkWriter interface {
	SetState(state WorkState)
	SetOutboundPayload(payload string)
	SetErrorMessage(err string)
}

type WorkController interface {
	Start()  // WorkRunning + timestamp
	Finish() // WorkDone + timestamp
	Fail(reason string)
}
