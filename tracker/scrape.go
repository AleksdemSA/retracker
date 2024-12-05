package tracker

import (
	"net/http"

	"github.com/zeebo/bencode"
)

func (t *Tracker) ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	infoHashes := query["info_hash"]

	t.mu.Lock()
	defer t.mu.Unlock()

	response := make(map[string]map[string]int)
	for _, infoHash := range infoHashes {
		torrent, exists := t.Torrents[infoHash]
		if !exists {
			continue
		}
		torrent.mu.Lock()
		response[infoHash] = map[string]int{
			"complete":   torrent.Seeders,
			"incomplete": torrent.Leechers,
			"downloaded": torrent.Completed,
		}
		torrent.mu.Unlock()
	}

	w.Header().Set("Content-Type", "text/plain")
	if err := bencode.NewEncoder(w).Encode(map[string]interface{}{
		"files": response,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
