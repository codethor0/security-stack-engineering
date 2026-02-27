package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/l6_response"
	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

func TestResponseEngineCreatesIncident(t *testing.T) {
	bus := orchestrator.NewInMemoryBus(nil, nil)
	persistence := orchestrator.NewInMemoryPersistence()

	incidentCh := bus.Subscribe("topic/incidents")

	roe := l6_response.NewResponseEngine(bus, persistence)
	go roe.Run()

	time.Sleep(50 * time.Millisecond)

	alert := l6_response.Alert{
		AlertID:      "alert-1",
		Timestamp:    "2025-01-01T12:00:00Z",
		RuleID:       "rule-1",
		TechniqueID:  "T1078",
		Severity:     "high",
		Entities:     map[string]interface{}{"task_id": "task-1"},
		EvidenceRefs: []string{},
		Status:       "new",
	}
	_ = bus.Publish("topic/alerts", alert)

	select {
	case msg := <-incidentCh:
		var incident l6_response.Incident
		if err := json.Unmarshal(msg.Payload, &incident); err != nil {
			t.Fatalf("Unmarshal incident: %v", err)
		}
		if incident.ID == "" {
			t.Error("Incident ID empty")
		}
		if incident.Severity != "high" {
			t.Errorf("Severity: got %s", incident.Severity)
		}
		if len(incident.Alerts) != 1 {
			t.Errorf("Alerts: got %d", len(incident.Alerts))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for incident")
	}
}
