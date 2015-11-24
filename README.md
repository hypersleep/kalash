# Kalash
Feel postgres like distributed database!

Kalash is auto-failover and cluster orchestration tool for PostgreSQL based on Consul.

> The AK-47 (also known as the Kalashnikov, AK, or in Russian slang, Kalash) is a selective-fire (semi-automatic and automatic), gas-operated 7.62×39mm assault rifle, developed in the Soviet Union by Mikhail Kalashnikov. It is officially known in the Soviet documentation as Avtomat Kalashnikova (Russian: Автомат Калашникова).

## Deploying & Operational map

![Deploying map](https://github.com/hypersleep/kalash/blob/master/map.png)

Kalash automatically elect a leader using consul.

Kalash automatically configure master and standby, syncs them and trying to start PostgreSQL as child process on all cluster nodes.
If postgres successfully started it registers a consul service with health check.
After health check passes, consul-template continuously changing a pgpool config and reload pgpool on changes.

## Usage

	$ kalash
	usage: kalash [--version] [--help] <command> [<args>]

	Available commands are:
	    join      Join kalash cluster
	    leave     Leave kalash cluster
	    status    Show kalash status
