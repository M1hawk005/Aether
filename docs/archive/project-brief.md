# Project Brief

## Problem

Today, most web content is delivered through centralized platforms, CDNs, cloud
storage providers, and large data aggregators. This makes delivery convenient, but
it concentrates control, metadata, availability, and failure modes in a small number
of places.

Fully permanent decentralized storage systems solve one part of the problem, but
often create another: once content is published, the owner may lose practical
control over deletion, expiry, or redistribution.

## Proposal

Aether is a decentralized content delivery layer based on signed content capsules.
Public capsules can be clear content-addressed data by default. Private, unlisted,
and privacy-conscious capsules can encrypt manifests and chunks.

When a publisher creates content, their client:

1. Builds a manifest describing the content and cache policy.
2. Chunks the content.
3. Encrypts chunks when the capsule visibility requires it.
4. Signs the capsule header and manifest reference.
5. Announces the capsule to one or more indexes, sites, forums, or app namespaces.

When a viewer opens the content:

1. Their client resolves the capsule through an index or namespace.
2. It discovers peers that have matching chunks.
3. It downloads chunks from nearby, fast, or trusted peers.
4. It requests decryption access according to the capsule policy.
5. It renders the content.
6. It may cache chunks and serve them to other viewers.

## Design Principles

- Private payloads are encrypted before entering the network.
- Public payloads can remain public and content-addressed.
- Content integrity is verified by hash.
- Ownership is verified by public-key signatures.
- Deletion is represented by signed revocation events.
- Peer caches are policy-bound and temporary by default.
- Indexes reference content; they do not own it.
- The protocol should support social posts, web pages, static sites, app bundles,
  documents, media, and API snapshots.

## Non-Goals For The First Prototype

- Global peer-to-peer networking.
- Financial incentives.
- Proof-of-storage.
- Browser extension integration.
- Production-grade anonymity.
- Perfect deletion guarantees against malicious clients.
- Full PGP implementation.

The first version should prove the capsule lifecycle locally before expanding.
*(Note: We have successfully proven multi-file directory publishing and multi-peer local networking, enabling fully decentralized static site delivery.)*

## First Product Surface

The first user-facing application should be a minimal Reddit-like forum built on top
of the network.

Users publish capsules they own, then choose which forums to submit those capsules
to. Forum operators enforce local rules by approving, rejecting, labeling, ranking,
or hiding submissions. The forum does not own the post; it owns the listing and
moderation decisions around that post.

Users should be able to connect to the Aether network with a clear app-level control,
similar in spirit to a "connect to Tor" button: disconnected mode for local/offline
use, connected mode for peer discovery, fetching, seeding, and forum participation.

Seeding should be voluntary. The early reward model should be community reputation,
badges, personal contribution stats, and scoped forum trust rather than money.
