# DukeDB

A distributed key-value database built from scratch in Go.

```
2026/06/04 13:27:56 Starting duke node on localhost:8002
2026/06/04 13:27:56 Starting duke node on localhost:8000
2026/06/04 13:27:56 Starting duke node on localhost:8001
tcp server listening on localhost:8001
tcp server listening on localhost:8002
tcp server listening on localhost:8000
Sending message to localhost:8000
2026/06/04 13:28:01 new connection
2026/06/04 13:28:01 parsed type=JOIN request_id=b93Msl1TUJdwmff3us6m
2026/06/04 13:28:01 [node=a] dispatching JOIN request_id=b93Msl1TUJdwmff3us6m
Sending message to localhost:8002
2026/06/04 13:28:01 new connection
2026/06/04 13:28:01 parsed type=JOIN_ACK request_id=b93Msl1TUJdwmff3us6m
2026/06/04 13:28:01 [node=c] dispatching JOIN_ACK request_id=b93Msl1TUJdwmff3us6m
2026/06/04 13:28:01 [node=c] pending request fulfilled request_id=b93Msl1TUJdwmff3us6m
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002

Sending message to localhost:8000
2026/06/04 13:28:04 new connection
2026/06/04 13:28:04 parsed type=JOIN request_id=OtJ88zBML796fTATzpw7
2026/06/04 13:28:04 [node=a] dispatching JOIN request_id=OtJ88zBML796fTATzpw7
Sending message to localhost:8001
2026/06/04 13:28:04 new connection
2026/06/04 13:28:04 parsed type=JOIN_ACK request_id=OtJ88zBML796fTATzpw7
2026/06/04 13:28:04 [node=b] dispatching JOIN_ACK request_id=OtJ88zBML796fTATzpw7
2026/06/04 13:28:04 [node=b] pending request fulfilled request_id=OtJ88zBML796fTATzpw7
Sending message to localhost:8002
2026/06/04 13:28:04 new connection
2026/06/04 13:28:04 parsed type=GOSSIP_MEMBERSHIP request_id=9ToPHl0XvhMIP1debiJv
2026/06/04 13:28:04 [node=c] dispatching GOSSIP_MEMBERSHIP request_id=9ToPHl0XvhMIP1debiJv
2026/06/04 13:28:04 [node=b] gossiped membership (1 peers) to c
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

Cluster Membership:
  c -> localhost:8002
  b -> localhost:8001

Cluster Membership:
  b -> localhost:8001
  a -> localhost:8000
  c -> localhost:8002

Sending message to localhost:8000
2026/06/04 13:28:14 new connection
2026/06/04 13:28:14 parsed type=GOSSIP_MEMBERSHIP request_id=6Slb1mu2A54zAE7VkGm6
2026/06/04 13:28:14 [node=a] dispatching GOSSIP_MEMBERSHIP request_id=6Slb1mu2A54zAE7VkGm6
2026/06/04 13:28:14 [node=b] gossiped membership (2 peers) to a
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

^Cmake: *** [Makefile:22: run-three-nodes] Interrupt
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
