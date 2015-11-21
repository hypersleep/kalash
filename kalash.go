package main

import (
	"log"
)

func (c ServerCommand) kalash() {
	log.Println("Starting kalash watcher")

	shutdownCh := makeShutdownCh()

	c.waitGroup.Add(1)

	defer c.waitGroup.Done()

	for {
		select {
		case <- shutdownCh:
			log.Println("Kalash watcher stopped")
			return
		}
	}
}
