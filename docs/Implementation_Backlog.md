---
title: Implementation Backlog
sidebar_position: 2
---

# Aether Implementation Backlog

This is the canonical backlog for the hybrid Aether architecture. It is formatted for copying into a Notion Kanban. Recommended Notion properties are `ID`, `Task`, `Priority`, `Status`, `Area`, `Milestone`, `Depends on`, and `Acceptance criteria`.

- **P0 — Blocker:** establish the architecture and security boundaries.
- **P1 — MVP:** deliver the first useful Aether Index release.
- **P2 — Beta:** make it reliable and safe for external users.
- **P3 — Scale:** grow the ecosystem and support advanced applications.

All tasks start in `Backlog`. Suggested Kanban states are `Backlog`, `Ready`, `In Progress`, `Review`, `Blocked`, and `Done`.

## P0 — Architecture and security blockers

| ID | Task | Area | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- |
| AET-001 | Freeze v1 scope and non-goals | Architecture | — | Review accepts Aether Index scope; private messaging, live streaming, and global search are explicitly deferred. |
| AET-002 | Define canonical signed-event serialization | Protocol | AET-001 | Go and TypeScript test vectors produce identical IDs and signatures. |
| AET-003 | Define v1 event schemas | Protocol | AET-002 | Schemas cover profile, app release, torrent listing, feed recommendation, label, key rotation, and tombstone events. |
| AET-004 | Implement strict event validation | Security | AET-002, AET-003 | Invalid signatures, sizes, sequences, predecessor links, required fields, and replays are rejected before persistence. |
| AET-005 | Replace first-seen usernames with public-key identities | Identity | AET-003 | Canonical APIs use public keys; aliases are non-consensus and cannot replace identity history. |
| AET-006 | Define torrent backend interface | Torrent | AET-001 | Interface covers create, add, stream, pause, remove, status, files, limits, trackers, web seeds, and seeding policy. |
| AET-007 | Implement qBittorrent proof-of-concept adapter | Torrent | AET-006 | Daemon can add a magnet, observe progress, select files, and seed through the local Web API with integration tests. |
| AET-008 | Define application manifest v1 | Runtime | AET-001, AET-003 | Manifest specifies app ID, version, entrypoint, content ID, publisher, capabilities, and resource limits. |
| AET-009 | Implement immutable app storage and safe paths | Security | AET-007, AET-008 | Traversal, device/absolute paths, escaping symlinks, and excessive size/file counts are rejected. |
| AET-010 | Implement per-application origin isolation | Runtime | AET-008, AET-009 | Apps cannot access another app, daemon cookies, gateway DOM, or privileged API origin. |
| AET-011 | Design capability API and permission model | Runtime | AET-008, AET-010 | Capabilities are deny-by-default, versioned, revocable, and classified by risk. |
| AET-012 | Remove hashless credit and storage paths | Security | AET-004, AET-009 | No payload is stored, executed, indexed, or credited before hashes and signatures are verified. |
| AET-013 | Define catalog provider API v1 | Catalog | AET-003 | API supports head resolution, incremental events, snapshots, submission, pagination, limits, and signed references. |
| AET-014 | Define local database schema and migrations | Storage | AET-003, AET-013 | Schema covers events, feeds, listings, releases, labels, blocks, grants, torrents, and full-text search. |
| AET-015 | Add architecture conformance fixtures | Quality | AET-002–AET-014 | Shared fixtures cover valid and hostile events, manifests, paths, responses, and release/torrent mismatches. |

## P1 — Aether Index MVP

