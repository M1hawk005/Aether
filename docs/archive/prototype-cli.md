# Prototype CLI

The current prototype is intentionally local-only. It proves the content capsule
lifecycle before networking, incentives, browser integration, or production key
management are introduced.

## Commands

Create an identity:

```sh
node src/aether.js init-user --username alice
```

Add another Ed25519 key to the same user:

```sh
node src/aether.js add-key --username alice --label laptop
```

Revoke a compromised key:

```sh
node src/aether.js revoke-key --username alice --key sha256:...
```

Create a forum controlled by an operator:

```sh
node src/aether.js create-forum --forum forum:aether --operator alice
```

Publish a file as an encrypted signed capsule:

```sh
node src/aether.js publish --username alice --input examples/hello.html --namespace site:hello
```

Submit a capsule to a forum. This expresses the publisher's intent to show the post
there, but does not list it until the forum accepts it:

```sh
node src/aether.js submit-forum --username alice --forum forum:aether --capsule sha256:...
```

Approve the submission as the forum operator:

```sh
node src/aether.js approve-forum --operator alice --forum forum:aether --capsule sha256:...
```

List forum state:

```sh
node src/aether.js list-forum --forum forum:aether
```

Publish private or unlisted content with encrypted chunks and manifest:

```sh
node src/aether.js publish --username alice --input secret.txt --visibility private
```

Inspect a capsule header:

```sh
node src/aether.js inspect --capsule sha256:...
```

Fetch and reconstruct a capsule:

```sh
node src/aether.js fetch --capsule sha256:... --out restored.html
```

Revoke a capsule:

```sh
node src/aether.js revoke --username alice --capsule sha256:...
```

## What It Proves

- Users have unique usernames backed by Ed25519 keysets.
- Users can own multiple active keys and revoke compromised keys.
- Username, key, and forum events are appended to a local registry chain.
- Forum listing is operator-approved instead of automatic.
- Public content can be cached without encryption by default.
- Private, unlisted, and opaque-public content is encrypted before being cached.
- Chunk integrity is verified by hash.
- Public manifests and encrypted manifests are both supported.
- A local namespace index can reference capsules.
- Revocation can delete compliant peer chunks and remove local readable access.

## What Is Still Fake

- Access grants are represented by a local key file.
- Peers are folders, not network processes.
- Cache selection does not yet consider proximity, speed, storage, or popularity.
- PGP identity is not implemented yet.
- Deletion only applies to compliant local peers.
