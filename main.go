package main

import (
	"log"
	"net/http"
	"retracker/tracker"
	"time"
)

func main() {
	tracker := tracker.NewTracker()

	go tracker.CleanupPeers()

	http.HandleFunc("/announce", tracker.AnnounceHandler)
	http.HandleFunc("/scrape", tracker.ScrapeHandler)
	http.HandleFunc("/", tracker.RootHandler)

	server := &http.Server{
		Addr:         ":80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Tracker running on port 80")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
