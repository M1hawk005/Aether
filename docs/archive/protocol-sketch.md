# Protocol Sketch

This is an early protocol shape. It should change as the prototype teaches us what
is too complex, too vague, or too expensive.

## Core Concepts

### Identity

An identity is a unique username bound to one or more Ed25519 public keys.

One user can own multiple active keys, similar to a crypto wallet with multiple
devices or signing authorities. A user profile declares the primary key, active
keys, and revoked keys. If a key is compromised, the user publishes a signed
key-revocation event and rotates to another active key.

PGP can still be supported later as an external verification or identity-binding
layer, but the base protocol uses Ed25519.

### Username Registry

Global usernames should be tracked by a lightweight bundled registry chain. The
registry stores identity and namespace control records, not user content.

Registry records include username registration, profile hash updates, key additions
and revocations, forum/site namespace ownership, operator changes, and dispute or
contest records.

Every app node can run the registry compute locally as part of the application. The
goal is not financial mining; it is a replicated, tamper-evident database for names,
keys, and namespace authority.

If all user keys are lost, recovery follows the same model as crypto: there is no
magic central reset. Recovery requires pre-created recovery keys, social recovery
guardians, multisig recovery, or accepting that the old identity is unrecoverable and
starting a new username.

### Capsule

A capsule is the unit of publication.

It contains:

- Public header.
- Public policy envelope.
- Manifest, either public or encrypted.
- Chunks, either public bytes or encrypted bytes.
- Owner signature.
- Cache policy.
- Optional namespace targets.

### Chunk

A chunk is a byte range addressed by hash.

Public content is not encrypted by default. Private, unlisted, opaque-public, and
high-risk content use encrypted chunks. Large encrypted media or files should use a
unique data key per file. Routine access rotation should rewrap the small file key,
not re-encrypt the whole file.

### Manifest

The manifest describes how to reconstruct content from chunks.

Manifest rules:

- Public sites can use public manifests.
- Private capsules must use encrypted manifests.
- Semi-public or unlisted content can use encrypted manifests with explicit access grants.
- Privacy-conscious public publishers can use opaque manifests where paths are hidden
  behind content IDs.
- A manifest can describe a single file or a full directory tree (using a `files` array), allowing a single capsule to package entire static websites or applications.

Public content should stay available through voluntary caching unless the uploader
deletes it. Viewers may download public content to their own device, but they can
also choose to stop serving their local copy back to the network.

### Revocation

A revocation is a signed lifecycle event from the owner.

Compliant peers must:

- Stop advertising the capsule.
- Stop serving chunks.
- Delete cached chunks.
- Forward the revocation event where appropriate.
- Keep only a tombstone if required for indexes or audit trails.

For public content, deletion means honest nodes stop serving the content, indexes
stop listing it, origin nodes remove it, clients treat it as unavailable, and caches
expire or purge it according to policy. For encrypted content, deletion also revokes
or destroys access keys and stops future access grants.

### Namespaces

Site:
A publisher-controlled namespace for web-like content. Usually has pages, assets,
manifests, routes, and versions.

Forum:
A shared writable namespace for many users. Usually has posts, replies, moderation
records, topic indexes, and permissions.

Users decide which forums they submit their capsules to. Forum operators decide
which submitted capsules are listed, ranked, hidden, labeled, or rejected in that
forum. The network proves authorship and availability; the forum enforces its own
listing rules.

App namespace:
A structured data namespace used by an application. Usually has schemas, records,
state updates, API-like operations, and permissions.

Search index:
A derived namespace that points to content from other namespaces. Usually stores
references, rankings, snippets, thumbnails, trust signals, and moderation metadata.

Indexes can store titles, snippets, thumbnails, safety labels, rankings, and trust
signals, but those are index-authored metadata unless signed by the publisher.

### Origins

A publisher can run an origin node on a VPS, home server, laptop, or overlay network.
The origin is the publisher-controlled source of truth. Peers cache and serve content,
but clients can fall back to the origin when peer capacity is missing, stale, or
overloaded.

Source preference:

