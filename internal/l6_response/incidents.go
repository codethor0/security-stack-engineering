package l6_response

// Incident aggregates alerts and tracks response lifecycle.
type Incident struct {
	ID             string           `json:"id"`
	Title          string           `json:"title"`
	Severity       string           `json:"severity"`
	Status         string           `json:"status"`
	Phase          string           `json:"phase"`
	CreatedAt      int64            `json:"created_at"`
	UpdatedAt      int64            `json:"updated_at"`
	Alerts         []string         `json:"alert_ids"`
	Entities       []Entity         `json:"entities"`
	PlaybookID     string           `json:"playbook_id"`
	Actions        []ResponseAction `json:"actions"`
	EvidenceRefs    []string         `json:"evidence_refs"`
	AssignedTo     string           `json:"assigned_to"`
	SLADeadline    int64            `json:"sla_deadline"`
	EngagementLink string           `json:"engagement_link"`
}

// Entity represents an affected asset, identity, or data.
type Entity struct {
	Type            string `json:"type"`
	ID              string `json:"id"`
	Criticality     string `json:"criticality"`
	CompromiseState string `json:"compromise_state"`
}

// ResponseAction records a single response step.
type ResponseAction struct {
	ActionID    string `json:"action_id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	ExecutedAt  int64  `json:"executed_at"`
	ExecutedBy  string `json:"executed_by"`
	Status      string `json:"status"`
	Result      string `json:"result"`
}
