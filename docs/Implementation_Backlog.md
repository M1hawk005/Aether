---
title: Implementation Backlog
sidebar_position: 2
---

# Aether Implementation Backlog

This is the canonical backlog for the hybrid Aether architecture. It is formatted for copying into a Notion Kanban. Recommended Notion properties are `ID`, `Task`, `Priority`, `Status`, `Area`, `Milestone`, `Depends on`, and `Acceptance criteria`.

- **P0 — Blocker:** establish the architecture and security boundaries.
- **P1 — MVP:** deliver the first useful Aether Index release.
- **P2 — Beta:** publish, run, and harden for external users.
- **P3 — Scale:** grow the ecosystem and support advanced applications.

Milestones sequence the work within those priorities:

| Milestone | Scope |
| --- | --- |
| M0 | Architecture contracts and security foundations |
| M1 | **Aether Index search client** — sync, verify, search, rank on-device, hand off |
| M2 | Publishing and catalog operation |
| M3 | Application runtime and capability model |
| M4 | Beta reliability and hardening |
| M5 | Scale and ecosystem |

Task IDs are stable. Where a task has moved milestone relative to earlier revisions of this backlog, the ID is unchanged so existing Notion records remain valid.

All tasks start in `Backlog`. Suggested Kanban states are `Backlog`, `Ready`, `In Progress`, `Review`, `Blocked`, and `Done`.

## P0 — Architecture and security blockers (M0)

| ID | Task | Area | Milestone | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- | --- |
| AET-001 | Freeze v1 scope and non-goals | Architecture | M0 | — | Review accepts a **search-client-only** v1. Application runtime, capability model, publishing, private messaging, live streaming, and global search are explicitly deferred to later milestones. |
| AET-002 | Define canonical signed-event serialization | Protocol | M0 | AET-001 | Go and TypeScript test vectors produce identical IDs and signatures. |
| AET-003 | Define v1 event schemas | Protocol | M0 | AET-002 | Schemas cover profile, app release, torrent listing, feed recommendation, label, key rotation, and tombstone events, and carry an `admission` field. |
| AET-004 | Implement strict event validation | Security | M0 | AET-002, AET-003 | Invalid signatures, sizes, sequences, predecessor links, admission proofs, required fields, and replays are rejected before persistence. |
| AET-005 | Replace first-seen usernames with public-key identities | Identity | M0 | AET-003 | Canonical APIs address identities by key. Names resolve through a provider/sequencer and cannot replace or rewrite identity history. |
| AET-006 | Define torrent backend interface | Torrent | M0 | AET-001 | Interface covers create, add, stream, pause, remove, status, files, limits, trackers, web seeds, and seeding policy. |
| AET-007 | Implement qBittorrent proof-of-concept adapter | Torrent | M0 | AET-006 | Daemon can add a magnet, observe progress, select files, and seed through the local Web API with integration tests. |
| AET-009 | Implement verified content storage and safe paths | Security | M0 | AET-007 | Traversal, device/absolute paths, escaping symlinks, and excessive size/file counts are rejected for any extracted content. |
| AET-012 | Remove hashless credit and storage paths | Security | M0 | AET-004, AET-009 | No payload is stored, executed, indexed, or credited before hashes and signatures are verified. |
| AET-013 | Define catalog provider API v1 | Catalog | M0 | AET-003 | API supports head resolution, incremental events, snapshots, submission, pagination, limits, and signed references. |
| AET-014 | Define local database schema and migrations | Storage | M0 | AET-003, AET-013 | Schema covers events, feeds, listings, releases, labels, blocks, grants, torrents, admission proofs, and full-text search. |
| AET-015 | Add architecture conformance fixtures | Quality | M0 | AET-002–AET-014 | Shared fixtures cover valid and hostile events, paths, responses, admission proofs, and release/torrent mismatches. |
| AET-016 | Define admission-proof interface and weighting | Identity | M0 | AET-003 | One pluggable proof field accepts attestation, onchain registration, fee/burn, and PoW forms. Methods carry distinct weights; catalogs set their own accepted set and thresholds; no method is protocol-mandatory. |
| AET-017 | Define partial-replication validation rule | Protocol | M0 | AET-002, AET-003 | An event is validatable without holding its full predecessor chain, via standalone proof or on-demand predecessor fetch. Documented with test vectors for gaps and out-of-order arrival. |
| AET-018 | Enforce public-by-design event content | Privacy | M0 | AET-003, AET-004 | Schema-level prohibition on personal data in events, with validation rejecting known-sensitive fields. Rationale documents that tombstones are delisting, not erasure. |
| AET-019 | Reserve reply and reference semantics | Protocol | M0 | AET-003 | Threading/reference fields are reserved and version-tolerant in v1 schemas even though v1 emits none. |
| AET-020 | Define aggregate and count provenance model | Protocol | M0 | AET-003, AET-013 | Votes, installs, and rankings are represented as provider-signed assertions over event sets. No aggregate is treated as a protocol fact anywhere in the stack. |
| AET-021 | Define data-ownership and opt-in collection policy | Privacy | M0 | AET-001 | Written policy: collection defaults off, grants are explicit/per-purpose/revocable, export is unconditional, and personalization must not require disclosure. Binding on all later tasks. |

