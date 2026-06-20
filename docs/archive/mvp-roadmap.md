# MVP Roadmap

## Phase 0: Local Capsule Prototype

Build a command-line prototype that runs on one machine.

Capabilities:

- Create a unique username with an Ed25519 keyset.
- Add and revoke user keys.
- Append username and key events to a local registry chain.
- Package a file or directory as a capsule.
- Store public content without encryption by default.
- Encrypt private, unlisted, and opaque-public chunks.
- Hash chunks.
- Sign the capsule manifest.
- Store chunks in a local cache directory.
- Resolve capsules from a local JSON index.
- Reconstruct public and private content from cached chunks.
- Publish and apply signed revocation events.
- Create a local forum namespace with an operator.
- Submit capsules to a forum.
- Approve forum submissions into the forum index.
- Open a local gateway shell with address bar, forum view, capsule view, and trust details.

Success criteria:

- A public local HTML file can be packaged, cached, fetched, and restored.
- A private local HTML file can be packaged, cached, fetched, decrypted, and restored.
- A revoked user key can be replaced by another active key.
- A replacement key can revoke older content signed by a compromised key.
- A revocation event causes the cache to remove chunks.
- Publishing to a forum does not automatically list the capsule.
- A forum operator can approve a submitted capsule for listing.
- Gateway can browse local accepted forum posts and render a verified capsule.
- Tampering with a chunk or manifest causes verification failure.

## Phase 1: Multi-Peer Local Network

Run multiple local peer processes.

Capabilities:

- Peer announces which capsule chunks it has.
- Viewer fetches chunks from more than one peer.
- Peer selection considers availability, latency, and cache policy.
- Cache TTL is enforced.
- Peers advertise a capacity envelope.
- Local seeding stats are recorded.

Success criteria:

- A capsule can be reconstructed from chunks spread across multiple local peers.
- Expired content is removed automatically.
- Revoked content stops being served.

## Phase 2: Namespaces And Indexes

Build the first minimal Reddit-like forum on top of the local network model.

Capabilities:

- A capsule declares intended namespaces.
- A namespace index lists accepted capsules.
- Indexes can reject, delist, or tombstone capsules without deleting owner data.
- Clients can subscribe to namespaces.
- Forum operators approve, reject, label, and rank submitted capsules.
- Users can choose which forums to submit their posts to.
- Community reputation and seeding badges are visible but non-financial.

Success criteria:

- A local "site" can publish and resolve a static page.
- A local "forum" can list several signed capsules.
- A forum page can show accepted posts, pending operator queue, and basic scores.
- Owner revocation removes readable content while preserving a tombstone.

## Phase 3: Web Gateway

Expose capsule content through a local HTTP gateway.

Capabilities:

- `http://localhost` routes resolve capsule IDs or site namespaces.
- HTML, CSS, JS, and media are served after verification and, when needed, decryption.
- Static web apps can be packaged and rendered.
- Browser shell supports Browse, Save, and Publish modes.
- Capsule pages show verified publisher, availability, offline status, and trust details.
- Failure states explain peer/origin/access fallback in plain language.

Success criteria:

- A small static site can be opened in a browser through the Aether gateway.
- Cached chunks can serve repeat requests without the origin.
- The UI feels like a reader first and exposes protocol details only on demand.

## Phase 4: Real Networking

Introduce network transport after the local model is stable.

Candidates:

- QUIC over UDP for native peers.
- DHT plus bootstrap or tracker nodes for discovery.
- Relay-assisted NAT fallback.
- HTTPS and/or QUIC publisher origins.
- Optional Tor/I2P origin addresses.

Success criteria:

- Two machines can exchange chunks according to visibility and cache policy.
- Peer discovery works without a central content host.
- Clients fall back from peer fetch to relay and then origin.
