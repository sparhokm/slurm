package models

type Event struct {
	EventID     int64  `json:"eventID,string"`
	EventType   int64  `json:"eventType,string"`
	UnixTime    int64  `json:"time,string"`
	FileID      string `json:"fileID"`
	OwnerID     int64  `json:"ownerID,string"`
	Filepath    string `json:"filepath"`
	Size        int64  `json:"size,string"`
	ContentType string `json:"contentType"`
	RequestID   string `json:"requestID"`
	TraceID     string `json:"traceID"`
	SpanID      string `json:"spanID"`
}
