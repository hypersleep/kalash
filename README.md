# kalash
Auto-failover orchestration tool for PostgreSQL based on Consul (prototyping)

## Deploying map



Kalash automatically elect a leader using consul.

Kalash configure master and standby nodes and try to start PostgreSQL service.

If service successfully started it registers a consul service with health check.

After health check passes, consul-template continuously changing a pgbonucer config and reload pgbonucer on changes. 
