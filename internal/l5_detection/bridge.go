package l5_detection

import (
	"encoding/json"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

// MessageBus for L5 bridge.
type MessageBus interface {
	Publish(topic string, payload interface{}) error
	Subscribe(topic string) <-chan orchestrator.Message
}

// DetectionBridge is a minimal Go bridge for Pathway A.
// It listens for L4 execution evidence and emits alerts when tasks complete.
type DetectionBridge struct {
	messageBus MessageBus
}

// Alert is the L5 alert structure.
type Alert struct {
	AlertID               string                 `json:"alert_id"`
	Timestamp             string                 `json:"timestamp"`
	RuleID                string                 `json:"rule_id"`
	TechniqueID           string                 `json:"technique_id"`
	Severity               string                 `json:"severity"`
	Entities               map[string]interface{} `json:"entities"`
	EvidenceRefs           []string               `json:"evidence_refs"`
	EngagementCorrelation  string                 `json:"engagement_correlation"`
	Status                 string                 `json:"status"`
}

// NewDetectionBridge creates a detection bridge.
func NewDetectionBridge(bus MessageBus) *DetectionBridge {
	return &DetectionBridge{messageBus: bus}
}

// Run listens for L4 evidence and emits alerts on task completion.
func (b *DetectionBridge) Run() {
	ch := b.messageBus.Subscribe("topic/evidence_records")
	for msg := range ch {
		var record orchestrator.EvidenceRecord
		if err := json.Unmarshal(msg.Payload, &record); err != nil {
			continue
		}
		if record.Layer != "L4" || record.Action != "task_completed" {
			continue
		}

		alert := Alert{
			AlertID:              orchestrator.GenerateID(),
			Timestamp:            time.Now().Format(time.RFC3339),
			RuleID:               "rule-rtea-simulated",
			TechniqueID:          record.TechniqueID,
			Severity:             "high",
			Entities:             map[string]interface{}{"task_id": record.TaskID},
			EvidenceRefs:          []string{record.TaskID},
			EngagementCorrelation: record.EngagementID,
			Status:                "new",
		}
		if alert.TechniqueID == "" {
			alert.TechniqueID = "T1078"
		}

		_ = b.messageBus.Publish("topic/alerts", alert)
	}
}
