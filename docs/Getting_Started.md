# Aether Protocol: Getting Started & Application Guide

> [!WARNING]
> This guide documents the runnable legacy capsule/libp2p prototype. It is retained for development and migration testing; it does not describe the accepted hybrid BitTorrent architecture. See [Aether Architectural Reference](./Aether_Architecture.md) and [Implementation Backlog](./Implementation_Backlog.md) before building new features.

Welcome to the Aether Protocol. This guide provides meticulous, step-by-step instructions for setting up the development environment, running the decentralized network locally, and utilizing the Aether Gateway UI.

---

## 1. Prerequisites

Before interacting with Aether, ensure your system has the following dependencies installed and configured in your system `PATH`:

- **Go (1.20 or higher):** Required to compile and run the LibP2P daemon and Wails backend.
- **Node.js (18 or higher):** Required for building the TypeScript SDK and running the Gateway frontend.
- **Wails CLI:** The framework used for the desktop Gateway UI.
  - Install via: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

---

## 2. Project Architecture Overview

The Aether repository is modularized into three distinct components:

1. **`src/daemon-go/`**: The core P2P LibP2P node. It runs as a background process (default `http://localhost:5000`) managing Kademlia DHT routing, cryptographic verification, and local storage.
2. **`packages/sdk/`**: A rigorously typed Node.js/TypeScript SDK (`aether-sdk`). It provides a robust, retry-enabled wrapper around the Daemon's API.
3. **`apps/delphi/`**: The Delphi decentralized Tracker Web App accessible at `http://localhost:5000/ui/`.
4. **`aether-gateway/`**: The frontend Desktop Application built using Wails (Go + React/Vite). It uses the local SDK to interact with the Daemon.

---

## 3. Step-by-Step Initialization

### Step 3.1: Start the P2P Daemon
The daemon must be running *before* the Gateway UI is opened, as the UI relies on its local HTTP server.

1. Open a terminal.
2. Navigate to the daemon directory:
   ```bash
   cd src/daemon-go
   ```
3. Run the daemon:
   ```bash
   go run .
   ```
4. **Expected Output:** You should see logs indicating that the local peer ID has been generated, local storage (`.aether`) has been initialized, and the HTTP server is listening on `:5000`. Leave this terminal open.

> [!TIP]
> All local state (keys, chunks, manifests, and reputation) is stored in a `.aether` directory within the daemon folder. If you ever need to perform a "factory reset" of your node, simply stop the daemon, delete the `.aether` folder, and restart.

### Step 3.2: Build the TypeScript SDK
The Gateway UI depends on the local SDK. We must build it so the UI can consume it.

1. Open a **second** terminal.
2. Navigate to the SDK directory:
   ```bash
   cd packages/sdk
   ```
3. Install dependencies and build:
   ```bash
   npm install
   npm run build
   ```
4. **Expected Output:** A `dist/` directory will be generated containing the compiled `index.js` and TypeScript declaration files.

### Step 3.3: Launch the Gateway UI
With the daemon running and the SDK compiled, we can now launch the user interface.

1. In your second terminal, navigate to the Gateway directory:
   ```bash
   cd aether-gateway
   ```
2. Launch the application in development mode using Wails:
   ```bash
   wails dev
   ```
3. **Expected Output:** Wails will automatically compile the frontend (Vite/React) and the Go backend, then open a native Desktop window displaying the Aether Gateway lock screen.

### Step 3.4: Launch Delphi
To view the live Delphi tracker feed:
1. Open your browser.
2. Navigate to `http://localhost:5000/ui/`.

---

## 4. Alpha Testing across the Internet
If you want to test Aether with a friend on another network:

1. **Build Cross-Platform Native Binaries**: 
   Run `.\build_alpha.ps1` from the root directory. This will compile optimized standalone binaries into the `\bin` directory (e.g. Windows amd64 and Mac M1 arm64). Send the correct executable to your friend.
