# DukeDB

A distributed key-value database built from scratch in Go.

DukeDB is an experiment in understanding distributed systems by implementing them from first principles rather than relying on existing frameworks or databases.

The project focuses on membership management, request routing, gossip-based state propagation, failure handling, and eventually data replication and recovery.

```
                    ┌─────────────────────┐
                    │     Client App      │
                    │  JS / curl / SDK    │
                    └──────────┬──────────┘
                               │ HTTP
                               ▼

                    ┌─────────────────────┐
                    │      Duke API       │
                    │                     │
                    └──────────┬──────────┘
                               │
                               ▼

                    ┌─────────────────────┐
                    │      Duke Node      │
                    │       Node A        │
                    │                     │
                    └──────────┬──────────┘
                               │
                ┌──────────────┼──────────────┐
                │              │              │
                ▼              ▼              ▼

        ┌────────────┐  ┌────────────┐  ┌────────────┐
        │  Node A    │  │  Node B    │  │  Node C    │
        │ localhost  │  │ localhost  │  │ localhost  │
        │            │  │            │  │            │
        └─────┬──────┘  └─────┬──────┘  └─────┬──────┘
              │               │               │
              └─────── Gossip / Membership ───┘

-----------------------------------------------------------
            
            ┌────────────────────────────────────┐
            │             Duke Node              │
            ├────────────────────────────────────┤
            │ HTTP/API Layer                     │
            ├────────────────────────────────────┤
            │ PUT() / GET()                      │
            ├────────────────────────────────────┤
            │ Routing                            │
            │   FindOwner(key)                   │
            ├────────────────────────────────────┤
            │ Pending Requests                   │
            │   RequestID → ResultChan           │
            ├────────────────────────────────────┤
            │ Membership State                   │
            │   Peers                            │
            │   MembershipVersion                │
            ├────────────────────────────────────┤
            │ Transport                          │
            │   TCP Messages                     │
            ├────────────────────────────────────┤
            │ Local KV Store                     │
            └────────────────────────────────────┘

------------------------------------------------------------
                  ┌───────────────────┐
                  │    Duke Client    │
                  │  JS SDK / curl    │
                  └─────────┬─────────┘
                            │
                            ▼

         ┌─────────────────────────────────────┐
         │           Duke API Layer            │
         │  HTTP / JSON Interface for Users    │
         └─────────────────┬───────────────────┘
                           │
                           ▼

         ┌─────────────────────────────────────┐
         │             Duke Cluster            │
         │                                     │
         │   Node A  ←→  Node B  ←→  Node C    │
         │                                     │
         │ Membership │ Routing │ Storage      │
         └─────────────────────────────────────┘
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

## Benchmarks (local):
- 5-node cluster
- 10,000 PUTs in ~0.9s
- 10,000 GETs in ~0.9s
- ~11k ops/sec

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

