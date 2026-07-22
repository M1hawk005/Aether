---
title: Development
sidebar_position: 3
---

# Development

## Requirements

- Go version declared in `src/daemon-go/go.mod`
- Node.js 20 or later
- npm
- Wails CLI for desktop gateway work

## Install dependencies

```powershell
npm ci --prefix packages/sdk
npm ci --prefix aether-docs
npm ci --prefix aether-gateway/frontend
```

## Run the current daemon prototype

```powershell
cd src/daemon-go
go run ./cmd/aether-daemon
```

The daemon listens on `http://127.0.0.1:5000`. Its current libp2p, capsule, tracker, reputation, and alias endpoints belong to the migration prototype. New work should follow [Aether Architectural Reference](./Aether_Architecture.md) and the P0 tasks in [Implementation Backlog](./Implementation_Backlog.md).

## Run the desktop gateway

```powershell
cd aether-gateway
wails dev
```

## Run checks

From the repository root:

```powershell
npm test
npm run test:docs
```

Run the gateway frontend checks separately:

```powershell
npm run build --prefix aether-gateway/frontend
```

## Documentation

The Markdown files in `docs/` are canonical. The Docusaurus project reads that directory directly.

```powershell
npm start --prefix aether-docs
```

Do not add generated documentation output, coverage reports, restored download output, editor state, or scratch migration scripts to Git.
