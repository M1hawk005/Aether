---
sidebar_position: 1
title: Quick Start
---

# Quick Start Guide

> **Migration notice:** This page exercises the legacy capsule API. The accepted target architecture uses signed events, standard torrents, catalog providers, and an isolated application runtime. Do not treat these calls as the stable v1 API.

Welcome to the Aether Protocol! This guide will get you building decentralized applications in minutes.

## 1. Prerequisites
You need Node.js installed, and the Aether Daemon running in the background.

```bash
# Clone the repository
git clone https://github.com/aether-network/aether.git
cd aether

# Install dependencies
npm install

# Start the daemon
npm run cli init-user --username my_app
```

## 2. Using the SDK
The `aether-sdk` is the easiest way to interact with the network.

### Publishing Data
```javascript
const AetherSDK = require('./src/sdk.js');

async function main() {
  const sdk = new AetherSDK({ daemonUrl: 'http://127.0.0.1:5000' });
  
  // Data is automatically Brotli-compressed, signed, and seeded!
  const capsuleId = await sdk.publish('my_app', 'Hello decentralized world!', {
    persistent: true
  });
  
  console.log('Capsule Published:', capsuleId);
}
main();
```

### Fetching Data
```javascript
const AetherSDK = require('./src/sdk.js');

async function main() {
  const sdk = new AetherSDK();
  const data = await sdk.fetch('sha256:your_capsule_id_here');
  
  console.log('Fetched data:', data);
}
main();
```
