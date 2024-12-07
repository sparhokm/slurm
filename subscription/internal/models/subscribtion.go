package models

import (
	"time"
)

type Subscription struct {
	ID           string
	UserID       int64
	FilesOwnerID *int64
	Prefix       string
	CreatedAt    time.Time
}

type SubscriptionShort struct {
	UserID       int64
	FilesOwnerID *int64
}
