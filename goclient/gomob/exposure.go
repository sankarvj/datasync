package gomob

import (
	"gitlab.com/vjopensrc/datasync/goclient/model"
)

var frontendAdapter IntfFrontendAdapter

type IntfFrontendAdapter interface {
	DatabasePath() string
}

func RegisterFrontendAdapter(f IntfFrontendAdapter) {
	frontendAdapter = f
	model.SetDBPath(frontendAdapter.DatabasePath())
}

type ClientCallback interface {
	OnResponseUpdated()
	OnResponseReceived(response string)
	OnError(errorCode int16, errorMsg string)
}
