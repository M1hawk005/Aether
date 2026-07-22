# Aether Documentation

Aether is an open, local-first application and media distribution layer built from signed records, interoperable BitTorrent content delivery, and replaceable hybrid service providers.

The canonical system decision is [Aether Architectural Reference](../../Aether_Architecture.md). Historical documents in `docs/archive/` explain earlier thinking but are not current decisions.

## Modules

- [**Daemon**](./Daemon.md) — trusted local boundary for identity, validation, catalog sync, torrents, policy, and the application runtime.
- [**DHT**](./DHT.md) — explains the distinction between Mainline BitTorrent DHT and the legacy libp2p prototype.
- [**Network**](./Network.md) — content, record, and provider planes in the hybrid architecture.
- [**Client SDK**](./SDK.md) — current TypeScript client boundary; APIs will evolve during migration.
- [**Desktop Gateway**](./Gateway.md) — Wails shell and future isolated application runtime.
- [**CLI Tools**](./CLI.md) — current prototype commands.
- [**Project Roadmap**](./Roadmap.md) — migration phases and release gates.
- [**Implementation Backlog**](../../Implementation_Backlog.md) — P0–P3 task list for delivery planning.

## Current versus target

The repository currently runs a Go/libp2p capsule prototype with a global GossipSub topic and custom tracker-like messages. The target uses Mainline BitTorrent infrastructure for immutable blobs and selected signed feeds/providers for mutable records. New features should target the accepted architecture rather than deepen the legacy transport.
