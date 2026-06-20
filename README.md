# Aether Protocol

Aether is a decentralized, peer-to-peer delivery layer for owner-controlled content. 

Rather than relying on central servers, publishers create cryptographically signed capsules. The network implements a "First-Seen" timestamp consensus for identity management, utilizes LibP2P Kademlia routing for peer discovery, and leverages GossipSub for peer-to-peer distribution. 

Nodes that fetch content temporarily cache data chunks and serve them to other peers, facilitating localized distribution through Proof-of-Storage mechanics.

## Architecture

Aether separates the core peer-to-peer network logic from the application layer interfaces. The primary cryptographic and network operations are handled by a dedicated Go daemon, which exposes an HTTP API for clients.

- **`src/daemon-go/`**: The core P2P daemon. Implements Kademlia DHT routing, GossipSub, Ed25519 cryptography, Brotli compression, and hosts the HTTP API.
- **`packages/sdk/`**: A TypeScript SDK providing typed payload resolution and retry mechanics for interacting with the daemon.
- **`apps/delphi/`**: The Delphi tracker interface, served by the daemon at `/ui/`, providing a feed of network activity.
- **`aether-gateway/`**: The primary desktop application built with Wails (Go and React).

## Documentation

For technical specifications and setup instructions, refer to the documentation:

- [Architectural Reference & Technical Decisions](docs/Aether_Architecture.md)
- [Getting Started & Application Guide](docs/Getting_Started.md)

*Note: Legacy design sketches and prototypes are retained in `docs/archive/`.*

## Quick Start

To run a local node and access the tracker or gateway:

1. **Start the Go daemon**:
   ```bash
   cd src/daemon-go
   go run .
   ```
2. **Build the SDK**:
   ```bash
   cd packages/sdk
   npm install && npm run build
   ```
3. **Access the Delphi Tracker interface**:
   The daemon serves the Delphi interface locally. Open a browser and navigate to:
   ```
   http://127.0.0.1:5000/ui/
   ```
4. **Launch the Gateway (Optional)**:
   ```bash
   cd aether-gateway
   wails dev
   ```

## Alpha Testing

To test the network across different environments (e.g., Windows and macOS):

1. Execute `.\build_alpha.ps1` from the root directory to compile native binaries for both platforms.
2. Start the daemon and exchange the printed public multiaddress with your peers to establish a connection.
3. Use the provided CLI scripts (`scripts/publish.ts` and `scripts/fetch.ts`) to transfer files across the network using the SDK.

## Contributing

Contributions are welcome. Please review the Architectural Reference to understand the system design before submitting pull requests. Current focus areas include extending the Web-of-Trust graph and implementing persistent storage via BadgerDB.
