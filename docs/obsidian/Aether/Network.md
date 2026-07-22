# Network Architecture

[Back to Home](./Home.md)

Aether composes existing networks rather than using one custom transport for every message and byte.

## Content plane

Immutable application bundles, media, catalog snapshots, and archives move through standard BitTorrent swarms. The torrent backend handles:

- Mainline DHT peer discovery;
- peer exchange and optional HTTP/UDP trackers;
- piece selection, verification, resume state, and streaming priority;
- NAT traversal and bandwidth controls;
- optional HTTP web seeds and mirrors.

Aether records only reference the torrent and its signed publication metadata. The legacy `BT_ANNOUNCE` message and in-memory tracker are not part of the target protocol.

## Record plane

Profiles, releases, listings, labels, recommendations, and later social activity are small signed events. They synchronize through selected catalog/feed providers and may be archived into torrents.

There is no mandatory global topic. Clients fetch only subscribed publisher, application, community, and moderation feeds.

## Service plane

Replaceable providers may supply search, ranking, recommendations, moderation, malware scanning, guaranteed availability, relays, or compute. HTTPS is acceptable here because provider responses are either subjective by definition or contain verifiable signed records and torrent references.

## Discovery flow

```text
known feed address
  -> resolve current signed head
  -> download catalog snapshot
  -> verify and index listings locally
  -> select known infohash
  -> use Mainline DHT/tracker/PEX to find peers
  -> download and verify pieces
```

BitTorrent discovers peers for a known infohash; it does not provide useful keyword search. Catalogs and local indexes solve that separate problem.

## Optional transports

libp2p, WebRTC, QUIC, or relays may be added later for scoped control messages, private interactions, notifications, or live media when an existing mechanism does not solve the requirement. They are not blob-transfer dependencies for Aether Index.