## P1 — Milestone 1: Aether Index search client

The first release synchronizes signed catalogs, verifies them, searches locally, ranks on-device, and hands verified torrents to an existing backend. It does **not** publish, install, or execute applications.

| ID | Task | Area | Milestone | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- | --- |
| AET-101 | Implement verified signed-event store | Protocol | M1 | AET-004, AET-014, AET-017 | Events deduplicate by ID, retain provenance, expose verified author chains, and accept valid events under partial replication. |
| AET-102 | Implement identity create/import/export/backup | Identity | M1 | AET-005 | Users can create and restore identities; apps never receive raw private keys. |
| AET-104 | Implement reference catalog provider | Catalog | M1 | AET-013, AET-101 | Provider accepts valid submissions, rejects invalid records, enforces an admission policy, and publishes a signed head. |
| AET-105 | Publish catalog snapshots as torrents | Catalog | M1 | AET-007, AET-104 | Current head and history can be reconstructed through the torrent backend. |
| AET-106 | Implement feed subscriptions and synchronization | Catalog | M1 | AET-101, AET-104, AET-105 | Client syncs two providers, resumes updates, retains provenance, and tolerates one unavailable provider. |
| AET-107 | Build local full-text catalog search | Search | M1 | AET-014, AET-106 | Search works offline over title, description, tags, publisher, and content type. |
| AET-108 | Evolve Delphi into Aether Index UI | Product | M1 | AET-107, AET-118 | Users can search, inspect provenance, admission strength, and security metadata, and hand a result to a torrent backend. Install and launch are out of scope for M1. |
| AET-113 | Implement storage and seeding controls | Storage | M1 | AET-007, AET-014 | Users set disk/bandwidth limits, pin, stop seeding, remove data, and see deletion consequences. |
| AET-114 | Add DHT, tracker, PEX, and web-seed source policy | Torrent | M1 | AET-006, AET-007 | Downloads use configured compatible sources and report successful paths without making one canonical. |
| AET-116 | Write Index operator and user guides | Documentation | M1 | AET-122 | Guides cover provider setup, subscription, admission policy, backup, recovery, seeding, export, and troubleshooting. |
| AET-117 | Implement attestation admission adapter | Identity | M1 | AET-016 | Client and provider can issue, carry, verify, and weight attestation credentials from more than one independent attester. Default admission path. |
| AET-118 | Implement on-device ranking engine | Search | M1 | AET-107, AET-020, AET-021 | Personalized ordering is computed locally against the synced corpus. No behavioral signal leaves the device. Provider ranking, where used, is labelled with provenance and is overridable. |
| AET-119 | Implement full local data export | Privacy | M1 | AET-021, AET-101, AET-102 | Identity, keys, subscriptions, grants, records, and local state export in a documented format without provider cooperation. |
| AET-120 | Enforce opt-in collection defaults | Privacy | M1 | AET-021 | All telemetry, diagnostics, and provider-side logging default off; enabling requires an explicit revocable grant; test asserts a clean client emits no behavioral data. |
| AET-121 | Implement verified hand-off to external torrent backend | Torrent | M1 | AET-009, AET-114 | Client verifies listing, publisher, and torrent metadata consistency, then hands a magnet/metainfo to the configured backend with verification state surfaced to the user. |
| AET-122 | Build two-provider search-client end-to-end demo | Quality | M1 | AET-101–AET-121 | A clean client subscribes to two independent providers, syncs, verifies, searches offline, ranks locally, and hands off a verified download while one provider is offline. |

## P2 — Publishing, runtime, and beta reliability

### Milestone 2 — Publishing and catalog operation

