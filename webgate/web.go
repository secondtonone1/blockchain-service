package main

import (
	"lbaas/webgate/webservice"
	"log"
)

func main() {
	webservice := webservice.GetWebServiceInst()

	if err := webservice.Init(); err != nil {
		log.Fatal(err)
	}

	if err := webservice.Run(); err != nil {
		log.Fatal(err)
	}
}
