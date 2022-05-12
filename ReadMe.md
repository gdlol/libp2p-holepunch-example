# libp2p-holepunch-example
A libp2p hole-punching example using Docker Compose, where 2 nodes behind NAT try to establish direct connection with the help of a relay server. It runs in a single machine, simulating NAT environments using container networks.

## Containers
### server
The server node provides the following services:
- DHT bootstrap node for **listener** and **dialer**.
- NAT service for **listener** and **dialer** to learn their own public addresses.
- Relay service

### listener-nat, dialer-nat
Simulating routers with **iptables** rules.

### listener
A listening node behind NAT, it will advertise itself to the DHT so that the **dialer** can discover its relayed address.

### dialer
A dialing node behind NAT, it will first connect to **listener** through the relay server.
Hole punching begins automatically after the connection is established, so that subsequent streams will be opened with a direct link.

## Build (export DOCKER_BUILDKIT=1)
```
docker compose build
```

## Run
```
docker compose up
```
