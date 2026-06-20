# Gateway UX

The first gateway should feel less like a crypto browser and more like a calm,
trustworthy reader. Protocol state should be visible only when useful.

## Core Job

Let a normal person open an Aether capsule, understand who published it, know whether
it is verified, cached, live, or offline, and share or save it without learning
protocol mechanics.

## First Screen

Start with a simple address/search bar:

```text
Open capsule, site, forum, or index...
```

Below it, show useful places to continue:

- Pinned.
- Recent.
- From followed indexes.
- Downloaded for offline.

Avoid a dashboard full of protocol stats. The first impression should be: "I can
browse things here."

## Capsule View

Each capsule, site, or forum should show a compact trust/status strip:

```text
alice.site
Verified publisher
Cached by 18 peers
Last updated 2 hours ago
Available offline
```

Use plain states:

- Verified.
- Unverified.
- Updated.
- Offline copy.
- Origin unavailable.
- Peer fallback active.
- Blocked by your policy.
- Encrypted, access required.

Avoid raw protocol terms like DHT, CID, provider record, epoch, or grant unless the
user opens details.

## Gateway Modes

Browse:
View public capsules through a normal web UI.

Save:
Keep a capsule available offline or help cache it.

Publish:
Run or connect an origin node and publish updates.

Browse should be primary. Save should feel like a natural bookmark/offline action.
Publish can be present but should not dominate the first version.

## Failure States

Failure messages should describe what the client is doing:

- Looking for peers.
- Trying publisher origin.
- Showing last verified offline copy.
- This capsule is encrypted. Request access?
- This content is hidden by one of your trust lists.

## Trust Popover

Each page should have one compact detail view:

```text
Publisher:
  Alice Example
  did:key:...

Content:
  Signed version
  Updated May 27, 2026
  Manifest visible/public

Sources:
  12 peers
  publisher origin
  1 relay fallback

Policy:
  cached for offline use
  allowed by your trust lists
```

Default view stays simple. Advanced details are one click away.

## Product Shape

The first product should be a local gateway plus lightweight browser shell.

Local gateway:
Handles protocol, caching, verification, grants, peer fetching, and trust policy.

Browser shell:
Provides address bar, identity/trust UI, offline controls, and publish entry points.

Public capsules can also be viewed from normal browsers through local or hosted
gateways:

```text
http://localhost:8787/capsule/...
https://gateway.aether.example/alice.site
```

The first UX worth proving:

> Open a thing, verify it, cache it, keep reading when the origin disappears.

