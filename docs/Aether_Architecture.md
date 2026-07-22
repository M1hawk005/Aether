---
title: Architecture
slug: /
sidebar_position: 1
---

# Aether Architectural Reference

**Status:** Accepted target architecture

**Implementation status:** Migration from the existing Go/libp2p capsule prototype

**First application:** Aether Index

## 1. Product definition

Aether is an open, local-first runtime and distribution layer for applications and public media. Users control their identity, subscriptions, installed applications, local catalog, and provider choices. Publishers receive signed releases and torrent-backed distribution without making a single catalog, frontend, or hosting company canonical.

Aether is a protocol-composition project: it reuses mature infrastructure instead of replacing every transport.

### Goals

- Interoperate with the existing BitTorrent ecosystem.
- Make application releases and public content immutable and independently verifiable.
- Keep identities and mutable records portable across providers.
- Let users choose catalog, ranking, moderation, availability, and compute providers.
- Provide fast local search and useful offline behavior.
- Run downloaded applications behind an explicit capability boundary.
- Use conventional infrastructure when it improves reliability without making it canonical.

### Non-goals

- Replacing every internet backend.
- Building a custom BitTorrent tracker or piece-exchange protocol.
- Storing full social records or search indexes in the BitTorrent DHT.
- Maintaining one global feed, username registry, reputation score, or moderation policy.
- Guaranteeing availability without seeders or an availability agreement.
- Treating torrent-delivered JavaScript as trusted code.

## 2. Architectural principles

1. **Reuse before invention.** Mainline DHT, PEX, standard trackers, web seeds, and mature torrent engines handle blobs.
2. **Immutable bytes, mutable pointers.** Torrents contain releases and snapshots; small signed records identify current versions.
3. **Canonical data is verifiable.** Providers may rank, omit, cache, or label records, but clients verify signatures and hashes.
4. **Providers are useful but replaceable.** Hybrid services are allowed when changing provider does not lose canonical identity or content.
5. **Local-first is a product feature.** Installed apps, catalogs, preferences, and retained content work offline where possible.
6. **No global firehose.** Clients synchronize selected publisher, community, app, and moderation feeds.
7. **Decentralization does not remove policy.** Catalogs curate, label providers assess, and clients enforce local rules.

## 3. System overview

```text
Aether desktop client
  +-- local daemon
      +-- identity and signing
      +-- signed-event store and verifier
      +-- catalog sync and local search
      +-- trust, moderation, and storage policy
      +-- BitTorrent backend
      +-- isolated application gateway
      +-- provider adapters

Open infrastructure
  +-- Mainline BitTorrent DHT
  +-- BitTorrent swarms, PEX, trackers, and web seeds
  +-- catalog and feed providers
  +-- optional moderation, search, availability, relay, and compute providers
```

| Plane | Responsibility | Mechanism |
| --- | --- | --- |
| Content | Large immutable app and media bytes | BitTorrent v2/hybrid torrents |
| Records | Small mutable publication and social history | Signed append-only events |
| Services | Search, ranking, moderation, relays, compute, availability | Replaceable providers |

## 4. Local daemon

The Go daemon is the trusted boundary between applications, keys, storage, transports, and providers. The desktop UI and installed applications use a versioned local API rather than accessing keys or the torrent engine directly.

Target services:

- **Identity:** create, import, rotate, and use Ed25519 keys; private keys never enter application JavaScript.
- **Events:** validate schemas, sequence numbers, predecessor links, signatures, sizes, and replay rules.
- **Catalog:** synchronize selected feeds and maintain a local full-text index.
- **Torrent:** create/add torrents, stream prioritized files, report availability, and seed within user policy.
- **Policy:** apply feed trust, moderation labels, blocks, quotas, and application permissions.
- **Application gateway:** serve verified bundles from isolated origins.
- **Provider adapters:** communicate with selected catalogs, labels, search, relay, availability, and compute services.

The first torrent integration may use qBittorrent's local Web API. The self-contained desktop target should embed libtorrent or another mature compatible engine. The backend is an implementation choice, not an Aether protocol dependency.

## 5. Identity and signed events

The canonical identity is a public key. Handles are optional aliases resolved through signed delegations, provider directories, DNS/WebFinger, or another non-consensus convenience layer.

The legacy first-seen username mechanism is rejected: message arrival order and publisher timestamps cannot produce global consensus.

```json
{
  "version": 1,
  "type": "app_release",
  "author": "ed25519:PUBLIC_KEY",
  "sequence": 42,
  "previous": "sha256:PREVIOUS_EVENT",
  "createdAt": "2026-07-22T00:00:00Z",
  "payload": {},
  "signature": "BASE64_SIGNATURE"
}
```

Before persistence, clients validate the schema and limits, canonical event ID, signature, author sequence and predecessor, replay status, and referenced identifiers. Initial event types are `profile_update`, `app_release`, `torrent_listing`, `feed_recommendation`, `moderation_label`, `key_rotation`, and `tombstone`. General social events follow after the release/catalog lifecycle is secure.

## 6. Immutable content and BitTorrent

BitTorrent distributes static web applications, media, releases, datasets, catalog snapshots, and event archives. A signed Aether record references its magnet URI, v2 root or hybrid infohashes, size, file count, metadata, and publisher.

The torrent engine—not Aether's event network—handles piece selection and verification, resume state, Mainline DHT, PEX, trackers, NAT traversal, bandwidth limits, and seeding.

