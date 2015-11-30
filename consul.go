package main

import (
	"os"
	"log"
	"net/http"
	consul "github.com/hashicorp/consul/api"
)

const LockKey = "kalash/leader"

type LeaderElection struct {
	client         *consul.Client
	lock           *consul.Lock
	gotLeaderCh    chan struct{}
	gotSlaveCh     chan struct{}
	failCh         <- chan struct{}
	leader         bool
	leaderHostname string
	nodeHostname     string
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

	go election.start()

	for {
		select {
		case <- election.gotLeaderCh:
			// if !node.leader {
			// 	log.Println("Promote this node to master")
			// }
			log.Println("Starting postgres as master")
			log.Println("Postgres successfully started as master")
			log.Println("Registring consul health check")
			log.Println("Consul health check successfully registred")
		case <- election.gotSlaveCh:
			// if node.leader {
			// 	log.Println("Promote this node to slave")
			// }
			log.Println("Stopping postgres")
			log.Println("Postgres succesfully stopped")
			log.Println("Try to sync with master")
			log.Println("Successfully synced with master")
			log.Println("Strarting postgres as slave")
			log.Println("Postgres successfully started as slave")
			log.Println("Registring consul health check")
			log.Println("Consul health check successfully registred")
		case <- election.failCh:
			log.Println("Restarting leader election")
			go election.start()
		case <- shutdownCh:
			log.Println("Consul watcher stopped")
			election.stop()
			return
		}
	}
}

func (e *LeaderElection) start() {
	log.Println("Looking for leader...")

	e.nodeHostname, _ = os.Hostname()

	kv := e.client.KV()
	pair, _, err := kv.Get(LockKey, nil)
	if err != nil {
		log.Println("Failed to fetch leader key from consul")
		gracefulShutdown()
		return
	}

	if pair != nil {
		e.leaderHostname = string(pair.Value)
		if e.leaderHostname != e.nodeHostname && pair.Session != "" {
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
			log.Println("Cannot acquire leadership")
			e.failCh = failCh
			return
		}

		log.Println("Election won, this node is master")
		e.failCh = failCh
		e.lock = lock
		e.leader = true
		e.gotLeaderCh <- struct{}{}
	}()
}

func (e *LeaderElection) stop() {
	log.Println("Deregister health checks")
	if e.leader {
		e.leader = false
		e.lock.Unlock()
		e.lock.Destroy()
	}
}
