package main

import (
	"flag"
	"os"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/download"
	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/status"
	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/upload"

	"fmt"
	"log"
	"net/http"
)

func setupRoutes(cfg *config.Config) {

	// Endpoint to init file upload
	http.HandleFunc("/start_uploads", upload.StartFileCarve)

	// Endpoint for healthchecks
	http.HandleFunc("/status", status.HealthCheck)

	// Endpoint for file upload destination
	// By default the if no other option is set to enabled
	// the option to write files to disk will be used.
	if cfg.Storage.Mongo.Enabled == true {
		log.Printf("[+] - Writing file uploads to Mongo")
		// Init Mongo Connectors for GridFS
		mongoBucketConnector, mongoClientConnector := upload.InitiateMongoClient(cfg)

		// Setup handler to use Mongo for file uploads
		http.HandleFunc("/upload_blocks", func(w http.ResponseWriter, r *http.Request) {
			upload.UploadFileCarveToMongo(w, r, mongoBucketConnector)
		})

		// Setup handler to use Mongo for file downloads
		http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
			download.FileRequestFromMongo(w, r, cfg, mongoClientConnector)
		})

	} else {
		log.Printf("[+] - Writing file uploads to disk")
		http.HandleFunc("/upload_blocks", upload.UploadFileCarveToDisk)
	}

	// If env debug listen on localhost and load SSL certs
	if os.Getenv("ENV") == "debug" {
		fmt.Println("################################ DEBUG MODE ################################")
		address := fmt.Sprintf(":%d", cfg.Webserver.Port)
		log.Fatal(http.ListenAndServeTLS(address, "conf/nginx/tls/snakeoil.crt", "conf/nginx/tls/snakeoil.key", nil))
	} else {
		address := fmt.Sprintf(":%d", cfg.Webserver.Port)
		log.Fatal(http.ListenAndServe(address, nil))
	}

}

func main() {
	// Read location of config from command line or load default location
	configLocationPtr := flag.String("config", "conf/osquery-file-carve.yml", "Set the file path for location of the config")
	flag.Parse()

	// Generate our config based on the config supplied
	cfg, err := config.NewConfig(*configLocationPtr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting osquery-file-carve-server on port %d\n", cfg.Webserver.Port)

	// Create GO ticker to clean up old entries
	//go session.CleanUpOldSessions(cfg)

	// Setup HTTP routes
	setupRoutes(cfg)

}
