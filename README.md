# DukeDB

A distributed key-value database built from scratch in Go.

DukeDB is an experiment in understanding distributed systems by implementing them from first principles rather than relying on existing frameworks or databases.

The project focuses on membership management, request routing, gossip-based state propagation, failure handling, and eventually data replication and recovery.

```
2026/06/05 14:13:12 Starting duke node on localhost:8000
2026/06/05 14:13:12 Starting duke node on localhost:8002
2026/06/05 14:13:12 Starting duke node on localhost:8001
tcp server listening on localhost:8002
tcp server listening on localhost:8000
tcp server listening on localhost:8001
Sending message to localhost:8000
2026/06/05 14:13:17 new connection
2026/06/05 14:13:17 parsed type=JOIN request_id=7ysXYsmkfZtyNDVhxCFa
2026/06/05 14:13:17 [node=a] dispatching JOIN request_id=7ysXYsmkfZtyNDVhxCFa
Sending message to localhost:8002
2026/06/05 14:13:17 new connection
2026/06/05 14:13:17 parsed type=JOIN_ACK request_id=7ysXYsmkfZtyNDVhxCFa
2026/06/05 14:13:17 [node=c] dispatching JOIN_ACK request_id=7ysXYsmkfZtyNDVhxCFa
2026/06/05 14:13:17 [node=c] pending request fulfilled request_id=7ysXYsmkfZtyNDVhxCFa
Sending message to localhost:8000
2026/06/05 14:13:17 new connection
2026/06/05 14:13:17 parsed type=GOSSIP_MEMBERSHIP request_id=pkammPsqBxWB2LZLzcaK
2026/06/05 14:13:17 [node=a] dispatching GOSSIP_MEMBERSHIP request_id=pkammPsqBxWB2LZLzcaK
2026/06/05 14:13:17 [node=c] gossiped membership (1 peers) to a
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002

Sending message to localhost:8000
2026/06/05 14:13:20 new connection
2026/06/05 14:13:20 parsed type=JOIN request_id=0JMuZ3f28NQI531hpeci
2026/06/05 14:13:20 [node=a] dispatching JOIN request_id=0JMuZ3f28NQI531hpeci
Sending message to localhost:8001
2026/06/05 14:13:20 new connection
2026/06/05 14:13:20 parsed type=JOIN_ACK request_id=0JMuZ3f28NQI531hpeci
2026/06/05 14:13:20 [node=b] dispatching JOIN_ACK request_id=0JMuZ3f28NQI531hpeci
2026/06/05 14:13:20 [node=b] pending request fulfilled request_id=0JMuZ3f28NQI531hpeci
Cluster Membership:
  b -> localhost:8001
  a -> localhost:8000
  c -> localhost:8002

Sending message to localhost:8002
2026/06/05 14:13:22 new connection
2026/06/05 14:13:22 [node=a] gossiped membership (1 peers) to c
Cluster Membership:
  c -> localhost:8002
  a -> localhost:8000
  b -> localhost:8001

2026/06/05 14:13:22 parsed type=GOSSIP_MEMBERSHIP request_id=WkHcXG0janRzRUywHokh
2026/06/05 14:13:22 [node=c] dispatching GOSSIP_MEMBERSHIP request_id=WkHcXG0janRzRUywHokh
2026/06/05 14:13:27 [TEST 0] PUT key=test-key-0 value=test-value-0 target=b
Sending message to localhost:8001
2026/06/05 14:13:27 new connection
2026/06/05 14:13:27 parsed type=PUT request_id=ScQHPWRZboMJYdJjwinE
2026/06/05 14:13:27 [node=b] dispatching PUT request_id=ScQHPWRZboMJYdJjwinE
2026/06/05 14:13:27 [node=b] PUT ENTER request_id=ScQHPWRZboMJYdJjwinE key=test-key-0 sender=a
2026/06/05 14:13:27 Putting test-key-0->test-value-0 to KV
Sending message to localhost:8000
2026/06/05 14:13:27 new connection
2026/06/05 14:13:27 parsed type=PUT_ACK request_id=ScQHPWRZboMJYdJjwinE
2026/06/05 14:13:27 [node=a] dispatching PUT_ACK request_id=ScQHPWRZboMJYdJjwinE
2026/06/05 14:13:27 [TEST 0] PUT response=PUT_ACK
Sending message to localhost:8001
2026/06/05 14:13:27 new connection
2026/06/05 14:13:27 parsed type=GET request_id=WA2scP0UnCNBkyKJ7wag
2026/06/05 14:13:27 [node=b] dispatching GET request_id=WA2scP0UnCNBkyKJ7wag
2026/06/05 14:13:27 Getting test-key-0->test-value-0 from KV
Sending message to localhost:8000
2026/06/05 14:13:27 new connection
2026/06/05 14:13:27 parsed type=GET_RESPONSE request_id=WA2scP0UnCNBkyKJ7wag
2026/06/05 14:13:27 [node=a] dispatching GET_RESPONSE request_id=WA2scP0UnCNBkyKJ7wag
2026/06/05 14:13:27 [node=a] pending request fulfilled request_id=WA2scP0UnCNBkyKJ7wag
2026/06/05 14:13:27 [TEST 0] SUCCESS key=test-key-0 value=test-value-0
Sending message to localhost:8000
2026/06/05 14:13:27 new connection
2026/06/05 14:13:27 parsed type=GOSSIP_MEMBERSHIP request_id=RI25HysbPXBvGp6IBIAo
2026/06/05 14:13:27 [node=c] gossiped membership (1 peers) to a
2026/06/05 14:13:27 [node=a] dispatching GOSSIP_MEMBERSHIP request_id=RI25HysbPXBvGp6IBIAo
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

Cluster Membership:
  a -> localhost:8000
  b -> localhost:8001
  c -> localhost:8002

Sending message to localhost:8001
2026/06/05 14:13:37 new connection
2026/06/05 14:13:37 [node=c] gossiped membership (1 peers) to b
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

2026/06/05 14:13:37 parsed type=GOSSIP_MEMBERSHIP request_id=SMzP6qNLnZirP02fSvEQ
2026/06/05 14:13:37 [node=b] dispatching GOSSIP_MEMBERSHIP request_id=SMzP6qNLnZirP02fSvEQ
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

Sending message to localhost:8002
2026/06/05 14:13:42 [TEST 1] PUT key=test-key-1 value=test-value-1 target=b
Sending message to localhost:8001
2026/06/05 14:13:42 new connection
2026/06/05 14:13:42 parsed type=PUT request_id=1VbnLyMpajrKRyTArWLe
2026/06/05 14:13:42 [node=b] dispatching PUT request_id=1VbnLyMpajrKRyTArWLe
2026/06/05 14:13:42 [node=b] PUT ENTER request_id=1VbnLyMpajrKRyTArWLe key=test-key-1 sender=a
2026/06/05 14:13:42 Putting test-key-1->test-value-1 to KV
Sending message to localhost:8000
2026/06/05 14:13:42 new connection
2026/06/05 14:13:42 [node=a] gossiped membership (2 peers) to c
2026/06/05 14:13:42 parsed type=GOSSIP_MEMBERSHIP request_id=J2Kf0km4HThha7OHcMWV
Sending message to localhost:8001
2026/06/05 14:13:42 [node=c] dispatching GOSSIP_MEMBERSHIP request_id=J2Kf0km4HThha7OHcMWV
2026/06/05 14:13:42 [node=a] gossiped membership (2 peers) to b
Cluster Membership:
  c -> localhost:8002
  a -> localhost:8000
  b -> localhost:8001

2026/06/05 14:13:42 new connection
2026/06/05 14:13:42 parsed type=GOSSIP_MEMBERSHIP request_id=J2Kf0km4HThha7OHcMWV
2026/06/05 14:13:42 [node=b] dispatching GOSSIP_MEMBERSHIP request_id=J2Kf0km4HThha7OHcMWV
2026/06/05 14:13:42 new connection
2026/06/05 14:13:42 parsed type=PUT_ACK request_id=1VbnLyMpajrKRyTArWLe
2026/06/05 14:13:42 [node=a] dispatching PUT_ACK request_id=1VbnLyMpajrKRyTArWLe
2026/06/05 14:13:42 [TEST 1] PUT response=PUT_ACK
Sending message to localhost:8001
2026/06/05 14:13:42 new connection
2026/06/05 14:13:42 parsed type=GET request_id=xq4w6BvFDbisNCtCZxLo
2026/06/05 14:13:42 [node=b] dispatching GET request_id=xq4w6BvFDbisNCtCZxLo
2026/06/05 14:13:42 Getting test-key-1->test-value-1 from KV
Sending message to localhost:8000
2026/06/05 14:13:42 new connection
2026/06/05 14:13:42 parsed type=GET_RESPONSE request_id=xq4w6BvFDbisNCtCZxLo
2026/06/05 14:13:42 [node=a] dispatching GET_RESPONSE request_id=xq4w6BvFDbisNCtCZxLo
2026/06/05 14:13:42 [node=a] pending request fulfilled request_id=xq4w6BvFDbisNCtCZxLo
2026/06/05 14:13:42 [TEST 1] SUCCESS key=test-key-1 value=test-value-1
Sending message to localhost:8001
2026/06/05 14:13:47 [node=c] gossiped membership (2 peers) to b
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

2026/06/05 14:13:47 new connection
2026/06/05 14:13:47 parsed type=GOSSIP_MEMBERSHIP request_id=KdjJE9e1k0CrIwIoGnfF
2026/06/05 14:13:47 [node=b] dispatching GOSSIP_MEMBERSHIP request_id=KdjJE9e1k0CrIwIoGnfF

```

