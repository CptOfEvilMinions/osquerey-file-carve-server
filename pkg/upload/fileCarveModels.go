package upload

import (
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// FilCarveSession struct to track state and data of file carve session
type FilCarveSession struct {
	Timestamp         time.Time
	ReceivedBlockIDs  []int
	totalBlocks       int
	lastBlockReceived int
	CarveID           string
	MongoUploadStream *gridfs.UploadStream
	FileStream        *os.File
}

// FileCarveBlock struct for incoming file carve block
type FileCarveBlock struct {
	BlockID   int    `json:"block_id"`
	SessionID string `json:"session_id"`
	RequestID string `json:"request_id"`
	BlockData string `json:"data"`
}

// FileCarveSessionMap Map of FilCarveSession structs
var FileCarveSessionMap = make(map[string]*FilCarveSession)

// Mutex for FileCarveSessionMap
var Mutex = &sync.Mutex{}
