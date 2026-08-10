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

Aether is an open, local-first runtime and distribution layer for applications and public media. It exists to reduce public reliance on centralized corporate infrastructure in two specific ways:

1. **Offload.** Move bytes to peers and move computation onto the user's own device, so the work a central server must perform shrinks toward coordination rather than delivery.
2. **Ownership.** Make identity, subscriptions, publication history, and social graph user-owned and portable, so leaving a provider costs nothing but time.

Aether is a protocol-composition project: it reuses mature infrastructure instead of replacing every transport.

### Goals

- Interoperate with the existing BitTorrent ecosystem.
- Reduce central-server bandwidth, storage, and computation for participating services.
- Keep identities, records, and social graphs portable across providers.
- Make data collection **strictly opt-in**, and make personalization work without it.
- Make application releases and public content immutable and independently verifiable.
- Let users choose catalog, ranking, moderation, availability, and compute providers.
- Provide fast local search and useful offline behavior.
- Run downloaded applications behind an explicit capability boundary.

### Non-goals

- **Eliminating servers.** A central or origin server is a permitted, expected component. The goal is to make it small, cheap, and replaceable — not absent.
- Building a general-purpose volunteer compute grid.
- Processing other users' private data on volunteer machines.
- Building a custom BitTorrent tracker or piece-exchange protocol.
- Storing full social records or search indexes in the BitTorrent DHT.
- Maintaining one global feed, reputation score, or moderation policy.
- Guaranteeing availability without seeders or an availability agreement.
- Treating torrent-delivered JavaScript as trusted code.

### The central claim

> The promise is not that servers cease to exist. It is that servers stop being the only copy, the only index, the only ranker, and the only owner — and that users and publishers can change servers without losing identities, applications, subscriptions, or verifiable history.

## 2. Architectural principles

1. **The server orders, but does not own.** Where a central service exists, it may sequence events, resolve names, and settle contention. It must not hold anything a user cannot export, re-sign, and carry to a competitor.
2. **Reuse before invention.** Mainline DHT, PEX, standard trackers, web seeds, and mature torrent engines handle blobs.
3. **Immutable bytes, mutable pointers.** Torrents contain releases and snapshots; small signed records identify current versions.
4. **Canonical data is verifiable.** Providers may rank, omit, cache, or label records, but clients verify signatures and hashes.
5. **Providers are useful but replaceable.** Hybrid services are allowed when changing provider does not lose canonical identity or content.
6. **Compute for yourself before computing for strangers.** Work performed by a user's device on that user's own data has no verification, privacy, or incentive problem. Prefer it.
7. **Local-first is a product feature.** Installed apps, catalogs, preferences, and retained content work offline where possible.
8. **No global firehose.** Clients synchronize selected publisher, community, app, and moderation feeds.
9. **Decentralization does not remove policy.** Catalogs curate, label providers assess, and clients enforce local rules.

## 3. System overview

```text
Aether desktop client
  +-- local daemon
      +-- identity, signing, and admission proofs
      +-- signed-event store and verifier
      +-- catalog sync and local search
      +-- on-device ranking and personalization
      +-- trust, moderation, and storage policy
      +-- BitTorrent backend
      +-- isolated application gateway
      +-- provider adapters

Open infrastructure
  +-- Mainline BitTorrent DHT
  +-- BitTorrent swarms, PEX, trackers, and web seeds
  +-- catalog and feed providers
  +-- optional moderation, search, availability, relay, and compute providers

Permitted central services (per deployment)
  +-- sequencer for ordering and name resolution
  +-- seeder of last resort / permanent web seed
  +-- submission, abuse handling, and legal interface
```

| Plane | Responsibility | Mechanism |
| --- | --- | --- |
| Content | Large immutable app and media bytes | BitTorrent v2/hybrid torrents |
| Records | Small mutable publication and social history | Signed append-only events |
| Ordering | Contention, naming, exclusive assignment | Permitted sequencer, exportable state |
| Services | Search, ranking, moderation, relays, compute, availability | Replaceable providers |

