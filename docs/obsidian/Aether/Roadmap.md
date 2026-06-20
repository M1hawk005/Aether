# Project Roadmap

This document outlines everything implemented so far in the Aether Network, and our strategic path ahead.

[Back to Home](./Home.md)

## What We've Implemented

1. **Core Network Architecture (Daemon):** A robust local Node.js daemon running a REST API that acts as a neutral routing pipe. It handles automatic 256KB chunking of large files and local storage ingestion.
2. **Kademlia DHT (`src/dht/kademlia.js`):** A custom Kademlia Distributed Hash Table implementation utilizing XOR distance geometry, k-buckets, and UDP-like asynchronous RPCs for decentralized routing.
3. **Transport Layer (`src/network/webrtc.js`):** 
   - **WebRTC DataChannels** for high-throughput, sequential peer-to-peer transfers.
   - **WebSockets** as a fallback and signaling layer for SDP offer/answer exchanges.
4. **LAN Discovery:** Zero-configuration `multicast-dns` (`_aether._tcp.local`) allowing nodes on the same local network to discover each other automatically without central bootstrap servers.
5. **Aether SDK (`packages/sdk`):** A TypeScript/JavaScript library providing a seamless API and exponential backoff resilience for third-party clients to interface with the local daemon.
6. **Desktop Gateway (`aether-gateway`):** A Wails-based desktop GUI wrapping a secure Go backend and React/Vite frontend. It manages Ed25519 identity generation securely.
7. **CLI Tools (`src/cli/index.js`):** Headless utilities for node orchestration, testing swarms, and interacting with the network directly.
8. **Extensive Documentation:** This deeply technical, modular Obsidian documentation vault.

## The Path Ahead

> [!TIP]
> The current focus is moving from a functional Alpha (LAN-capable) to a robust Beta (WAN-capable and production-ready).

### 1. Comprehensive Unit Testing
- **Status:** Planned
- **Details:** Implement rigorous Jest test coverage for the SDK and Daemon to ensure stability across edge cases and network drops.

### 2. Wide Area Network (WAN) Traversal
- **Status:** Planned
- **Details:** Currently, WebRTC relies on standard ICE. For nodes behind symmetric NATs over the public internet, we will need to integrate STUN/TURN servers or implement ICE TCP fallbacks to ensure connectivity outside of a LAN.

### 3. Chunk Garbage Collection & Storage Quotas
- **Status:** Planned
- **Details:** Implement LRU (Least Recently Used) caching mechanisms or user-defined storage quotas to prevent the `.aether` cache directory from filling the user's hard drive indefinitely as they browse the network.

### 4. Data Pinning
- **Status:** Planned
- **Details:** Allow users to explicitly "pin" specific capsules (like IPFS pinning) so they are exempt from garbage collection, ensuring long-term availability.

### 5. Automated CI/CD
- **Status:** Planned
- **Details:** Set up GitHub Actions to automatically run Jest tests and build the standalone Wails `.exe`, `.app`, and `.AppImage` binaries for Windows, Mac, and Linux on every commit.
