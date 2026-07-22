# Aether Project Todo

The canonical prioritized backlog is maintained in [docs/Implementation_Backlog.md](docs/Implementation_Backlog.md). Tasks are grouped for a Notion Kanban using:

- **P0:** architecture or security blockers;
- **P1:** required for the first useful Aether Index release;
- **P2:** required for a credible public beta;
- **P3:** scale, ecosystem, and advanced applications.

## Current focus

- [ ] Complete all P0 architecture contracts and security foundations.
- [ ] Prove publish -> catalog -> search -> download -> verify -> isolated-run end to end.
- [ ] Do not extend the custom tracker, global GossipSub blob transfer, or first-seen username system.

## Completed prototype work retained for reference

- [x] Go daemon and local HTTP boundary.
- [x] Ed25519 identity/signature prototype.
- [x] TypeScript SDK boundary.
- [x] Wails desktop gateway shell.
- [x] Delphi tracker/catalog proof-of-concept UI.
- [x] Local content publication and retrieval prototype.

These completed items prove product concepts but do not imply production readiness or conformance with the target architecture.