| ID | Task | Area | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- |
| AET-101 | Implement verified signed-event store | Protocol | AET-004, AET-014 | Events deduplicate by ID, retain provenance, and expose verified author chains. |
| AET-102 | Implement identity create/import/export/backup | Identity | AET-005 | Users can create and restore identities; apps never receive raw private keys. |
| AET-103 | Implement app bundle publisher | Publishing | AET-007–AET-009 | A static directory becomes validated torrent metadata plus a signed release in one workflow. |
| AET-104 | Implement reference catalog provider | Catalog | AET-013, AET-101 | Provider accepts valid submissions, rejects invalid records, and publishes a signed head. |
| AET-105 | Publish catalog snapshots as torrents | Catalog | AET-007, AET-104 | Current head and history can be reconstructed through the torrent backend. |
| AET-106 | Implement feed subscriptions and synchronization | Catalog | AET-101, AET-104, AET-105 | Client syncs two providers, resumes updates, retains provenance, and tolerates one unavailable provider. |
| AET-107 | Build local full-text catalog search | Search | AET-014, AET-106 | Search works offline over title, description, tags, publisher, and content type. |
| AET-108 | Evolve Delphi into Aether Index UI | Product | AET-107 | Users can search, inspect provenance/security metadata, install, launch, pin, and remove an app. |
| AET-109 | Implement verified install pipeline | Runtime | AET-009–AET-012, AET-103 | Install verifies torrent, release, publisher, manifest, limits, and permissions before activation. |
| AET-110 | Implement capability bridge MVP | Runtime | AET-011, AET-109 | Index can search and request torrent operations without direct daemon or key access. |
| AET-111 | Implement daemon-controlled permission prompts | Security | AET-011, AET-110 | Prompt shows app, publisher, capability, scope, and supports deny, once, or persistent grant. |
| AET-112 | Implement signed application updates | Runtime | AET-101, AET-103, AET-109 | Client detects successors, highlights permission changes, and supports update, ignore, and pin. |
| AET-113 | Implement storage and seeding controls | Storage | AET-007, AET-014 | Users set disk/bandwidth limits, pin, stop seeding, remove data, and see deletion consequences. |
| AET-114 | Add DHT, tracker, PEX, and web-seed source policy | Torrent | AET-006, AET-007 | Downloads use configured compatible sources and report successful paths without making one canonical. |
| AET-115 | Build two-publisher/two-provider end-to-end demo | Quality | AET-101–AET-114 | Clean clients publish, discover, install, run, update, and reseed while one provider is offline. |
| AET-116 | Write Index operator and user guides | Documentation | AET-115 | Guides cover publishing, provider setup, permissions, backup, recovery, seeding, and troubleshooting. |

## P2 — Public beta reliability

| ID | Task | Area | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- |
| AET-201 | Embed a production torrent engine | Torrent | AET-115 | Desktop package joins Mainline DHT, downloads, streams, resumes, and seeds without a separate client. |
| AET-202 | Implement signed mutable catalog-head pointers | Catalog | AET-105, AET-201 | Heads resolve without a mandatory provider endpoint; rollback and expiration have interoperability tests. |
| AET-203 | Add catalog head/archive compaction | Catalog | AET-105, AET-106 | Recent sync stays bounded while complete history remains recoverable. |
| AET-204 | Add feed admission and anti-spam controls | Abuse | AET-104 | Providers support rate, size, age, proof, allowlist, and review policies with clear rejection reasons. |
| AET-205 | Implement signed moderation labels | Moderation | AET-003, AET-106 | Users combine multiple label providers with local allow, warn, hide, and block rules. |
| AET-206 | Add malware-scanning provider hooks | Security | AET-109, AET-205 | Install pipeline displays multi-provider verdict provenance/age and applies local policy. |
| AET-207 | Add publisher continuity and key-rotation UX | Identity | AET-003, AET-102 | Valid rotations preserve history; unexpected publisher changes produce a blocking warning. |
| AET-208 | Implement availability-provider contracts | Availability | AET-114 | Publishers attach replaceable seed/web-seed sources and clients verify all bytes against torrent hashes. |
| AET-209 | Add streaming and file-priority support | Media | AET-201 | Media starts before full download, seeks within a defined target, and falls back across sources. |
| AET-210 | Harden local API authentication | Security | AET-010, AET-110 | Other processes and arbitrary websites cannot invoke privileged daemon APIs without authorization. |
| AET-211 | Add runtime resource containment | Runtime | AET-109, AET-210 | Apps have enforced storage, request, CPU/time, and network policies with usable failures. |
| AET-212 | Add reproducible-release attestations | Supply chain | AET-103, AET-112 | UI distinguishes publisher signatures from source/build reproducibility evidence. |
| AET-213 | Build cross-platform signed installers and updates | Release | AET-201, AET-210 | Windows, macOS, and Linux packages update safely and preserve identity/database state. |
| AET-214 | Add privacy and exposure disclosures | Privacy | AET-201 | Users understand peer-IP exposure, seeding, storage, and provider requests before enabling them. |
| AET-215 | Run adversarial and load testing | Quality | AET-203–AET-214 | Tests cover malicious torrents, event floods, equivocation, DHT failure, NAT diversity, disk pressure, and large catalogs. |
| AET-216 | Add opt-in diagnostics and operational dashboards | Operations | AET-215 | Measure availability, sync latency, rejections, and failures without collecting content history by default. |

