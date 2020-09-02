package session

import (
	"fmt"
	"time"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/upload"
)

// CleanUpOldSessions if timestamp is greater than 8 hours remove entry
// https://stackoverflow.com/questions/23838857/how-to-get-hours-difference-between-two-dates
func CleanUpOldSessions(cfg *config.Config) {
	// Create ticker
	ticker := time.NewTicker(time.Duration(cfg.Cleanup.TimeInvterval) * time.Second)

	for range ticker.C {

		upload.Mutex.Lock()
		for sessionID, fileCarveSession := range upload.FileCarveSessionMap {
			// If timestamp of CarveSession hasn't had a new block
			// in X amount of time delete it
			// https://www.geeksforgeeks.org/time-since-function-in-golang-with-examples/
			// https://www.golangprograms.com/convert-int-to-float-in-golang.html
			if time.Since(fileCarveSession.Timestamp).Seconds() >= float64(cfg.Cleanup.ExpireInterval) {
				delete(upload.FileCarveSessionMap, sessionID)
				fmt.Printf("[+] - %s - Cleaned up old session ID -s %s", time.Now().Format("2006-01-02 15:04:05"), sessionID)
			}
		}
		upload.Mutex.Unlock()

	}
}
