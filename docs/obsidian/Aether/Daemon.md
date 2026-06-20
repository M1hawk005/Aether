# Core Daemon (`src/daemon.js`)

The Daemon is the backbone of an Aether node. It is a long-running background process that orchestrates the internal sub-systems of the network. It does not contain any user interface; instead, it exposes a local HTTP REST API for clients (like the [Gateway](./Gateway.md) or [CLI](./CLI.md)) to interact with.

[Back to Home](./Home.md)

## Architecture & Neutrality

In accordance with Aether's [Base Layer Philosophy](./Home.md), the daemon is designed as a **neutral routing pipe**. 
- It **does not** perform any content moderation. 
- It **does not** issue rejection receipts.
- Its primary responsibility is the ingress, validation, chunking, storage, and egress of data blobs (capsules).

## Key Sub-Systems

When the daemon starts, it initializes several dependencies in a strict order:
1. **Local Storage:** The local `.aether/` directory is mapped for caching blobs.
2. **Keypair Identity:** An Ed25519 cryptographic keypair is loaded from disk.
3. **Kademlia DHT:** The routing table is initialized using the node's NodeID (hash of the public key). See [DHT](./DHT.md).
4. **WebRTC Swarm:** The P2P transport layer is started, listening for signaling and LAN discovery. See [Network](./Network.md).

## Local REST API Endpoints

The daemon exposes a local HTTP server (default port `5000`) for the [SDK](./SDK.md) to consume.

### Publishing Data
`POST /api/publish`
- Accepts a `payload` (base64 or string) and an `alias/username`.
- **Chunking Logic:** The daemon splits large files into highly available 256KB chunks. 
- Each chunk is hashed (SHA-256) to produce a `CapsuleID`.
- The daemon instructs the DHT to announce these chunks to the closest nodes.

### Fetching Data
`GET /api/fetch/:capsuleId`
1. Checks the local `.aether` cache.
2. If missing, queries the [DHT](./DHT.md) for peers holding the capsule.
3. Instructs the [Network](./Network.md) layer to initiate a WebRTC transfer to download the chunk.
4. Validates the hash against the chunk content to ensure data integrity.

### Identity Management
`POST /api/claim-username`
`POST /api/release-username`
- Broadcasts a cryptographic proof linking the Ed25519 public key to a human-readable username across the DHT.

> [!TIP]
> The daemon binds to `0.0.0.0` to ensure it is accessible locally and across the LAN, which is critical for [mDNS discovery](./Network.md).
