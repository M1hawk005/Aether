# Aether Network Documentation

Welcome to the internal documentation vault for the **Aether Network**. This vault is designed to provide deep, modular, and technically rigorous explanations of every component within the Aether ecosystem.

Aether is a decentralized, peer-to-peer file sharing and communication network built around the philosophy of a **neutral base routing layer**. 

## Modules

Navigate through the architectural modules below:

- [**Daemon**](./Daemon.md) 
  The core process (`src/daemon.js`). Explains the HTTP REST API, local storage chunking mechanisms, and the orchestration of network sub-systems.
- [**DHT (Distributed Hash Table)**](./DHT.md) 
  The routing engine (`src/dht/kademlia.js`). Explains Kademlia XOR metrics, routing tables, k-buckets, and network lookup RPCs.
- [**Network & WebRTC**](./Network.md)
  The transport layer (`src/network/webrtc.js`). Details the use of WebRTC DataChannels for high-throughput P2P transfers, fallback WebSockets for signaling, and mDNS zero-configuration discovery on Local Area Networks.
- [**Client SDK**](./SDK.md) 
  The JavaScript/TypeScript library (`packages/sdk`). Documentation for the `AetherClient` used by third-party applications to interact with a local Aether daemon.
- [**Desktop Gateway**](./Gateway.md) 
  The user interface (`aether-gateway`). Explains the Wails application that wraps a local daemon with a React/Vite frontend to provide the graphical experience.
- [**CLI Tools**](./CLI.md) 
  The command-line interface (`src/cli/index.js`). Details headless usage, identity generation, and direct publishing commands.
- [**Project Roadmap**](./Roadmap.md)
  A high-level overview of everything implemented so far and the strategic path forward for the network.

---
*Generated for architectural review and community onboarding.*