1. Local cache.
2. Nearby peers.
3. Wider peers.
4. Trusted relays or bootstrap caches.
5. Publisher origin.

Static assets are content-addressed. Dynamic APIs are identity-addressed service
capabilities with schemas, auth rules, operations, and cache policies. Signed dynamic
responses can become content-addressed snapshots.

### Transport

The first real network should be QUIC-native for peers, HTTPS-compatible for origins,
relay-assisted for NAT, and privacy-overlay-optional rather than privacy-overlay-required.

Recommended first shape:

- Primary peer transport: QUIC over UDP. *(Note: The current Phase 1 prototype uses lightweight local HTTP daemons to prove peer-to-peer chunk and metadata proxying.)*
- Peer discovery and routing: DHT plus bootstrap or tracker nodes.
- NAT fallback: Aether relay nodes.
- Origin fallback: HTTPS and/or QUIC endpoints.
- Browser-compatible bridge: HTTPS relay first, WebRTC later.
- Privacy overlay: optional Tor/I2P-compatible origin addresses, not default.

Delay full VPN or mesh networking until the content model is proven.

### Metadata Visibility

Use visibility levels instead of one global metadata rule.

Level 0, opaque private:
Expose object hash, size band, expiry/cache class, and proof parameters. Hide
publisher, capsule ID, paths, media type, and access policy detail.

Level 1, private but cacheable:
Expose object hash, exact or rounded size, encrypted capsule ID, cache TTL,
serve/storage permissions, and broad content class. Hide plaintext names, paths,
titles, membership, and manifest contents.

Level 2, public capsule:
Expose publisher ID, namespace/capsule ID, object hashes, exact size, media type,
cache policy, optional public manifest path, and content version.

Level 3, indexable public:
Expose title, description, thumbnail, language, tags, canonical name, publisher-signed
preview metadata, and index permissions.

The rule is to expose only the metadata needed for routing, caching, safety, and
discovery at that visibility level.

### Access Grants

Access grants are signed, scoped capability tokens.

An access grant states who issued it, who can use it, what it grants access to, which
key material it wraps or references, when it expires, whether it can be delegated,
and which actions are allowed.

Identity keys sign and authorize. Access keys unlock capsules, groups, threads, or
epochs. Data keys encrypt one object or chunk set.

Revocation uses access epochs. When membership changes, increment the epoch. Future
objects use the new epoch key. Large objects are not re-encrypted during routine
rotation; their small data keys are rewrapped.

### Trust Claims

Moderation, blocklists, rankings, and reputation are signed claims, not universal
protocol truth.

Claims can target:

- Object hash.
- Capsule ID.
- Namespace ID.
- Publisher identity.
- Origin endpoint.
- Index ID.

Claims should include scope, such as personal, community, index, legal jurisdiction,
or network default. Users decide which claims to follow.

### Seeding Reputation (Subjective Web-of-Trust)

To prevent Sybil attacks and gamification, Aether abandons global network-wide scores. Instead, reputation relies on a **Subjective Web-of-Trust** combined with **Proof of Useful Work**.

1. **Reputation is Local and Earned, Not Global**
   Your Aether node doesn't ask the network "What is Bob's score?" Instead, your node maintains its own private scoreboard. Bob only gets a point on your scoreboard if your node personally asks Bob for a file, and Bob successfully delivers the correct, cryptographically verified bytes. If attackers spin up 1,000 fake nodes and trade fake data among themselves, they get zero points on your scoreboard because they never helped you.

2. **Expanding to "Network-Wide" via Trusted Peers**
   If you want to find good seeds you haven't interacted with yet, your node asks the peers you already trust. By asking "Who are your best seeds?" and following this chain of trust, your node builds a "network-wide" view anchored entirely in real, verified work.

