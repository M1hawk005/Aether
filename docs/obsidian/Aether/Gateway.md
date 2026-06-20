# Desktop Gateway (`aether-gateway`)

The Gateway is the official desktop graphical user interface (GUI) for the Aether Network. While the [Daemon](./Daemon.md) operates headlessly in the background, the Gateway provides a rich, user-friendly interface for browsing, publishing, and managing identity.

[Back to Home](./Home.md)

## Framework Architecture

The Gateway is built using **Wails**, an Electron alternative that wraps native OS webviews.
- **Backend:** Go (Golang). Handles OS-level window management, system tray, and secure file I/O.
- **Frontend:** React & TypeScript, compiled via Vite. 

## Communication Flow

1. **User Interaction:** The user interacts with the React frontend.
2. **SDK Integration:** The frontend uses the [Aether SDK](./SDK.md) (`AetherClient`) to send HTTP REST commands directly to the local Node.js daemon running on port `5000`.
3. **Go Bridge:** For actions requiring high security (like saving passwords or retrieving the local identity cryptographic keys), the frontend uses Wails' IPC bridge to call native Go functions (`App.go`).

## Security & Identity

Aether utilizes Ed25519 cryptography for decentralized identity.
- When a user first opens the Gateway, they are prompted to create a Master Password.
- The Go backend securely encrypts the Ed25519 private key to disk using this password.
- The UI handles the decryption in-memory to sign payloads before pushing them through the local daemon.

> [!CAUTION]
> The private key never leaves the user's local machine, and the master password is never transmitted to the daemon or the network.
