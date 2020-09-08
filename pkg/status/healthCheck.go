package status

import (
	"encoding/json"
	"net/http"
)

// HealthCheck is an API endpoint to query to see if this
// seervie is still alive
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Create map for JSON and set vaule
	resp := map[string]bool{"success": true}

	// Marshal map into JSON
	// Return 404 if JOSN can't be marshalled
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return sucess to client
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
