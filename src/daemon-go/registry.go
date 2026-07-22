package daemon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type RegistryMessage struct {
	Type      string `json:"type"` // "claim" or "release"
	Username  string `json:"username"`
	KeyID     string `json:"key_id"`
	PublicKey string `json:"public_key"`
	Timestamp string `json:"timestamp"` // RFC3339
	Signature string `json:"signature"`
}

type registryMessagePayload struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	KeyID     string `json:"key_id"`
	PublicKey string `json:"public_key"`
	Timestamp string `json:"timestamp"`
}

func registryPayload(msg RegistryMessage) registryMessagePayload {
	return registryMessagePayload{
		Type:      msg.Type,
		Username:  strings.ToLower(strings.TrimSpace(msg.Username)),
		KeyID:     msg.KeyID,
		PublicKey: msg.PublicKey,
		Timestamp: msg.Timestamp,
	}
}

func CreateSignedRegistryMessage(msgType string, username string, keyId string, timestamp string) (RegistryMessage, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	if timestamp == "" {
		timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	publicKeyPath := PathInAether("users", SafeName(keyId), "keys", SafeName(keyId), "public.key")
	privateKeyPath := PathInAether("users", SafeName(keyId), "keys", SafeName(keyId), "private.key")
	publicKeyBytes, err := ReadFileBytes(publicKeyPath)
	if err != nil {
		return RegistryMessage{}, fmt.Errorf("failed to load registry public key: %w", err)
	}
	privateKey, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return RegistryMessage{}, fmt.Errorf("failed to load registry private key: %w", err)
	}

	msg := RegistryMessage{
		Type:      msgType,
		Username:  username,
		KeyID:     keyId,
		PublicKey: base64.StdEncoding.EncodeToString(publicKeyBytes),
		Timestamp: timestamp,
	}

	signature, err := SignJson(privateKey, registryPayload(msg))
	if err != nil {
		return RegistryMessage{}, err
	}
	msg.Signature = signature
	return msg, nil
}

func VerifyRegistryMessage(msg RegistryMessage) error {
	payload := registryPayload(msg)
	if payload.Type != "claim" && payload.Type != "release" {
		return fmt.Errorf("unsupported registry message type: %s", payload.Type)
	}
	if payload.Username == "" || payload.KeyID == "" || payload.PublicKey == "" || payload.Timestamp == "" || msg.Signature == "" {
		return fmt.Errorf("registry message is missing required fields")
	}
	if _, err := time.Parse(time.RFC3339, payload.Timestamp); err != nil {
		return fmt.Errorf("invalid registry timestamp: %w", err)
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(payload.PublicKey)
	if err != nil {
		return fmt.Errorf("invalid registry public key encoding: %w", err)
	}
	if Sha256Bytes(publicKeyBytes) != payload.KeyID {
		return fmt.Errorf("registry key id does not match public key")
	}

	publicKey, err := PublicKeyFromBytes(publicKeyBytes)
	if err != nil {
		return fmt.Errorf("invalid registry public key: %w", err)
	}
	valid, err := VerifyJson(publicKey, payload, msg.Signature)
	if err != nil {
		return fmt.Errorf("invalid registry signature encoding: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid registry signature")
	}
	return nil
}

func ProcessRegistryMessage(data []byte) {
	var msg RegistryMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	username := strings.ToLower(msg.Username)
	msg.Username = username
	if err := VerifyRegistryMessage(msg); err != nil {
		log.Printf("[Registry] Rejected %s for '%s': %v", msg.Type, username, err)
		return
	}

	file := PathInAether("usernames", SafeName(username)+".json")
	var reg Registry

	// Check existing registry file
	exists := true
	if err := ReadJson(file, &reg); err != nil {
		exists = false
	}

	if msg.Type == "claim" {
		// First-Seen Rule
		if exists && len(reg.Epochs) > 0 {
			// Check if currently active
			now := time.Now().UTC().Format(time.RFC3339)
			isActive := false
			for _, e := range reg.Epochs {
				if e.ValidFrom <= now && (e.ValidUntil == "" || e.ValidUntil >= now) {
					isActive = true
					break
				}
			}
			if isActive {
				// We already have an active claim, reject this new claim
				log.Printf("[Registry] Rejected claim for '%s': already active", username)
				return
			}
		}

		newEpoch := Epoch{
			BoundIdentity: msg.KeyID,
			ValidFrom:     msg.Timestamp,
		}

		if !exists {
			reg = Registry{
				PrimaryKeyID: msg.KeyID,
				CreatedAt:    msg.Timestamp,
				Epochs:       []Epoch{newEpoch},
			}
		} else {
			reg.Epochs = append(reg.Epochs, newEpoch)
		}

		WriteJson(file, reg)
		log.Printf("[Registry] Accepted claim for '%s' by '%s'", username, msg.KeyID)

	} else if msg.Type == "release" {
		if !exists || len(reg.Epochs) == 0 {
			return
		}

		// Find active epoch
		now := time.Now().UTC().Format(time.RFC3339)
		for i, e := range reg.Epochs {
			if e.ValidFrom <= now && (e.ValidUntil == "" || e.ValidUntil >= now) {
				if e.BoundIdentity == msg.KeyID {
					// Verify this release is requested by the current owner
					reg.Epochs[i].ValidUntil = msg.Timestamp
					WriteJson(file, reg)
					log.Printf("[Registry] Released claim for '%s' by '%s'", username, msg.KeyID)
					return
				}
			}
		}
	}
}

// Ensure the local node has its own identity registered correctly in the global format
func InitializeLocalRegistryClaim(username string, keyId string) error {
	username = strings.ToLower(username)
	file := PathInAether("usernames", SafeName(username)+".json")

	now := time.Now().UTC().Format(time.RFC3339)
	reg := Registry{
		PrimaryKeyID: keyId,
		CreatedAt:    now,
		Epochs: []Epoch{
			{
				BoundIdentity: keyId,
				ValidFrom:     now,
			},
		},
	}

	EnsureDir(filepath.Dir(file))
	err := WriteJson(file, reg)
	if err != nil {
		return err
	}

	msg, err := CreateSignedRegistryMessage("claim", username, keyId, now)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(msg)
	BroadcastRegistryMessage(b)

	return nil
}

func BroadcastRegistryMessage(data []byte) {
	if GlobalPubSub != nil {
		// Compatibility path for the legacy global topic.
		BroadcastMessage(map[string]interface{}{
			"type": "REGISTRY_SYNC",
			"data": string(data),
		})
	}
}
