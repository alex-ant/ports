package api

import (
	"net/http"
	"time"
)

func (a *API) fetchPortInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Register start time.
	startTime := time.Now()

	// Fetch port info slice.
	res, resErr := a.portInfoFetcher()
	if resErr != nil {
		respond("", nil, "failed to fetch port info: "+resErr.Error(), http.StatusInternalServerError, startTime, w)
		return
	}

	respond("port-info", res, "ok", http.StatusOK, startTime, w)
}
