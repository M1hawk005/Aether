---
sidebar_position: 2
title: Core Concepts
---

# Theoretical Concepts

To build safely on Aether, it is important to understand the underlying mathematics and cryptography.

## Distributed Hash Table (DHT)
Aether uses a **Kademlia DHT** (via Libp2p) for peer discovery.
Instead of routing by IP addresses, nodes are routed by the hash of their Public Key (Node ID). 

When you publish a capsule, Aether finds the nodes mathematically closest to you (using XOR distance) and connects to them to seed the data.

## Cryptography & Data Integrity
Every piece of data published to Aether is wrapped in a **Capsule**.
1. **Payload:** The raw data, compressed losslessly using **Brotli**.
2. **Manifest:** A JSON index of the payload chunks, hashed using **SHA-256**.
3. **Envelope:** A signed object proving ownership, signed using an **Ed25519** Private Key.

Because every layer is cryptographically signed, the data cannot be tampered with by intermediary nodes.

## Decentralized Identifiers (DIDs)
Usernames on Aether are mapped to Ed25519 Public Keys via **Epochs**. 
This prevents identity theft: even if a username is transferred to a new user (a new Epoch), historical posts retain the cryptographic signature of the original owner.