## 4. Workload fit model

Not every workload benefits from this architecture. Score a workload on four axes before assuming it does.

**Aether fits data that is large, immutable, public, and requested by many.** It fits poorly when data is small, mutable, private, or real-time.

| Workload | Large | Immutable | Public | Many readers | Verdict |
| --- | --- | --- | --- | --- | --- |
| App bundles, media, datasets | yes | yes | yes | yes | Ideal |
| Catalog snapshots and archives | yes | yes | yes | yes | Ideal |
| Map tiles, routing graphs, model weights | yes | yes | yes | yes | Ideal |
| Forum post text and metadata | no | no | yes | yes | Sign and own it; do not distribute by torrent |
| Personalized feeds | no | no | no | no | Compute on-device |
| Direct messages | no | no | no | no | Out of scope for v1 |
| Location pings, live dispatch | no | no | no | no | Not applicable |

Two consequences follow:

- **Offload value is proportional to the byte size of the median object.** For video, images, software, and archives, egress dominates infrastructure cost and peers absorb it. For text-heavy workloads the bandwidth saving is near zero, and the value is entirely in ownership and portability.
- **Browsers cannot join BitTorrent swarms directly.** WebTorrent peers only with other WebRTC peers, which fragments the swarm. The offload story is native/desktop-first by construction. A browser client receives a materially weaker version and must fall back to web seeds.

## 5. Identity, naming, and sybil resistance

The canonical identity is a public key. Private keys never enter application JavaScript.

### Naming

The legacy first-seen username mechanism is rejected: message arrival order and publisher timestamps cannot produce global consensus in a gossip network.

However, a permitted sequencer *can* resolve names, because ordering is exactly what a sequencer provides. Names are therefore a **provider-supplied convenience** resolved through signed delegations, provider directories, DNS/WebFinger, or an onchain registry — never a protocol-level consensus claim. A name may be lost or contested; a key never is. All canonical APIs address identities by key.

### Admission proofs

Sybil resistance does not come from any particular ledger. It comes from a scarce credential that is hard to mass-produce. Aether therefore defines a single pluggable **admission proof** carried on events, not a single mechanism.

| Method | Cost profile | Notes |
| --- | --- | --- |
| Attestation credential | Low friction, no chain required | Default path. Signed personhood/uniqueness claims from one or more attesters. |
| One-time onchain registration | Genuinely compute-light | Registry of identities, not a ledger of activity. Proven pattern: onchain registry, offchain messages. |
| Fee or burn per action | Effective against spam-scale sybils | **Regressive.** Prices out low-income participants. Reserve for high-abuse or high-value contexts. |
| Proof of work stamp | No payment rails needed | Weak against funded adversaries; useful as a rate limiter. |

Rules:

1. **Communities set their own admission policy.** Sybil resistance is per-catalog configuration, not a global protocol constant.
2. **Methods carry different weights, not just accept/reject.** A menu is only as strong as its weakest accepted method, so an identity attested by several independent attesters must outrank one that paid a trivial fee.
3. **No chain may be canonical.** An onchain registry is *one attester among several*. If a chain's fees, liveness, or governance degrade, identity must survive. This follows directly from principle 5.
4. **Per-action records never go onchain.** Vote- and post-volume workloads exceed what any compute-light chain can absorb. Chains register identities; they do not witness activity.
5. **Global sybil resistance is often unnecessary.** Within-community reputation delivers most of the ranking quality at a fraction of the cost.

## 6. Signed events

```json
{
  "version": 1,
  "type": "app_release",
  "author": "ed25519:PUBLIC_KEY",
  "sequence": 42,
  "previous": "sha256:PREVIOUS_EVENT",
  "createdAt": "2026-07-22T00:00:00Z",
  "admission": {},
  "payload": {},
  "signature": "BASE64_SIGNATURE"
}
```

Before persistence, clients validate the schema and limits, canonical event ID, signature, author sequence and predecessor, admission proof, replay status, and referenced identifiers.

