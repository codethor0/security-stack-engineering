package l6_response

// Playbook defines a response workflow.
type Playbook struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	IncidentType    string          `json:"incident_type"`
	Phases          []PlaybookPhase `json:"phases"`
	AutomationLevel string          `json:"automation_level"`
	GovernanceReqs  []string        `json:"governance_requirements"`
}

// PlaybookPhase is a phase in a playbook.
type PlaybookPhase struct {
	Name         string           `json:"name"`
	Actions      []ActionTemplate `json:"actions"`
	ExitCriteria []string         `json:"exit_criteria"`
	SLAMinutes   int              `json:"sla_minutes"`
}

// ActionTemplate defines an action in a phase.
type ActionTemplate struct {
	ID               string            `json:"id"`
	Type             string            `json:"type"`
	Automation       bool              `json:"automation"`
	ApprovalRequired bool              `json:"approval_required"`
	Parameters       map[string]string `json:"parameters"`
}

// loadDefaultPlaybooks returns default incident response playbooks.
func loadDefaultPlaybooks() map[string]*Playbook {
	return map[string]*Playbook{
		"pb-credential-abuse": {
			ID:              "pb-credential-abuse",
			Name:            "Credential abuse response",
			IncidentType:    "credential_abuse",
			AutomationLevel: "partial",
			Phases: []PlaybookPhase{
				{
					Name:       "detection_analysis",
					SLAMinutes: 15,
					Actions: []ActionTemplate{
						{
							ID:               "act-triage",
							Type:             "triage",
							Automation:       false,
							ApprovalRequired: false,
						},
						{
							ID:               "act-disable-identity",
							Type:             "disable_identity",
							Automation:       true,
							ApprovalRequired: true,
						},
					},
					ExitCriteria: []string{"identity_disabled", "containment_confirmed"},
				},
			},
		},
		"pb-default": {
			ID:              "pb-default",
			Name:            "Default incident response",
			IncidentType:    "unknown",
			AutomationLevel: "manual",
			Phases: []PlaybookPhase{
				{
					Name:       "detection_analysis",
					SLAMinutes: 30,
					Actions: []ActionTemplate{
						{
							ID:               "act-triage",
							Type:             "triage",
							Automation:       false,
							ApprovalRequired: false,
						},
					},
					ExitCriteria: []string{"triage_complete"},
				},
			},
		},
	}
}
