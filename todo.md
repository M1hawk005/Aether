# Aether Project Todo

## V2 Architecture
- [ ] **Replace WebRTC:** Rip out the WebRTC dependency in the native daemon and replace it with a lightweight, custom raw UDP/TCP + UPnP protocol (similar to BitTorrent's uTP) to make the daemon ultra-lean.

## Decentralized Tracker App
- [x] Implement mutable state / signed feeds for tracker swarm announcements.
- [x] Build the Tracker UI (Web App interfacing with the Aether Daemon SDK).
