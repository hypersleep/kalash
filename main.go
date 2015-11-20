package main

import (
	"os"
	"log"

	"github.com/mitchellh/cli"
)

func serverCommandFactory() {
	go consul()
	go postgres()
	go kalash()
}

func statusCommandFactory() {
	log.Println("Status string")
}

func main() {
	log.Println("Starting kalash instance")

	c := cli.NewCLI("kalash", "0.0.1")
    c.Args = os.Args[1:]
    c.Commands = map[string]cli.CommandFactory{
        "server": serverCommandFactory,
        "status": statusCommandFactory,
    }

	exitCode, err := cli.Run()
	if err != nil {
		log.Println("Error executing CLI:", err)
		os.Exit(exitCode)
	}
}