2. **Exchange Multiaddresses**: 
   Start the daemon. It will output your Public/Local IP Multiaddresses. Send the public Multiaddress to your friend.
3. **Connect as a Bootstrap Node**: 
   Your friend starts their daemon using the environment variable: 
   `AETHER_BOOTSTRAP_PEERS=<your_multiaddress> ./aether-daemon-mac-arm64`
4. **Publish and Fetch Data**: 
   Use the testing scripts `scripts/publish.ts` and `scripts/fetch.ts` to seamlessly transfer a file!

## 4. Application User Guide

Once the Gateway is open, you will be interacting with the true decentralized network. Here is how to navigate the CTOS interface:

### 4.1 Initial System Unlock
When you first launch the app, you will be greeted by the lock screen.
- **Setup Mode:** If you do not have a master password set, enter a secure `NEW_PASSWORD` and click **INITIALIZE_SYSTEM**. This encrypts your local Ed25519 keyring.
- **Login Mode:** On subsequent boots, enter your master password to decrypt the system.

### 4.2 [01] Identity.Mgr (Generating your Alias)
Before publishing data, you need a decentralized identity.
1. Enter your desired human-readable alias (e.g., `alice`).
2. Click **Generate**. The daemon will generate an Ed25519 keypair and associate it locally with this alias.
3. Click **Broadcast_Claim**. This utilizes the "First-Seen" consensus mechanism. It sends a cryptographically signed envelope to the network claiming the username. If no older claim exists, the alias is yours.

### 4.3 [02] Data.Uplink (Publishing to the DHT)
You can now inject JSON payloads into the global network.
1. Write a valid JSON payload in the text area (e.g., `{ "message": "Hello Decentralized World" }`).
2. Click **Execute_Uplink**. 
3. The SDK routes this to the daemon, which Brotli-compresses the data, generates a manifest, signs a Cryptographic Envelope, and broadcasts a `DELIVER_USER_CAPSULES` packet to the swarm.
4. Note the resulting `CAPSULE_ID` (the hash of the envelope).

### 4.4 [03] Net.Trace (Exploring the Network)
Trace allows you to look up other users and read their data.
1. Enter a `TARGET_ALIAS` (e.g., `bob`) and click **Trace_Target**.
2. The daemon opens a 3-second context window, broadcasts a `WANT_USER_CAPSULES` request, deduplicates identical network responses, and returns the unique capsules to the UI.
3. You will see Capsule Cards showing the Publisher, Timestamp, and **Swarm Health** (the number of peers globally seeding this data).
4. Click **Decode_Payload**. The daemon will perform a P2P fetch for the underlying chunks, decompress them, cache them locally (turning you into a seeder), and display the raw JSON.

### 4.5 [04] Net.Swarm (Reputation Scoreboard)
Aether operates on a local Proof-of-Storage reputation system.
- Click **Refresh_Scores** to view your local scoreboard.
- This panel ranks peers based on how many bytes of data they have successfully seeded on your behalf. There are no global leaderboards—this is a subjective, gamification-resistant Web of Trust metric.

### 4.6 [05] Sys.Control (Settings)
Configure your local node parameters:
- **Disk Allocation:** Set the gigabyte limit for your `.aether/global_cache` directory.
- **Network Mode:** Toggle between Online (Active Kademlia routing) and Offline (Read-only local cache).
- **Lock Session:** Instantly drops the decryption keys from RAM and returns you to the lock screen.

---

## 5. Troubleshooting Common Issues

- **Daemon fails to start with "Address already in use":** Port 5000 is occupied. Kill any lingering Node or Go processes running on that port.
- **Gateway hangs on "Decoding...":** The capsule data may have been dropped from the network if no seeds are available. The robust SDK will automatically timeout and return a `404` without crashing the application.
- **Unresponsive UI during `wails dev`:** Ensure you ran `npm install` inside `aether-gateway/frontend` so the React dependencies (and the locally linked `aether-sdk`) are present.
