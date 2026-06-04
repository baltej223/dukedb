# DukeDB

A distributed key-value database built from scratch in Go.

```
2026/06/04 17:06:59 Starting duke node on localhost:8002
2026/06/04 17:06:59 Starting duke node on localhost:8001
2026/06/04 17:06:59 Starting duke node on localhost:8000
2026/06/04 17:06:59 [TEST] starting put/get test
tcp server listening on localhost:8000
tcp server listening on localhost:8001
tcp server listening on localhost:8002
Sending message to localhost:8000
2026/06/04 17:07:04 new connection
2026/06/04 17:07:04 parsed type=JOIN request_id=rHbs0jj2pUxv8fHczOZX
2026/06/04 17:07:04 [node=a] dispatching JOIN request_id=rHbs0jj2pUxv8fHczOZX
Sending message to localhost:8002
2026/06/04 17:07:04 new connection
2026/06/04 17:07:04 parsed type=JOIN_ACK request_id=rHbs0jj2pUxv8fHczOZX
2026/06/04 17:07:04 [node=c] dispatching JOIN_ACK request_id=rHbs0jj2pUxv8fHczOZX
2026/06/04 17:07:04 [node=c] pending request fulfilled request_id=rHbs0jj2pUxv8fHczOZX
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002

Sending message to localhost:8000
2026/06/04 17:07:07 new connection
2026/06/04 17:07:07 parsed type=JOIN request_id=XpMf0dGwxEMhwimfxErk
2026/06/04 17:07:07 [node=a] dispatching JOIN request_id=XpMf0dGwxEMhwimfxErk
Sending message to localhost:8001
2026/06/04 17:07:07 new connection
2026/06/04 17:07:07 parsed type=JOIN_ACK request_id=XpMf0dGwxEMhwimfxErk
2026/06/04 17:07:07 [node=b] dispatching JOIN_ACK request_id=XpMf0dGwxEMhwimfxErk
2026/06/04 17:07:07 [node=b] pending request fulfilled request_id=XpMf0dGwxEMhwimfxErk
Sending message to localhost:8000
2026/06/04 17:07:07 new connection
2026/06/04 17:07:07 [node=b] gossiped membership (2 peers) to a
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

2026/06/04 17:07:07 parsed type=GOSSIP_MEMBERSHIP request_id=T7GnZcyzcFfTQ5w4PGN5
2026/06/04 17:07:07 [node=a] dispatching GOSSIP_MEMBERSHIP request_id=T7GnZcyzcFfTQ5w4PGN5
Sending message to localhost:8001
2026/06/04 17:07:09 new connection
2026/06/04 17:07:09 parsed type=GOSSIP_MEMBERSHIP request_id=YlQ0DOvHfabTXMUqrafY
2026/06/04 17:07:09 [node=b] dispatching GOSSIP_MEMBERSHIP request_id=YlQ0DOvHfabTXMUqrafY
2026/06/04 17:07:09 [node=a] gossiped membership (1 peers) to b
Cluster Membership:
  c -> localhost:8002
  b -> localhost:8001
  a -> localhost:8000

2026/06/04 17:07:14 [TEST] woke up after sleep
2026/06/04 17:07:14 [TEST] cluster has 3 peers
2026/06/04 17:07:14 [TEST] selected target node=a addr=localhost:8000
PUT
REQUEST_ID 6Jd5vhoQnFZGaxkpJWKO
KEY Name
NODE_ID a
VALUE_BASE64 QmFsdGVq

2026/06/04 17:07:14 [TEST] sending PUT request_id=6Jd5vhoQnFZGaxkpJWKO
Sending message to localhost:8000
2026/06/04 17:07:14 new connection
2026/06/04 17:07:14 [PARSE PUT] NODE_ID="a" KEY="Name"
2026/06/04 17:07:14 parsed type=PUT request_id=6Jd5vhoQnFZGaxkpJWKO
2026/06/04 17:07:14 [node=a] dispatching PUT request_id=6Jd5vhoQnFZGaxkpJWKO
2026/06/04 17:07:14 [node=a] PUT ENTER request_id=6Jd5vhoQnFZGaxkpJWKO key=Name sender=a
2026/06/04 17:07:14 [node=a] PUT calling FindOwner(Name)
2026/06/04 17:07:14 [node=a] PUT owner=b self=a
2026/06/04 17:07:14 [node=a] PUT rejected owner=b self=a
Sending message to localhost:8000
2026/06/04 17:07:14 new connection
2026/06/04 17:07:14 [node=a] PUT_REJ sent successfully
2026/06/04 17:07:14 parsed type=PUT_REJ request_id=6Jd5vhoQnFZGaxkpJWKO
2026/06/04 17:07:14 [node=a] dispatching PUT_REJ request_id=6Jd5vhoQnFZGaxkpJWKO
2026/06/04 17:07:14 [TEST] PUT response received type=PUT_REJ request_id=6Jd5vhoQnFZGaxkpJWKO
2026/06/04 17:07:14 [TEST] sending GET request_id=1pR4DrWwbPYBkg9Vd58b
Sending message to localhost:8000
2026/06/04 17:07:14 new connection
2026/06/04 17:07:14 parsed type=GET request_id=1pR4DrWwbPYBkg9Vd58b
2026/06/04 17:07:14 [node=a] dispatching GET request_id=1pR4DrWwbPYBkg9Vd58b
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002

Sending message to localhost:8002
2026/06/04 17:07:17 [node=b] gossiped membership (1 peers) to c
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

2026/06/04 17:07:17 new connection
2026/06/04 17:07:17 parsed type=GOSSIP_MEMBERSHIP request_id=keE8ZYT8qZ0qjyvIPEXJ
2026/06/04 17:07:17 [node=c] dispatching GOSSIP_MEMBERSHIP request_id=keE8ZYT8qZ0qjyvIPEXJ
Sending message to localhost:8002
2026/06/04 17:07:19 new connection
2026/06/04 17:07:19 [node=a] gossiped membership (2 peers) to c
Cluster Membership:
  b -> localhost:8001
  a -> localhost:8000
  c -> localhost:8002

2026/06/04 17:07:19 parsed type=GOSSIP_MEMBERSHIP request_id=xpSAFFLipH6wBXVMaEwk
2026/06/04 17:07:19 [node=c] dispatching GOSSIP_MEMBERSHIP request_id=xpSAFFLipH6wBXVMaEwk
Sending message to localhost:8001
2026/06/04 17:07:24 new connection
2026/06/04 17:07:24 [node=c] gossiped membership (1 peers) to b
Cluster Membership:
  a -> localhost:8000
  c -> localhost:8002
  b -> localhost:8001

2026/06/04 17:07:24 parsed type=GOSSIP_MEMBERSHIP request_id=PsjEIf4p3UDbR6fa6Eld
2026/06/04 17:07:24 [node=b] dispatching GOSSIP_MEMBERSHIP request_id=PsjEIf4p3UDbR6fa6Eld
```

## Why?

To learn:

- Distributed systems
- Database internals
- Networking
- Systems programming

## Current Features

- Custom protocol
- Gossip protocol
- Routing

## Planned

- Replication
- Persistent storage

## Status

🚧 Work in progress.
