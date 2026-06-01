# DukeDB

A distributed key-value database built from scratch in Go.

```
./main -selfAddr "localhost:8000" -selfNodeID "a" -peerAddr "localhost:8001" -peerNodeID "b" -delay 3 & \
./main -selfAddr "localhost:8001" -selfNodeID "b" -peerAddr "localhost:8000" -peerNodeID "a" -delay 8 & \
./main1 -selfAddr "localhost:8002" -selfNodeID "c" -peerAddr "localhost:8000" -peerNodeID "a" -delay 5
2026/06/01 20:03:23 Starting duke node on localhost:8000
2026/06/01 20:03:23 Starting duke node on localhost:8002
2026/06/01 20:03:23 Starting duke node on localhost:8001
tcp server listening on localhost:8000
tcp server listening on localhost:8002
tcp server listening on localhost:8001
Sending message to localhost:8001
2026/06/01 20:03:26 new connection
2026/06/01 20:03:26 raw message:
PING
REQUEST_ID W5hklmQI6ZylVoxXtUms
NODE_ID a
2026/06/01 20:03:26 parsed type=PING request_id=W5hklmQI6ZylVoxXtUms
2026/06/01 20:03:26 [node=b] dispatching PING request_id=W5hklmQI6ZylVoxXtUms
2026/06/01 20:03:26 [node=b] received PING request_id=W5hklmQI6ZylVoxXtUms from node=a
2026/06/01 20:03:26 [node=b] sending PONG request_id=W5hklmQI6ZylVoxXtUms to node=a
Sending message to localhost:8000
2026/06/01 20:03:26 new connection
2026/06/01 20:03:26 raw message:
PONG
REQUEST_ID W5hklmQI6ZylVoxXtUms
NODE_ID b
2026/06/01 20:03:26 parsed type=PONG request_id=W5hklmQI6ZylVoxXtUms
2026/06/01 20:03:26 [node=a] dispatching PONG request_id=W5hklmQI6ZylVoxXtUms
2026/06/01 20:03:26 [node=a] received PONG request_id=W5hklmQI6ZylVoxXtUms from node=b
2026/06/01 20:03:26 [node=a] fulfilling pending request request_id=W5hklmQI6ZylVoxXtUms
2026/06/01 20:03:26 [node=b] PONG sent successfully request_id=W5hklmQI6ZylVoxXtUms
2026/06/01 20:03:26 [node=a] pending request fulfilled request_id=W5hklmQI6ZylVoxXtUms
Done PONG
Sending message to localhost:8000
2026/06/01 20:03:28 new connection
2026/06/01 20:03:28 raw message:
JOIN
REQUEST_ID 64oKa16UlwMFAcpSyFAf
ADDR localhost:8002
NODE_ID c
2026/06/01 20:03:28 parsed type=JOIN request_id=64oKa16UlwMFAcpSyFAf
2026/06/01 20:03:28 [node=a] dispatching JOIN request_id=64oKa16UlwMFAcpSyFAf
Sending message to localhost:8002
2026/06/01 20:03:28 new connection
2026/06/01 20:03:28 raw message:
JOIN_ACK
REQUEST_ID 64oKa16UlwMFAcpSyFAf
PEER_0_ADDR localhost:8002
PEER_0_NODE_ID c
PEER_1_ADDR localhost:8001
PEER_1_NODE_ID b
PEER_COUNT 2
2026/06/01 20:03:28 parsed type=JOIN_ACK request_id=64oKa16UlwMFAcpSyFAf
2026/06/01 20:03:28 [node=c] dispatching JOIN_ACK request_id=64oKa16UlwMFAcpSyFAf
2026/06/01 20:03:28 [node=c] pending request fulfilled request_id=64oKa16UlwMFAcpSyFAf
Done JOIN_ACK
Sending message to localhost:8000
2026/06/01 20:03:31 new connection
2026/06/01 20:03:31 raw message:
PING
REQUEST_ID HMzV58iLKhjh76MrPCIG
NODE_ID b
2026/06/01 20:03:31 parsed type=PING request_id=HMzV58iLKhjh76MrPCIG
2026/06/01 20:03:31 [node=a] dispatching PING request_id=HMzV58iLKhjh76MrPCIG
2026/06/01 20:03:31 [node=a] received PING request_id=HMzV58iLKhjh76MrPCIG from node=b
2026/06/01 20:03:31 [node=a] sending PONG request_id=HMzV58iLKhjh76MrPCIG to node=b
Sending message to localhost:8001
2026/06/01 20:03:31 new connection
2026/06/01 20:03:31 raw message:
PONG
REQUEST_ID HMzV58iLKhjh76MrPCIG
NODE_ID a
2026/06/01 20:03:31 parsed type=PONG request_id=HMzV58iLKhjh76MrPCIG
2026/06/01 20:03:31 [node=b] dispatching PONG request_id=HMzV58iLKhjh76MrPCIG
2026/06/01 20:03:31 [node=b] received PONG request_id=HMzV58iLKhjh76MrPCIG from node=a
2026/06/01 20:03:31 [node=b] fulfilling pending request request_id=HMzV58iLKhjh76MrPCIG
2026/06/01 20:03:31 [node=b] pending request fulfilled request_id=HMzV58iLKhjh76MrPCIG
Done PONG
2026/06/01 20:03:31 [node=a] PONG sent successfully request_id=HMzV58iLKhjh76MrPCIG
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
