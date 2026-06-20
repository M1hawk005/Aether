package daemon

import (
	"sync"
)

var (
	reputationMutex sync.RWMutex
)

// RecordTransfer increments the bytes delivered for a specific peer.
func RecordTransfer(peerId string, bytes int) {
	reputationMutex.Lock()
	defer reputationMutex.Unlock()

	file := PathInAether("reputation.json")
	var rep map[string]int
	if err := ReadJson(file, &rep); err != nil || rep == nil {
		rep = make(map[string]int)
	}

	rep[peerId] += bytes
	WriteJson(file, rep)
}

// GetReputation returns the recorded bytes delivered by a peer.
func GetReputation(peerId string) int {
	reputationMutex.RLock()
	defer reputationMutex.RUnlock()

	file := PathInAether("reputation.json")
	var rep map[string]int
	if err := ReadJson(file, &rep); err != nil || rep == nil {
		return 0
	}

	return rep[peerId]
}

// GetAllReputation returns the full local scoreboard.
func GetAllReputation() map[string]int {
	reputationMutex.RLock()
	defer reputationMutex.RUnlock()

	file := PathInAether("reputation.json")
	var rep map[string]int
	if err := ReadJson(file, &rep); err != nil || rep == nil {
		return make(map[string]int)
	}

	return rep
}
