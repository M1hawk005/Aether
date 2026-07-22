# Project Roadmap

[Back to Home](./Home.md) · [Detailed P0–P3 backlog](../../Implementation_Backlog.md)

## Starting point

The repository already proves several useful concepts: a Go daemon and local API, Ed25519 identity/signing, a TypeScript SDK, a Wails desktop shell, a Delphi catalog/tracker UI, and local publication/retrieval flows.

The current custom libp2p/GossipSub blob transport, in-memory BitTorrent-style tracker, first-seen username claims, and JSON/base64 chunks are prototype components to be replaced or isolated—not production foundations.

## Phase P0 — Foundation

Define and test the boundaries before extending product behavior:

- canonical event serialization and schemas;
- public-key identity and key rotation;
- torrent backend abstraction and qBittorrent proof of concept;
- catalog provider API and local database schema;
- application manifest, origin isolation, and capabilities;
- mandatory hash/signature/path/size validation;
- shared hostile and valid conformance fixtures.

**Exit gate:** Go and TypeScript agree on signed records; an untrusted bundle cannot reach privileged APIs; all network content is verified before use.

## Phase P1 — Aether Index MVP

Deliver the complete first application lifecycle:

- publish a static app as a torrent and signed release;
- submit it to two independent catalog providers;
- synchronize catalogs and search locally;
- download through DHT/tracker/PEX/web seed;
- install and run from an isolated origin;
- approve explicit capabilities;
- publish and install a signed update;
- retain operation with one provider unavailable.

**Exit gate:** the end-to-end flow works on clean machines without hidden database edits.

## Phase P2 — Public beta

Make the system suitable for external users:

- embed a production torrent engine;
- add signed mutable catalog-head pointers and archive compaction;
- implement anti-spam, signed moderation labels, and malware hooks;
- add publisher continuity and key-rotation UX;
- support guaranteed availability providers and media streaming;
- harden local API authentication and runtime containment;
- ship signed cross-platform installers;
- run adversarial, NAT, failure, and load tests.

**Exit gate:** the system fails safely under hostile content and provider/network outages and remains usable on ordinary consumer connections.

## Phase P3 — Scale and ecosystem

Add optional capabilities after the core is proven:

- provider discovery and federated search;
- community pinning and provider marketplace abstractions;
- privacy-preserving creator analytics;
- mobile consumption client;
- social events and notification relays;
- provider-based transcoding;
- research for BEP 51 indexing, private groups, and live streaming;
- public conformance suite and protocol governance.

These are not allowed to block or weaken the first release's security and portability guarantees.
