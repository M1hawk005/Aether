---
sidebar_position: 3
title: API Reference
---

# API Reference

> **Migration notice:** These endpoints describe the legacy capsule prototype. The target v1 API will be versioned around identities, events, feeds, catalogs, torrents, applications, permissions, storage, and providers.

The `AetherSDK` class provides simple, promise-based methods to interact with the Aether Daemon.

## `publish(publisher, payload, options)`
Publishes data to the network.

*   **`publisher`** (string): The username of the local identity publishing the data.
*   **`payload`** (string): The raw string data to publish.
*   **`options.visibility`** (string): `'public'` or `'private'`.
*   **`options.persistent`** (boolean): Whether to retain the file locally forever.

**Returns:** A `Promise<string>` containing the unique `Capsule ID` (SHA-256 hash).

---

## `fetch(capsuleId)`
Retrieves a capsule's payload from the local cache or the P2P network.

*   **`capsuleId`** (string): The SHA-256 hash of the capsule.

**Returns:** A `Promise<string>` containing the uncompressed payload.

---

## `availability(capsuleId)`
Checks the global DHT for how many seeders are currently hosting this capsule.

*   **`capsuleId`** (string): The SHA-256 hash of the capsule.

**Returns:** A `Promise<number>` containing the active seed count.

---

## `deleteLocal(capsuleId)`
Deletes the capsule from your device's hard drive. It remains on the network if seeded by others.

*   **`capsuleId`** (string): The SHA-256 hash of the capsule.

**Returns:** A `Promise<boolean>` indicating success.
