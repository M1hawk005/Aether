# Aether

Aether is an open, local-first application and media distribution layer. It combines signed, user-owned records with the existing BitTorrent ecosystem so applications and public content remain verifiable, portable, and available across replaceable service providers.

## Project status

Aether is an architecture-migration prototype, not a production network. The repository currently contains a Go/libp2p capsule prototype, a TypeScript SDK, the Delphi proof-of-concept interface, and a Wails desktop gateway. The target architecture replaces the custom blob transport and in-memory tracker with interoperable BitTorrent components.

## Target architecture

- **Identity and mutable state:** Ed25519 identities and signed append-only events for releases, listings, profiles, follows, comments, labels, and tombstones.
- **Content distribution:** BitTorrent v2 or hybrid torrents for immutable application bundles, media, catalog snapshots, and archives.
- **Peer discovery:** Mainline BitTorrent DHT, peer exchange, optional standard trackers, and HTTP web seeds. Aether does not implement another tracker network.
- **Discovery and search:** Multiple independently operated catalog/feed providers. Clients verify records and merge selected feeds into a local search database.
- **Availability:** Volunteer peers, publisher seeders, community pinning, and optional paid seed/web-seed providers.
- **Application runtime:** Torrent-delivered web applications run from isolated origins and receive only user-approved daemon capabilities.
- **Hybrid services:** Search, recommendation, moderation, transcoding, relays, and other expensive or real-time work may be supplied by replaceable providers.

The first target application is **Aether Index**: a locally searchable, multi-provider, signed catalog of torrent-addressed applications and content.

See the [architectural reference](docs/Aether_Architecture.md) and [implementation backlog](docs/Implementation_Backlog.md).

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

The first end-to-end milestone is:

```text
build static app
  -> create interoperable torrent
  -> seed through a torrent backend
  -> sign an Aether release
  -> publish it through two catalogs
  -> synchronize and search locally
  -> download and verify it
  -> run it from an isolated local origin
```

Contributions should follow the P0 tasks in the [implementation backlog](docs/Implementation_Backlog.md). Do not extend the legacy custom tracker, global GossipSub blob delivery, or first-seen username mechanism unless the work is explicitly part of removing or migrating them.