| ID | Task | Area | Milestone | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- | --- |
| AET-103 | Implement app bundle publisher | Publishing | M2 | AET-007, AET-009 | A static directory becomes validated torrent metadata plus a signed release in one workflow. |
| AET-115 | Build two-publisher/two-provider publish demo | Quality | M2 | AET-103, AET-122 | Two publishers sign and submit releases; both providers verify and republish; clients converge with one provider offline. |
| AET-123 | Implement origin seeder of last resort | Availability | M2 | AET-105, AET-114 | A permitted origin serves as permanent web seed for the cold tail. Origin bytes are verified against the same torrent hashes as peer bytes, and the client reports which source path succeeded. |
| AET-124 | Implement sequencer and name resolution service | Identity | M2 | AET-005, AET-016 | An optional central service resolves names and settles contention. Names are exportable and non-canonical; losing the service loses names but never identity, records, or graph. |
| AET-125 | Implement optional onchain identity registry adapter | Identity | M2 | AET-016, AET-117 | One-time identity registration only; no per-action records onchain. Registry is one attester among several, and identity survives chain unavailability or abandonment. |
| AET-126 | Implement optional fee/burn admission adapter | Abuse | M2 | AET-016 | Available to catalogs that choose it, with documented regressiveness tradeoff and a required attestation-based alternative path for the same community. |
| AET-204 | Add feed admission and anti-spam controls | Abuse | M2 | AET-104, AET-016 | Providers support rate, size, age, admission-weight, allowlist, and review policies with clear rejection reasons. |
| AET-205 | Implement signed moderation labels | Moderation | M2 | AET-003, AET-106 | Users combine multiple label providers with local allow, warn, hide, and block rules. |
| AET-207 | Add publisher continuity and key-rotation UX | Identity | M2 | AET-003, AET-102 | Valid rotations preserve history; unexpected publisher changes produce a blocking warning. |

### Milestone 3 — Application runtime

Deferred until M1 ships. This is a browser security model built in-house; each item carries real exploit history and should not be rushed to reach a demo.

| ID | Task | Area | Milestone | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- | --- |
| AET-008 | Define application manifest v1 | Runtime | M3 | AET-001, AET-003 | Manifest specifies app ID, version, entrypoint, content ID, publisher, capabilities, and resource limits. |
| AET-010 | Implement per-application origin isolation | Runtime | M3 | AET-008, AET-009 | Apps cannot access another app, daemon cookies, gateway DOM, or privileged API origin. |
| AET-011 | Design capability API and permission model | Runtime | M3 | AET-008, AET-010 | Capabilities are deny-by-default, versioned, revocable, and classified by risk. |
| AET-109 | Implement verified install pipeline | Runtime | M3 | AET-010, AET-011, AET-103 | Install verifies torrent, release, publisher, manifest, limits, and permissions before activation. |
| AET-110 | Implement capability bridge MVP | Runtime | M3 | AET-011, AET-109 | Index can search and request torrent operations without direct daemon or key access. |
| AET-111 | Implement daemon-controlled permission prompts | Security | M3 | AET-011, AET-110 | Prompt shows app, publisher, capability, scope, and supports deny, once, or persistent grant. |
| AET-112 | Implement signed application updates | Runtime | M3 | AET-101, AET-103, AET-109 | Client detects successors, highlights permission changes, and supports update, ignore, and pin. |
| AET-206 | Add malware-scanning provider hooks | Security | M3 | AET-109, AET-205 | Install pipeline displays multi-provider verdict provenance/age and applies local policy. |
| AET-211 | Add runtime resource containment | Runtime | M3 | AET-109, AET-210 | Apps have enforced storage, request, CPU/time, and network policies with usable failures. |

### Milestone 4 — Beta reliability

| ID | Task | Area | Milestone | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- | --- |
| AET-201 | Embed a production torrent engine | Torrent | M4 | AET-122 | Desktop package joins Mainline DHT, downloads, streams, resumes, and seeds without a separate client. |
| AET-202 | Implement signed mutable catalog-head pointers | Catalog | M4 | AET-105, AET-201 | Heads resolve without a mandatory provider endpoint; rollback and expiration have interoperability tests. |
| AET-203 | Add catalog head/archive compaction | Catalog | M4 | AET-105, AET-106 | Recent sync stays bounded while complete history remains recoverable. |
| AET-208 | Implement availability-provider contracts | Availability | M4 | AET-114, AET-123 | Publishers attach replaceable seed/web-seed sources and clients verify all bytes against torrent hashes. |
| AET-209 | Add streaming and file-priority support | Media | M4 | AET-201 | Media starts before full download, seeks within a defined target, and falls back across sources. |
| AET-210 | Harden local API authentication | Security | M4 | AET-010 | Other processes and arbitrary websites cannot invoke privileged daemon APIs without authorization. |
| AET-212 | Add reproducible-release attestations | Supply chain | M4 | AET-103, AET-112 | UI distinguishes publisher signatures from source/build reproducibility evidence. |
| AET-213 | Build cross-platform signed installers and updates | Release | M4 | AET-201, AET-210 | Windows, macOS, and Linux packages update safely and preserve identity/database state. |
| AET-214 | Add privacy and exposure disclosures | Privacy | M4 | AET-201, AET-021 | Users understand peer-IP exposure, that installs are observable by strangers, seeding, storage, and provider requests before enabling them. |
| AET-215 | Run adversarial and load testing | Quality | M4 | AET-203–AET-214 | Tests cover malicious torrents, event floods, equivocation, sybil admission bypass, DHT failure, NAT diversity, disk pressure, and large catalogs. |
| AET-216 | Add opt-in diagnostics and operational dashboards | Operations | M4 | AET-215, AET-120 | Measure availability, sync latency, rejections, and failures without collecting content history, and only under an explicit grant. |
| AET-127 | Document abuse and legal posture | Governance | M4 | AET-205, AET-124 | Written position on infringing use, delisting versus deletion, moderation provider replaceability, reporting, and the legal interface at the permitted central service. Completed before any public launch. |

