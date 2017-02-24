package main

import (
	"fmt"
	"gitlab.com/vjopensrc/datasync/goclient"
)

func main() {
	frontendAdapter := goclient.FrontendAdapter{}
	goclient.RegisterFrontendAdapter(frontendAdapter)

	fmt.Println("__________________________Ticket Create____________________________")
	goclient.TicketCreateHandler("subject 1", "description 1")
}
