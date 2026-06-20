package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/andybalholm/brotli"
)

var aetherDir = ".aether"

func PathInAether(parts ...string) string {
	wd, _ := os.Getwd()
	allParts := append([]string{wd, aetherDir}, parts...)
	return filepath.Join(allParts...)
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func ReadJson(file string, target interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func ReadFileBytes(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func WriteJson(file string, data interface{}) error {
	EnsureDir(filepath.Dir(file))
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, bytes, 0644)
}

func SafeName(val string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9.-]`)
	return re.ReplaceAllString(val, "_")
}

func CompressBrotli(data []byte) []byte {
	var b bytes.Buffer
	w := brotli.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

func DecompressBrotli(data []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(data))
	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	return b.Bytes(), err
}

type NetworkConfig struct {
	Peers map[string]string `json:"peers"`
}

func LoadNetworkConfig() error {
	netFile := PathInAether("network.json")
	var net NetworkConfig
	if err := ReadJson(netFile, &net); err != nil {
		net.Peers = make(map[string]string)
	}
	if net.Peers == nil {
		net.Peers = make(map[string]string)
	}
	net.Peers[PeerId] = fmt.Sprintf("http://127.0.0.1:%d", HttpPort)
	return WriteJson(netFile, net)
}

func CleanupNetworkConfig() {
	netFile := PathInAether("network.json")
	var net NetworkConfig
	if err := ReadJson(netFile, &net); err == nil && net.Peers != nil {
		delete(net.Peers, PeerId)
		WriteJson(netFile, net)
	}
}

// Identity and Epoch logic
type Epoch struct {
	BoundIdentity string `json:"bound_identity"`
	ValidFrom     string `json:"valid_from"`
	ValidUntil    string `json:"valid_until,omitempty"`
}

type Registry struct {
	PrimaryKeyID string  `json:"primary_key_id"`
	CreatedAt    string  `json:"created_at"`
	Epochs       []Epoch `json:"epochs,omitempty"`
}

type UserProfile struct {
	PrimaryKeyID string                 `json:"primary_key_id"`
	Keys         map[string]interface{} `json:"keys"`
}

func ResolveCurrentKeyId(username string) (string, error) {
	file := PathInAether("usernames", username+".json")
	var reg Registry
	if err := ReadJson(file, &reg); err != nil {
		return "", fmt.Errorf("user %s not found in registry", username)
	}

	now := time.Now().Format(time.RFC3339)
	if len(reg.Epochs) > 0 {
		for _, e := range reg.Epochs {
			if e.ValidFrom <= now && (e.ValidUntil == "" || e.ValidUntil >= now) {
				return e.BoundIdentity, nil
			}
		}
	} else {
		return reg.PrimaryKeyID, nil
	}
	return "", fmt.Errorf("no active epoch for %s", username)
}

func LoadUserProfile(keyId string) (UserProfile, error) {
	file := PathInAether("users", SafeName(keyId), "profile.json")
	var p UserProfile
	err := ReadJson(file, &p)
	return p, err
}
