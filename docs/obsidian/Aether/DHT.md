# Distributed Hash Table (`src/dht/kademlia.js`)

The routing engine of the Aether Network is a bespoke implementation of the **Kademlia Distributed Hash Table (DHT)** algorithm. It governs how nodes locate each other and where data is stored across the decentralized swarm.

[Back to Home](./Home.md)

## Node IDs and XOR Distance

- **NodeID:** Upon first boot, the [Daemon](./Daemon.md) generates an Ed25519 keypair. The SHA-256 hash of the public key acts as the 160-bit NodeID.
- **XOR Metric:** The distance between any two nodes, or a node and a piece of data (CapsuleID), is calculated using the bitwise XOR operator. 
- In Kademlia, distance is geometric. The closer the XOR result is to zero, the "closer" the two entities are.

## Routing Table & K-Buckets

A node cannot memorize the IP address of every participant in the network. Instead, it maintains a routing table organized into **k-buckets**.
- Each k-bucket covers a specific segment of the XOR distance space.
- A k-bucket stores a maximum of `k` contacts (in Aether, `k=20`).
- This mathematical structure ensures that any node can locate any piece of data in $O(\log n)$ hops, where $n$ is the total number of network peers.

## Remote Procedure Calls (RPCs)

The DHT communicates asynchronously via UDP (or simulated UDP over WebSockets/WebRTC) using four standard Kademlia RPCs:

1. **PING:** Verifies that a peer in a k-bucket is still alive.
2. **STORE:** Instructs a peer to store a `<Key, Value>` pair (e.g., `CapsuleID -> ChunkData`).
3. **FIND_NODE:** Requests a list of the `k` closest peers to a given NodeID. Used heavily to bootstrap the routing table.
4. **FIND_VALUE:** Requests a specific `CapsuleID`. If the receiving peer has the data, it returns it; otherwise, it returns the `k` closest peers to the ID, allowing the search to iterate closer to the target.
