package main

import (
	"os"
	"log"
	"net"
	"errors"
	"syscall"
	"net/rpc"
	"net/http"
)

type (
	KalashRPC struct {}
)


func (c JoinCommand) kalash() {
	log.Println("Starting kalash watcher")

	shutdownCh := makeShutdownCh()

	c.waitGroup.Add(1)

	defer c.waitGroup.Done()

	log.Println("Starting RPC server on:", c.rpcAddr)

	kalashRPC := new(KalashRPC)
	rpc.Register(kalashRPC)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", c.rpcAddr)
	if e != nil {
		log.Println("RPC listen error:", e)
		c.watchersErrorCh <- 2
		return
	}

	go http.Serve(l, nil)

	for {
		select {
		case <- shutdownCh:
			log.Println("Kalash watcher stopped")
			return
		}
	}
}

func (kalash *KalashRPC) Leave(args *int, reply *string) error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		*reply = "ERROR"
		return errors.New("Can't find process: " + err.Error())
	}

	p.Signal(syscall.SIGINT)

	*reply = "OK"
	return nil
}
