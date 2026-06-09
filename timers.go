package main

import (
	"sync"
	"time"
)

type pingTimer struct {
	stop chan struct{}
}

var (
	pingMu     sync.Mutex
	playerPing = map[int]*pingTimer{}
)

func startPlayerPing(playerID int) bool {
	pingMu.Lock()
	defer pingMu.Unlock()
	if _, ok := playerPing[playerID]; ok {
		return false
	}
	stop := make(chan struct{})
	playerPing[playerID] = &pingTimer{stop: stop}
	go func(id int, done <-chan struct{}) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if !bridgeIsConnected(id) {
					stopPlayerPing(id)
					return
				}
				bridgeSendClientMessage(id, ColourTimer, "PING!")
			}
		}
	}(playerID, stop)
	return true
}

func stopPlayerPing(playerID int) bool {
	pingMu.Lock()
	defer pingMu.Unlock()
	t, ok := playerPing[playerID]
	if !ok {
		return false
	}
	close(t.stop)
	delete(playerPing, playerID)
	return true
}

func clearPlayerPing(playerID int) {
	pingMu.Lock()
	defer pingMu.Unlock()
	if t, ok := playerPing[playerID]; ok {
		close(t.stop)
		delete(playerPing, playerID)
	}
}
