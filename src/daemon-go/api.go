package daemon

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func SendJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func HandlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendJson(w, 405, map[string]string{"error": "Method Not Allowed"})
		return
	}

	var req struct {
		Publisher  string `json:"publisher"`
		Payload    string `json:"payload"`
		Visibility string `json:"visibility"`
		Persistent bool   `json:"persistent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJson(w, 400, map[string]string{"error": "Invalid JSON"})
		return
	}

	keyId, err := ResolveCurrentKeyId(req.Publisher)
	if err != nil {
		SendJson(w, 404, map[string]string{"error": err.Error()})
		return
	}

	privKeyPath := PathInAether("users", SafeName(keyId), "keys", SafeName(keyId), "private.key")
	privKey, err := LoadPrivateKey(privKeyPath)
	if err != nil {
		SendJson(w, 500, map[string]string{"error": "Failed to load private key"})
		return
	}

	// Compress
	rawBytes := []byte(req.Payload)
	compressed := CompressBrotli(rawBytes)
	chunkHash := Sha256Bytes(compressed)

	// Create Manifest
	manifest := map[string]interface{}{
		"version":      1,
		"created_at":   time.Now().UTC().Format(time.RFC3339),
		"is_directory": false,
		"chunks": []map[string]interface{}{
			{
				"index": 0,
				"hash":  chunkHash,
				"size":  len(compressed),
			},
		},
	}
	manifestHash, _ := Sha256Json(manifest)

	visibility := req.Visibility
	if visibility == "" {
		visibility = "public"
	}

	// Create Envelope
	envelope := map[string]interface{}{
		"version":          1,
		"created_at":       manifest["created_at"],
		"publisher":        req.Publisher,
		"publisher_key_id": keyId,
		"object_hash":      manifestHash,
		"visibility":       visibility,
		"persistent":       req.Persistent,
	}
	envSig, _ := SignJson(privKey, envelope)
	envelope["signature"] = envSig
	envelopeHash, _ := Sha256Json(envelope)

	// Create Header
	header := map[string]interface{}{
		"version":       1,
		"created_at":    manifest["created_at"],
		"publisher":     req.Publisher,
		"owner_key_id":  keyId,
		"visibility":    visibility,
		"manifest_hash": manifestHash,
		"envelope_hash": envelopeHash,
		"capsule_id":    envelopeHash,
	}
	hdrSig, _ := SignJson(privKey, header)
	header["signature"] = hdrSig

	// Write to disk
	capsDir := PathInAether("peers", PeerId, "capsules", SafeName(envelopeHash))
	EnsureDir(capsDir)
	WriteJson(filepath.Join(capsDir, "header.json"), header)
	WriteJson(filepath.Join(capsDir, "envelope.json"), envelope)
	WriteJson(filepath.Join(capsDir, "manifest.json"), manifest)

	// Write chunk
	chunkDir := PathInAether("peers", PeerId, "chunks")
	EnsureDir(chunkDir)
	WriteJson(filepath.Join(chunkDir, SafeName(chunkHash)+".json"), map[string]interface{}{
		"hash":       chunkHash,
		"encoding":   "base64",
		"bytes":      compressed,
		"compressed": "brotli",
	})

	// Broadcast
	BroadcastMessage(map[string]string{
		"type":      "NEW_CAPSULE",
		"capsuleId": envelopeHash,
	})

	SendJson(w, 200, map[string]interface{}{
		"ok":         true,
		"capsule_id": envelopeHash,
	})
}

func HandleResolve(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/api/resolve/")
	username = strings.TrimSpace(strings.ToLower(username))

	var localResults []map[string]interface{}

	dirsToCheck := []string{
		PathInAether("peers", PeerId, "capsules"),
		PathInAether("global_cache", "capsules"),
	}

	for _, d := range dirsToCheck {
		entries, err := os.ReadDir(d)
		if err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				headerFile := filepath.Join(d, e.Name(), "header.json")
				var header map[string]interface{}
				if err := ReadJson(headerFile, &header); err == nil {
					pub, _ := header["publisher"].(string)
					if strings.ToLower(pub) == username {
						localResults = append(localResults, header)
					}
				}
			}
		}
	}

	// Now ask the network
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	ch := make(chan []map[string]interface{}, 50)
	ResolveMutex.Lock()
	PendingResolves[username] = ch
	ResolveMutex.Unlock()

	defer func() {
		ResolveMutex.Lock()
		delete(PendingResolves, username)
		ResolveMutex.Unlock()
	}()

	BroadcastMessage(map[string]interface{}{
		"type":     "WANT_USER_CAPSULES",
		"username": username,
	})

	// Collect responses
	aggregated := make(map[string]map[string]interface{})
	for _, r := range localResults {
		if sig, ok := r["signature"].(string); ok {
			aggregated[sig] = r
		}
	}

	for {
		select {
		case <-ctx.Done():
			// Timeout reached, return aggregated
			var finalResults []map[string]interface{}
			for _, v := range aggregated {
				finalResults = append(finalResults, v)
			}

			if len(finalResults) == 0 {
				SendJson(w, 404, map[string]string{"error": "No capsules found for " + username})
				return
			}

			SendJson(w, 200, map[string]interface{}{
				"ok":       true,
				"username": username,
				"capsules": finalResults,
			})
			return
		case networkResults := <-ch:
			for _, r := range networkResults {
				if sig, ok := r["signature"].(string); ok {
					aggregated[sig] = r
				}
			}
		}
	}
}

func HandleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	capsuleId := strings.TrimPrefix(r.URL.Path, "/api/fetch/")
	capsDir := PathInAether("peers", PeerId, "capsules", SafeName(capsuleId))

	var manifest map[string]interface{}
	if err := ReadJson(filepath.Join(capsDir, "manifest.json"), &manifest); err == nil {
		chunks := manifest["chunks"].([]interface{})
		if len(chunks) > 0 {
			chunkMeta := chunks[0].(map[string]interface{})
			chunkHash := chunkMeta["hash"].(string)

			var chunkData map[string]interface{}
			if err := ReadJson(PathInAether("peers", PeerId, "chunks", SafeName(chunkHash)+".json"), &chunkData); err == nil {
				b64Str := chunkData["bytes"].(string)
				// Decode base64
				compressed, _ := base64.StdEncoding.DecodeString(b64Str)
				// Decompress brotli
				b, _ := DecompressBrotli(compressed)
				SendJson(w, 200, map[string]interface{}{
					"ok":   true,
					"data": string(b),
				})
				return
			}
		}
	}
	// If not found locally, broadcast a P2P request
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ch := make(chan []byte, 1)
	FetchMutex.Lock()
	PendingFetches[capsuleId] = ch
	FetchMutex.Unlock()

	defer func() {
		FetchMutex.Lock()
		delete(PendingFetches, capsuleId)
		FetchMutex.Unlock()
	}()

	BroadcastMessage(map[string]interface{}{
		"type":      "WANT_CAPSULE",
		"capsuleId": capsuleId,
	})

	select {
	case <-ctx.Done():
		SendJson(w, 404, map[string]string{"error": "Not found locally and P2P fetch timed out"})
	case b := <-ch:
		decompressed, err := DecompressBrotli(b)
		if err != nil {
			SendJson(w, 500, map[string]string{"error": "Failed to decompress P2P chunk"})
			return
		}

		// Retain the fetched chunk for the legacy seeding path.
		capsDir := PathInAether("global_cache", "capsules", SafeName(capsuleId))
		EnsureDir(capsDir)
		chunkHash := capsuleId + "_chunk"

		// The legacy responder expects a manifest containing the cached chunk key.
		WriteJson(filepath.Join(capsDir, "manifest.json"), map[string]interface{}{
			"chunks": []map[string]string{{"hash": chunkHash}},
		})

		chunkDir := PathInAether("global_cache", "chunks")
		EnsureDir(chunkDir)
		WriteJson(filepath.Join(chunkDir, SafeName(chunkHash)+".json"), map[string]interface{}{
			"bytes": base64.StdEncoding.EncodeToString(b),
		})

		SendJson(w, 200, map[string]interface{}{
			"ok":   true,
			"data": string(decompressed),
		})
	}
}

func HandleAvailability(w http.ResponseWriter, r *http.Request) {
	capsuleId := strings.TrimPrefix(r.URL.Path, "/api/availability/")
	seeds := 0
	peersDir := PathInAether("peers")
	entries, err := os.ReadDir(peersDir)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				headerFile := filepath.Join(peersDir, e.Name(), "capsules", SafeName(capsuleId), "header.json")
				if _, err := os.Stat(headerFile); err == nil {
					seeds++
				}
			}
		}
	}

	// Add Global P2P Seeds
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	ch := make(chan string, 50)
	AvailabilityMutex.Lock()
	PendingAvailability[capsuleId] = ch
	AvailabilityMutex.Unlock()

	defer func() {
		AvailabilityMutex.Lock()
		delete(PendingAvailability, capsuleId)
		AvailabilityMutex.Unlock()
	}()

	BroadcastMessage(map[string]interface{}{
		"type":      "CHECK_AVAILABILITY",
		"capsuleId": capsuleId,
	})

	seenPeers := make(map[string]bool)
	for {
		select {
		case <-ctx.Done():
			SendJson(w, 200, map[string]interface{}{
				"ok":         true,
				"capsule_id": capsuleId,
				"seeds":      seeds + len(seenPeers),
			})
			return
		case peerId := <-ch:
			seenPeers[peerId] = true
		}
	}
}

func HandleDeleteLocal(w http.ResponseWriter, r *http.Request) {
	capsuleId := strings.TrimPrefix(r.URL.Path, "/api/delete/")
	capsDir := PathInAether("peers", PeerId, "capsules", SafeName(capsuleId))
	if _, err := os.Stat(capsDir); err == nil {
		os.RemoveAll(capsDir)
		SendJson(w, 200, map[string]interface{}{"ok": true, "message": "Deleted locally"})
	} else {
		SendJson(w, 404, map[string]string{"error": "Not found locally"})
	}
}

func HandleRevokeGlobal(w http.ResponseWriter, r *http.Request) {
	capsuleId := strings.TrimPrefix(r.URL.Path, "/api/revoke/")
	capsDir := PathInAether("peers", PeerId, "capsules", SafeName(capsuleId))
	os.RemoveAll(capsDir)

	BroadcastMessage(map[string]interface{}{
		"type":      "REVOKE_CAPSULE",
		"capsuleId": capsuleId,
	})
	SendJson(w, 200, map[string]interface{}{"ok": true, "message": "Revocation broadcasted"})
}

func HandleRegistryClaim(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendJson(w, 405, map[string]string{"error": "Method Not Allowed"})
		return
	}
	var req struct {
		Username string `json:"username"`
		KeyId    string `json:"key_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJson(w, 400, map[string]string{"error": "Invalid JSON"})
		return
	}
	err := InitializeLocalRegistryClaim(req.Username, req.KeyId)
	if err != nil {
		SendJson(w, 500, map[string]string{"error": err.Error()})
		return
	}
	SendJson(w, 200, map[string]interface{}{"ok": true, "message": "Claim broadcasted"})
}

