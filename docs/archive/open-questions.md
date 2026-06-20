# Open Questions And Current Decisions

This file now separates decisions we have made from areas that still need design or
experimentation.

## Decided

### Identity

- The base protocol uses Ed25519.
- Every user has a unique username.
- A username can own multiple active keys.
- Key compromise follows a crypto-wallet-style model: revoke the compromised key,
  rotate to another active key, and publish a signed key-revocation event.
- Username registration should use a lightweight bundled registry chain that tracks
  names, profile hashes, key events, and namespace ownership without storing content.
- If all keys are lost, recovery is the same as crypto: recovery keys, guardians,
  multisig recovery, or the identity is unrecoverable.
- PGP can be added later as an identity verification bridge, not the base signing
  primitive.

### Public, Private, And Opaque Content

- Public content is public by default and does not need encryption.
- Public sites can use public manifests.
- Private capsules must use encrypted manifests.
- Semi-public or unlisted content can use encrypted manifests with explicit grants.
- Privacy-conscious public publishers can use opaque manifests where paths are hidden
  behind content IDs.
- Public content should remain available through voluntary caching until the uploader
  deletes it.
- Users can download public content to their device and decide whether their copy is
  shared back to the network.

### Keying

- Private capsules and groups use one data key per object.
- Private group object keys are wrapped by the current group access key.
- Group access keys rotate when membership changes.
- Normal private groups can rotate access keys on a 90-day cadence.
- High-risk content should use narrower access keys per thread, collection, or event.
- Large media and files always use a unique data key per file.
- Routine access rotation rewraps file keys instead of re-encrypting large files.

### Deletion

- Honest nodes stop serving deleted content.
- Indexes stop listing it.
- Origin nodes remove it.
- Access grants are revoked.
- Future clients treat it as unavailable.
- Caches expire or purge it according to policy.
- Encrypted payloads should be deleted where possible, and keys should be revoked or
  destroyed.

### Cache Policy

- Caching is voluntary, like a torrent network.
- Individual viewers have total control over their local cache and serving policy.
- Cache duration should be adjusted by storage, proximity, network quality, popularity,
  and local device policy.
- Mobile should default to download-only unless charging, on Wi-Fi, and explicitly
  opted in.
- Publishers express intent; peers make final local storage decisions.

### Overload Control

- Clients should fetch from many peers, not one preferred peer.
- Peers advertise serving limits.
- Routing should be load-aware.
- Popular capsules should replicate wider automatically.
- Chunking lets different peers serve different pieces.
- Clients fall back to origin or relays when peer capacity is low.

### Namespaces And Indexes

- A site is a publisher-controlled namespace for web-like content.
- A forum is a shared writable namespace for conversation.
- Users choose which forums they submit their posts to.
- Forum operators enforce the forum's rules by approving, rejecting, labeling,
  ranking, or hiding submissions.
- An app namespace stores structured application state and operations.
- A search index is a derived namespace that points to content it usually does not own.
- Indexes can store titles, snippets, thumbnails, rankings, labels, and trust signals.
- Index metadata is index-authored unless signed by the publisher.
- Moderation records are signed claims with scope, not edits to the original content.

### Origins

- Users should be able to deploy a VPS, home server, or overlay-connected node as an
  always-on serving origin.
- Origins can live behind overlays so users do not need public IP or port-forwarding
  expertise.
- Clients try local cache, nearby peers, wider peers, relays, then publisher origin.
- Static assets are content-addressed.
- Dynamic APIs are identity-addressed capabilities with schemas, auth rules, allowed
  operations, and cache policies.

### Transport

- First real peer transport should be QUIC over UDP.
- Peer discovery should use a DHT plus bootstrap or tracker nodes.
- NAT fallback should use relays.
- Origin fallback should support HTTPS and/or QUIC.
- Browser compatibility can use HTTPS relays first and WebRTC later.
- Tor/I2P should be optional privacy transports, not defaults.
- Full VPN/mesh networking should wait until the content model is proven.

### Metadata Visibility

- Use Level 0 opaque private, Level 1 private-but-cacheable, Level 2 public capsule,
  and Level 3 indexable public metadata.
- Expose only the metadata needed for routing, caching, safety, and discovery at each
  visibility level.

### Access Grants

- Access grants are signed scoped capability tokens.
- Identity keys sign and authorize.
- Access keys unlock capsules, groups, threads, or epochs.
- Data keys encrypt one object or chunk set.
- Revocation uses access epochs.

### Useful Storage

- Useful storage means bytes are stored, retrievable, and useful to demand.
- Peers can prove storage by responding to random encrypted chunk challenges.
- Prototype verification should use random chunk challenge plus Merkle proof.
- First real network verification should use erasure-coded proof of retrievability.
- Usefulness can include retrieval success, availability, latency, rarity, policy fit,
  and integrity.

### Seeding Rewards

- Rewards should start as non-financial reputation, badges, personal stats, publisher
  thanks, and community-scoped karma.
- Seed score should reward verified retrievals, rarity, availability, policy
  compliance, and anti-spam weighting.
- Rewards need daily caps, diminishing returns, self-fetch loop protection, and
  suspicious traffic filtering.

### Abuse And Safety

- Abuse and safety policy is chosen by users, communities, and indexes.
- Peers inspect a signed public policy envelope without decrypting private payloads.
- Blocklists and trust signals are signed scoped claims, not global protocol commands.

### Content Packaging (Phase 0)

- A capsule can represent a single file or a full directory tree (e.g., a static site).
- Directory manifests contain a `files` array, each with its own path, content type, and chunk list.
- The gateway natively supports directory capsules by resolving relative paths and rendering `index.html` within a sandboxed iframe.

### Multi-Peer Architecture (Phase 1 Prototype)

- The local prototype uses lightweight HTTP daemons to represent independent peers.
- Service discovery for the prototype relies on a local `.aether/network.json` registry.
- Daemons proxy chunk and metadata requests across the network when files are missing locally.
- Infinite request loops between peers are prevented using `x-no-forward` HTTP headers.

### Registry Consensus

- Start with a Proof-of-Authority (PoA) Consortium.
- Initial trusted validators vote on-chain to add new ones. 
- The data structure remains identical, so it can seamlessly swap to a Proof-of-Stake validator set later without breaking names.

### Seeding Reputation

- Abandon global network-wide scores in favor of a Subjective Web-of-Trust and Proof of Useful Work.
- Nodes score peers locally based on verified bytes delivered. Network discovery expands through the trust graph (asking trusted peers for their trusted seeds).
- The UI uses qualitative relative badges (🟢 Highly Reliable, 🔵 Trusted by peers, ⚪ Unknown) rather than raw numbers to prevent gamification.

### Identity Deletion, Reissue, and History

- Enforce Time-Aware Epochs mapped to underlying DIDs.
- The registry preserves a history of leases. Clients resolve post signatures against these historical epochs.
- The UI visually distinguishes historical owners from current owners (e.g., using a "Historical" tag and a broken link icon) and routes them to an archived view. Hovering over names exposes the raw DID to prevent power-user ambiguity.

## Still Open

- How should the local gateway evolve into a packaged browser shell?
- What exact proof-of-stake validator and slashing model fits a name registry?
- What should the minimal Reddit-like forum extension API look like?
