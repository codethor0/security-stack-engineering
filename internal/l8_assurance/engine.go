package l8_assurance

import (
	"time"

	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

// MessageBus for L8.
type MessageBus interface {
	Publish(topic string, payload interface{}) error
	ConsumeLatest(topic string) interface{}
}

// EvidenceStore provides evidence queries.
type EvidenceStore interface {
	AllEvidence() []orchestrator.EvidenceRecord
}

// AssuranceEngine generates reports and attestations.
type AssuranceEngine struct {
	evidenceStore EvidenceStore
	messageBus    MessageBus
}

// AssuranceReport is the canonical L8 output.
type AssuranceReport struct {
	ReportID        string                 `json:"report_id"`
	ReportType      string                 `json:"report_type"`
	PeriodStart     string                 `json:"period_start"`
	PeriodEnd       string                 `json:"period_end"`
	GeneratedAt     int64                  `json:"generated_at"`
	GeneratedBy     string                 `json:"generated_by"`
	Sections        []ReportSection        `json:"sections"`
	EvidenceSummary EvidenceSummary        `json:"evidence_summary"`
}

// ReportSection is a section of a report.
type ReportSection struct {
	Title        string      `json:"title"`
	Type         string      `json:"type"`
	Content      interface{} `json:"content"`
	EvidenceRefs []string     `json:"evidence_refs"`
}

// EvidenceSummary summarizes evidence for a period.
type EvidenceSummary struct {
	TotalRecords int            `json:"total_records"`
	ByLayer      map[string]int `json:"by_layer"`
	TimeRange    TimeRange      `json:"time_range"`
}

// TimeRange is a time window.
type TimeRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// NewAssuranceEngine creates an L8 engine.
func NewAssuranceEngine(store EvidenceStore, bus MessageBus) *AssuranceEngine {
	return &AssuranceEngine{
		evidenceStore: store,
		messageBus:    bus,
	}
}

// GenerateReport produces an assurance report for the given period.
func (e *AssuranceEngine) GenerateReport(periodStart, periodEnd time.Time) AssuranceReport {
	evidence := e.evidenceStore.AllEvidence()

	byLayer := make(map[string]int)
	for _, r := range evidence {
		byLayer[r.Layer]++
	}

	report := AssuranceReport{
		ReportID:    orchestrator.GenerateID(),
		ReportType:  "operational",
		PeriodStart: periodStart.Format(time.RFC3339),
		PeriodEnd:   periodEnd.Format(time.RFC3339),
		GeneratedAt: time.Now().Unix(),
		GeneratedBy: "L8_ASSURANCE_AUDIT",
		Sections: []ReportSection{
			{
				Title:   "Executive Summary",
				Type:    "summary",
				Content: map[string]interface{}{
					"engagement_summary": "Governed adversary simulation completed with full attribution.",
					"detection_status":   "Alerts generated for simulated techniques.",
					"evidence_count":     len(evidence),
				},
			},
			{
				Title:   "Evidence by Layer",
				Type:    "metrics",
				Content: byLayer,
			},
		},
		EvidenceSummary: EvidenceSummary{
			TotalRecords: len(evidence),
			ByLayer:      byLayer,
			TimeRange: TimeRange{
				Start: periodStart.Unix(),
				End:   periodEnd.Unix(),
			},
		},
	}

	_ = e.messageBus.Publish("topic/assurance_reports", report)
	return report
}
