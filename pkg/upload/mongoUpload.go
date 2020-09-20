package upload

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo"
)

// FileCarveToMongo allows file carves to be uploaded to Mongo
func FileCarveToMongo(w http.ResponseWriter, r *http.Request, mongoBucketConnector *gridfs.Bucket, mongoClientConnector *mongo.Client, mongoCollectionConnector *mongo.Collection) {
	// Declare a new FileCarveBlock obj
	var fileCarveBlock FileCarveBlock

	// Decode JSON blob to fileCarveBlock obj
	JSONDecode(w, r.Body, &fileCarveBlock)

	// Check if Sesssion ID exists.
	result := CheckSessionIDexists(fileCarveBlock.SessionID, FileCarveSessionMap)
	if result == false {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, "Session ID does not exist"))
		return
	}

	// If sessionID exists
	// Update values for File Carve
	Mutex.Lock()                                                                                                                                                    // Lock access to FileCarveSessionMap
	FileCarveSessionMap[fileCarveBlock.SessionID].Timestamp = time.Now()                                                                                            // Set timestamp to the latest time a block was received
	FileCarveSessionMap[fileCarveBlock.SessionID].lastBlockReceived = fileCarveBlock.BlockID                                                                        // Set this to the latest block ID received
	FileCarveSessionMap[fileCarveBlock.SessionID].ReceivedBlockIDs = append(FileCarveSessionMap[fileCarveBlock.SessionID].ReceivedBlockIDs, fileCarveBlock.BlockID) // Add current block ID to blocks received list

	// If Mongo Upload Stream is nil create one
	if FileCarveSessionMap[fileCarveBlock.SessionID].MongoUploadStream == nil {
		// Create Mongo Upload Stream
		mStream, err := createMongoUploadStream(mongoBucketConnector, FileCarveSessionMap[fileCarveBlock.SessionID].CarveID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}
		FileCarveSessionMap[fileCarveBlock.SessionID].MongoUploadStream = mStream
	}
	Mutex.Unlock() // UNlock access to FileCarveSessionMap

	// Extract data block from JSON payload
	// Write the current data block to the file stream
	if err := writeDataToMongoStream(FileCarveSessionMap[fileCarveBlock.SessionID].MongoUploadStream, fileCarveBlock.BlockData); err != nil {
		log.Panicln(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	// Check if all blocks have been received
	// If NOT all blocks have been received return 200 for sucessful block upload
	if len(FileCarveSessionMap[fileCarveBlock.SessionID].ReceivedBlockIDs) < FileCarveSessionMap[fileCarveBlock.SessionID].totalBlocks {
		if err := SucessfulUpload(w, false); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}
		return
	}

	// Close File Stream
	if err := closeMongoUploadStream(FileCarveSessionMap[fileCarveBlock.SessionID].MongoUploadStream, fileCarveBlock.SessionID, FileCarveSessionMap); err != nil {
		log.Panicln("Close Mongo file stream error: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	// Instruct client sucessful upload
	if err := SucessfulUpload(w, true); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

}



func closeMongoUploadStream(mStream *gridfs.UploadStream, sessionID string, FileCarveSessionMap map[string]*FilCarveSession) error {
	// Delete session from FileCarveSessionMap
	fmt.Println(FileCarveSessionMap)
	Mutex.Lock()                           // Lock access to FileCarveSessionMap
	delete(FileCarveSessionMap, sessionID) // Delete session from FileCarveSessionMap
	Mutex.Unlock()                         // UNlock access to FileCarveSessionMap

	fmt.Println(FileCarveSessionMap)

	// Close Mongo stream
	if err := mStream.Close(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func createMongoUploadStream(mongoBucketConnector *gridfs.Bucket, CarveID string) (*gridfs.UploadStream, error) {
	log.Println("Create Mongo GridFS file object: ", CarveID)
	// Create uploadStream
	println("File CarveID: ", CarveID)	
	uploadStream, err := mongoBucketConnector.OpenUploadStream(CarveID)
	return uploadStream, err
}

func writeDataToMongoStream(mStream *gridfs.UploadStream, currentDataBlock string) error {
	// Decode Base64 block of data
	rawDataBlock, err := base64.StdEncoding.DecodeString(currentDataBlock)
	if err != nil {
		return err
	}

	// Write raw data block to file
	if _, err := mStream.Write(rawDataBlock); err != nil {
		return err
	}
	return nil

}