package upload

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitiateMongoClient this function inits the Mongo connector
// https://www.mongodb.com/blog/post/quick-start-golang--mongodb--a-quick-look-at-gridfs
func InitiateMongoClient(cfg *config.Config) (*gridfs.Bucket, *mongo.Client) {
	log.Printf("[+] - Init Mongo connector")

	var mongoClientConnector *mongo.Client
	var err error
	mongoURI := fmt.Sprintf("mongodb://%s:%d", cfg.Storage.Mongo.Host, cfg.Storage.Mongo.Port)
	opts := options.Client()
	opts.ApplyURI(mongoURI)
	opts.SetMaxPoolSize(5)
	if mongoClientConnector, err = mongo.Connect(context.Background(), opts); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Init Mongo GridFS bucket
	mongoBucketConnector, err := gridfs.NewBucket(mongoClientConnector.Database(cfg.Storage.Mongo.Database))
	log.Printf("[+] - Created GridFS bucket for file uploads")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Return Mongo Bucket Connector
	return mongoBucketConnector, mongoClientConnector
}

// UploadFileCarveToMongo
func UploadFileCarveToMongo(w http.ResponseWriter, r *http.Request, mongoBucketConnector *gridfs.Bucket) {
	fmt.Println("######################################### Uploading Block to GridFS #########################################")
	// Declare a new FileCarveBlock obj
	var fileCarveBlock FileCarveBlock

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&fileCarveBlock)
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

	// Write File Carve to GridFS
	writeDataToGridFS(FileCarveSessionMap[fileCarveBlock.SessionID], w, mongoBucketConnector)

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

func writeDataToGridFS(fileCarveSession *FilCarveSession, w http.ResponseWriter, mongoBucketConnector *gridfs.Bucket) {
	fmt.Println("######################################### Create GridFS #########################################")
	grifsFileName := fmt.Sprintf("%s", fileCarveSession.CarveID)

	// Create GridFS file
	// https://golang.org/pkg/os/
	// this is the name of the file which will be saved in the database
	uploadStream, err := mongoBucketConnector.OpenUploadStream(grifsFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer uploadStream.Close()

	// IN ORDER of each data block
	var fullFileBytes []byte
	for i := 0; i <= len(fileCarveSession.blockData); i++ {
		// Decode Base64 block of data
		raw, err := base64.StdEncoding.DecodeString(fileCarveSession.blockData[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Append raw bytes
		fullFileBytes = append(fullFileBytes, raw...)

	}

	// Upload data block
	fileSize, err := uploadStream.Write(fullFileBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(fileSize)

	fmt.Println("######################################### Write file to DB was successful. #########################################")

}
