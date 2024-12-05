package tracker

import (
	"sync"
	"time"
)

type Peer struct {
	IP       string
	Port     int
	LastSeen time.Time
}

type TorrentData struct {
	Peers     map[string]Peer
	Seeders   int
	Leechers  int
	Completed int
	mu        sync.Mutex
}

type Tracker struct {
	Torrents map[string]*TorrentData
	mu       sync.Mutex
}

func NewTracker() *Tracker {
	return &Tracker{
		Torrents: make(map[string]*TorrentData),
	}
}
