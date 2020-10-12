package main

import (
	"log"
	"os"

	"eatask.com/eatask"
)

func main() {

	var (
		err      error
		exitCh   chan bool
		eaServer *eatask.EaServer
	)

	log.Println("Starting EA Task")

	exitCh = make(chan bool)

	if eaServer, err = eatask.NewEaServer(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	if err = eaServer.Run(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	if err = eaServer.ProcessInputFromFile("/opt/eatask/data/input.json"); err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	<-exitCh

}
