# Aether Protocol: Architectural Reference & Technical Decisions

This document details the comprehensive architectural, technical, and design decisions made during the development of the Aether protocol and its ecosystem.

---

## 1. Core Philosophy
The guiding principle of Aether is **True Decentralization from Day 1**. 
- **No Central Bottlenecks:** There are no master nodes, hardcoded genesis authorities, or central databases.
- **Thin Clients:** The UI and SDK are kept intentionally "dumb" and thin. All heavy lifting, cryptographic verification, data deduplication, and P2P aggregation are handled by the local Go daemon.
- **Local-First Reputation:** Reputation and trust are not global scores managed by a blockchain, but rather subjective, localized metrics based on actual data seeding.

## 2. System Components

### 2.1 The Go Daemon (`daemon-go`)
The heart of the Aether protocol. It runs a local HTTP server (`localhost:5000`) for the client apps, and a LibP2P listener for the global network.
- **Language:** Go (chosen for concurrency, speed, and excellent LibP2P support).
- **Networking:** `go-libp2p`. Uses Kademlia DHT for peer routing and discovery, and `GossipSub` for pub/sub message broadcasting.
- **Concurrency:** Thread-safe state maps utilizing `sync.Mutex` (`PendingFetches`, `PendingResolves`, `PendingAvailability`) handle asynchronous P2P pubsub events and map them back to synchronous HTTP requests.
- **Storage Strategy:** Uses a file-system based approach for MVP (`peers/` and `global_cache/` directories).
- **Data Compression:** Payloads are **Brotli-compressed** before network transmission to save bandwidth, and base64 encoded for JSON-safe transport.

### 2.2 The Node.js SDK (`packages/sdk`)
A strictly typed TypeScript wrapper that allows web and Node applications to interface with the local Go daemon.
- **Resilience:** Integrates `axios-retry` with exponential backoff.
- **Selective Retries:** Retries are strictly limited to `5xx` internal server errors (daemon crash/restart) and `ECONNREFUSED` (network drops). It **does not** retry `404` or `408` errors resulting from P2P fetch timeouts, as missing network data shouldn't freeze the UI.
- **Test Coverage:** Achieved 100% test coverage using `jest` and `axios-mock-adapter`, mocking out the daemon responses to ensure edge cases are handled elegantly.

### 2.3 The Gateway UI (`aether-gateway`)
A desktop application that acts as the primary user interface to the network.
- **Framework:** Built with **Wails** (Go backend + React/Vite frontend). Allows the UI to compile to a native desktop executable while retaining modern web technologies.
- **Aesthetics:** Minimalistic, "CTOS Hacker" aesthetic. Translucent glassmorphism was proposed but rejected in favor of a sleek, bare-bones, terminal-inspired style.
- **Integration:** Bypasses raw HTTP calls in favor of importing the local `aether-sdk`, inheriting all robustness and type-safety mechanisms.

---

## 3. Key Architectural Decisions

### 3.1 Cryptographic Identity & Envelopes
- **Keys:** Ed25519 for high-speed, secure digital signatures.
- **Encryption:** AES-GCM for payload encryption.
- **Capsule Structure:** Data is packaged into "Capsules". Each Capsule contains a `manifest.json` and a series of `chunks`. 
- **Envelopes:** The header metadata (`header.json`) is cryptographically signed by the publisher. The ID of a capsule is the hash of its Envelope. This ensures immutability; if the data changes, the hash changes.

### 3.2 Decentralized Name Registry ("First-Seen" Consensus)
- **Problem:** How to map human-readable usernames (e.g., `@alice`) to Ed25519 public keys without a centralized DNS server or expensive blockchain.
- **Solution:** A timestamp-based "First-Seen" consensus. When resolving a username, the network broadcasts a `WANT_USER_CAPSULES` packet. The daemon receives all signed claims for that username. If multiple users claim the same name, the cryptographic signature with the oldest valid timestamp wins. 
- **Tombstones:** To release a username, users publish a Cryptographic Tombstone. The daemon checks for this tombstone and frees the alias for the next user.

