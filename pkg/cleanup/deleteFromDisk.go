package cleanup

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
)

func fileTimestampCheck(fTimestamp time.Time, expireInterval int) bool {
	return time.Now().Sub(fTimestamp) > time.Duration(expireInterval)*time.Hour
}

func getListOfFilesInDir(storageDirectory string) []os.FileInfo {
	// Scan all files in directory
	files, err := ioutil.ReadDir(storageDirectory)
	if err != nil {
		log.Println(err.Error())
	}
	return files
}

func deleteFile(storageDirectory string, f os.FileInfo, expireInterval int) {
	if fileTimestampCheck(f.ModTime(), expireInterval) {
		filePath := fmt.Sprintf("%s/%s", storageDirectory, f.Name())
		log.Printf("Deleting %s\n", filePath)

		err := os.Remove(filePath)
		if err != nil {
			log.Printf("Coulld NOT delete %s - Error: %s\n", filePath, err.Error())
		}

	}
}

// DeleteOldFilesOnDisk is a golang ticker that runs
// perodically to clean up old files
func DeleteOldFilesOnDisk(cfg *config.Config) {
	// Creat ticker
	ticker := time.NewTicker(time.Duration(cfg.Cleanup.TimeInvterval) * time.Second)

	for range ticker.C {
		// Get list of files in directory
		fileList := getListOfFilesInDir(cfg.Storage.File.Location)

		// Delete old files
		for _, f := range fileList {
			deleteFile(cfg.Storage.File.Location, f, cfg.Cleanup.ExpireInterval)
		}
		log.Printf("[+] - %s - Cleaned up old files on disk", time.Now().Format("2006-01-02 15:04:05"))

	}
}
