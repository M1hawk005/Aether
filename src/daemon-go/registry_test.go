package daemon

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

func withRegistryTestDir(t *testing.T) {
	t.Helper()
	oldDir := aetherDir
	dir := ".aether-registry-test"
	os.RemoveAll(dir)
	aetherDir = dir
	t.Cleanup(func() {
		aetherDir = oldDir
		os.RemoveAll(dir)
	})
}

func createRegistryTestIdentity(t *testing.T) string {
	t.Helper()
	priv, pub, err := p2pcrypto.GenerateKeyPair(p2pcrypto.Ed25519, 256)
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	pubBytes, err := p2pcrypto.MarshalPublicKey(pub)
	if err != nil {
		t.Fatalf("MarshalPublicKey failed: %v", err)
	}
	privBytes, err := p2pcrypto.MarshalPrivateKey(priv)
	if err != nil {
		t.Fatalf("MarshalPrivateKey failed: %v", err)
	}

	keyID := Sha256Bytes(pubBytes)
	keyDir := PathInAether("users", SafeName(keyID), "keys", SafeName(keyID))
	if err := EnsureDir(keyDir); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(keyDir, "public.key"), pubBytes, 0644); err != nil {
		t.Fatalf("write public key failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(keyDir, "private.key"), privBytes, 0600); err != nil {
		t.Fatalf("write private key failed: %v", err)
	}
	return keyID
}

func processMessage(t *testing.T, msg RegistryMessage) {
	t.Helper()
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal message failed: %v", err)
	}
	ProcessRegistryMessage(data)
}

func TestRegistryAcceptsSignedClaim(t *testing.T) {
	withRegistryTestDir(t)
	keyID := createRegistryTestIdentity(t)
	msg, err := CreateSignedRegistryMessage("claim", "alice", keyID, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("CreateSignedRegistryMessage failed: %v", err)
	}

	processMessage(t, msg)

	var reg Registry
	if err := ReadJson(PathInAether("usernames", "alice.json"), &reg); err != nil {
		t.Fatalf("expected registry file to be written: %v", err)
	}
	if reg.PrimaryKeyID != keyID {
		t.Fatalf("expected primary key %s, got %s", keyID, reg.PrimaryKeyID)
	}
}

func TestRegistryRejectsUnsignedClaim(t *testing.T) {
	withRegistryTestDir(t)
	keyID := createRegistryTestIdentity(t)
	msg, err := CreateSignedRegistryMessage("claim", "alice", keyID, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("CreateSignedRegistryMessage failed: %v", err)
	}
	msg.Signature = "not-a-valid-signature"

	processMessage(t, msg)

	if _, err := os.Stat(PathInAether("usernames", "alice.json")); !os.IsNotExist(err) {
		t.Fatalf("forged claim should not create registry file")
	}
}

func TestRegistryRejectsPublicKeyMismatch(t *testing.T) {
	withRegistryTestDir(t)
	aliceKeyID := createRegistryTestIdentity(t)
	bobKeyID := createRegistryTestIdentity(t)
	msg, err := CreateSignedRegistryMessage("claim", "alice", aliceKeyID, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("CreateSignedRegistryMessage failed: %v", err)
	}
	msg.KeyID = bobKeyID

	processMessage(t, msg)

	if _, err := os.Stat(PathInAether("usernames", "alice.json")); !os.IsNotExist(err) {
		t.Fatalf("mismatched public key claim should not create registry file")
	}
}

func TestRegistryReleaseRequiresActiveOwnerSignature(t *testing.T) {
	withRegistryTestDir(t)
	aliceKeyID := createRegistryTestIdentity(t)
	attackerKeyID := createRegistryTestIdentity(t)

	claim, err := CreateSignedRegistryMessage("claim", "alice", aliceKeyID, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("CreateSignedRegistryMessage claim failed: %v", err)
	}
	processMessage(t, claim)

	forgedRelease, err := CreateSignedRegistryMessage("release", "alice", attackerKeyID, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("CreateSignedRegistryMessage release failed: %v", err)
	}
	forgedRelease.KeyID = aliceKeyID
	processMessage(t, forgedRelease)

	var reg Registry
	if err := ReadJson(PathInAether("usernames", "alice.json"), &reg); err != nil {
		t.Fatalf("expected registry file: %v", err)
	}
	if reg.Epochs[0].ValidUntil != "" {
		t.Fatalf("forged release should not close active epoch")
	}

	release, err := CreateSignedRegistryMessage("release", "alice", aliceKeyID, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("CreateSignedRegistryMessage release failed: %v", err)
	}
	processMessage(t, release)

	if err := ReadJson(PathInAether("usernames", "alice.json"), &reg); err != nil {
		t.Fatalf("expected registry file: %v", err)
	}
	if reg.Epochs[0].ValidUntil == "" {
		t.Fatalf("valid release should close active epoch")
	}
}
