# DukeDB

A distributed key-value database built from scratch in Go.

```
make kill-pro ; make compile && make run-two-nodes
make[1]: Entering directory '/home/baltej/Desktop/projects/DukeDB'
sudo lsof -ti:8000,8001 | xargs kill -9
make[1]: Leaving directory '/home/baltej/Desktop/projects/DukeDB'
make[1]: Entering directory '/home/baltej/Desktop/projects/DukeDB'
go build ./cmd/duke/main.go
make[1]: Leaving directory '/home/baltej/Desktop/projects/DukeDB'
make[1]: Entering directory '/home/baltej/Desktop/projects/DukeDB'
./main -selfAddr "localhost:8000" -selfNodeID "a" -peerAddr "localhost:8001" -peerNodeID "b" -delay 3 & \
./main -selfAddr "localhost:8001" -selfNodeID "b" -peerAddr "localhost:8000" -peerNodeID "a" -delay 8 & \
wait
2026/06/01 14:03:21 Starting duke node on 8000
2026/06/01 14:03:21 Starting duke node on 8000
tcp server listening on localhost:8001
tcp server listening on localhost:8000
[localhost:8000] Sending message to localhost:8001
2026/06/01 14:03:24 [node=b] received PING request_id=O1ge8uJA4U8zyJ48XWvL from node=a
2026/06/01 14:03:24 [node=b] sending PONG request_id=O1ge8uJA4U8zyJ48XWvL to node=a
[localhost:8001] Sending message to localhost:8000
2026/06/01 14:03:24 [node=b] PONG sent successfully request_id=O1ge8uJA4U8zyJ48XWvL
[localhost:8001] Sending message to localhost:8000
2026/06/01 14:03:29 [node=a] received PING request_id=vAd4TR3v1z7JGS4p3yjZ from node=b
2026/06/01 14:03:29 [node=a] sending PONG request_id=vAd4TR3v1z7JGS4p3yjZ to node=b
[localhost:8000] Sending message to localhost:8001
2026/06/01 14:03:29 [node=a] PONG sent successfully request_id=vAd4TR3v1z7JGS4p3yjZ
^Cmake[1]: *** [Makefile:8: run-two-nodes] Interrupt
make: *** [Makefile:16: restart] Interrupt
```

## Why?

To learn:

- Distributed systems
- Database internals
- Networking
- Systems programming

## Current Features

- TCP transport
- Custom protocol
- Multi-node communication

## Planned

- Gossip protocol
- Routing
- Replication
- Persistent storage

## Status

🚧 Work in progress.
