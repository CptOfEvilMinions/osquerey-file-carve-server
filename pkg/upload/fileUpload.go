package upload

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
)

func writeDataToFileStream(fStream *os.File, currentDataBlock string) error {
	// Decode Base64 block of data
	rawDataBlock, err := base64.StdEncoding.DecodeString(currentDataBlock)
	if err != nil {
		return err
	}

	// Write raw data block to file
	if _, err := fStream.Write(rawDataBlock); err != nil {
		return err
	}
	return nil
}

func closeFileStream(fStream *os.File, sessionID string, FileCarveSessionMap map[string]*FilCarveSession) {
	// Close file stream
	fStream.Close()

	// Delete session from FileCarveSessionMap
	Mutex.Lock()                           // Lock access to FileCarveSessionMap
	delete(FileCarveSessionMap, sessionID) // Delete session from FileCarveSessionMap
	Mutex.Unlock()                         // UNlock access to FileCarveSessionMap
}

func createFileStream(storageLocation string, carveID string) (*os.File, error) {
	outFileName := fmt.Sprintf("%s/%s.tar", storageLocation, carveID)
	log.Println("Create file stream: ", outFileName)
	fo, err := os.Create(outFileName)
	return fo, err
}

// FileCarveToDisk takes in file carve data blocks and processes them
//https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func FileCarveToDisk(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	// Declare a new FileCarveBlock obj
	var fileCarveBlock FileCarveBlock

	// Decode JSON blob to fileCarveBlock obj
	JSONDecode(w, r.Body, &fileCarveBlock)

	// Check if Sesssion ID exists.
	result := CheckSessionIDexists(fileCarveBlock.SessionID, FileCarveSessionMap)
	if result == false {
		http.Error(w, "Session ID does not exist", http.StatusBadRequest)
		return
	}

	// If sessionID exists
	// Update values for File Carve
	Mutex.Lock()                                                                                                                                                    // Lock access to FileCarveSessionMap
	FileCarveSessionMap[fileCarveBlock.SessionID].Timestamp = time.Now()                                                                                            // Set timestamp to the latest time a block was received
	FileCarveSessionMap[fileCarveBlock.SessionID].lastBlockReceived = fileCarveBlock.BlockID                                                                        // Set this to the latest block ID received
	FileCarveSessionMap[fileCarveBlock.SessionID].ReceivedBlockIDs = append(FileCarveSessionMap[fileCarveBlock.SessionID].ReceivedBlockIDs, fileCarveBlock.BlockID) // Add current block ID to blocks received list

	// If File Stream is nil create one
	if FileCarveSessionMap[fileCarveBlock.SessionID].FileStream == nil {
		fStream, err := createFileStream(cfg.Storage.File.Location, FileCarveSessionMap[fileCarveBlock.SessionID].CarveID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		FileCarveSessionMap[fileCarveBlock.SessionID].FileStream = fStream
	}

	Mutex.Unlock() // UNlock access to FileCarveSessionMap

	// Extract data block from JSON payload
	// Write the current data block to the file stream
	if err := writeDataToFileStream(FileCarveSessionMap[fileCarveBlock.SessionID].FileStream, fileCarveBlock.BlockData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if all blocks have been received
	// If NOT all blocks have been received return 200 for sucessful block upload
	if len(FileCarveSessionMap[fileCarveBlock.SessionID].ReceivedBlockIDs) < FileCarveSessionMap[fileCarveBlock.SessionID].totalBlocks {
		if err := SucessfulUpload(w, false); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	// Close File Stream
	closeFileStream(FileCarveSessionMap[fileCarveBlock.SessionID].FileStream, fileCarveBlock.SessionID, FileCarveSessionMap)

	// Instruct client sucessful upload
	if err := SucessfulUpload(w, true); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
