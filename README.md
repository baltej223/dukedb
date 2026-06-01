# DukeDB

A distributed key-value database built from scratch in Go.

```
2026/06/01 18:28:18 Starting duke node on localhost:8000
2026/06/01 18:28:18 Starting duke node on localhost:8001
tcp server listening on localhost:8001
tcp server listening on localhost:8000
Sending message to
%!(EXTRA string=localhost:8001)2026/06/01 18:28:21 new connection
2026/06/01 18:28:21 raw message:
PING
REQUEST_ID q3KzsxWtkBC4hNrSTEVb
NODE_ID a
2026/06/01 18:28:21 parsed type=PING request_id=q3KzsxWtkBC4hNrSTEVb
2026/06/01 18:28:21 [node=b] dispatching PING request_id=q3KzsxWtkBC4hNrSTEVb
2026/06/01 18:28:21 [node=b] received PING request_id=q3KzsxWtkBC4hNrSTEVb from node=a
2026/06/01 18:28:21 [node=b] sending PONG request_id=q3KzsxWtkBC4hNrSTEVb to node=a
Sending message to
%!(EXTRA string=localhost:8000)2026/06/01 18:28:21 new connection
2026/06/01 18:28:21 [node=b] PONG sent successfully request_id=q3KzsxWtkBC4hNrSTEVb
2026/06/01 18:28:21 raw message:
PONG
REQUEST_ID q3KzsxWtkBC4hNrSTEVb
NODE_ID b
2026/06/01 18:28:21 parsed type=PONG request_id=q3KzsxWtkBC4hNrSTEVb
2026/06/01 18:28:21 [node=a] dispatching PONG request_id=q3KzsxWtkBC4hNrSTEVb
2026/06/01 18:28:21 [node=a] received PONG request_id=q3KzsxWtkBC4hNrSTEVb from node=b
2026/06/01 18:28:21 [node=a] fulfilling pending request request_id=q3KzsxWtkBC4hNrSTEVb
2026/06/01 18:28:21 [node=a] pending request fulfilled request_id=q3KzsxWtkBC4hNrSTEVb
Done PONG
Sending message to
%!(EXTRA string=localhost:8000)2026/06/01 18:28:26 new connection
2026/06/01 18:28:26 raw message:
PING
REQUEST_ID 7uXVgfIfvo7d7bvAlxIf
NODE_ID b
2026/06/01 18:28:26 parsed type=PING request_id=7uXVgfIfvo7d7bvAlxIf
2026/06/01 18:28:26 [node=a] dispatching PING request_id=7uXVgfIfvo7d7bvAlxIf
2026/06/01 18:28:26 [node=a] received PING request_id=7uXVgfIfvo7d7bvAlxIf from node=b
2026/06/01 18:28:26 [node=a] sending PONG request_id=7uXVgfIfvo7d7bvAlxIf to node=b
Sending message to
%!(EXTRA string=localhost:8001)2026/06/01 18:28:26 new connection
2026/06/01 18:28:26 raw message:
PONG
REQUEST_ID 7uXVgfIfvo7d7bvAlxIf
NODE_ID a
2026/06/01 18:28:26 parsed type=PONG request_id=7uXVgfIfvo7d7bvAlxIf
2026/06/01 18:28:26 [node=b] dispatching PONG request_id=7uXVgfIfvo7d7bvAlxIf
2026/06/01 18:28:26 [node=a] PONG sent successfully request_id=7uXVgfIfvo7d7bvAlxIf
2026/06/01 18:28:26 [node=b] received PONG request_id=7uXVgfIfvo7d7bvAlxIf from node=a
2026/06/01 18:28:26 [node=b] fulfilling pending request request_id=7uXVgfIfvo7d7bvAlxIf
2026/06/01 18:28:26 [node=b] pending request fulfilled request_id=7uXVgfIfvo7d7bvAlxIf
Done PONG

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
