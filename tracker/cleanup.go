package tracker

import (
	"log"
	"time"
)

func (t *Tracker) CleanupPeers() {
	for {
		time.Sleep(10 * time.Minute)

		t.mu.Lock()
		for infoHash, torrent := range t.Torrents {
			torrent.mu.Lock()
			for peerKey, peer := range torrent.Peers {
				if time.Since(peer.LastSeen) > 10*time.Minute {
					log.Printf("Peer timed out: %s (info_hash: %s)", peerKey, infoHash)
					delete(torrent.Peers, peerKey)
				}
			}
			if len(torrent.Peers) == 0 {
				log.Printf("No peers left, removing torrent: %s", infoHash)
				delete(t.Torrents, infoHash)
			}
			torrent.mu.Unlock()
		}
		t.mu.Unlock()
	}
}