Initial event types are `profile_update`, `app_release`, `torrent_listing`, `feed_recommendation`, `moderation_label`, `key_rotation`, and `tombstone`.

Four schema rules are load-bearing and must be settled before publishers sign against the format:

- **Events must be validatable under partial replication.** Requiring the full predecessor chain to accept event *n* does not survive social volume. An event must be acceptable standalone with a proof, or with predecessors fetched on demand.
- **Events carry only public-by-design content.** Append-only signatures distributed through torrent archives make deletion impossible; tombstones are delisting, not erasure. Personal data must therefore never enter an event. This is a schema-level prohibition, not a convention.
- **Counts are provider-computed opinions, not protocol facts.** Votes, likes, install counts, and rankings are assertions by a named provider over a set of events, carrying that provider's provenance. Nothing in the protocol treats an aggregate as canonical. This is what makes forum-shaped products possible later without a consensus layer.
- **Reply and reference semantics are reserved now, even though v1 does not use them.** Retrofitting threading onto a release-only schema is expensive; reserving the fields is nearly free.

## 7. Immutable content and BitTorrent

BitTorrent distributes static web applications, media, releases, datasets, catalog snapshots, and event archives. A signed Aether record references its magnet URI, v2 root or hybrid infohashes, size, file count, metadata, and publisher.

The torrent engine—not Aether's event network—handles piece selection and verification, resume state, Mainline DHT, PEX, trackers, NAT traversal, bandwidth limits, and seeding.

Mainline BitTorrent DHT is distinct from the current libp2p Kademlia DHT. Aether must use a compatible torrent backend to join the existing network.

Availability is a policy, not a consequence of content addressing. Sources may include publishers, consumer seeders, community pinning, paid seeders, trackers, HTTP web seeds, mirrors, and archives. Every source is verified against the same torrent hashes.

**A permitted origin server acts as seeder of last resort.** Peers absorb popular content; the origin guarantees the cold tail as a permanent web seed. This is the P2P-CDN model, and it resolves the long-tail availability gap that a purely volunteer swarm cannot close. Verification is unchanged: origin bytes are checked against the same hashes as peer bytes.

## 8. Catalogs, feeds, and search

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

The DHT must not store full listings or search indexes. Search is local or provider-assisted. Full-text indexing over a global corpus is expensive and therefore structurally centralizing; a search provider is a legitimate power center and must remain replaceable and provenance-labelled.

BEP 51 may later help specialized indexers sample public infohashes, but it is not a search protocol and its output is untrusted.

## 9. Data ownership and opt-in collection

Data ownership is a primary product commitment, not a side effect of decentralization.

- **Collection is strictly opt-in.** No behavioral, install, search, or content-history data leaves the device without an explicit, revocable, per-purpose grant. Defaults are off. Diagnostics are aggregate and opt-in.
- **Export is unconditional.** Identity, keys, subscriptions, grants, publication history, and local records export in a documented format at any time, without provider cooperation.
- **Personalization must not require disclosure.** This is the hard part, and it is the reason surveillance won elsewhere: ranking normally requires the ranker to observe behavior. Moving ranking to a provider relocates surveillance rather than removing it.

The architectural answer is that **the client already holds the corpus metadata it needs**. Because catalogs sync into a local full-text index, personalized ranking can run entirely on-device: the provider ships data and never receives behavior. On-device ranking simultaneously delivers the privacy guarantee, removes server ranking cost, and avoids every trust problem described in the next section.

## 10. Compute model

Volunteer computing works at scale — distributed science projects have exceeded the aggregate throughput of the largest supercomputers. But those workloads share properties most application workloads do not: embarrassingly parallel, latency-tolerant, redundancy-verifiable, and operating on non-private input.

The distinction that governs Aether's design is **who the computation is for**.

### Computing for yourself — preferred

When a user's device works on that user's own data for that user's benefit, there is nothing to verify, nothing to leak, and no incentive to engineer. Local search, on-device ranking and personalization, deduplication, preview and thumbnail generation, and integrity scanning all fall here.

