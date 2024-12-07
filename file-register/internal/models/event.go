package models

import (
	"time"
)

const (
	EventFileAdd = iota
	EventFileUpdate
	EventFileDelete
)

type FileEvent struct {
	ID          int64
	Type        string
	RequestID   *string
	TraceID     *string
	SpanID      *string
	File        File
	CreatedAt   *time.Time
	ProcessedAt *time.Time
}
