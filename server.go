package main

import (
	"fmt"
	"gitlab.com/vjopensrc/datasync/goclient/gomob"
)

func main() {
	frontendAdapter := FrontendAdapter{}
	gomob.RegisterFrontendAdapter(frontendAdapter)
	callback := CallBack{}

	//fmt.Println("__________________________Ticket Create____________________________")
	gomob.TicketCreateHandler(callback, "offfy subject 1", "description 1")

	//fmt.Println("__________________________Note Create____________________________")
	//gomob.NoteCreateHandler(callback, "lnote name", "note description", 296)

	//fmt.Println("__________________________Ticket List____________________________")
	//gomob.TicketListHandler(callback)

	//fmt.Println("__________________________Note List____________________________")
	//gomob.NoteListHandler(callback, 297)

	//fmt.Println("__________________________Ticket Edit____________________________")
	//gomob.TicketEditHandler(callback, "layss", "sankar", 296)

	//fmt.Println("____________________________Periodic Sync____________________________")
	//gomob.PeriodicSync()

}

type FrontendAdapter struct {
}

func (f FrontendAdapter) DatabasePath() string {
	return "datasync.db"
}

type CallBack struct {
}

func (c CallBack) OnResponseReceived(response string) {
	fmt.Println("......................................................")
	fmt.Println("OnResponseReceived", response)
}

func (c CallBack) OnResponseUpdated() {
	fmt.Println("......................................................")
	fmt.Println("OnResponseUpdated")
}

func (c CallBack) OnError(errorCode int16, errorMsg string) {
	fmt.Println("......................................................")
	fmt.Println("OnError errorCode ::: ", errorCode)
	fmt.Println("OnError errorMsg ::: ", errorMsg)
}
