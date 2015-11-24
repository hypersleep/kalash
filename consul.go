package main

import (
	"log"
	"net/http"
	consul "github.com/hashicorp/consul/api"
)

const LockKey = "kalash/leader"

type LeaderElection struct {
	client          *consul.Client
	lock            *consul.Lock
	gotLeaderCh     chan struct{}
	leaderLostCh    <- chan struct{}
	gotSlaveCh      chan struct{}
	leaderChangedCh <- chan struct{}
	leader          bool
	leaderHostname  string
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
		client:          client,
		gotLeaderCh:     make(chan struct{}),
		leaderLostCh:    make(chan struct{}),
		gotSlaveCh:      make(chan struct{}),
		leaderChangedCh: make(chan struct{}),
	}

	election.start()

	for {
		select {
		case <- election.gotLeaderCh:
			if !election.leader {
				log.Println("Promote this node to master")
			}
			log.Println("Starting postgres as master")
			log.Println("Postgres successfully started as master")
			log.Println("Registring consul health check")
			log.Println("Consul health check successfully registred")
		case <- election.leaderLostCh:
			log.Println("Leadership lost")
			election.start()
		case <- election.gotSlaveCh:
			log.Println("Stopping postgres")
			log.Println("Postgres succesfully stopped")
			log.Println("Try to sync with master")
			log.Println("Successfully synced with master")
			log.Println("Strarting postgres as slave")
			log.Println("Postgres successfully started as slave")
			log.Println("Registring consul health check")
			log.Println("Consul health check successfully registred")
		case <- election.leaderChangedCh:
			log.Println("Leader changed, try to sync")
			election.start()
		case <- shutdownCh:
			log.Println("Consul watcher stopped")
			election.stop()
			return
		}
	}
}

func (e *LeaderElection) start() {
	log.Println("Looking for leader...")

	lock, err := e.client.LockKey(LockKey)
	if err != nil {
		log.Println("Leader not found, try to acquire leadership")
	}

	go e.acquire(lock)
}

func (e *LeaderElection) stop() {
	if e.leader {
		e.leader = false
		e.lock.Unlock()
		e.lock.Destroy()
	}
}

func (e *LeaderElection) acquire(lock *consul.Lock) {
	leaderLostCh, err := lock.Lock(make(chan struct{}))
	if err != nil {
		log.Println("Cannot acquire leadership")

		kv := e.client.KV()
		pair, _, err := kv.Get(LockKey, nil)
		if err != nil {
			log.Println("Failed to fetch leader key from consul")
			gracefulShutdown()
			return
		}

		e.leaderHostname = string(pair.Value)
		log.Println("Leader found, this node is slave")
		e.gotSlaveCh <- struct{}{}
		return
	}

	log.Println("Election win, this node is master")
	e.lock = lock
	e.leaderLostCh = leaderLostCh
	e.leader = true
	e.gotLeaderCh <- struct{}{}
}
