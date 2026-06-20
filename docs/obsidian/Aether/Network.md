# Network Transport Layer (`src/network/webrtc.js`)

The network transport layer is responsible for the actual byte-transfer of data between peers. It employs an Adapter Pattern to support both WebRTC DataChannels and WebSockets.

[Back to Home](./Home.md)

## WebRTC DataChannels (Primary)

Aether utilizes `node-datachannel` for high-throughput, low-latency, and decentralized data transfer.
- **Rarest-First & Sequential Downloading:** WebRTC transfers are optimized for large media formats, fetching chunks efficiently.
- **Direct Peer-to-Peer:** Unlike centralized client-server models, WebRTC establishes direct TCP/UDP connections between nodes.

## Signaling (WebSockets)

Because WebRTC requires nodes to exchange connection metadata (SDP offers/answers and ICE candidates) *before* a direct connection can be established, Aether uses WebSockets as a fallback signaling layer.
- Nodes connect to a set of known bootstrap servers (or LAN peers).
- The `join-swarm` signal broadcasts the node's intent to connect.
- `webrtc-offer`, `webrtc-answer`, and `ice-candidate` signals are relayed over WebSockets to bootstrap the direct WebRTC DataChannel.

## Local Area Network (LAN) Discovery

Aether implements **Zero-Configuration Networking** using `multicast-dns` (mDNS).
- **Service Name:** `_aether._tcp.local`
- **Behavior:** The daemon continuously queries the local subnet for other Aether nodes.
- **Offline Capabilities:** This allows Aether nodes on the same Wi-Fi or LAN to discover each other, sync routing tables, and exchange data entirely without internet access or centralized tracker servers.

> [!WARNING]
> While mDNS is excellent for LAN discovery, Wide Area Network (WAN) routing relies entirely on the [DHT](./DHT.md).
