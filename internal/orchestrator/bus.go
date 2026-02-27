package orchestrator

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SchemaValidator validates message payloads against JSON schemas.
type SchemaValidator interface {
	Validate(topic string, payload json.RawMessage) error
}

// Persistence stores evidence and supports queries.
type Persistence interface {
	StoreEvidence(record EvidenceRecord) error
	QueryEvidence(layer string, since int64) ([]EvidenceRecord, error)
}

// MessageBus provides publish/subscribe and consume-latest semantics.
type MessageBus interface {
	Publish(topic string, payload interface{}) error
	Subscribe(topic string) <-chan Message
	ConsumeLatest(topic string) interface{}
}

// NoopValidator accepts all payloads (PoC stub).
type NoopValidator struct{}

func (NoopValidator) Validate(topic string, payload json.RawMessage) error {
	return nil
}

// InMemoryBus is a simple in-memory message bus.
type InMemoryBus struct {
	mu         sync.RWMutex
	subscribers map[string][]chan Message
	messageLog []Message
	logMu      sync.RWMutex
	validator  SchemaValidator
	persistence Persistence
	maxLog     int
}

// NewInMemoryBus creates a bus with optional validator and persistence.
func NewInMemoryBus(validator SchemaValidator, persistence Persistence) *InMemoryBus {
	if validator == nil {
		validator = NoopValidator{}
	}
	return &InMemoryBus{
		subscribers: make(map[string][]chan Message),
		messageLog:  make([]Message, 0, 1000),
		validator:   validator,
		persistence: persistence,
		maxLog:      10000,
	}
}

// Publish sends a message to all subscribers of the topic.
func (b *InMemoryBus) Publish(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}
	raw := json.RawMessage(data)
	if err := b.validator.Validate(topic, raw); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	msg := Message{
		Topic:         topic,
		Payload:       raw,
		SchemaVersion: "1.0",
		Timestamp:     time.Now().Unix(),
		MessageID:     GenerateID(),
		Producer:      producerFromTopic(topic),
	}

	b.mu.RLock()
	subs := append([]chan Message(nil), b.subscribers[topic]...)
	b.mu.RUnlock()

	for _, ch := range subs {
		select {
		case ch <- msg:
		default:
		}
	}

	b.logMu.Lock()
	b.messageLog = append(b.messageLog, msg)
	if len(b.messageLog) > b.maxLog {
		b.messageLog = b.messageLog[1000:]
	}
	b.logMu.Unlock()

	if topic == "topic/evidence_records" {
		var record EvidenceRecord
		if err := json.Unmarshal(data, &record); err == nil && b.persistence != nil {
			_ = b.persistence.StoreEvidence(record)
		}
	}

	return nil
}

// Subscribe returns a channel that receives messages for the topic.
func (b *InMemoryBus) Subscribe(topic string) <-chan Message {
	ch := make(chan Message, 100)
	b.mu.Lock()
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	b.mu.Unlock()
	return ch
}

// ConsumeLatest returns the most recent message payload for a topic, or nil.
func (b *InMemoryBus) ConsumeLatest(topic string) interface{} {
	b.logMu.RLock()
	defer b.logMu.RUnlock()
	for i := len(b.messageLog) - 1; i >= 0; i-- {
		if b.messageLog[i].Topic == topic {
			var result interface{}
			_ = json.Unmarshal(b.messageLog[i].Payload, &result)
			return result
		}
	}
	return nil
}

func producerFromTopic(topic string) string {
	switch topic {
	case "topic/governance_tokens":
		return "L0"
	case "topic/environment_snapshots":
		return "L1"
	case "topic/offensive_plan", "topic/offensive_findings":
		return "L4"
	case "topic/alerts":
		return "L5"
	case "topic/incidents":
		return "L6"
	case "topic/engineering_changes":
		return "L7"
	case "topic/assurance_reports":
		return "L8"
	case "topic/evidence_records":
		return "ORCH"
	default:
		return "UNKNOWN"
	}
}
