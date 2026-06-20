# Aether SDK (`packages/sdk`)

The Aether Node.js SDK provides a strongly typed TypeScript interface for third-party applications and the official [Gateway](./Gateway.md) to interact with a local Aether [Daemon](./Daemon.md).

[Back to Home](./Home.md)

## Installation and Initialization

The SDK is published as `aether-sdk`. It encapsulates the `Axios` HTTP client to communicate with the local daemon's REST API.

```typescript
import { AetherClient } from 'aether-sdk';

// Initializes connection to local daemon on port 5000
const client = new AetherClient('http://localhost:5000'); 
```

## Resiliency and Retries

The `AetherClient` is designed to be highly resilient against transient local network failures. It utilizes `axios-retry` to implement **exponential backoff**.
- **Retry Conditions:** It automatically retries requests on `5xx` server errors or network drops (`ECONNREFUSED`).
- **P2P Exceptions:** It explicitly *does not* retry on `404` or `408` errors, as these represent expected Distributed Hash Table ([DHT](./DHT.md)) misses or P2P fetch timeouts.

## Core API Methods

### `publish(publisher: string, payload: string, options?: PublishOptions)`
Submits a data payload to the local daemon, which handles chunking and propagation to the swarm.

### `fetch(capsuleId: string)`
Requests the daemon to download a specific chunk by its hash. If not found locally, the daemon will query the network.

### `resolve(targetAlias: string)`
Looks up a human-readable alias in the DHT to resolve its associated cryptographic public key and latest published capsules.

### Identity Methods
- `claimUsername(username: string, keyId: string)`: Broadcasts a cryptographically signed claim over the DHT.
- `releaseUsername(username: string, keyId: string)`: Revokes a previously claimed username.
- `getReputation()`: Returns the local trust metrics for the node.