## Why?

The goal is not to compete with production systems. The goal is to build and understand the machinery that makes distributed systems work.

## Current Features

### Cluster Membership

Nodes can join an existing cluster using a seed node.

New members receive the current cluster view and become part of the membership list.

### Gossip-Based Membership Propagation

Membership information is propagated between nodes using gossip.

Nodes gradually converge toward a consistent view of the cluster without requiring a central coordinator.

### Consistent Ownership Routing

Keys are deterministically mapped to owning nodes.

Requests sent to the wrong node can be rejected and redirected toward the correct owner.

### Distributed PUT / GET

Clients can:

- Store values
- Retrieve values
- Route requests across the cluster

### Membership Versioning

Each node maintains a membership version.

Version numbers are used to detect stale routing information and support future cluster state synchronization.

### Failure Detection (Work In Progress)

Nodes track suspected failures and timeouts.

This lays the groundwork for cluster healing and recovery.

---

## Architecture

```text
Client
   |
   v
Any Node
   |
   +---- Owner Node ----> Storage
```

Requests may enter through any node.

Ownership is determined using cluster membership information.

If a node receives a request for data it does not own, it can redirect the request toward the appropriate owner.

---

## Protocol

DukeDB uses a custom text-based protocol over TCP.

Examples:

```text
PUT
REQUEST_ID abc123
NODE_ID node-a
KEY user:42
VALUE_BASE64 SGVsbG8=
```

```text
GET
REQUEST_ID xyz456
NODE_ID node-b
KEY user:42
```

Membership propagation is performed using gossip messages exchanged between nodes.

---

## Current Status

Implemented:

- Node join protocol
- Membership propagation
- Gossip loop
- Request routing
- PUT
- GET
- Request/response correlation
- Membership versioning
- Redirect hints for stale routing

In Progress:

- Membership synchronization protocol
- Stale route repair
- Failure recovery
- Replication
- Data migration
- Persistence

---

## Goals

The long-term goal is to explore:

- Gossip protocols
- Membership management
- Consistent routing
- Replication
- Failure detection
- Distributed consensus
- Cluster recovery

while keeping the implementation understandable and built from first principles.

---

Built as a learning project to understand how distributed systems actually work beneath the abstractions.

# TODO

- Implement membership sync request/response
  like

```
MEMBERSHIP_SYNC_REQ

and

MEMBERSHIP_SYNC_RES
```

- Add redirect count.
