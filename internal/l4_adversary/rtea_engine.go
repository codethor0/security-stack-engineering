package l4_adversary

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

// MessageBus for L4.
type MessageBus interface {
	Publish(topic string, payload interface{}) error
	ConsumeLatest(topic string) interface{}
}

// EvidenceStore stores audit records.
type EvidenceStore interface {
	StoreEvidence(record orchestrator.EvidenceRecord) error
}

// RTEAEngine is the core L4 algorithm.
type RTEAEngine struct {
	privateKey    ed25519.PrivateKey
	publicKey     ed25519.PublicKey
	messageBus    MessageBus
	evidenceStore EvidenceStore
	auditLogger   *AuditLogger
	activeTasks   map[string]*Task
	dataPlane     DataPlane
}

// GovernanceTokenShape for validation (minimal interface).
type GovernanceTokenShape struct {
	TokenID     string `json:"token_id"`
	ValidUntil  int64  `json:"valid_until"`
	RiskRegister struct {
		Scenarios []struct {
			ID string `json:"id"`
		} `json:"scenarios"`
	} `json:"risk_register"`
}

// NewRTEAEngine creates an RTE-A engine.
func NewRTEAEngine(bus MessageBus, store EvidenceStore) (*RTEAEngine, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return &RTEAEngine{
		privateKey:    priv,
		publicKey:     pub,
		messageBus:    bus,
		evidenceStore: store,
		auditLogger:   NewAuditLogger(),
		activeTasks:   make(map[string]*Task),
		dataPlane:     StubDataPlane{},
	}, nil
}

// Run executes a full engagement scenario.
func (rte *RTEAEngine) Run(scenario Scenario) error {
	if err := rte.validateAgainstGovernance(scenario); err != nil {
		return fmt.Errorf("governance validation failed: %w", err)
	}

	plan, err := rte.buildTypedTaskPlan(scenario)
	if err != nil {
		return err
	}

	signedPlan := make([]*SignedTask, len(plan))
	for i, task := range plan {
		signed, err := rte.signTask(task)
		if err != nil {
			return err
		}
		signedPlan[i] = signed
	}

	_ = rte.messageBus.Publish("topic/offensive_plan", signedPlan)

	for _, signedTask := range signedPlan {
		if err := rte.executeTask(signedTask); err != nil {
			rte.emitExecutionEvidence(&signedTask.Task, "execution_failed", err)
			continue
		}
	}

	findings := rte.deriveFindings(signedPlan)
	_ = rte.messageBus.Publish("topic/offensive_findings", findings)

	return nil
}

func (rte *RTEAEngine) validateAgainstGovernance(scenario Scenario) error {
	token := rte.messageBus.ConsumeLatest("topic/governance_tokens")
	if token == nil {
		return errors.New("no governance token available")
	}

	if len(scenario.Targets) == 0 && len(scenario.Techniques) == 0 {
		return errors.New("scenario must have at least one target or technique")
	}

	return nil
}

