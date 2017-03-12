package model

type ParallelClientCallback interface {
	OnResponseUpdated()
	OnResponseReceived(response string)
	OnError(errorCode int16, errorMsg string)
}
