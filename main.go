package main

import (
	"os"
	"log"
	"os/signal"
	"syscall"

	"github.com/mitchellh/cli"
)

func main() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	c := cli.NewCLI("kalash", "0.0.1")
    c.Args = os.Args[1:]
    c.Commands = map[string]cli.CommandFactory{
        "join": func() (cli.Command, error) {
			return JoinCommand{
				Ui:                ui,
				ShutdownCh:        makeShutdownCh(),
			}, nil
		},
        "status": func() (cli.Command, error) {
			return StatusCommand{
				Ui: ui,
			}, nil
		},
		"leave": func() (cli.Command, error) {
			return LeaveCommand{
				Ui: ui,
			}, nil
		},
    }

	exitCode, err := c.Run()
	if err != nil {
		log.Println("Error executing CLI:", err)
		os.Exit(exitCode)
	}
}

func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		resultCh <- struct{}{}
	}()

	return resultCh
}
