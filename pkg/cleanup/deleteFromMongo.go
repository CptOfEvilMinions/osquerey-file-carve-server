package cleanup

import (
	"context"
	"log"
	"time"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func deleteMongoDoc(mongoDatabase string, bucket *gridfs.Bucket, expireInterval int) {
	// Generate expired timestamp
	expiredTime := time.Now().Add(time.Duration(-expireInterval) * time.Second)

	// Generate filter
	mFilter := bson.M{"uploadDate": bson.M{"$lt": expiredTime}}

	// Find all the documents that need to be deleted
	cursor, err := bucket.Find(mFilter)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err := cursor.Close(context.TODO()); err != nil {
			log.Println(err)
		}
	}()

	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	type gridfsFile struct {
		Name       string             `bson:"filename"`
		ObjectID   primitive.ObjectID `bson:"_id"`
		Length     int64              `bson:"length"`
		UploadDate primitive.DateTime `bson:"uploadDate"`
	}

	var results []gridfsFile
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// Print all files being deleted
	for _, result := range results {
		log.Printf("Deleting Filename: %s - ObjectID: %s - Upload Date: %s From Mongo\n", result.Name, result.ObjectID, result.UploadDate.Time())
		if err := bucket.Delete(result.ObjectID); err != nil {
			log.Println(err)
		}
	}
}

// DeleteOldFilesFromMongo is a golang ticker that runs
// perodically to clean up old files
func DeleteOldFilesFromMongo(cfg *config.Config, mongoBucketConnector *gridfs.Bucket) {

	// Creat ticker
	ticker := time.NewTicker(time.Duration(cfg.Cleanup.TimeInvterval) * time.Second)

	for range ticker.C {
		log.Println("Start Mongo cleanup of old files")
		deleteMongoDoc(cfg.Storage.Mongo.Database, mongoBucketConnector, cfg.Cleanup.ExpireInterval)
	}
}
