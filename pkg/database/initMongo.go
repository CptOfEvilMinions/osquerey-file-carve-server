package database

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// mongodb://[<username>:<password>@]<mongo hostname>[:port][/[defaultauthdb]?authSource=admin[&ssl=<bool>&option1&option2]
	// https://docs.mongodb.com/manual/reference/connection-string/
	var mongoURI string
	if cfg.Storage.Mongo.Username != "" && cfg.Storage.Mongo.Password != "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin&ssl=%t&%s",
			cfg.Storage.Mongo.Username,
			cfg.Storage.Mongo.Password,
			cfg.Storage.Mongo.Host,
			cfg.Storage.Mongo.Port,
			cfg.Storage.Mongo.Database,
			cfg.Storage.Mongo.SSL,
			cfg.Storage.Mongo.Options,
		)
	} else {
		mongoURI = fmt.Sprintf("mongodb://@%s:%d/%s?authSource=admin&ssl=%t&%s",
			cfg.Storage.Mongo.Host,
			cfg.Storage.Mongo.Port,
			cfg.Storage.Mongo.Database,
			cfg.Storage.Mongo.SSL,
			cfg.Storage.Mongo.Options,
		)
	}

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
