package auth

import (
	"errors"
	"sync"
	"time"
)

// This map contains all valid tokens for download requests
// key: vault token ACCESSOR - NOT secret token
// Value: expiration date and time
// mutex lock for map
var tokenSessionMap = make(map[string]time.Time)
var mutex = &sync.Mutex{}

func localTokenLookup(tokenAccessor string) bool {
	// First check to see if the token exist in map
	if tokenExpiration, ok := tokenSessionMap[tokenAccessor]; ok {
		// Make sure the token is still valid
		// https://stackoverflow.com/questions/20924303/how-to-do-date-time-comparison
		// Conver both timestamps to UTC then compare
		// https://socketloop.com/tutorials/golang-convert-unix-timestamp-to-utc-timestamp#:~:text=Solution%20%3A,Coordinated%20Universal%20Time)%20time%20stamp.
		if tokenExpiration.UTC().After(time.Now().UTC()) {
			return true
		}
	}

	// Remove item from map
	mutex.Lock()
	delete(tokenSessionMap, tokenAccessor)
	mutex.Unlock()

	return false
}

// TokenValdiation performs token validation checks
func TokenValdiation(tokenAccessor string, token string) error {
	// Make sure JSON payload has token if not reject
	if token == "" || tokenAccessor == "" {
		return errors.New(`{"auth_error":"No token accessor OR token provided"}`)
	}

	// Check if token accessor is known AND valid
	if localTokenLookup(tokenAccessor) {
		return nil
	}

	// If token accessor does exit do a lookup on the token
	if err := vaultTokenLookup(token); err == nil {
		return nil
	}
	// Reject request if it did not pass the checks
	return errors.New(`{"auth_error": "Token AND token accessor did not pass validation checks"}`)
}
