package time

import (
	"time"
)

// Get the current timestamp for s
func GetCurrentTSForS() int64 {
	return time.Now().Unix()
}

// Get the current timestamp for ms
func GetCurrentTSForMS() int64 {
	return time.Now().UnixNano() / 1e6
}