### 3.3 Server-Side P2P Deduplication
- **Problem:** In a GossipSub network, resolving `@alice` might yield 50 identical responses from 50 different peers. Sending an array of 50 identical JSON objects over HTTP to the UI is a massive waste of IPC bandwidth and memory.
- **Solution (The "Thin Client" philosophy):** Deduplication happens entirely in the Go daemon. The `/api/resolve` endpoint opens a 3-second context window, aggregates all incoming capsules into a map keyed by their cryptographic `signature`, and flattens the map into a unique array before sending it to the SDK. 

### 3.4 Subjective Proof-of-Storage (Reputation)
- **Problem:** Gamifiable, global reputation scores lead to botnets and centralization.
- **Solution:** Reputation is local and earned. A node's `/api/reputation` scoreboard simply measures how many bytes of data each peer has verifiably seeded on its behalf. When a node fetches a chunk from a peer, it adds the byte size to that peer's local score. 
- **UX Implementation:** Instead of showing raw numbers (e.g., "Score: 10,432"), the Gateway translates these scores into qualitative badges and leaderboard tiers, encouraging organic, reciprocal seeding.

### 3.5 Global Caching & Swarm Health
- **The Seeding Fix:** When a user fetches a capsule via `HandleFetch`, the daemon automatically decompresses the brotli chunks, serves the raw data to the UI, and simultaneously writes the compressed chunks to the `.aether/global_cache/` directory.
- **Swarm Availability:** When tracing a capsule, the daemon broadcasts a `CHECK_AVAILABILITY` ping. Peers checking their local `global_cache` reply, allowing the daemon to aggregate unique Peer IDs over a 2-second window and report an accurate "Global Seeds" count to the UI.

---

## 4. Phase 2 Scalability Architecture

To transition the Aether MVP into a production-ready, highly scalable global network, two major architectural shifts are required:

### 4.1 Dual-Layer Storage (Metadata vs. Blobs)
Currently, Aether writes chunks as raw `.json` files to the OS filesystem, which limits I/O scaling and large file support.
- **Layer 1 (The "Hot" Metadata Layer):** We will implement **BadgerDB**, an insanely fast Log-Structured Merge (LSM) key-value store. BadgerDB will handle all highly dynamic, tiny data: Usernames, Signatures, Reputation Scores, and Short JSON Posts. This eliminates OS file-descriptor bottlenecks by flushing tiny metadata writes in large contiguous blocks.
- **Layer 2 (The "Cold" Blob Layer):** To support massive static files (like 50GB videos), Aether will natively integrate **Bitswap / BitTorrent** block-exchange logic over our existing LibP2P stack. The raw bytes will be saved as contiguous files, while only their cryptographic hash (e.g., an IPFS CID) is stored inside the BadgerDB Capsule. The daemon will automatically route metadata requests to BadgerDB and stream massive blobs via Bitswap.

### 4.2 Subjective Web-of-Trust Graph
Currently, Proof-of-Storage reputation is strictly **local**—your node only trusts peers it has personally fetched data from.
- **The Upgrade:** To defend against Sybil attacks and spam without a centralized authority or global blockchain, nodes will securely gossip their local scoreboards to trusted peers.
- **Mechanism:** If Node A trusts Node B (due to high local seeding scores), and Node B trusts Node C, Node A can safely assign a proxy-trust score to Node C before ever interacting with them. This creates a highly resilient, subjective "Web of Trust" graph capable of organically isolating malicious botnets.

### 4.3 Production Bootstrapping
- **Bootstrap Nodes:** In production, the daemon will drop the localhost fallback and instead read an `AETHER_BOOTSTRAP_PEERS` environment variable to connect to dedicated, high-availability entry points (e.g., DigitalOcean/AWS droplets) to initialize the Kademlia DHT routing tables.
