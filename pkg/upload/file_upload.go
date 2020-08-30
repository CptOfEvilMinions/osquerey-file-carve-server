package upload

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// UploadFileCarveToDisk takes in file carve data blocks and processes them
//https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func UploadFileCarveToDisk(w http.ResponseWriter, r *http.Request) {
	fmt.Println("######################################### Uploading Block #########################################")
	// Declare a new FileCarveBlock obj
	var fileCarveBlock FileCarveBlock

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(fileCarveBlock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if Sesssion ID exists, if not return 404
	// https://www.quora.com/In-Go-how-do-I-use-a-map-with-a-string-key-and-a-struct-as-value
	// https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
	if _, exist := FileCarveSessionMap[fileCarveBlock.SessionID]; !exist {
		for k := range FileCarveSessionMap {
			fmt.Println(k)
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// LOCK Mutex
	// Add block to map
	// Set last block
	Mutex.Lock()
	FileCarveSessionMap[fileCarveBlock.SessionID].Timestamp = time.Now()
	FileCarveSessionMap[fileCarveBlock.SessionID].lastBlock = fileCarveBlock.BlockID
	FileCarveSessionMap[fileCarveBlock.SessionID].blockData[fileCarveBlock.BlockID] = fileCarveBlock.BlockData
	Mutex.Unlock()

	// Check if received if received all data blocks
	if len(FileCarveSessionMap[fileCarveBlock.SessionID].blockData) < FileCarveSessionMap[fileCarveBlock.SessionID].totalBlocks {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Write File Carve to disk
	writeDataToFile(FileCarveSessionMap[fileCarveBlock.SessionID], w)

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

func writeDataToFile(fileCarveSession *FilCarveSession, w http.ResponseWriter) {
	fmt.Println("######################################### Creaete file #########################################")
	outFileName := fmt.Sprintf("%s/%s.tar", "/tmp", fileCarveSession.CarveID)
	fmt.Println(outFileName)

	// Create file
	// https://golang.org/pkg/os/
	f, err := os.Create(outFileName) // For read access.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

	// IN ORDER of each data block
	for i := 0; i <= len(fileCarveSession.blockData); i++ {
		raw, err := base64.StdEncoding.DecodeString(fileCarveSession.blockData[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		f.Write([]byte(raw))
	}
	f.Close()
	fmt.Println("######################################### Close file #########################################")
}
