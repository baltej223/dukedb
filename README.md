# Duke

A distributed key-value database built from scratch in Go.
DukeDB is an simple and small distributed database implementing membership, **gossip**, routing, failure detection, and replication without relying on 
existing distributed systems frameworks. Its properly distributed using **Gossip Protocol** with the idea of eventual consistency.

<p>
<img alt="GitHub Actions Workflow Status" src="https://img.shields.io/github/actions/workflow/status/baltej223/dukedb/ci.yml">
</p>

## Current Capabilities

-  [x] Distributed key-value storage
-  [x]  Gossip-based cluster membership
-  [x]  Membership synchronization
-  [x] Membership versioning
-  [x] Consistent ownership routing
-  [x] Distributed PUT / GET
-  [x] Data replication
-  [x] Automatic request forwarding
-  [x] Stale route detection and repair
-  [x] Request/response correlation
-  [x] HTTP API
-  [x] Offical JavaScript SDK
-  [x] Custom TCP protocol
-  [x] End-to-end CI testing
-  [x] Replication

In Progress:

-  [ ] Failure recovery
-  [ ] Data migration / rebalancing
-  [ ] Persistence

# Quick Start
- Clone this repositry.
```bash
git clone https://github.com/baltej223/dukedb.git dukedb
cd dukedb
```
- Then build duke. (make sure you have `make` installed)
```bash
make compile
```
- Then run few duke nodes.
```bash
make run-five-nodes
```
- One liner:
```sh
git clone https://github.com/baltej223/dukedb.git dukedb && cd dukedb && make compile && make run-five-nodes ;
```
- For running custom number of nodes, git clone the [`duke-orchestrator`](https://github.com/baltej223/duke-orchestrator). 

- If cluster started by  `make run-five-nodes` the nodes will expose a client facing http API at ports `9000`, `9001`, `9002`, `9003`, `9004`. Then using it from the client side is super-duper simple.
## Example

Store a value:

```bash
curl -X PUT \
http://localhost:9000/put \
-H "Content-Type: application/json" \
-d '{"key":"name","value":"Duke"}'
```

Retrieve it:

```bash
curl "http://localhost:9003/get?key=name"
```

Requests may be sent to **any node** in the cluster. Duke automatically routes them to the node responsible for the key.

> **NOTE**:
DukeDB's [offical Javascript client](https://github.com/baltej223/duke-client) is available, and is recommended.
Install it using `npm install duke-client`

## How It Works

- The duke node is divided into layers, like tranport layer, routing layer, api layer, node runtime layer, storing layer, routing layer and cluster layer where each layer servers its very specific purpose.
- Every node maintains a view of the cluster through periodic gossip.
- Keys are deterministically mapped to owner nodes using the routing layer.
If a request reaches a node that does not own the key (at the client interface), it is transparently forwarded to the appropriate owner.

Writes are replicated to the configured replica set, allowing data to remain available even if individual nodes fail.

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
# Running a Duke node manually
- There can be two types of duke nodes at the node startup, either a seed node, or a non-seed node.
For starting a node as a seed node at the time of startup:
> (Assuming main is the duke compiled executable)
```sh
./main -self-addr "localhost:8000" -self-node-id "a" -seed-node=true -api-at ":9000" -replication-factor 3 
```
For starting a node as a non seed node, it needs node id and address of a already running seed node, to which it can connect to form the cluster.
```sh
./main -self-addr "localhost:8001" -self-node-id "b" -peer-addr "localhost:8000" -peer-node-id "a" -delay 2 -api-at ":9001" -replication-factor 3 & \
```
## Parameters definition:

## Contributing
Its a rather new project, any types of contribution, bug reports, features, and changes are accepted.
