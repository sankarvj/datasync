package main

import (
	"fmt"
	"gitlab.com/vjopensrc/datasync/goclient/controller"
	"gitlab.com/vjopensrc/datasync/goclient/model"
)

func main() {
	frontendAdapter := model.FrontendAdapter{}
	model.RegisterFrontendAdapter(frontendAdapter)

	fmt.Println("__________________________Ticket Create____________________________")
	controller.TicketCreateHandler("subject 1", "description 1")
}