This is genuine offload. A million clients ranking their own feeds is a million ranking jobs the origin never runs — the aggregation happens because everyone does their own work, not because a scheduler farms units out.

### Computing for strangers — narrowly scoped

Three problems apply, and none are solved by the fact that BitTorrent works:

1. **Verification.** Hashing a torrent piece is free; checking an arbitrary computation is often as expensive as performing it. Redundancy-and-vote costs 2–3x and fails against collusion. Verifiable-computation proofs remain orders of magnitude over native execution. Trusted execution requires specific hardware with a poor vulnerability record.
2. **Privacy.** A volunteer machine processing another user's data reconstitutes surveillance with less accountable operators. This directly contradicts the goals in section 9. Volunteer compute may touch **public data only**.
3. **Economics.** Residential electricity is frequently more expensive per compute-hour than cloud spot pricing, so naive volunteer compute can destroy value. Bandwidth sharing succeeds precisely because a flat-rate residential link has near-zero marginal cost and free verification; compute has neither.

Aether therefore accepts stranger-compute only where **the output is cheap to check**, the input is public, and latency is tolerant. Segment-parallel media transcoding is the qualifying workload: output verifiable by perceptual comparison, public source material, no deadline, and a demonstrated market. Everything else stays on the origin until the verification story improves.

### Summary

| Contribution | Marginal cost to contributor | Verification | Status |
| --- | --- | --- | --- |
| Bandwidth and storage | ~zero | free (hashes) | Core to v1 |
| Compute on your own data | negligible | not required | Core to v1 |
| Transcoding-class work | real but bounded | cheap output check | Later milestone |
| General-purpose compute | real | unsolved | Rejected |

## 11. Hybrid provider model

Providers intentionally improve latency, availability, safety, and computation without owning canonical identities or bytes. Provider types include catalog ingestion, search/ranking, moderation/malware labels, recommendations, guaranteed seeding/web seeds, notification relays, attestation, sequencing, and transcoding.

Clients can configure, disable, or replace providers. Responses claiming canonical facts contain or reference verifiable signed records. Subjective ranking results expose their provider provenance.

Two value propositions are separable and should be tracked separately, because they have different buyers:

- **Bandwidth and compute offload** — sells to operators, measured in infrastructure cost avoided, immediately provable.
- **Data ownership and portability** — sells to users and regulators, measured in exit rights.

An adopter may take either without the other. Offload is the likely wedge because it appears on a cost line; ownership is the reason the project matters.

## 12. Application runtime

Torrent-delivered applications are untrusted. They do not share an origin with the daemon, gateway, other applications, or other identities' storage.

Each release is served read-only from an origin derived from its verified app ID and release hash. Its manifest declares an entrypoint, publisher, version, content ID, limits, and requested capabilities.

Example capability classes:

- low-risk: `catalog.search`, `catalog.read`, `torrent.status`;
- user-mediated: `torrent.open`, `content.pin`, `event.publish`;
- sensitive: `identity.sign`, `filesystem.read`, `provider.configure`.

Sensitive actions use a daemon-controlled confirmation surface showing the exact operation. The gateway enforces CSP, path normalization, MIME policy, resource limits, and network policy.

**Sequencing note.** The runtime is the most novel component and the least battle-tested — it is a browser security model built in-house, and each requirement above carries a real exploit history. It is therefore sequenced *after* the search client ships, not cut. See section 13.

## 13. Delivery sequence

The first release is **Aether Index as a search client**: signed multi-catalog sync, local full-text search, on-device ranking, and verify-and-hand-off to a torrent backend. No application runtime, no capability model, no gateway.

This is deliberately narrow. It exercises identity, admission proofs, event validation, catalog sync, and local ranking under real load, and it answers the two assumptions everything else rests on: *will anyone operate a catalog, and will a publisher sign a release?*

