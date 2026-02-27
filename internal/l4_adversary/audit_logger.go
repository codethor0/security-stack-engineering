package l4_adversary

import (
	"encoding/json"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

// AuditLogger produces tamper-evident audit records.
type AuditLogger struct {
	records []orchestrator.EvidenceRecord
}

// NewAuditLogger creates an audit logger.
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{records: make([]orchestrator.EvidenceRecord, 0)}
}

// LogEvent records a task execution event.
func (a *AuditLogger) LogEvent(task *Task, result interface{}) orchestrator.EvidenceRecord {
	data, _ := json.Marshal(map[string]interface{}{
		"task_id":   task.ID,
		"state":     task.State,
		"result":    result,
		"timestamp": time.Now().Unix(),
	})
	record := orchestrator.EvidenceRecord{
		Layer:        "L4",
		Action:       "task_execution",
		EngagementID: task.Engagement,
		TaskID:       task.ID,
		Operator:     task.Operator,
		TechniqueID:  task.TechniqueID,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		Attribution:  "L4_RTE_A",
	}
	a.records = append(a.records, record)
	return record
}
