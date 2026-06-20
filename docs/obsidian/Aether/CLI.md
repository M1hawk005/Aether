# Command Line Interface (`src/cli/index.js`)

Aether provides a robust CLI built with `commander.js`. This allows developers to run headless nodes, orchestrate test swarms, and script interactions with the network without requiring the visual [Gateway](./Gateway.md).

[Back to Home](./Home.md)

## Core Commands

### `start-daemon`
Initializes and runs the core [Daemon](./Daemon.md) background process.
- **Usage:** `node src/cli/index.js start-daemon [options]`
- **Options:** 
  - `--port <number>`: Override the default REST API port.
  - `--dir <path>`: Specify a custom storage directory (useful for testing multiple local nodes).

### `init-user`
Generates a new Ed25519 identity keypair and stores it in the local `.aether` configuration directory.
- **Usage:** `node src/cli/index.js init-user <username>`

### `publish`
Publishes a payload to the network via the local daemon.
- **Usage:** `node src/cli/index.js publish <payload>`

### `fetch`
Attempts to fetch a specific chunk from the network using its cryptographic hash.
- **Usage:** `node src/cli/index.js fetch <capsule-id>`

### `resolve`
Looks up an alias in the [DHT](./DHT.md) to retrieve its associated records.
- **Usage:** `node src/cli/index.js resolve <alias>`

## Headless Swarm Testing

The CLI is heavily utilized in internal integration tests (e.g., `test-swarm.js`). By utilizing the `--dir` flag, developers can spawn multiple daemon instances on a single machine, each with their own isolated local storage, keypairs, and ports, to simulate a localized Kademlia swarm.
