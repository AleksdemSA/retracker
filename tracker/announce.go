package tracker

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/zeebo/bencode"
)

func (t *Tracker) AnnounceHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	infoHash := query.Get("info_hash")
	peerPortStr := query.Get("port")
	event := query.Get("event")

	if infoHash == "" || peerPortStr == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	peerIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Invalid remote address", http.StatusBadRequest)
		return
	}

	peerPort, err := strconv.Atoi(peerPortStr)
	if err != nil {
		http.Error(w, "Invalid port", http.StatusBadRequest)
		return
	}

	t.mu.Lock()
	torrent, exists := t.Torrents[infoHash]
	if !exists {
		torrent = &TorrentData{Peers: make(map[string]Peer)}
		t.Torrents[infoHash] = torrent
	}
	t.mu.Unlock()

	torrent.mu.Lock()
	defer torrent.mu.Unlock()

	peerKey := fmt.Sprintf("%s:%d", peerIP, peerPort)

	switch event {
	case "stopped":
		delete(torrent.Peers, peerKey)
		log.Printf("Peer disconnected: %s (info_hash: %s)", peerKey, infoHash)
	case "completed":
		torrent.Completed++
		fallthrough
	case "started":
		torrent.Seeders++
		log.Printf("Peer connected: %s (info_hash: %s)", peerKey, infoHash)
	default:
		torrent.Leechers++
	}

	torrent.Peers[peerKey] = Peer{
		IP:       peerIP,
		Port:     peerPort,
		LastSeen: time.Now(),
	}

	peers := make([]byte, 0)
	for _, peer := range torrent.Peers {
		ipBytes := net.ParseIP(peer.IP).To4()
		if ipBytes != nil {
			peers = append(peers, ipBytes...)
			peers = append(peers, byte(peer.Port>>8), byte(peer.Port&0xFF))
		}
	}

	response := map[string]interface{}{
		"interval": 1800,
		"peers":    peers,
	}

	w.Header().Set("Content-Type", "text/plain")
	err = bencode.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
