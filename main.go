package main

import (
	"os"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/status"
	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/upload"

	"fmt"
	"log"
	"net/http"
)

func setupRoutes(cfg *config.Config) {
	http.HandleFunc("/start_uploads", upload.StartFileCarve)
	http.HandleFunc("/upload_blocks", upload.UploadFileCarve)
	http.HandleFunc("/status", status.HealthCheck)

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
	// Generate our config based on the config supplied
	cfg, err := config.NewConfig("conf/osquery-file-carve.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting osquery-file-carve-server on port %d\n", cfg.Webserver.Port)

	// Create GO ticker to clean up old entries
	//go session.CleanUpOldSessions(cfg)

	setupRoutes(cfg)

}
