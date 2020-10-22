package models

import (
	"time"
)

const (
	StatusEnable  = "0"
	StatusDisable = "1"
)

type Model struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
