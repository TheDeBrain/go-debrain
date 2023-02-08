package utils

import (
	uuid "github.com/satori/go.uuid"
	"log"
)

func CrtUUID() string {
	u := uuid.NewV1()
	if len(u) == 0 {
		log.Println("Failed to create uuid")
		return ""
	}
	return u.String()
}
