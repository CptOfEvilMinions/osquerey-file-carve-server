package download

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/auth"
	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
)

func checkFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FileRequestFromDisk this function allows authenticated clients to download
// Osquery uploads by requesting the files GUID
func FileRequestFromDisk(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	// Declare a new FileRequest obj,
	var fileRequest FileRequest

	// Decode JSON file request
	err := json.NewDecoder(r.Body).Decode(&fileRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	// Validate token
	err = auth.TokenValdiation(fileRequest.TokenAccessor, fileRequest.Token, cfg.Vault.Policy)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	// Log request
	log.Printf("[*] - File GUID request: %s - FROM: %s - Token accessor: %s", fileRequest.FileCarveGUID, extractXFordwardedFor(r), fileRequest.TokenAccessor)

	// Check if file exists
	filePath := fmt.Sprintf("%s/%s.tar", cfg.Storage.File.Location, fileRequest.FileCarveGUID)
	if result := checkFileExists(filePath); result == false {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, "File with that GUID does not exist"))
		return
	}

	// Open file
	fromFile, err := os.Open(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	defer fromFile.Close()

	// Get the Content-Type of the file
	// Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)

	// Copy the headers into the FileHeader buffer
	fromFile.Read(FileHeader)

	// Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	// Get the file size
	FileStat, _ := fromFile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.tar", fileRequest.FileCarveGUID))
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	fromFile.Seek(0, 0)
	io.Copy(w, fromFile) //'Copy' the file to the client
	return
}
