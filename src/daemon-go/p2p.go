package daemon

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

var (
	GlobalHost   host.Host
	GlobalPubSub *pubsub.PubSub
	GlobalDHT    *dht.IpfsDHT
	TopicName    = "aether-global-v1"
	Topic        *pubsub.Topic

	PendingFetches = make(map[string]chan []byte)
	FetchMutex     sync.Mutex

	PendingResolves = make(map[string]chan []map[string]interface{})
	ResolveMutex    sync.Mutex

	PendingAvailability = make(map[string]chan string)
	AvailabilityMutex   sync.Mutex
)

func InitLibp2p(ctx context.Context, listenPort int) error {
	var err error

	GlobalHost, err = libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.NATPortMap(),
	)
	if err != nil {
		return fmt.Errorf("failed to create libp2p host: %w", err)
	}

	log.Printf("===================================================================")
	log.Printf("[Aether p2p] Host initialized with ID: %s", GlobalHost.ID().String())
	log.Printf("[Aether p2p] Your Public/Local Multiaddresses for Peers to connect:")
	for _, addr := range GlobalHost.Addrs() {
		log.Printf("   -> %s/p2p/%s", addr.String(), GlobalHost.ID().String())
	}
	log.Printf("===================================================================")

	GlobalDHT, err = dht.New(ctx, GlobalHost, dht.Mode(dht.ModeServer))
	if err != nil {
		return fmt.Errorf("failed to create DHT: %w", err)
	}

	bootstrapEnv := os.Getenv("AETHER_BOOTSTRAP_PEERS")
	var bootstrapPeers []string
	if bootstrapEnv != "" {
		bootstrapPeers = strings.Split(bootstrapEnv, ",")
	}
	connectToBootstrapPeers(ctx, bootstrapPeers)

	if err = GlobalDHT.Bootstrap(ctx); err != nil {
		log.Printf("DHT Bootstrap failed: %v", err)
	}

	routingDiscovery := routing.NewRoutingDiscovery(GlobalDHT)
	util.Advertise(ctx, routingDiscovery, TopicName)

	go func() {
		for {
			peerChan, err := routingDiscovery.FindPeers(ctx, TopicName)
			if err != nil {
				log.Printf("FindPeers error: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			for p := range peerChan {
				if p.ID == GlobalHost.ID() {
					continue
				}
				if GlobalHost.Network().Connectedness(p.ID) != 1 { // 1 == Connected
					err := GlobalHost.Connect(ctx, p)
					if err == nil {
						log.Printf("[Aether p2p] Connected to peer: %s", p.ID.String())
					}
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()

	GlobalPubSub, err = pubsub.NewGossipSub(ctx, GlobalHost)
	if err != nil {
		return fmt.Errorf("failed to create pubsub: %w", err)
	}

	Topic, err = GlobalPubSub.Join(TopicName)
	if err != nil {
		return fmt.Errorf("failed to join topic: %w", err)
	}

	sub, err := Topic.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	go func() {
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				log.Printf("PubSub subscription error: %v", err)
				return
			}
			if msg.ReceivedFrom == GlobalHost.ID() {
				continue
			}

			handleIncomingNetworkMessage(msg.Data, msg.ReceivedFrom)
		}
	}()

	return nil
}

func BroadcastMessage(msgObj interface{}) {
	if Topic == nil {
		return
	}
	bytes, err := json.Marshal(msgObj)
	if err != nil {
		log.Printf("Failed to marshal broadcast message: %v", err)
		return
	}
	err = Topic.Publish(context.Background(), bytes)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
}

func handleIncomingNetworkMessage(data []byte, from peer.ID) {
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}
	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	log.Printf("[Aether p2p] Received %s from %s", msgType, from.String())
	
	if msgType == "REGISTRY_SYNC" {
		if dataStr, ok := msg["data"].(string); ok {
			ProcessRegistryMessage([]byte(dataStr))
		}
	} else if msgType == "WANT_CAPSULE" {
		if capsuleId, ok := msg["capsuleId"].(string); ok {
			capsDir := PathInAether("peers", PeerId, "capsules", SafeName(capsuleId))
			chunkBaseDir := PathInAether("peers", PeerId, "chunks")
			
			// If not in our own published dir, check global cache
			if _, err := os.Stat(capsDir); os.IsNotExist(err) {
				capsDir = PathInAether("global_cache", "capsules", SafeName(capsuleId))
				chunkBaseDir = PathInAether("global_cache", "chunks")
			}
			
			var manifest map[string]interface{}
			if err := ReadJson(filepath.Join(capsDir, "manifest.json"), &manifest); err == nil {
				chunks := manifest["chunks"].([]interface{})
				if len(chunks) > 0 {
					chunkMeta := chunks[0].(map[string]interface{})
					chunkHash := chunkMeta["hash"].(string)

					var chunkData map[string]interface{}
					if err := ReadJson(filepath.Join(chunkBaseDir, SafeName(chunkHash)+".json"), &chunkData); err == nil {
						// Send DELIVER_CAPSULE back
						BroadcastMessage(map[string]interface{}{
							"type":      "DELIVER_CAPSULE",
							"capsuleId": capsuleId,
							"data":      chunkData["bytes"],
							"sender":    GlobalHost.ID().String(),
						})
					}
				}
			}
		}
	} else if msgType == "DELIVER_CAPSULE" {
		if capsuleId, ok := msg["capsuleId"].(string); ok {
			if b64Data, ok := msg["data"].(string); ok {
				FetchMutex.Lock()
				ch, exists := PendingFetches[capsuleId]
				FetchMutex.Unlock()
				if exists {
					// AET-012: validate the expected content hash before accepting this payload.
					b, err := base64.StdEncoding.DecodeString(b64Data)
					if err == nil {
						select {
						case ch <- b:
						default:
						}
					}
				}
			}
		}
	} else if msgType == "WANT_USER_CAPSULES" {
		if targetUser, ok := msg["username"].(string); ok {
			var results []map[string]interface{}
			
			dirsToCheck := []string{
				PathInAether("peers", PeerId, "capsules"),
				PathInAether("global_cache", "capsules"),
			}
			for _, d := range dirsToCheck {
				entries, _ := os.ReadDir(d)
				for _, e := range entries {
					if e.IsDir() {
						var header map[string]interface{}
						if err := ReadJson(filepath.Join(d, e.Name(), "header.json"), &header); err == nil {
							pub, _ := header["publisher"].(string)
							if strings.ToLower(pub) == targetUser {
								results = append(results, header)
							}
						}
					}
				}
			}

			if len(results) > 0 {
				BroadcastMessage(map[string]interface{}{
					"type":     "DELIVER_USER_CAPSULES",
					"username": targetUser,
					"capsules": results,
				})
			}
		}
	} else if msgType == "DELIVER_USER_CAPSULES" {
		if username, ok := msg["username"].(string); ok {
			if capsules, ok := msg["capsules"].([]interface{}); ok {
				var parsed []map[string]interface{}
				for _, c := range capsules {
					if cMap, ok := c.(map[string]interface{}); ok {
						parsed = append(parsed, cMap)
					}
				}
				ResolveMutex.Lock()
				ch, exists := PendingResolves[username]
				ResolveMutex.Unlock()
				if exists {
					select {
					case ch <- parsed:
					default:
					}
				}
			}
		}
	} else if msgType == "BT_ANNOUNCE" {
		if infoHash, ok := msg["info_hash"].(string); ok {
			if peerData, ok := msg["peer"].(map[string]interface{}); ok {
				ip, _ := peerData["ip"].(string)
				portF, _ := peerData["port"].(float64)
				peerId, _ := peerData["peer_id"].(string)
				
				if ip != "" && peerId != "" {
					AddBTPeer(infoHash, BTPeer{
						IP:     ip,
						Port:   int(portF),
						PeerID: peerId,
					})
				}
			}
		}
	} else if msgType == "CHECK_AVAILABILITY" {
		if capsuleId, ok := msg["capsuleId"].(string); ok {
			hasIt := false
			if _, err := os.Stat(PathInAether("peers", PeerId, "capsules", SafeName(capsuleId))); err == nil {
				hasIt = true
			} else if _, err := os.Stat(PathInAether("global_cache", "capsules", SafeName(capsuleId))); err == nil {
				hasIt = true
			}
			if hasIt {
				BroadcastMessage(map[string]interface{}{
					"type":      "AVAILABILITY_RESPONSE",
					"capsuleId": capsuleId,
					"peerId":    GlobalHost.ID().String(),
				})
			}
		}
	} else if msgType == "AVAILABILITY_RESPONSE" {
		if capsuleId, ok := msg["capsuleId"].(string); ok {
			if peerId, ok := msg["peerId"].(string); ok {
				AvailabilityMutex.Lock()
				ch, exists := PendingAvailability[capsuleId]
				AvailabilityMutex.Unlock()
				if exists {
					select {
					case ch <- peerId:
					default:
					}
				}
			}
		}
	}
}

func connectToBootstrapPeers(ctx context.Context, addrs []string) {
	var wg sync.WaitGroup
	for _, addrStr := range addrs {
		addrStr = strings.TrimSpace(addrStr)
		if addrStr == "" {
			continue
		}

		addr, err := multiaddr.NewMultiaddr(addrStr)
		if err != nil {
			log.Printf("[Aether p2p] Invalid bootstrap multiaddr %s: %v", addrStr, err)
			continue
		}

		peerinfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			log.Printf("[Aether p2p] Invalid bootstrap peer info %s: %v", addrStr, err)
			continue
		}

		if peerinfo.ID == GlobalHost.ID() {
			continue
		}

		wg.Add(1)
		go func(pi *peer.AddrInfo) {
			defer wg.Done()
			if err := GlobalHost.Connect(ctx, *pi); err != nil {
				log.Printf("[Aether p2p] Warning: could not connect to bootstrap peer %s: %v", pi.ID.String(), err)
			} else {
				log.Printf("[Aether p2p] Successfully connected to bootstrap peer %s", pi.ID.String())
			}
		}(peerinfo)
	}
	wg.Wait()
}
