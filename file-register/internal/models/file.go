package models

import (
	"time"
)

type File struct {
	ID          string     `json:"id"`
	OwnerID     int64      `json:"owner_id"`
	ContentType string     `json:"conten_type"`
	Size        int64      `json:"size"`
	Filepath    string     `json:"filepath"`
	Version     int64      `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
