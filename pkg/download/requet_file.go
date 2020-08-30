package download

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// FileRequest
type FileRequest struct {
	FileCarveGUID string `json:"file_carve_guid"`
}

type nope struct {
	FileCarveGUID string `json:"file_carve_guid"`
}

// retrieveFile checks if the file by GGUID exists in Mongo.
// If it exists it returns result and nil, if not returns nil and err
func retrieveFile(fileCarveGUIDgo string, mongoCollection *mongo.Collection) (primitive.M, error) {
	// Generate filter to search for document
	mongoFilter := bson.M{"filename": fileCarveGUIDgo}

	// Search for document
	var result bson.M
	err := mongoCollection.FindOne(context.Background(), mongoFilter).Decode(&result)

	// If error is nil a document was found
	// Else return false
	if err == nil {
		return result, nil
	}
	return nil, err
}

// FileRequestFromMongo this function will take in download requests
func FileRequestFromMongo(w http.ResponseWriter, r *http.Request, cfg *config.Config, mongoClientConnector *mongo.Client) {
	// Declare a new FileRequest obj,
	var fileRequest FileRequest

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&fileRequest)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create Mongo DB connector
	db := mongoClientConnector.Database(cfg.Storage.Mongo.Database)
	mongoCollection := db.Collection("fs.files")

	// Check file exists
	// If file does NOT exist return error to user
	result, err := retrieveFile(fileRequest.FileCarveGUID, mongoCollection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Print result
	fmt.Println(result)

	// // Create outfile name
	fileName := fmt.Sprintf("%s/mongo-%s.tar", "/tmp", fileRequest.FileCarveGUID)
	fmt.Println(fileName)

	// Connect files bucket
	b, err := gridfs.NewBucket(mongoClientConnector.Database(cfg.Storage.Mongo.Database))
	if err != nil {
		log.Fatalln(err)
	}

	// Create buffer
	var buf bytes.Buffer

	// Retrieve file from Mongo by GGUID and write to buffer
	// Write buffer to disk
	dStream, err := b.DownloadToStreamByName(fileRequest.FileCarveGUID, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v\n", dStream)
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)

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
