package upload

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func writeDataToFileStream(w http.ResponseWriter, fStream *os.File, currentDataBlock string) error {
	fmt.Println("######################################### Writing current data block to file stream #########################################")
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
	fmt.Println("######################################### Close file stream #########################################")
	// Close file stream
	fStream.Close()

	// Delete session from FileCarveSessionMap
	Mutex.Lock()                           // Lock access to FileCarveSessionMap
	delete(FileCarveSessionMap, sessionID) // Delete session from FileCarveSessionMap
	Mutex.Unlock()                         // UNlock access to FileCarveSessionMap
}

func createFileStream(carveID string) (*os.File, error) {
	fmt.Println("######################################### Create file stream #########################################")
	outFileName := fmt.Sprintf("%s/%s.tar", "/tmp", carveID)
	fo, err := os.Create(outFileName)
	return fo, err
}

// FileCarveToDisk takes in file carve data blocks and processes them
//https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func FileCarveToDisk(w http.ResponseWriter, r *http.Request) {
	fmt.Println("######################################### Uploading Block #########################################")
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
	if result == true {
		Mutex.Lock()                                                                             // Lock access to FileCarveSessionMap
		FileCarveSessionMap[fileCarveBlock.SessionID].Timestamp = time.Now()                     // Set timestamp to the latest time a block was received
		FileCarveSessionMap[fileCarveBlock.SessionID].lastBlockReceived = fileCarveBlock.BlockID // Set this to the latest block ID received

		// If File Stream is nil create one
		if FileCarveSessionMap[fileCarveBlock.SessionID].FileStream == nil {
			fStream, err := createFileStream(FileCarveSessionMap[fileCarveBlock.SessionID].CarveID)
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
		if err := writeDataToFileStream(w, FileCarveSessionMap[fileCarveBlock.SessionID].FileStream, fileCarveBlock.BlockData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}

	// Check if all blocks have been received
	// If NOT all blocks have been received return 200 for sucessful block upload
	if len(FileCarveSessionMap[fileCarveBlock.SessionID].blockData) < FileCarveSessionMap[fileCarveBlock.SessionID].totalBlocks {
		SucessfulUpload(w, false)
		return
	}

	// Close File Stream
	closeFileStream(FileCarveSessionMap[fileCarveBlock.SessionID].FileStream, fileCarveBlock.SessionID, FileCarveSessionMap)

	// Instruct client sucessful upload
	SucessfulUpload(w, true)
}
