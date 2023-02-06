package utils

import "time"

// get current timestamp by ms
func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
