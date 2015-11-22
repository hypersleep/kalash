package main

import (
	"log"
)

func (c JoinCommand) consul() {
	log.Println("Starting consul watcher")

	shutdownCh := makeShutdownCh()

	c.waitGroup.Add(1)

	defer c.waitGroup.Done()

	for {
		select {
		case <- shutdownCh:
			log.Println("Consul watcher stopped")
			return
		}
	}
}
