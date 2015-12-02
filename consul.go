package main

import (
	"os"
	"log"
	"net/http"
	consul "github.com/hashicorp/consul/api"
)

const LockKey = "kalash/master"

type LeaderElection struct {
	client         *consul.Client
	lock           *consul.Lock
	gotLeaderCh    chan struct{}
	gotSlaveCh     chan struct{}
	failCh         <- chan struct{}
	master         bool
	masterHostname string
	nodeHostname   string
}

func (c JoinCommand) consul() {
	log.Println("Starting consul watcher")

	shutdownCh := makeShutdownCh()

	c.waitGroup.Add(1)

	defer c.waitGroup.Done()

	consulConfig := &consul.Config{
		Address:    "127.0.0.1:8500",
		Scheme:     "http",
		HttpClient: http.DefaultClient,
	}

	client, err := consul.NewClient(consulConfig)
	if err != nil {
		log.Fatal("Failed to connect consul", err)
	}

	election := &LeaderElection{
		client:       client,
		gotLeaderCh:  make(chan struct{}),
		gotSlaveCh:   make(chan struct{}),
		failCh:       make(chan struct{}),
	}

	go election.Start()

	postgresActor := &PostgresActor{
		PgData: "/usr/local/var/postgres",
		PgBin:  "/usr/local/Cellar/postgresql/9.4.0/bin",
	}

	for {
		select {
		case <- election.gotLeaderCh:
			postgresActor.Master = true

			log.Println("Starting postgres as master")
			err = postgresActor.Start()
			if err != nil {
				log.Println("Failed to start postgres as master:", err)
				gracefulShutdown()
				return
			}

			log.Println("Postgres successfully started as master")

			log.Println("Registring consul health check")
			log.Println("Consul health check successfully registred")
		case <- election.gotSlaveCh:
			postgresActor.Master = false

			log.Println("Stopping postgres")
			err = postgresActor.Stop()
			if err != nil {
				log.Println("Failed to stop postgres:", err)
			}

			log.Println("Postgres succesfully stopped")

			log.Println("Try to sync with master")
			err = postgresActor.Sync(election.masterHostname)
			if err != nil {
				log.Println("Failed to sync with postgres master:", err)
				gracefulShutdown()
				return
			}

			log.Println("Successfully synced with master")


			log.Println("Starting postgres as slave")
			err = postgresActor.Start()
			if err != nil {
				log.Println("Failed to start postgres as slave:", err)
				gracefulShutdown()
				return
			}

			log.Println("Postgres successfully started as slave")

			log.Println("Registring consul health check")
			log.Println("Consul health check successfully registred")
		case <- election.failCh:
			log.Println("Restarting master election")
			go election.Start()
		case <- shutdownCh:
			log.Println("Consul watcher stopped")
			postgresActor.Stop()
			election.Stop()
			return
		}
	}
}

func (e *LeaderElection) Start() {
	log.Println("Looking for master...")

	e.nodeHostname, _ = os.Hostname()

	kv := e.client.KV()
	pair, _, err := kv.Get(LockKey, nil)
	if err != nil {
		log.Println("Failed to fetch master key from consul")
		gracefulShutdown()
		return
	}

	if pair != nil {
		e.masterHostname = string(pair.Value)
		if e.masterHostname != e.nodeHostname && pair.Session != "" {
			log.Println("Leader found, this node is slave")
			e.gotSlaveCh <- struct{}{}
		}
	}

	lockOptions := &consul.LockOptions{
		Key:   LockKey,
		Value: []byte(e.nodeHostname),
	}

	go func() {
		lock, _ := e.client.LockOpts(lockOptions)
		failCh, err := lock.Lock(make(chan struct{}))
		if err != nil {
			log.Println("Cannot acquire leadership: ", err)
			e.failCh = failCh
			return
		}

		log.Println("Election won, this node is master")
		e.failCh = failCh
		e.lock = lock
		e.master = true
		e.gotLeaderCh <- struct{}{}
	}()
}

func (e *LeaderElection) Stop() {
	log.Println("Deregister health checks")
	if e.master {
		e.master = false
		e.lock.Unlock()
		e.lock.Destroy()
	}
}
