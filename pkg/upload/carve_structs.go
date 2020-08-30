package upload

import (
	"sync"
	"time"
)

// FilCarveSession struct to track state and data of file carve session
type FilCarveSession struct {
	Timestamp   time.Time
	blockData   map[int]string
	totalBlocks int
	lastBlock   int
	CarveID     string
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
