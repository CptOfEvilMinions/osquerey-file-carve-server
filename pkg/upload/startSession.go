package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	guuid "github.com/google/uuid"
)

// Generate a UUID 4 as a sessionID
func generateSessionID() string {
	return fmt.Sprintf("%s", guuid.New())
}

type fileCarve struct {
	BlockCount int    `json:"block_count"`
	BlockSize  int    `json:"block_size"`
	CarveSize  int    `json:"carve_size"`
	CarveID    string `json:"carve_id"`
	RequestID  string `json:"request_id"`
	NodeKey    string `json:"node_key"`
}

// Osquery client calls this function to receive a session ID
// for the file upload. This function generates a UUID4 for a
// session ID, adds sesssion ID to global map and returns it
// to the client.
func StartFileCarve(w http.ResponseWriter, r *http.Request) {
	// Lock map
	Mutex.Lock()

	// Declare a new StartFileCarve obj
	var startFileCarve fileCarve

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&startFileCarve)
	if err != nil {
		log.Println("startSession - StartFileCarve - 1 -", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	if r.Header.Get("X-FORWARDED-FOR") != "" {
		log.Printf("Host %s starting a file upload for Mongo File GUID: %s\n", r.Header.Get("X-FORWARDED-FOR"), startFileCarve.CarveID)
	} else {
		log.Printf("Host %s starting a file upload for Mongo File GUID: %s\n", r.RemoteAddr, startFileCarve.CarveID)
	}

	// Generate sessionID
	sessionID := generateSessionID()

	// Add sessionID to map
	FileCarveSessionMap[sessionID] = &FilCarveSession{
		Timestamp:        time.Now(),
		totalBlocks:      startFileCarve.BlockCount,
		CarveID:          startFileCarve.CarveID,
		ReceivedBlockIDs: []int{},
	}

	fmt.Println(FileCarveSessionMap)

	// Unlock map
	Mutex.Unlock()

	// Create map for JSON and set vaule
	resp := map[string]string{"session_id": sessionID}

	// Marshal map into JSON
	// Return 404 if JOSN can't be marshalled
	js, err := json.Marshal(resp)
	if err != nil {
		log.Println("startSession - StartFileCarve - 2 -", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	// Return session ID
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
