package main

import (
	"log"
	"sync"
	"time"
	"net/rpc"

	"github.com/mitchellh/cli"
)

type (
	JoinCommand struct {
		Ui              cli.Ui
		ShutdownCh      <- chan struct{}
		watchersErrorCh chan int
		waitGroup       sync.WaitGroup
		rpcAddr         string

	}

	StatusCommand struct {
		Ui cli.Ui
		rpcAddr string
	}

	LeaveCommand struct {
		Ui cli.Ui
		rpcAddr string
	}
)

func (c JoinCommand) Run(args []string) int {
	log.Println("Starting kalash instance")
	c.rpcAddr = "127.0.0.1:8543"
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

func (c JoinCommand) Help() string {
	return "Help string"
}

func (c JoinCommand) Synopsis() string {
	return "Join kalash cluster"
}

func (c StatusCommand) Run(args []string) int {
	c.rpcAddr = "127.0.0.1:8543"
	client, err := rpc.DialHTTP("tcp", c.rpcAddr)
	if err != nil {
		log.Println("Dialing error:", err)
		return 2
	}

	var reply string
	err = client.Call("KalashRPC.Leave", 0, &reply)
	if err != nil {
		log.Println("Leave error:", err)
		return 2
	}

	log.Println("Leave status:", reply)

	return 0
}

func (c StatusCommand) Help() string {
	return "Help string"
}

func (c StatusCommand) Synopsis() string {
	return "Show kalash status"
}

func (c LeaveCommand) Run(args []string) int {
	c.rpcAddr = "127.0.0.1:8543"
	client, err := rpc.DialHTTP("tcp", c.rpcAddr)
	if err != nil {
		log.Println("Dialing error:", err)
		return 2
	}

	var reply string
	err = client.Call("KalashRPC.Leave", 0, &reply)
	if err != nil {
		log.Println("Leave error:", err)
		return 2
	}

	log.Println("Leave status:", reply)

	return 0
}

func (c LeaveCommand) Help() string {
	return "Help string"
}

func (c LeaveCommand) Synopsis() string {
	return "Leave kalash cluster"
}
