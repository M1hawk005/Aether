# Distributed Hash Tables

[Back to Home](./Home.md)

The target architecture primarily uses the **BitTorrent Mainline DHT** through a compatible torrent backend. The repository's existing go-libp2p Kademlia DHT is a different network with different messages and peers; participating in it does not join Mainline DHT.

## Mainline DHT responsibilities

Mainline DHT stores short-lived peer locations associated with known torrent infohashes. Torrent clients use `get_peers` to find peers and `announce_peer` to advertise participation.

It is appropriate for:

- locating peers for a known torrent;
- trackerless torrent operation;
- small interoperable extensions implemented by the selected backend.

It is not appropriate for:

- keyword search;
- complete torrent metadata catalogs;
- social timelines or comments;
- global usernames;
- recommendation indexes;
- durable large records.

## Mutable pointers

BEP 44-style signed mutable items and BEP 46-style update pointers may identify the current torrent for a catalog or publisher feed. Values remain small; catalog data itself is stored in torrents. Publishers or mirrors must republish pointers because DHT values are not permanent storage.

## Bootstrap

The embedded torrent engine retains its routing table and uses multiple ordinary BitTorrent bootstrap mechanisms. Bootstrap nodes introduce a client to peers but do not become canonical providers.

## Legacy libp2p use

The current global libp2p DHT/GossipSub implementation is migration-era code. libp2p may remain later for narrowly scoped control, rendezvous, private, or real-time functions only when BitTorrent and provider APIs do not meet the requirement.
