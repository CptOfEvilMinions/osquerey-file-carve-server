package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// JSONDecode attempts to decode the request body into the struct. If there is an error,
// respond to the client with the error message and a 400 status code.
func JSONDecode(w http.ResponseWriter, rBody io.ReadCloser, fileCarveBlock *FileCarveBlock) {
	err := json.NewDecoder(rBody).Decode(fileCarveBlock)
	if err != nil {
		log.Println("uploadHelper - JSONDecode -", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
}

// CheckSessionIDexists check if Sesssion ID exists. If sessionID exists return true
// else return false
// https://www.quora.com/In-Go-how-do-I-use-a-map-with-a-string-key-and-a-struct-as-value
// https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
func CheckSessionIDexists(sessionID string, FileCarveSessionMap map[string]*FilCarveSession) bool {
	if _, exist := FileCarveSessionMap[sessionID]; !exist {
		for k := range FileCarveSessionMap {
			log.Println(k)
		}
		return false
	}
	return true
}

// SucessfulUpload returns status codes (200) to the client for each
// sucessful data block upload or a status code of 200 AND {"success": true}
// for complete file uploads
func SucessfulUpload(w http.ResponseWriter, uploadComplete bool) error {
	// Let the client know ALL data block have been recevied sucessfully
	if uploadComplete == true {
		// Create map for JSON and set vaule
		resp := map[string]bool{"success": true}

		// Marshal map into JSON
		// Return 404 if JOSN can't be marshalled
		js, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		// Return sucess to client
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

	// Let the client know the data block was uploaded sucessfully
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	return nil
}
