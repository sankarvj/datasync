package main

import (
	"fmt"
	"gitlab.com/vjopensrc/datasync/goclient"
)

func main() {
	frontendAdapter := FrontendAdapter{}
	goclient.RegisterFrontendAdapter(frontendAdapter)
	callback := CallBack{}

	//fmt.Println("__________________________Ticket Create____________________________")
	goclient.TicketCreateHandler(callback, "subject 1", "description 1")

	//fmt.Println("__________________________Note Create____________________________")
	//goclient.NoteCreateHandler("note name", "note description", 245)

	//fmt.Println("__________________________Ticket List____________________________")
	//goclient.TicketListHandler(callback)

	//fmt.Println("__________________________Note List____________________________")
	//goclient.NoteListHandler(callback, 234)

	//fmt.Println("__________________________Ticket Edit____________________________")
	//goclient.TicketEditHandler("vijay", "sankar", 233)

	fmt.Println("____________________________Generic Sync____________________________")
	//goclient.GenericSync()

}

//Adapter talks with frontend and get back info needed for goclient
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
