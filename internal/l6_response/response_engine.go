package l6_response

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

// Alert is the L5 alert format (matches l5_detection.Alert).
type Alert struct {
	AlertID              string                 `json:"alert_id"`
	Timestamp            string                 `json:"timestamp"`
	RuleID               string                 `json:"rule_id"`
	TechniqueID          string                 `json:"technique_id"`
	Severity              string                 `json:"severity"`
	Entities              map[string]interface{} `json:"entities"`
	EvidenceRefs          []string               `json:"evidence_refs"`
	EngagementCorrelation string                 `json:"engagement_correlation"`
	Status                string                 `json:"status"`
}

// MessageBus for L6.
type MessageBus interface {
	Publish(topic string, payload interface{}) error
	Subscribe(topic string) <-chan orchestrator.Message
	ConsumeLatest(topic string) interface{}
}

// EvidenceStore for L6.
type EvidenceStore interface {
	StoreEvidence(record orchestrator.EvidenceRecord) error
}

// ResponseOrchestrationEngine implements the L6 algorithm.
type ResponseOrchestrationEngine struct {
	mu              sync.RWMutex
	incidents       map[string]*Incident
	playbooks       map[string]*Playbook
	activePlaybooks map[string]string
	bus             MessageBus
	store           EvidenceStore
}

// NewResponseEngine creates an L6 engine.
func NewResponseEngine(bus MessageBus, store EvidenceStore) *ResponseOrchestrationEngine {
	return &ResponseOrchestrationEngine{
		incidents:       make(map[string]*Incident),
		playbooks:       loadDefaultPlaybooks(),
		activePlaybooks: make(map[string]string),
		bus:             bus,
		store:           store,
	}
}

// Run starts the L6 loops.
func (roe *ResponseOrchestrationEngine) Run() {
	go roe.alertProcessingLoop()
	go roe.incidentLifecycleLoop()
}

func (roe *ResponseOrchestrationEngine) alertProcessingLoop() {
	ch := roe.bus.Subscribe("topic/alerts")
	for msg := range ch {
		var alert Alert
		if err := json.Unmarshal(msg.Payload, &alert); err != nil {
			continue
		}

		incident := roe.findOrCreateIncident(alert)
		playbook := roe.selectPlaybook(incident)
		roe.mu.Lock()
		roe.activePlaybooks[incident.ID] = playbook.ID
		roe.mu.Unlock()

		go roe.executePlaybook(incident, playbook)
	}
}

func (roe *ResponseOrchestrationEngine) findOrCreateIncident(alert Alert) *Incident {
	roe.mu.Lock()
	defer roe.mu.Unlock()

	for _, inc := range roe.incidents {
		if inc.Status != "closed" && inc.Severity == alert.Severity {
			inc.Alerts = append(inc.Alerts, alert.AlertID)
			inc.UpdatedAt = time.Now().Unix()
			return inc
		}
	}

	incident := roe.createIncident(alert)
	roe.incidents[incident.ID] = incident
	return incident
}

func (roe *ResponseOrchestrationEngine) createIncident(alert Alert) *Incident {
	slaMinutes := roe.getSLAForSeverity(alert.Severity)

	incident := &Incident{
		ID:          orchestrator.GenerateID(),
		Title:       fmt.Sprintf("Incident from alert %s", alert.AlertID),
		Severity:    alert.Severity,
		Status:      "new",
		Phase:       "detection_analysis",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		Alerts:      []string{alert.AlertID},
		Entities:    roe.extractEntitiesFromAlert(alert),
		SLADeadline: time.Now().Add(time.Duration(slaMinutes) * time.Minute).Unix(),
	}

	if alert.EngagementCorrelation != "" {
		incident.EngagementLink = alert.EngagementCorrelation
		incident.Title = fmt.Sprintf("[RED TEAM] %s", incident.Title)
	}

	roe.emitEvidence("incident_created", incident)
	_ = roe.bus.Publish("topic/incidents", incident)
	return incident
}

