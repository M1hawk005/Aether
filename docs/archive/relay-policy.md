# Relay Policy

Relays are constrained fallback infrastructure. They must not become the default
content path or a central choke point.

## Hard Rules

1. Relay only when direct peer or origin connection fails.
2. Relay capacity is capped per peer, per capsule, and per relay.
3. Relay selection is multi-provider and client-randomized.
4. Relayed traffic is short-lived and migrates back to direct peers when possible.

## Request Order

Clients should try:

1. Local cache.
2. Direct peers.
3. Relay-assisted peer connection.
4. Origin.
5. Relay-served origin/cache as last resort.

## Relay Capacity Envelope

```json
{
  "relay": "did:key:relay",
  "region": "ap-southeast",
  "max_concurrent_streams": 2000,
  "max_bytes_per_peer_per_day": 524288000,
  "max_bytes_per_capsule_per_day": 10737418240,
  "max_stream_duration_seconds": 300,
  "supports": ["nat-assist", "rendezvous", "limited-forwarding"],
  "price": "free-tier-or-paid",
  "signature": "sig:..."
}
```

## Selection

Clients should use weighted random selection:

```text
candidate_relays =
  trusted relays
  + relays near client
  + relays near publisher/peer
  + relays under capacity
```

Pick randomly from the top N, not always the top 1.

## Diversity

- Minimum 3 relay providers per region.
- No single relay gets more than a configured share of a capsule's relayed traffic.
- Clients rotate relays between sessions.
- Publishers cannot require one exclusive relay for public content.
- Indexes should list multiple relay options.

## Allowed

- NAT hole-punch assistance.
- Short-lived stream forwarding.
- Rendezvous.
- Temporary bootstrap cache.

## Not Default

- Becoming canonical origin.
- Long-term exclusive hosting.
- Rewriting content.
- Requiring account login for public reads.
- Seeing plaintext for encrypted/private content.

## Popular Content

If relayed traffic rises:

- Increase peer discovery fanout.
- Encourage more caching.
- Notify publisher origin.
- Rate-limit relay use.
- Shift clients to direct peers as soon as possible.

## Refusals

Relays may refuse blocked hashes, abusive peers, traffic above envelope, or capsules
exceeding fair-use limits. Where safe, they should publish signed refusal reasons.

```json
{
  "target": "capsule://example",
  "action": "relay_refused",
  "reason_code": "capacity_exceeded",
  "scope": "this-relay",
  "signature": "sig:..."
}
```

The strongest anti-choke-point rule:

> A relay must never be required to verify, discover, fetch, or publish public content.