```text
Milestone 1 — Index search client
  synchronize catalog heads
  -> download snapshots
  -> verify provider and publisher events
  -> merge into local search
  -> rank on-device
  -> hand verified torrent off to a backend

Milestone 2 — Publishing
  build artifact
  -> validate metadata
  -> create and seed torrent
  -> sign app_release
  -> submit to selected catalogs
  -> catalogs verify and publish new heads

Milestone 3 — Application runtime
  verify torrent, release, manifest, and permissions
  -> run from isolated origin
  -> enforce capabilities
  -> seed according to user policy
```

An update is a new torrent and signed event referencing its predecessor. Clients verify publisher continuity and capability changes and support update, ignore, or pin-old-version choices.

## 14. Security and abuse requirements

Before storage, execution, indexing, or credit, verify:

1. Torrent pieces and final roots.
2. Canonical event IDs and signatures.
3. Admission proofs and their weights.
4. Release-to-torrent metadata consistency.
5. File, size, piece, and nesting limits.
6. Safe relative paths without traversal, device paths, or escaping symlinks.
7. Allowed MIME types and entrypoints.
8. Replay, sequence, and oversized-message protections.
9. Origin isolation and capability grants.

The platform also needs spam admission, local blocks, signed labels, publisher history, reporting, storage quotas, and malware-scanner hooks. Delisting is not global deletion.

Public torrent participation exposes peer-network metadata including IP addresses. For an application distribution channel this means installs are observable by strangers — a materially worse privacy posture than an HTTPS download, and it must be disclosed plainly rather than footnoted. Private messaging/groups require a separate threat model for encryption, membership, key rotation, delivery, and metadata protection and are outside the first release.

A verifiable, portable catalog of torrent-addressed content will attract infringing use. The posture — per-client labels, replaceable moderation providers, delisting rather than deletion, and a legal interface at the permitted central service — should be chosen deliberately and documented before launch rather than discovered afterward.

## 15. Scaling model

- Popular immutable content scales through consumer seeding.
- Cold content is guaranteed by the origin acting as permanent web seed, supplemented by publishers, seedboxes, and availability providers.
- Feeds partition by publisher, topic, community, and time; there is no global event firehose.
- Clients retain subscribed metadata and user-selected content, not the whole network.
- Ranking and personalization scale by running on each client's own device.
- Historical records move into immutable archives.
- Mobile clients seed conservatively on suitable power and network conditions.
- Global full-text search and expensive stranger-compute may be delegated to competing providers.

## 16. Migration from the prototype

Keep the Go daemon, local API boundary, Ed25519 implementation, TypeScript SDK, Wails gateway, local-first storage concepts, and Delphi UI as the seed for Aether Index.

| Current prototype | Target |
| --- | --- |
| Custom capsule transfer | Standard BitTorrent backend |
| In-memory `BTTrackerSwarm` | Mainline DHT, PEX, standard trackers |
| `BT_ANNOUNCE` over GossipSub | Torrent-client announces |
| Global `aether-global-v1` topic | Selected scoped feeds/providers |
| First-seen usernames | Public-key identity with provider-resolved names |
| No sybil resistance | Pluggable weighted admission proofs |
| JSON/base64 chunks | Torrent piece/resume storage |
| Transfer leaderboard | Local availability metrics and signed labels |
| Custom transport/NAT roadmap | Mature torrent engine and optional relays |

Legacy code may remain behind compatibility boundaries while the new path is built, but new features must not depend on it.

## 17. Decision summary

Aether is hybrid by design. A central service may sequence, resolve names, guarantee the cold tail, and handle abuse and legal process — because those are the functions that genuinely require a coordinator. What it may not do is own the user's identity, records, graph, or history, or collect behavior without an explicit grant.

The offload is real but specific: bytes move to peers, and personalization moves to the user's own device. General-purpose volunteer compute is not part of the design, because verification and residential energy economics do not support it. Transcoding-class work is the one exception worth pursuing.

The result is not the absence of servers. It is servers that are small, cheap, contestable, and exit-safe.
