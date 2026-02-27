package orchestrator

import "encoding/json"

// Message is the unit of exchange between algorithms on the bus.
type Message struct {
	Topic         string          `json:"topic"`
	Payload       json.RawMessage `json:"payload"`
	SchemaVersion string          `json:"schema_version"`
	Timestamp     int64           `json:"timestamp"`
	Producer      string          `json:"producer"`
	MessageID     string          `json:"message_id"`
	CorrelationID  string          `json:"correlation_id"`
}

// EvidenceRecord is the canonical audit evidence structure.
type EvidenceRecord struct {
	Layer        string      `json:"layer"`
	Action       string      `json:"action"`
	Timestamp    int64       `json:"timestamp"`
	Attribution  string      `json:"attribution,omitempty"`
	EngagementID string      `json:"engagement_id,omitempty"`
	TaskID       string      `json:"task_id,omitempty"`
	Operator     string      `json:"operator,omitempty"`
	TechniqueID  string      `json:"technique_id,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Error        string      `json:"error,omitempty"`
}
