package gomob

import (
	"gitlab.com/vjopensrc/datasync/goclient/model"
	"gitlab.com/vjopensrc/datasync/goclient/network"
)

var frontendAdapter IntfFrontendAdapter

type IntfFrontendAdapter interface {
	DatabasePath() string
	Domain() string
}

func RegisterFrontendAdapter(f IntfFrontendAdapter) {
	frontendAdapter = f
	model.SetDBPath(frontendAdapter.DatabasePath())
	network.SetDomain(frontendAdapter.Domain())
}

type ClientCallback interface {
	OnResponseUpdated()
	OnResponseReceived(response string)
	OnError(errorCode int16, errorMsg string)
}
