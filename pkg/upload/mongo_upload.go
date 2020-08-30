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

	// Decode JSON blob to fileCarveBlock obj
	JSONDecode(w, r.Body, &fileCarveBlock)

	// Check if Sesssion ID exists.
	result := CheckSessionIDexists(fileCarveBlock.SessionID, FileCarveSessionMap)
	if result == false {
		http.Error(w, "Session ID does not exist", http.StatusBadRequest)
		return
	}

	// Create uploadStream
	mongoUploadStream := createMongoUploadStream(mongoBucketConnector)

	////////////////////////////////////////////// LOCK Mutex //////////////////////////////////////////////
	// Add block to map
	// Set last block
	Mutex.Lock()
	FileCarveSessionMap[fileCarveBlock.SessionID].Timestamp = time.Now()
	FileCarveSessionMap[fileCarveBlock.SessionID].lastBlockReceived = fileCarveBlock.BlockID
	FileCarveSessionMap[fileCarveBlock.SessionID].MongoUploadStream = mongoUploadStream
	//FileCarveSessionMap[fileCarveBlock.SessionID].blockData[fileCarveBlock.BlockID] = fileCarveBlock.BlockData
	Mutex.Unlock()
	////////////////////////////////////////////// LOCK Mutex //////////////////////////////////////////////

	// Check if all data blocks were received
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

func createMongoUploadStream(mongoBucketConnector *gridfs.Bucket) *gridfs.UploadStream {
	// Create uploadStream
	uploadStream, err := mongoBucketConnector.OpenUploadStream("filename")
	if err != nil {
		log.Fatal(err)
	}

	// Defer closing stream unless there is an error
	defer func() {
		if err = uploadStream.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	return uploadStream

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