func HandleRegistryRelease(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendJson(w, 405, map[string]string{"error": "Method Not Allowed"})
		return
	}
	var req struct {
		Username string `json:"username"`
		KeyId    string `json:"key_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJson(w, 400, map[string]string{"error": "Invalid JSON"})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	msg, err := CreateSignedRegistryMessage("release", req.Username, req.KeyId, now)
	if err != nil {
		SendJson(w, 500, map[string]string{"error": err.Error()})
		return
	}
	b, _ := json.Marshal(msg)

	// Apply locally first
	ProcessRegistryMessage(b)

	// Broadcast to network
	BroadcastRegistryMessage(b)

	SendJson(w, 200, map[string]interface{}{"ok": true, "message": "Release broadcasted"})
}

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendJson(w, 405, map[string]string{"error": "Method Not Allowed"})
		return
	}

	peerCount := 0
	if GlobalHost != nil {
		peerCount = len(GlobalHost.Network().Peers())
	}

	SendJson(w, 200, map[string]interface{}{
		"ok":      true,
		"peers":   peerCount,
		"peer_id": PeerId,
	})
}

func HandleReputation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendJson(w, 405, map[string]string{"error": "Method Not Allowed"})
		return
	}

	rep := GetAllReputation()

	SendJson(w, 200, map[string]interface{}{
		"ok":         true,
		"reputation": rep,
	})
}

func HandleAnnounce(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	q := r.URL.Query()
	infoHashBytes := []byte(q.Get("info_hash"))
	if len(infoHashBytes) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("Missing info_hash"))
		return
	}
	infoHashHex := fmt.Sprintf("%x", infoHashBytes)

	peerIdStr := q.Get("peer_id")
	portStr := q.Get("port")

	port, _ := strconv.Atoi(portStr)

	ip := "127.0.0.1"
	if GlobalHost != nil {
		for _, addr := range GlobalHost.Addrs() {
			addrStr := addr.String()
			if strings.Contains(addrStr, "127.0.0.1") || strings.Contains(addrStr, "::1") {
				continue
			}
			parts := strings.Split(addrStr, "/")
			if len(parts) >= 3 && parts[1] == "ip4" {
				ip = parts[2]
				break
			}
		}
	}

	if peerIdStr != "" && port > 0 {
		me := BTPeer{
			IP:     ip,
			Port:   port,
			PeerID: peerIdStr,
		}
		AddBTPeer(infoHashHex, me)

		msg := map[string]interface{}{
			"type":      "BT_ANNOUNCE",
			"info_hash": infoHashHex,
			"peer":      me,
		}
		BroadcastMessage(msg)
	}

	peers := GetBTPeers(infoHashHex)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)

	var sb strings.Builder
	sb.WriteString("d8:intervali1800e5:peersl")
	for _, p := range peers {
		sb.WriteString("d")
		sb.WriteString(fmt.Sprintf("2:ip%d:%s", len(p.IP), p.IP))

		pid := p.PeerID
		if len(pid) > 20 {
			pid = pid[:20]
		}
		sb.WriteString(fmt.Sprintf("7:peer_id%d:%s", len(pid), pid))
		sb.WriteString(fmt.Sprintf("4:porti%de", p.Port))
		sb.WriteString("e")
	}
	sb.WriteString("ee")

	w.Write([]byte(sb.String()))
}

func HandleFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	var results []map[string]interface{}

	dirsToCheck := []string{
		PathInAether("peers", PeerId, "capsules"),
		PathInAether("global_cache", "capsules"),
	}

	seen := make(map[string]bool)

	for _, d := range dirsToCheck {
		entries, err := os.ReadDir(d)
		if err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				headerFile := filepath.Join(d, e.Name(), "header.json")
				var header map[string]interface{}
				if err := ReadJson(headerFile, &header); err == nil {
					hashStr, ok := header["capsule_id"].(string)
					if ok && !seen[hashStr] {
						seen[hashStr] = true
						results = append(results, header)
					}
				}
			}
		}
	}

	SendJson(w, 200, map[string]interface{}{
		"ok":       true,
		"capsules": results,
	})
}
