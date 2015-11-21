package main

import (
	"log"
	"sync"
	"time"

	"github.com/mitchellh/cli"
)

type (
	ServerCommand struct {
		Ui              cli.Ui
		ShutdownCh      <- chan struct{}
		watchersErrorCh <- chan int
		waitGroup       sync.WaitGroup

	}

	StatusCommand struct {
		Ui cli.Ui
	}
)

func (c ServerCommand) Run(args []string) int {
	log.Println("Starting kalash instance")
	c.watchersErrorCh = make(chan int, 3)

	go c.consul()
	go c.postgres()
	go c.kalash()

	select {
	case <- c.ShutdownCh:
		log.Println("Wait until all watchers stop")
		c.waitGroup.Wait()
		time.Sleep(time.Second)
	case exitCode := <- c.watchersErrorCh:
		return exitCode
	}

	return 0
}

func (c ServerCommand) Help() string {
	return "Help string"
}

func (c ServerCommand) Synopsis() string {
	return "Start kalash server"
}

func (c StatusCommand) Run(args []string) int {
	log.Println("Status string")
	return 0
}

func (c StatusCommand) Help() string {
	return "Help string"
}

func (c StatusCommand) Synopsis() string {
	return "Show kalash status"
}