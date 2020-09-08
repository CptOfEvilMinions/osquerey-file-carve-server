package download

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/auth"
	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// retrieveFile checks if the file by GGUID exists in Mongo.
// If it exists it returns result and nil, if not returns nil and err
func retrieveFile(fileCarveGUID string, mongoCollection *mongo.Collection) (primitive.M, error) {
	// Generate filter to search for document
	mongoFilter := bson.M{"filename": fileCarveGUID}

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

func checkMongoDocExists(fileCarveGUID string, mongoCollection *mongo.Collection) bool {
	// Generate filter to search for document
	mongoFilter := bson.M{"filename": fileCarveGUID}

	// Search for document
	// If error is nil a document was found
	// Else return false
	var result bson.M
	if err := mongoCollection.FindOne(context.Background(), mongoFilter).Decode(&result); err == nil {
		return true
	}
	return false
}

// FileRequestFromMongo this function will take in download requests
func FileRequestFromMongo(w http.ResponseWriter, r *http.Request, cfg *config.Config, mongoClientConnector *mongo.Client) {
	// Declare a new FileRequest obj,
	var fileRequest FileRequest

	// Decode JSON file request
	err := json.NewDecoder(r.Body).Decode(&fileRequest)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate token
	err = auth.TokenValdiation(fileRequest.TokenAccessor, fileRequest.Token)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log request
	log.Printf("[*] - File GUID request: %s - FROM: %s - Token accessor: %s", fileRequest.FileCarveGUID, extractXFordwardedFor(r), fileRequest.TokenAccessor)

	// Create Mongo DB connector
	db := mongoClientConnector.Database(cfg.Storage.Mongo.Database)
	mongoCollection := db.Collection("fs.files")

	// Check file exists
	if result := checkMongoDocExists(fileRequest.FileCarveGUID, mongoCollection); result == false {
		err := fmt.Errorf("File does not exist: %s", fileRequest.FileCarveGUID)
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	// Get content type of file
	FileContentType := "application/octet-stream"

	// Get the file size
	FileSize := strconv.FormatInt(dStream, 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.tar", fileRequest.FileCarveGUID))
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	io.Copy(w, &buf) //'Copy' the file to the client
	return
}