## P3 — Milestone 5: Scale and ecosystem

| ID | Task | Area | Milestone | Depends on | Acceptance criteria |
| --- | --- | --- | --- | --- | --- |
| AET-301 | Define provider discovery events | Ecosystem | M5 | AET-202, AET-205 | Users import, compare, verify, and remove provider recommendations without a canonical directory. |
| AET-302 | Implement federated search providers | Search | M5 | AET-203, AET-216 | Clients merge signed references, expose ranking provenance, and fall back to local search and on-device ranking. |
| AET-303 | Evaluate BEP 51 infohash indexing | Research | M5 | AET-201, AET-204 | Prototype measures coverage, cost, metadata quality, abuse, and operational risk before commitment. |
| AET-304 | Add privacy-preserving creator analytics | Publishing | M5 | AET-208, AET-216 | Creators receive useful aggregate availability/install data with no mandatory user tracking and no per-user records. |
| AET-305 | Add provider marketplace/payment abstraction | Ecosystem | M5 | AET-208, AET-301 | Users can buy availability/compute from competing providers without protocol consensus. |
| AET-306 | Add community pinning and retention policies | Availability | M5 | AET-208 | Communities publish signed policies; clients opt in with explicit disk/bandwidth budgets. |
| AET-307 | Build mobile consumption client | Mobile | M5 | AET-213, AET-214 | Mobile defaults to conservative Wi-Fi/charging seeding and preserves identity/catalog portability. |
| AET-308 | Add post, reply, follow, and reaction events | Social | M5 | AET-019, AET-020, AET-204, AET-205 | Social state syncs through selected feeds with signature, admission, spam, moderation, and offline tests. Counts remain provider assertions. |
| AET-309 | Add replaceable notification relays | Realtime | M5 | AET-308 | Relays are authenticated and rate-limited and cannot forge underlying events. |
| AET-310 | Add verifiable transcoding contribution | Compute | M5 | AET-209, AET-301 | The one accepted stranger-compute workload: segment-parallel, public input, latency-tolerant, with cheap output verification and provider/publisher authorization through signed records. General-purpose volunteer compute remains a non-goal. |
| AET-311 | Research private groups and messaging | Research | M5 | AET-210, AET-309 | Threat model covers encryption, membership, rotation, forward secrecy, delivery, and metadata leakage. |
| AET-312 | Research low-latency live streaming | Research | M5 | AET-209, AET-309 | Prototype compares WebRTC/QUIC/relay options and torrents completed-stream archives. |
| AET-313 | Publish protocol conformance suite | Ecosystem | M5 | AET-215 | Independent clients can validate events, manifests, feeds, admission proofs, permissions, and provider behavior. |
| AET-314 | Establish protocol governance/versioning | Governance | M5 | AET-313 | Extensions, compatibility, deprecation, ownership, and security response are documented. |
| AET-315 | Measure realized offload ratio | Quality | M5 | AET-208, AET-216 | Instrumented deployment reports the share of bytes served by peers versus origin, and the share of ranking performed on-device, so the offload claim is evidenced rather than asserted. |

## Definition of the first useful release

P0 and M1 are complete only when this works without hidden database edits:

```text
clean client subscribes to two independent providers
  -> synchronizes signed catalog heads and snapshots
  -> verifies provider and publisher events and admission proofs
  -> merges accepted metadata into local full-text search
  -> searches offline
  -> ranks results on-device with no behavioral data leaving the machine
  -> hands a verified torrent to the configured backend
  -> exports identity, subscriptions, and records in full
  -> flow still works with one provider unavailable
```

Publishing (M2) and the isolated application runtime (M3) are separately gated and must not be pulled forward to make a demo look complete.