## P3 — Scale and ecosystem

| ID | Task | Area | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- |
| AET-301 | Define provider discovery events | Ecosystem | AET-202, AET-205 | Users import, compare, verify, and remove provider recommendations without a canonical directory. |
| AET-302 | Implement federated search providers | Search | AET-203, AET-216 | Clients merge signed references, expose ranking provenance, and fall back to local search. |
| AET-303 | Evaluate BEP 51 infohash indexing | Research | AET-201, AET-204 | Prototype measures coverage, cost, metadata quality, abuse, and operational risk before commitment. |
| AET-304 | Add privacy-preserving creator analytics | Publishing | AET-208, AET-216 | Creators receive useful aggregate availability/install data without mandatory user tracking. |
| AET-305 | Add provider marketplace/payment abstraction | Ecosystem | AET-208, AET-301 | Users can buy availability/compute from competing providers without protocol consensus. |
| AET-306 | Add community pinning and retention policies | Availability | AET-208 | Communities publish signed policies; clients opt in with explicit disk/bandwidth budgets. |
| AET-307 | Build mobile consumption client | Mobile | AET-213, AET-214 | Mobile defaults to conservative Wi-Fi/charging seeding and preserves identity/catalog portability. |
| AET-308 | Add post, reply, follow, and reaction events | Social | AET-203–AET-205 | Social state syncs through selected feeds with signature, spam, moderation, and offline tests. |
| AET-309 | Add replaceable notification relays | Realtime | AET-308 | Relays are authenticated and rate-limited and cannot forge underlying events. |
| AET-310 | Add provider-based transcoding | Compute | AET-208, AET-209, AET-301 | Derived media links to source and publisher authorization through signed records. |
| AET-311 | Research private groups and messaging | Research | AET-210, AET-309 | Threat model covers encryption, membership, rotation, forward secrecy, delivery, and metadata leakage. |
| AET-312 | Research low-latency live streaming | Research | AET-209, AET-309 | Prototype compares WebRTC/QUIC/relay options and torrents completed-stream archives. |
| AET-313 | Publish protocol conformance suite | Ecosystem | AET-215 | Independent clients can validate events, manifests, feeds, permissions, and provider behavior. |
| AET-314 | Establish protocol governance/versioning | Governance | AET-313 | Extensions, compatibility, deprecation, ownership, and security response are documented. |

## Definition of the first useful release

P0 and P1 are complete only when this works without hidden database edits:

```text
publisher builds static app
  -> daemon creates and seeds torrent
  -> publisher signs release
  -> two providers publish the listing
  -> clean client synchronizes both providers
  -> client searches locally
  -> client downloads and verifies app
  -> app runs in isolated origin with explicit capabilities
  -> publisher releases an update
  -> client verifies and installs update
  -> flow still works with one provider unavailable
```