3. **De-gamifying the UI**
   Showing raw numbers (like "Score: 10,432") encourages farming. Aether should never show raw numbers. Instead, display qualitative, relative badges:
   - 🟢 **"Highly Reliable Seed"** (Top 10% of nodes in your trust graph)
   - 🔵 **"Trusted by your peers"** (You haven't interacted with them, but your trusted friends have)
   - ⚪ **"Unknown/New"**

By making reputation a locally calculated metric based on actual bytes delivered, and restricting second-hand reputation to nodes you already trust, it becomes mathematically impossible to farm reputation.

### Verification

Prototype verification can use random chunk retrieval plus Merkle proofs.

First real network verification should use erasure-coded proof of retrievability.
Avoid starting with heavy zero-knowledge proofs unless financial incentives require
them later.

## Example Public Header

```json
{
  "version": 1,
  "capsule_id": "sha256:...",
  "publisher": "alice",
  "owner_key_id": "sha256:...",
  "visibility": "public",
  "created_at": "2026-05-26T00:00:00Z",
  "manifest_hash": "sha256:...",
  "cache_policy_hash": "sha256:...",
  "signature": "sig:..."
}
```

## Example Cache Policy

```json
{
  "max_ttl_seconds": 604800,
  "allow_peer_cache": true,
  "allow_persistent_pin": false,
  "delete_on_revoke": true,
  "serve_while_offline": true,
  "required_storage_class": "temporary",
  "mobile_default": "download-only"
}
```

Peers combine publisher policy with local policy. A simple first formula is:

```text
cache_duration = base_duration * usefulness_score * affordability_score
```

Signals:

- Storage: more free disk keeps content longer; low disk expires aggressively.
- Proximity: same region, low latency, same community, or same overlay neighborhood.
- Network quality: stable, fast, unmetered peers keep longer.
- Popularity: frequently requested capsules stay longer.
- Device policy: mobile defaults to download-only unless charging, on Wi-Fi, and opted in.

Publisher policy says "please cache this for up to 7 days." Peer policy says "given
my situation, I will cache it for 2 days."

## Example Capacity Envelope

```json
{
  "max_concurrent_streams": 8,
  "max_daily_egress_mb": 500,
  "max_capsule_egress_share": 0.25,
  "metered": false,
  "serve_when_on_battery": false
}
```

Popular capsules should spread across more peers instead of concentrating traffic on
the first available sources.

## Example Revocation Event

```json
{
  "version": 1,
  "type": "capsule.revoked",
  "capsule_id": "sha256:...",
  "publisher": "alice",
  "owner_key_id": "sha256:...",
  "created_at": "2026-05-26T00:00:00Z",
  "reason": "owner_deleted",
  "signature": "sig:..."
}
```

## Example Policy Envelope

```json
{
  "version": 1,
  "object_hash": "sha256:...",
  "publisher": "alice",
  "capsule_id": "sha256:...",
  "size": 1048576,
  "media_type_hint": "page",
  "visibility": "public",
  "cache_policy_hash": "sha256:...",
  "signature": "sig:..."
}
```

## Example Access Grant

```json
{
  "type": "aether.access_grant.v1",
  "issuer": "alice",
  "subject": "bob",
  "scope": {
    "namespace": "capsule://alice/private",
    "objects": ["sha256:..."],
    "version_range": {
      "from": "v42",
      "to": "v48"
    }
  },
  "permissions": ["read", "cache", "serve"],
  "key_wrap": {
    "algorithm": "X25519+HKDF+AES-KW",
    "wrapped_key": "base64..."
  },
  "not_before": "2026-05-27T00:00:00Z",
  "expires": "2026-06-27T00:00:00Z",
  "delegable": false,
  "revocation_epoch": 12,
  "signature": "sig:..."
}
```

## Local Prototype Storage Layout

```text
.aether/
  network.json
  users/
    alice/
      profile.json
      keys/
  indexes/
  forums/
  registry/
  usernames/
  peers/
    peer-a/
      chunks/
      capsules/
      revocations/
    peer-b/
      chunks/
      capsules/
      revocations/
```

## Important Constraint

Aether can enforce deletion only among compliant peers. It cannot prevent a viewer
from saving plaintext after viewing, taking screenshots, modifying their client, or
re-publishing content elsewhere.

The honest promise is:

> Official peers obey signed revocations, stop serving deleted content, remove
> cache entries where possible, and stop granting readable access for encrypted
> content when owners revoke it.
