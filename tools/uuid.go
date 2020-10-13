package tools

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func GetUUID() string {
	u := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}