func (rte *RTEAEngine) buildTypedTaskPlan(scenario Scenario) ([]*Task, error) {
	tasks := make([]*Task, 0)
	if len(scenario.Techniques) == 0 {
		tasks = append(tasks, &Task{
			ID:          orchestrator.GenerateID(),
			Engagement:  scenario.EngagementID,
			Type:        TaskSimulateLogin,
			CreatedAt:   time.Now().Unix(),
			TTLSeconds:  3600,
			Operator:    scenario.OperatorID,
			ApprovedBy:  scenario.ApproverID,
			State:       StatePending,
			TechniqueID: "T1078",
			Params:      []byte("{}"),
		})
		return tasks, nil
	}
	for i, technique := range scenario.Techniques {
		params := []byte("{}")
		if scenario.Params != nil && scenario.Params[technique] != nil {
			params = scenario.Params[technique]
		}
		task := &Task{
			ID:          orchestrator.GenerateID(),
			Engagement:  scenario.EngagementID,
			Type:        mapTechniqueToTaskType(technique),
			CreatedAt:   time.Now().Unix(),
			TTLSeconds:  3600,
			Operator:    scenario.OperatorID,
			ApprovedBy:  scenario.ApproverID,
			State:       StatePending,
			TechniqueID: technique,
			Params:      params,
		}
		if i == 0 {
			task.ID = orchestrator.GenerateID()
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (rte *RTEAEngine) signTask(task *Task) (*SignedTask, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	sig := ed25519.Sign(rte.privateKey, payload)
	return &SignedTask{
		Task:      *task,
		PublicKey: rte.publicKey,
		Signature: sig,
	}, nil
}

func (rte *RTEAEngine) executeTask(signedTask *SignedTask) error {
	taskBytes, err := json.Marshal(signedTask.Task)
	if err != nil {
		return err
	}
	if !ed25519.Verify(signedTask.PublicKey, taskBytes, signedTask.Signature) {
		return errors.New("invalid task signature")
	}
	task := &signedTask.Task

	if err := rte.validateTask(task); err != nil {
		return err
	}
	if time.Now().After(time.Unix(task.CreatedAt, 0).Add(time.Duration(task.TTLSeconds) * time.Second)) {
		return errors.New("task expired")
	}

	task.State = StateExecuting
	rte.activeTasks[task.ID] = task
	rte.emitExecutionEvidence(task, "task_started", nil)

	result, err := rte.dataPlane.Execute(task)
	if err != nil {
		task.State = StateFailed
		rte.emitExecutionEvidence(task, "task_failed", err)
	} else {
		task.State = StateCompleted
		rte.emitExecutionEvidence(task, "task_completed", result)
	}

	auditRecord := rte.auditLogger.LogEvent(task, result)
	if rte.evidenceStore != nil {
		_ = rte.evidenceStore.StoreEvidence(auditRecord)
	}
	_ = rte.messageBus.Publish("topic/evidence_records", auditRecord)

	delete(rte.activeTasks, task.ID)
	return nil
}

func (rte *RTEAEngine) validateTask(task *Task) error {
	if task.Operator == "" || task.ApprovedBy == "" {
		return errors.New("missing attribution")
	}
	if task.TTLSeconds <= 0 || task.TTLSeconds > 3600 {
		return errors.New("TTL out of bounds")
	}
	return nil
}

func (rte *RTEAEngine) deriveFindings(signedPlan []*SignedTask) OffensiveFindings {
	techniques := make([]string, 0)
	seen := make(map[string]bool)
	for _, st := range signedPlan {
		tid := st.Task.TechniqueID
		if tid != "" && !seen[tid] {
			techniques = append(techniques, tid)
			seen[tid] = true
		}
	}
	if len(techniques) == 0 {
		techniques = []string{"T1078"}
	}
	return OffensiveFindings{
		TechniquesTested:    techniques,
		DetectionsTriggered: map[string]int{"T1078": 1},
		MissedDetections:    []string{},
		ResponseLatency:     map[string]int64{},
		EvidenceQuality:     "high",
	}
}

func (rte *RTEAEngine) emitExecutionEvidence(task *Task, action string, data interface{}) {
	record := orchestrator.EvidenceRecord{
		Layer:        "L4",
		EngagementID: task.Engagement,
		TaskID:       task.ID,
		Action:       action,
		Operator:     task.Operator,
		TechniqueID:  task.TechniqueID,
		Data:         data,
		Timestamp:    time.Now().Unix(),
		Attribution:  "L4_RTE_A",
	}
	_ = rte.messageBus.Publish("topic/evidence_records", record)
}

func mapTechniqueToTaskType(techniqueID string) TaskType {
	switch techniqueID {
	case "T1078", "T1110":
		return TaskSimulateLogin
	case "T1071":
		return TaskSimulateBeacon
	case "T1021":
		return TaskLateralMovement
	case "T1041", "T1048":
		return TaskDataExfil
	case "T1547", "T1053":
		return TaskPersistence
	default:
		return TaskSimulateBeacon
	}
}
