package main

import (
	"fmt"
	"github.com/sankarvj/sample_syncadapter_client/goclient/gomob"
	"sync"
)

func main() {
	frontendAdapter := FrontendAdapter{}
	gomob.RegisterFrontendAdapter(frontendAdapter)
	callback := CallBack{}

	//fmt.Println("__________________________Ticket Create____________________________")
	//gomob.TicketCreateHandler(callback, "offfy subject 1", "description 1")

	//fmt.Println("__________________________Note Create____________________________")
	//gomob.NoteCreateHandler(callback, "lnote name", "note description", 313)

	//fmt.Println("__________________________Ticket List____________________________")
	//gomob.TicketListHandler(callback)

	//fmt.Println("__________________________Note List____________________________")
	gomob.NoteListHandler(callback, 314)

	//fmt.Println("__________________________Ticket Edit____________________________")
	//gomob.TicketEditHandler(callback, "layss", "sankar", 296)

	//fmt.Println("____________________________Periodic Sync____________________________")
	//gomob.PeriodicSync()

	//Add it here so that the program will wait for go routines. Like normal scenerios
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}

type FrontendAdapter struct {
}

func (f FrontendAdapter) DatabasePath() string {
	return "datasync.db"
}

func (f FrontendAdapter) Domain() string {
	return "192.168.61.13:8080"
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
