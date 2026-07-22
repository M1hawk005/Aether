# Core Daemon (`src/daemon-go`)

[Back to Home](./Home.md)

The daemon is Aether's trusted local boundary. It exposes a versioned API to the Wails gateway, SDK, CLI, and isolated Aether applications. It is not merely a neutral router: it must verify content, enforce local security and storage policy, and mediate privileged actions.

## Target services

1. **Identity:** securely load and use Ed25519 keys; applications never receive raw private keys.
2. **Events:** canonicalize, sign, verify, order, deduplicate, and store Aether events.
3. **Catalog:** synchronize selected feeds and maintain local full-text search.
4. **Torrent:** create, add, stream, resume, and seed standard torrents through a backend interface.
5. **Policy:** apply labels, local blocks, quotas, retention, bandwidth, and provider choices.
6. **Runtime:** serve verified app bundles from isolated origins and enforce capabilities.
7. **Providers:** communicate with catalog, label, search, relay, availability, and compute services.

## Startup sequence

```text
open secure identity and database
  -> load policy and capability grants
  -> initialize torrent backend and Mainline DHT
  -> resume retained torrents
  -> start authenticated local API
  -> synchronize selected feeds
  -> expose gateway and isolated app origins
```

## API direction

The legacy capsule endpoints remain prototype interfaces during migration. The v1 target API should be resource-oriented and versioned around:

- `/api/v1/identity` and mediated signing;
- `/api/v1/events`;
- `/api/v1/feeds` and `/api/v1/catalog/search`;
- `/api/v1/torrents`;
- `/api/v1/apps`, releases, and updates;
- `/api/v1/permissions`;
- `/api/v1/storage` and seeding policy;
- `/api/v1/providers`.

Localhost is not an authentication mechanism. The production API must prevent arbitrary websites and unrelated local processes from invoking privileged actions.

## Migration note

Current code still contains libp2p/GossipSub message routing, custom capsule chunks, first-seen aliases, an in-memory tracker table, and unverified delivery comments. New features should be implemented behind target interfaces so those paths can be removed without another rewrite.
