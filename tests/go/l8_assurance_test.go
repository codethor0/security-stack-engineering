package tests

import (
	"testing"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/l8_assurance"
	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

func TestAssuranceEngineGenerateReport(t *testing.T) {
	persistence := orchestrator.NewInMemoryPersistence()
	bus := orchestrator.NewInMemoryBus(nil, persistence)

	_ = persistence.StoreEvidence(orchestrator.EvidenceRecord{
		Layer: "L4", Action: "task_completed", Timestamp: time.Now().Unix(),
	})

	engine := l8_assurance.NewAssuranceEngine(persistence, bus)
	report := engine.GenerateReport(
		time.Now().Add(-1*time.Hour),
		time.Now(),
	)

	if report.ReportID == "" {
		t.Error("ReportID empty")
	}
	if report.ReportType != "operational" {
		t.Errorf("ReportType: got %s", report.ReportType)
	}
	if report.EvidenceSummary.TotalRecords != 1 {
		t.Errorf("EvidenceSummary.TotalRecords: got %d", report.EvidenceSummary.TotalRecords)
	}
}
