package tracker

import (
	"fmt"
	"net/http"
)

func (t *Tracker) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Open Torrent Tracker (add this retracker)")
}
