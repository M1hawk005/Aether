# Registry Policy

The username registry should be a lightweight, bundled registry chain. It tracks
names, keys, namespace authority, leases, and disputes. It does not store user
content.

## Consensus Direction

Proof of stake can make sense if the registry becomes a shared global database, but
the first implementation should keep the registry small and application-bundled.

Because we cannot launch with a massive, permissionless Proof-of-Stake network on Day 1,
the registry will start as a **Proof-of-Authority (PoA) Consortium**. The initial trusted
validators can vote on-chain to add new validators as the community grows. 

Recommended progression:

1. Local append-only registry chain for the prototype.
2. **Proof-of-Authority (PoA) Consortium** for early network tests and controlled growth.
3. Because the underlying data structure remains the same, seamlessly swap the consensus mechanism to an efficient **Proof-of-Stake (PoS)** validator set as the network matures, without breaking usernames or causing a chaotic migration.
4. Slashing or governance only for registry misbehavior, not content moderation.

The registry is for name and authority state, not for routing every content request.

## Identity vs Username

A cryptographic identity is stable:

```text
did:key:z6Mk...
```

A username is a human-readable pointer:

```text
alice -> did:key:z6Mk...
```

Disputes are about who controls a name, not who cryptographically authored old
content.

## Name Classes

Personal name:
`alice`

Community name:
`photography.club`

Verified external name:
`alice.com`, `alice@example.com`, `@alice` on another service

Reserved/protected name:
`bank`, `police`, `aether`, `support`

Do not treat all names equally. A random unused handle and a real-world brand need
different dispute rules.

## First Version Policy

Avoid launching with universal scarce usernames.

Use:

- Cryptographic identities as source of truth.
- Non-unique display names.
- Scoped names from indexes or communities.
- Domain-backed names for stronger claims.
- Lease-based names only where a registry exists.

## Deletion And Reissue (Time-Aware Epochs)

Users should be able to delete or release a username from the registry and reissue it
with different keys. Because Aether uses cryptographic keys (DIDs) as the true identity and the username is just a pointer, the registry handles this by maintaining a historical ledger of "leases" known as **Time-Aware Epochs**.

Example of Epochs:
- **Epoch 1:** `@alice` belonged to Key A (Jan 2024 — Dec 2025)
- **Epoch 2:** `@alice` belongs to Key B (Jan 2026 — Present)

Every binding needs time/version history:

```json
{
  "name": "alice",
  "bound_identity": "did:key:old...",
  "valid_from": "2024-01-01T00:00:00Z",
  "valid_until": "2025-12-31T23:59:59Z",
  "registry": "aether://registry/main",
  "signature": "sig:..."
}
```

### Time-Aware Client Resolution
When an Aether client downloads a post, it checks the signature on the post (which includes a timestamp) and checks the registry. If it downloads a post signed by Key A in 2024, it confirms: "Yes, Key A legally owned the name `@alice` in 2024."

### Visual Distinction in the UI
The UI must treat historical epochs as entirely different accounts from the current owner to prevent getting credit or blame for an old user's posts:
- **The Current Owner:** Displays normally as `@alice`, which links to the new owner's profile (Key B).
- **The Historical Owner:** Visually breaks the link to the current owner. Displays as `@alice (Historical, 2024-2025)`, with a greyed-out badge or broken-link icon. Clicking does not route to the new owner, but to an archived view for Key A.
- **Zero Ambiguity for Power Users:** Hovering over the username in the UI displays a tooltip with the raw cryptographic ID (`did:key:z6Mk...`), proving unequivocally that the past author is mathematically distinct from the current author.

## Abandoned Names

Names can expire only when the registry explicitly uses leases.

Suggested defaults:

- Free personal names expire after 12-24 months of no signed activity.
- Paid or verified names expire at lease end with a renewal grace period.
- Community names expire according to community registry policy.
- External-domain names expire when the domain proof no longer validates.

Expiry flow:

1. Mark inactive.
2. Notify known recovery contacts/devices.
3. Enter grace period.
4. Enter reclaimable state.
5. Allow new claim.
6. Preserve historical binding records.

## Disputes

Disputes are signed claims with evidence and registry-scoped outcomes.

Reasons:

- Impersonation.
- Trademark or brand conflict.
- Abandoned name.
- Account compromise.
- Harassment or squatting.
- Verified external identity conflict.

Outcomes:

- No action.
- Warning label.
- Suspension.
- Transfer after grace period.
- Reservation/protection.
- Split namespace recommendation.

