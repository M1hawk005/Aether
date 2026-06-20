package daemon

import (
	"sync"
	"time"
)

type BTPeer struct {
	IP       string    `json:"ip"`
	Port     int       `json:"port"`
	PeerID   string    `json:"peer_id"`
	LastSeen time.Time `json:"-"`
}

var (
	BTTrackerSwarm = make(map[string][]BTPeer)
	BTTrackerMutex sync.Mutex
)

func AddBTPeer(infoHash string, peer BTPeer) {
	BTTrackerMutex.Lock()
	defer BTTrackerMutex.Unlock()
	
	peers := BTTrackerSwarm[infoHash]
	updated := false
	for i, p := range peers {
		if p.PeerID == peer.PeerID {
			peers[i].IP = peer.IP
			peers[i].Port = peer.Port
			peers[i].LastSeen = time.Now()
			updated = true
			break
		}
	}
	if !updated {
		peer.LastSeen = time.Now()
		peers = append(peers, peer)
	}
	BTTrackerSwarm[infoHash] = peers
}

func GetBTPeers(infoHash string) []BTPeer {
	BTTrackerMutex.Lock()
	defer BTTrackerMutex.Unlock()
	
	// Filter out peers older than 30 minutes
	var activePeers []BTPeer
	now := time.Now()
	for _, p := range BTTrackerSwarm[infoHash] {
		if now.Sub(p.LastSeen) < 30*time.Minute {
			activePeers = append(activePeers, p)
		}
	}
	BTTrackerSwarm[infoHash] = activePeers
	return activePeers
}
