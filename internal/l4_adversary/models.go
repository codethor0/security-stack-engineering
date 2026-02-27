package l4_adversary

// TaskType identifies the kind of offensive action.
type TaskType string

const (
	TaskSimulateLogin   TaskType = "simulate_login"
	TaskSimulateBeacon  TaskType = "simulate_beacon"
	TaskLateralMovement TaskType = "lateral_movement"
	TaskDataExfil       TaskType = "data_exfiltration"
	TaskPersistence     TaskType = "persistence"
)

// TaskState represents task lifecycle state.
type TaskState string

const (
	StatePending   TaskState = "pending"
	StateExecuting TaskState = "executing"
	StateCancelled TaskState = "cancelled"
	StateCompleted TaskState = "completed"
	StateFailed    TaskState = "failed"
)

// Task represents one atomic offensive action.
type Task struct {
	ID          string          `json:"id"`
	Engagement  string          `json:"engagement_id"`
	Type        TaskType        `json:"type"`
	CreatedAt   int64           `json:"created_at"`
	TTLSeconds  int64           `json:"ttl_seconds"`
	Operator    string          `json:"operator"`
	ApprovedBy  string          `json:"approved_by"`
	State       TaskState       `json:"state"`
	TechniqueID string          `json:"technique_id"`
	Params      []byte          `json:"params"`
	CancelToken string          `json:"cancel_token,omitempty"`
}

// SignedTask provides cryptographic attribution.
type SignedTask struct {
	Task      Task   `json:"task"`
	PublicKey []byte `json:"public_key"`
	Signature []byte `json:"signature"`
}

// Scenario describes an engagement to run.
type Scenario struct {
	EngagementID string              `json:"engagement_id"`
	OperatorID   string              `json:"operator_id"`
	ApproverID   string              `json:"approver_id"`
	Targets      []string            `json:"targets"`
	Techniques   []string            `json:"techniques"`
	Params       map[string][]byte   `json:"params"`
}

// OffensiveFindings summarizes engagement results.
type OffensiveFindings struct {
	TechniquesTested    []string               `json:"techniques_tested"`
	DetectionsTriggered map[string]int          `json:"detections_triggered"`
	MissedDetections    []string                `json:"missed_detections"`
	ResponseLatency     map[string]int64        `json:"response_latency"`
	EvidenceQuality     string                  `json:"evidence_quality"`
}
