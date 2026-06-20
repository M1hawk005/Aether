package daemon

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
)

func TestSha256Bytes(t *testing.T) {
	hash := Sha256Bytes([]byte("hello world"))
	expected := "sha256:b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	if hash != expected {
		t.Errorf("Expected %s, got %s", expected, hash)
	}
}

func TestSignAndVerify(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	payload := map[string]string{"message": "test payload"}
	
	sigBase64, err := SignJson(priv, payload)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	valid, err := VerifyJson(pub, payload, sigBase64)
	if err != nil || !valid {
		t.Fatalf("Failed to verify signature: %v", err)
	}
}