Mainline BitTorrent DHT is distinct from the current libp2p Kademlia DHT. Aether must use a compatible torrent backend to join the existing network.

Availability is a policy, not a consequence of content addressing. Sources may include publishers, consumer seeders, community pinning, paid seeders, trackers, HTTP web seeds, mirrors, and archives. Every source is verified against the same torrent hashes.

## 7. Catalogs, feeds, and search

A tracker finds peers for a known infohash. A catalog describes content. A search engine queries catalog metadata. **Aether Index is a decentralized catalog and local search application, not another tracker.**

Each catalog is independently operated and signed. Clients subscribe to multiple catalogs and label providers, verify entries, deduplicate by event/content ID, apply local policy, and import accepted metadata into a local SQLite full-text index.

Catalog history is distributed as torrent snapshots. A small signed mutable pointer identifies the current head. The initial implementation may use HTTPS provider APIs; BEP 44/46-style DHT mutable pointers can follow after interoperability and republishing behavior are proven.

```text
signed catalog head pointer
  -> catalog-head torrent
      -> recent signed listings
      -> moderation labels
      -> previous head
      -> archive torrent references
```

The DHT must not store full listings or search indexes. Search is local or provider-assisted. BEP 51 may later help specialized indexers sample public infohashes, but it is not a search protocol and its output is untrusted.

## 8. Hybrid provider model

Providers intentionally improve latency, availability, safety, and computation without owning canonical identities or bytes. Provider types include catalog ingestion, search/ranking, moderation/malware labels, recommendations, guaranteed seeding/web seeds, notification relays, and transcoding.

Clients can configure, disable, or replace providers. Responses claiming canonical facts contain or reference verifiable signed records. Subjective ranking results expose their provider provenance.

## 9. Application runtime

Torrent-delivered applications are untrusted. They do not share an origin with the daemon, gateway, other applications, or other identities' storage.

Each release is served read-only from an origin derived from its verified app ID and release hash. Its manifest declares an entrypoint, publisher, version, content ID, limits, and requested capabilities.

Example capability classes:

- low-risk: `catalog.search`, `catalog.read`, `torrent.status`;
- user-mediated: `torrent.open`, `content.pin`, `event.publish`;
- sensitive: `identity.sign`, `filesystem.read`, `provider.configure`.

Sensitive actions use a daemon-controlled confirmation surface showing the exact operation. The gateway enforces CSP, path normalization, MIME policy, resource limits, and network policy.

## 10. Aether Index flow

```text
Publish
  build static app
  -> validate manifest
  -> create and seed torrent
  -> sign app_release
  -> submit to selected catalogs
  -> catalogs verify and publish new heads

Discover and run
  synchronize catalog heads
  -> download snapshots
  -> verify provider and publisher events
  -> merge into local search
  -> download selected torrent through DHT/tracker/PEX/web seed
  -> verify torrent, release, manifest, and permissions
  -> run from isolated origin
  -> seed according to user policy
```

An update is a new torrent and signed event referencing its predecessor. Clients verify publisher continuity and capability changes and support update, ignore, or pin-old-version choices.

## 11. Security and abuse requirements

Before storage, execution, indexing, or credit, verify:

1. Torrent pieces and final roots.
2. Canonical event IDs and signatures.
3. Release-to-torrent metadata consistency.
4. File, size, piece, and nesting limits.
5. Safe relative paths without traversal, device paths, or escaping symlinks.
6. Allowed MIME types and entrypoints.
7. Replay, sequence, and oversized-message protections.
8. Origin isolation and capability grants.

The platform also needs spam admission, local blocks, signed labels, publisher history, reporting, storage quotas, and malware-scanner hooks. Delisting is not global deletion.

Public torrent participation exposes peer-network metadata including IP addresses. Private messaging/groups require a separate threat model for encryption, membership, key rotation, delivery, and metadata protection and are outside the first release.

## 12. Scaling model

- Popular immutable content scales through consumer seeding.
- Cold content depends on publishers, seedboxes, web seeds, or availability providers.
- Feeds partition by publisher, topic, community, and time; there is no global event firehose.
- Clients retain subscribed metadata and user-selected content, not the whole network.
- Historical records move into immutable archives.
- Mobile clients seed conservatively on suitable power and network conditions.
- Search and expensive computation may be local or delegated to competing providers.

## 13. Migration from the prototype

Keep the Go daemon, local API boundary, Ed25519 implementation, TypeScript SDK, Wails gateway, local-first storage concepts, and Delphi UI as the seed for Aether Index.

| Current prototype | Target |
| --- | --- |
| Custom capsule transfer | Standard BitTorrent backend |
| In-memory `BTTrackerSwarm` | Mainline DHT, PEX, standard trackers |
| `BT_ANNOUNCE` over GossipSub | Torrent-client announces |
| Global `aether-global-v1` topic | Selected scoped feeds/providers |
| First-seen usernames | Public-key identity with optional aliases |
| JSON/base64 chunks | Torrent piece/resume storage |
| Transfer leaderboard | Local availability metrics and signed labels |
| Custom transport/NAT roadmap | Mature torrent engine and optional relays |

Legacy code may remain behind compatibility boundaries while the new path is built, but new features must not depend on it.

## 14. Decision summary

Aether will be hybrid where the benefits are material. HTTP providers and web seeds may accelerate discovery and availability; professional services may provide moderation, search, relays, and compute. They remain replaceable, while clients verify canonical events and content.

The promise is not that servers cease to exist. It is that users and publishers can change servers without losing identities, applications, subscriptions, or verifiable content history.
