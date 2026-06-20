package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"aether-daemon"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func sha256Bytes(data []byte) string {
	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:])
}

func sha256Json(value interface{}) string {
	bytes, _ := json.Marshal(value)
	return sha256Bytes(bytes)
}

func safeName(val string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9.-]`)
	return reg.ReplaceAllString(val, "_")
}

func getAetherDir() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, ".aether")
}

func ensureDir(dir string) {
	os.MkdirAll(dir, 0755)
}

func nowIso() string {
	return time.Now().UTC().Format(time.RFC3339)
}

type Identity struct {
	KeyID string `json:"key_id"`
	Label string `json:"label"`
}

type UserProfile struct {
	Version      int                    `json:"version"`
	Username     string                 `json:"username"`
	CreatedAt    string                 `json:"created_at"`
	PrimaryKeyID string                 `json:"primary_key_id"`
	Keys         map[string]interface{} `json:"keys"`
	RevokedKeys  map[string]interface{} `json:"revoked_keys"`
}

func (a *App) InitUser(username string) (string, error) {
	username = strings.TrimSpace(strings.ToLower(username))
	if matched, _ := regexp.MatchString(`^[a-z0-9][a-z0-9_-]{2,31}$`, username); !matched {
		return "", fmt.Errorf("invalid username format")
	}

	aetherDir := getAetherDir()
	registryDir := filepath.Join(aetherDir, "usernames")
	ensureDir(registryDir)
	registryFile := filepath.Join(registryDir, username+".json")

	if _, err := os.Stat(registryFile); err == nil {
		return "", fmt.Errorf("username is currently active and taken: %s", username)
	}

	// Generate Ed25519 Key
	priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519, 256)
	if err != nil {
		return "", err
	}

	pubBytes, _ := crypto.MarshalPublicKey(pub)
	privBytes, _ := crypto.MarshalPrivateKey(priv)

	keyID := sha256Bytes(pubBytes)
	userDir := filepath.Join(aetherDir, "users", safeName(keyID))
	keyDir := filepath.Join(userDir, "keys", safeName(keyID))
	ensureDir(keyDir)

	// Optionally encrypt
	var finalBytes []byte
	if daemon.MasterKey != nil && len(daemon.MasterKey) == 32 {
		enc, err := daemon.EncryptBytes(privBytes)
		if err != nil {
			return "", err
		}
		finalBytes = enc
	} else {
		finalBytes = privBytes
	}

	os.WriteFile(filepath.Join(keyDir, "public.key"), pubBytes, 0644)
	os.WriteFile(filepath.Join(keyDir, "private.key"), finalBytes, 0600)

	publicRecord := map[string]interface{}{
		"key_id":          keyID,
		"algorithm":       "ed25519",
		"label":           "primary",
		"status":          "active",
		"created_at":      nowIso(),
		"public_key_file": "public.key",
	}

	profile := UserProfile{
		Version:      1,
		Username:     username,
		CreatedAt:    nowIso(),
		PrimaryKeyID: keyID,
		Keys:         map[string]interface{}{keyID: publicRecord},
		RevokedKeys:  map[string]interface{}{},
	}

	profileBytes, _ := json.MarshalIndent(profile, "", "  ")
	os.WriteFile(filepath.Join(userDir, "profile.json"), profileBytes, 0644)

	err = daemon.InitializeLocalRegistryClaim(username, keyID)
	if err != nil {
		return "", err
	}

	return keyID, nil
}

func (a *App) GetIdentity(username string) (string, error) {
	registryFile := filepath.Join(getAetherDir(), "usernames", username+".json")
	bytes, err := os.ReadFile(registryFile)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}
	var registry map[string]interface{}
	json.Unmarshal(bytes, &registry)

	epochs := registry["epochs"].([]interface{})
	lastEpoch := epochs[len(epochs)-1].(map[string]interface{})
	return lastEpoch["bound_identity"].(string), nil
}

func (a *App) ReleaseIdentity(username string) error {
	username = strings.TrimSpace(strings.ToLower(username))
	keyId, err := a.GetIdentity(username)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	msg, err := daemon.CreateSignedRegistryMessage("release", username, keyId, now)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(msg)

	daemon.ProcessRegistryMessage(b)
	daemon.BroadcastRegistryMessage(b)

	return nil
}