func (roe *ResponseOrchestrationEngine) extractEntitiesFromAlert(alert Alert) []Entity {
	entities := make([]Entity, 0)
	if actor, ok := alert.Entities["actor"]; ok {
		if am, ok := actor.(map[string]interface{}); ok {
			if id, ok := am["id"].(string); ok {
				entities = append(entities, Entity{
					Type: "identity", ID: id, Criticality: "high", CompromiseState: "suspected",
				})
			}
		}
	}
	if taskID, ok := alert.Entities["task_id"]; ok {
		if s, ok := taskID.(string); ok {
			entities = append(entities, Entity{
				Type: "asset", ID: s, Criticality: "medium", CompromiseState: "suspected",
			})
		}
	}
	if len(entities) == 0 {
		entities = append(entities, Entity{
			Type: "identity", ID: "unknown", Criticality: "medium", CompromiseState: "suspected",
		})
	}
	return entities
}

func (roe *ResponseOrchestrationEngine) selectPlaybook(incident *Incident) *Playbook {
	if incident.EngagementLink != "" {
		if pb, ok := roe.playbooks["pb-credential-abuse"]; ok && pb != nil {
			return pb
		}
	}
	if pb, ok := roe.playbooks["pb-default"]; ok && pb != nil {
		return pb
	}
	for _, pb := range roe.playbooks {
		if pb != nil {
			return pb
		}
	}
	return &Playbook{ID: "fallback", Name: "Fallback", Phases: []PlaybookPhase{}}
}

func (roe *ResponseOrchestrationEngine) executePlaybook(incident *Incident, playbook *Playbook) {
	roe.mu.Lock()
	incident.Phase = playbook.Phases[0].Name
	roe.mu.Unlock()

	for _, phase := range playbook.Phases {
		for _, tmpl := range phase.Actions {
			action := ResponseAction{
				ActionID:    orchestrator.GenerateID(),
				Type:        tmpl.Type,
				Description: fmt.Sprintf("Execute %s", tmpl.Type),
				ExecutedAt:  time.Now().Unix(),
				ExecutedBy:  "L6_RESPONSE_ORCHESTRATION",
				Status:      "completed",
				Result:      "simulated",
			}
			if tmpl.Automation {
				action.Result = "automated_simulated"
			}

			roe.mu.Lock()
			incident.Actions = append(incident.Actions, action)
			incident.UpdatedAt = time.Now().Unix()
			roe.mu.Unlock()
		}
	}

	roe.mu.Lock()
	incident.Status = "contained"
	incident.UpdatedAt = time.Now().Unix()
	roe.mu.Unlock()

	roe.emitEvidence("playbook_completed", incident)
	_ = roe.bus.Publish("topic/incidents", incident)
}

func (roe *ResponseOrchestrationEngine) incidentLifecycleLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		roe.mu.RLock()
		incidents := make([]*Incident, 0, len(roe.incidents))
		for _, inc := range roe.incidents {
			incidents = append(incidents, inc)
		}
		roe.mu.RUnlock()

		now := time.Now()
		for _, incident := range incidents {
			if incident.Status == "closed" {
				continue
			}
			if now.After(time.Unix(incident.SLADeadline, 0)) {
				roe.escalateIncident(incident, "sla_breach")
			}
			if (incident.Severity == "critical" || incident.Severity == "high") &&
				len(incident.Actions) == 0 &&
				time.Since(time.Unix(incident.CreatedAt, 0)) > 15*time.Minute {
				roe.escalateIncident(incident, "no_response_action")
			}
		}
	}
}

func (roe *ResponseOrchestrationEngine) escalateIncident(incident *Incident, reason string) {
	roe.emitEvidence("incident_escalated", map[string]string{
		"incident_id": incident.ID,
		"reason":      reason,
	})
}

func (roe *ResponseOrchestrationEngine) getSLAForSeverity(severity string) int {
	switch severity {
	case "critical":
		return 15
	case "high":
		return 30
	case "medium":
		return 60
	default:
		return 120
	}
}

func (roe *ResponseOrchestrationEngine) emitEvidence(action string, data interface{}) {
	record := orchestrator.EvidenceRecord{
		Layer:       "L6",
		Action:      action,
		Data:        data,
		Timestamp:   time.Now().Unix(),
		Attribution: "L6_RESPONSE_ORCHESTRATION",
	}
	_ = roe.bus.Publish("topic/evidence_records", record)
	if roe.store != nil {
		_ = roe.store.StoreEvidence(record)
	}
}
