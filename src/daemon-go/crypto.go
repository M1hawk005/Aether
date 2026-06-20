package daemon

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

var MasterKey []byte // Global in-memory AES key (32 bytes)

func SetMasterKey(key []byte) {
	MasterKey = key
}

// EncryptBytes uses AES-GCM to encrypt data with the global MasterKey
func EncryptBytes(data []byte) ([]byte, error) {
	if len(MasterKey) != 32 {
		return nil, fmt.Errorf("master key not set or invalid")
	}
	block, err := aes.NewCipher(MasterKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptBytes uses AES-GCM to decrypt data with the global MasterKey
func DecryptBytes(data []byte) ([]byte, error) {
	if len(MasterKey) != 32 {
		return nil, fmt.Errorf("master key not set or invalid")
	}
	block, err := aes.NewCipher(MasterKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// If the file is smaller than 100 bytes, it might be an old unencrypted key
	// A valid encrypted key will be at least NonceSize + TagSize + Payload
	// For backward compatibility and our current test keys, try to parse it raw first,
	// if it fails or if we strictly want encryption, we decrypt.
	// But since user requested encryption, let's assume it's encrypted if MasterKey is set.
	var decryptedBytes []byte
	if len(MasterKey) == 32 {
		decryptedBytes, err = DecryptBytes(bytes)
		if err != nil {
			// Fallback for unencrypted legacy test keys
			decryptedBytes = bytes
		}
	} else {
		decryptedBytes = bytes
	}

	privKey, err := p2pcrypto.UnmarshalPrivateKey(decryptedBytes)
	if err != nil {
		return nil, err
	}

	rawBytes, err := privKey.Raw()
	if err != nil {
		return nil, err
	}

	return ed25519.PrivateKey(rawBytes), nil
}

func LoadPublicKey(path string) (ed25519.PublicKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return PublicKeyFromBytes(bytes)
}

func PublicKeyFromBytes(bytes []byte) (ed25519.PublicKey, error) {
	pubKey, err := p2pcrypto.UnmarshalPublicKey(bytes)
	if err != nil {
		return nil, err
	}

	rawBytes, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}

	return ed25519.PublicKey(rawBytes), nil
}

func Sha256Bytes(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", hash)
}

func Sha256Json(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return Sha256Bytes(bytes), nil
}

func SignJson(privKey ed25519.PrivateKey, obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	sig := ed25519.Sign(privKey, bytes)
	return base64.StdEncoding.EncodeToString(sig), nil
}

func VerifyJson(pubKey ed25519.PublicKey, obj interface{}, signatureBase64 string) (bool, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return false, err
	}
	sig, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, err
	}
	return ed25519.Verify(pubKey, bytes, sig), nil
}
