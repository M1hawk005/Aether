---
sidebar_position: 2
title: Core Concepts
---

# Core Concepts

Aether combines independently verifiable records with existing content-distribution infrastructure.

## Identity and signed events

An Ed25519 public key is the canonical identity. Profiles, releases, listings, labels, and later social activity form per-author signed event chains with sequence numbers and predecessor links. Handles are optional aliases, not global consensus identifiers.

## Immutable content

Static applications, media, catalog snapshots, and archives are standard BitTorrent v2 or hybrid torrents. Torrent hashes verify downloaded bytes; the Aether event signature separately verifies which publisher endorsed those bytes.

## Mainline BitTorrent DHT

The Mainline DHT locates peers for known torrent infohashes. It is not the repository's legacy libp2p DHT and is not a keyword-search or social database. A compatible torrent engine joins the existing BitTorrent network.

## Catalogs and local search

Independent catalog providers distribute signed listing events and torrent-backed snapshots. Users subscribe to several providers; the local daemon verifies, deduplicates, filters, and indexes accepted records for offline search.

## Hybrid providers

Replaceable services can improve search, ranking, moderation, availability, relays, or computation. They may be centralized operationally without becoming canonical owners of identity or content.

## Application isolation

Downloaded web applications are untrusted. Each verified release runs from an isolated origin and receives only explicit daemon capabilities. Private keys remain in the daemon, and sensitive signing or filesystem actions require daemon-controlled consent.
