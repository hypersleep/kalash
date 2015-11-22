package main

import (
	"log"
)

func (c JoinCommand) postgres() {
	log.Println("Starting postgres watcher")

	shutdownCh := makeShutdownCh()

	c.waitGroup.Add(1)

	defer c.waitGroup.Done()

	for {
		select {
		case <- shutdownCh:
			log.Println("Postgres watcher stopped")
			return
		}
	}
}
