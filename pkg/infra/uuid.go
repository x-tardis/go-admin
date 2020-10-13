package infra

import (
	"github.com/google/uuid"
	uuid2 "github.com/satori/go.uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateUUID2() string {
	return uuid2.NewV4().String()
}
