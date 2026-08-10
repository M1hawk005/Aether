# Aether Project Todo

The canonical prioritized backlog is maintained in [docs/Implementation_Backlog.md](docs/Implementation_Backlog.md). Tasks are grouped for a Notion Kanban using:

- **P0:** architecture or security blockers;
- **P1:** required for the first useful Aether Index release;
- **P2:** publishing, application runtime, and beta reliability;
- **P3:** scale, ecosystem, and advanced applications.

Milestones sequence the work: **M0** foundations, **M1** search client, **M2** publishing, **M3** runtime, **M4** beta, **M5** scale.

## Thesis

Aether reduces reliance on centralized corporate infrastructure through **offload** (bytes to peers, computation to the user's own device) and **ownership** (portable identity and records, strictly opt-in data collection). A central server is permitted. The rule is *the server orders, but does not own*.

## Current focus

- [ ] Complete all P0/M0 architecture contracts and security foundations.
- [ ] Ship **M1: Aether Index as a search client** — sync two signed catalogs, verify, search offline, rank on-device, hand off a verified torrent, export everything.
- [ ] Settle the four load-bearing schema rules before publishers sign against the format: partial-replication validation, public-by-design content, counts as provider assertions, reserved reply/reference semantics.
- [ ] Define the pluggable admission-proof interface with per-method weights; attestation is the default path.
- [ ] Answer the two assumptions the project rests on: will anyone operate a catalog, and will a publisher sign a release?
- [ ] Do not extend the custom tracker, global GossipSub blob transfer, or first-seen username system.
- [ ] Do not pull publishing (M2) or the application runtime (M3) forward to make a demo look complete.

## Deliberate non-goals

- Eliminating servers. The goal is small, cheap, contestable, exit-safe servers.
- General-purpose volunteer compute. Verification cost and residential energy economics do not support it; transcoding-class work is the one exception.
- Processing other users' private data on volunteer machines.
- Any chain as a canonical dependency, or per-action records onchain.

## Completed prototype work retained for reference

- [x] Go daemon and local HTTP boundary.
- [x] Ed25519 identity/signature prototype.
- [x] TypeScript SDK boundary.
- [x] Wails desktop gateway shell.
- [x] Delphi tracker/catalog proof-of-concept UI.
- [x] Local content publication and retrieval prototype.

These completed items prove product concepts but do not imply production readiness or conformance with the target architecture.
