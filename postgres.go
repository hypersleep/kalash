package main

import (
	"os/exec"
)

type PostgresActor struct {
	Master         bool
	MasterHostname string
	PgData         string
	PgBin          string
}

func (postgres *PostgresActor) Start() error {
	cmd := exec.Command(postgres.PgBin + "/pg_ctl", "start", "-D", postgres.PgBin)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (postgres *PostgresActor) Stop() error {
	cmd := exec.Command(postgres.PgBin + "/pg_ctl", "start", "-D", postgres.PgBin)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (postgres *PostgresActor) Sync(masterHostname string) error {
	return nil
}
