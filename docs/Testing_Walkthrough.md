# Aether Network: Detailed Testing Walkthrough Plan

> [!WARNING]
> This walkthrough targets a superseded Node/WebRTC prototype and is preserved only as historical test context. It is not a validation plan for the hybrid BitTorrent architecture. The replacement end-to-end acceptance flow is defined in [Implementation Backlog](./Implementation_Backlog.md).

This guide outlines exactly how to test the Aether Network as a neutral, decentralized file-sharing platform across two devices on the same Local Area Network (LAN). 

Because we integrated **mDNS (Zero-Configuration Discovery)**, you no longer need to worry about manually typing IP addresses or editing `network.json`. The nodes will find each other automatically.

---

## Prerequisites
1. **Devices:** One Windows PC and one Mac.
2. **Network:** Both devices must be connected to the same LAN (e.g., the same WiFi network).
3. **Environment:** Node.js must be installed on both machines.
4. **Codebase:** Ensure the latest version of the Aether codebase is copied to both machines.

> [!IMPORTANT]  
> If you have a strict local firewall (like Windows Defender), ensure that Node.js is allowed to communicate over public/private networks, specifically for UDP port `5353` (used by mDNS for discovery) and the daemon ports (`5000` HTTP, `6000` WebSocket, `7000` Kademlia UDP).

---

## Phase 1: Start the Initial Node (Windows)

First, we will spin up the daemon on the Windows machine. It will bind to `0.0.0.0` (all network interfaces), allowing it to accept connections from the Mac.

1. Open PowerShell and navigate to the Aether project directory.
2. Start the daemon using the CLI:
   ```bash
   node src/cli/index.js start-daemon --host 0.0.0.0 --port 5000 --peer node-windows
   ```
3. You should see logs confirming it is listening on HTTP `5000`, Gossip WS `6000`, and DHT `7000`. You will also see it initialize the mDNS discovery engine.

---

## Phase 2: Start the Second Node (Mac)

Now, we will spin up the daemon on the Mac. It will also bind to `0.0.0.0`. Thanks to mDNS, the Mac should immediately detect the Windows node without any manual configuration.

1. Open Terminal on the Mac and navigate to the Aether project directory.
2. Start the daemon:
   ```bash
   node src/cli/index.js start-daemon --host 0.0.0.0 --port 5000 --peer node-mac
   ```
3. **Verification:** Watch the console output on both machines. 
   - On the Mac, you should see a log like: `[mDNS] Discovered local Aether node: http://<Windows-IP>:5000`
   - Both terminals should log that a `HANDSHAKE` was exchanged and a `WebRTC DataChannel` was successfully established between them.

---

## Phase 3: Publish Content (Windows)

We will now create a user identity on the Windows machine and publish a file into the network.

1. Open a **second** PowerShell tab on the Windows machine (leave the daemon running in the first tab).
2. Initialize a new user identity:
   ```bash
   node src/cli/index.js init-user --username win-user
   ```
3. Create a test file, for example, a text file or an image named `test-file.jpg`.
4. Publish the file to the network:
   ```bash
   node src/cli/index.js publish --username win-user --input ./test-file.jpg
   ```
5. **The Output:** The CLI will output a long string like `sha256:...`. **This is the Capsule ID.** Copy this ID. 
   *Behind the scenes, the Windows daemon has shredded the file, built a Merkle Tree, and announced to the Mac (via Kademlia DHT) that it is seeding this content.*

---

## Phase 4: Fetch Content (Mac)

Finally, we will request the file from the Mac, proving that the trackerless WebRTC swarm fetching works seamlessly across devices.

1. Open a **second** Terminal tab on the Mac (leave the daemon running).
2. Initialize a new user identity for the Mac:
   ```bash
   node src/cli/index.js init-user --username mac-user
   ```
3. Fetch the file using the Capsule ID you copied from the Windows machine:
   ```bash
   node src/cli/index.js fetch --capsule <PASTE_CAPSULE_ID_HERE> --out ./received-file.jpg
   ```
4. **Verification:** 
   - The Mac daemon will query the DHT, discover the Windows node holds the file, and request the chunks over the high-speed WebRTC connection.
   - The command will exit successfully, and `received-file.jpg` will appear on your Mac!
   - Open `received-file.jpg` to verify it perfectly matches the original file.

> [!TIP]  
> To thoroughly test the architecture, reverse the roles for Phase 3 and 4! Publish a file on the Mac and fetch it on Windows to ensure bi-directional NAT traversal and routing are working flawlessly.
