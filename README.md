# Aether

Aether is an open, local-first application and media distribution layer. It combines signed, user-owned records with the existing BitTorrent ecosystem so applications and public content remain verifiable, portable, and available across replaceable service providers.

Aether exists to reduce public reliance on centralized corporate infrastructure in two specific ways: **offload** — moving bytes to peers and computation onto the user's own device — and **ownership** — making identity, subscriptions, and history portable, with data collection strictly opt-in.

Aether is not an attempt to eliminate servers. A central or origin server is a permitted, expected component. The goal is to make it small, cheap, and replaceable. The guiding rule is *the server orders, but does not own*.

## Project status

Aether is an architecture-migration prototype, not a production network. The repository currently contains a Go/libp2p capsule prototype, a TypeScript SDK, the Delphi proof-of-concept interface, and a Wails desktop gateway. The target architecture replaces the custom blob transport and in-memory tracker with interoperable BitTorrent components.

## Target architecture

- **Identity and mutable state:** Ed25519 identities and signed append-only events for releases, listings, profiles, follows, comments, labels, and tombstones. Names are provider-resolved conveniences; keys are canonical.
- **Sybil resistance:** a pluggable, weighted admission proof — attestation credentials by default, optional one-time onchain registration, and fee/burn or proof-of-work where the friction is justified. No chain is canonical, and per-action records never go onchain.
- **Content distribution:** BitTorrent v2 or hybrid torrents for immutable application bundles, media, catalog snapshots, and archives.
- **Peer discovery:** Mainline BitTorrent DHT, peer exchange, optional standard trackers, and HTTP web seeds. Aether does not implement another tracker network.
- **Discovery and search:** Multiple independently operated catalog/feed providers. Clients verify records and merge selected feeds into a local search database.
- **Ranking and personalization:** computed on-device against the locally synced corpus, so providers ship data and never receive behavior.
- **Availability:** volunteer peers, publisher seeders, community pinning, optional paid providers, and a permitted origin server acting as seeder of last resort for the cold tail.
- **Application runtime:** torrent-delivered web applications run from isolated origins and receive only user-approved daemon capabilities. Sequenced after the search client ships.
- **Hybrid services:** search, recommendation, moderation, transcoding, relays, sequencing, and other expensive or real-time work may be supplied by replaceable providers.

The first target application is **Aether Index**: a locally searchable, multi-provider, signed catalog of torrent-addressed applications and content.

See the [architectural reference](docs/Aether_Architecture.md) and [implementation backlog](docs/Implementation_Backlog.md).

## What Aether is good at

Aether fits data that is **large, immutable, public, and requested by many** — app bundles, media, datasets, catalog snapshots, map tiles, model weights. It fits poorly when data is small, mutable, private, or real-time.

Offload value is proportional to the byte size of the median object: transformative for media and software, near zero for text-heavy workloads, where the value is entirely ownership and portability. Browsers cannot join BitTorrent swarms directly, so the offload story is native/desktop-first by construction.

On computation, the governing distinction is *who the work is for*. Computing on your own data — local search, on-device ranking, previews, integrity scanning — has no verification, privacy, or incentive problem and is core to v1. Computing for strangers faces unsolved verification cost, a privacy contradiction, and residential energy economics that frequently make it value-destroying; it is accepted only where output is cheap to check, input is public, and latency is tolerant. Transcoding-class work is the qualifying case.

## Repository layout

- **`src/daemon-go/`**: Current Go daemon and legacy libp2p/capsule prototype; migration target for identity, feed synchronization, torrent integration, storage, policy, and the local API.
- **`packages/sdk/`**: TypeScript client SDK for applications communicating with the local daemon.
- **`apps/delphi/`**: Existing tracker/catalog proof of concept; intended to evolve into Aether Index.
- **`aether-gateway/`**: Wails desktop gateway and future secure application runtime.
- **`docs/`**: Canonical architecture, implementation backlog, and developer guide.
- **`aether-docs/`**: Docusaurus site that renders the canonical files from `docs/`.

## Running the current prototype

The commands below run the existing prototype and do not yet demonstrate the target BitTorrent architecture.

```powershell
cd src/daemon-go
go run ./cmd/aether-daemon
```

In another terminal:

```powershell
cd packages/sdk
npm install
npm run build
```

The current Delphi UI is served at `http://127.0.0.1:5000/ui/`. The Wails gateway can be started separately from `aether-gateway/` with `wails dev` when the Wails toolchain is installed.

## Near-term objective

The first release is **Aether Index as a search client only** — no application runtime, no capability model, no gateway. It is deliberately narrow: it exercises identity, admission proofs, event validation, catalog sync, and on-device ranking under real load, and it answers the two assumptions everything else depends on — *will anyone operate a catalog, and will a publisher sign a release?*

```text
synchronize catalog heads
  -> download snapshots
  -> verify provider and publisher events
  -> merge into local search
  -> rank on-device
  -> hand verified torrent off to a backend
```

Publishing is milestone 2; the isolated application runtime is milestone 3.

Contributions should follow the P0 tasks in the [implementation backlog](docs/Implementation_Backlog.md). Do not extend the legacy custom tracker, global GossipSub blob delivery, or first-seen username mechanism unless the work is explicitly part of removing or migrating them.
