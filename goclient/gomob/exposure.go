package gomob

import (
	"github.com/sankarvj/sample_syncadapter_client/goclient/model"
	"github.com/sankarvj/sample_syncadapter_client/goclient/network"
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
