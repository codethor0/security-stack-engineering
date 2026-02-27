package l0_govern

import "time"

// GovernanceToken is the canonical output of L0, consumed by all downstream layers.
type GovernanceToken struct {
	TokenID      string        `json:"token_id"`
	IssuedAt     int64         `json:"issued_at"`
	ValidUntil   int64         `json:"valid_until"`
	Governance   GovernanceSet `json:"governance"`
	RiskRegister RiskRegister  `json:"risk_register"`
	Objectives   []Objective   `json:"objectives"`
	TargetMetrics Metrics      `json:"target_metrics"`
	Signature    []byte        `json:"signature"`
	Issuer       string        `json:"issuer"`
}

// GovernanceSet captures external obligations and internal policy.
type GovernanceSet struct {
	Regulations      []Regulation `json:"regulations"`
	InternalPolicies []Policy     `json:"internal_policies"`
	Contracts        []Contract   `json:"contracts"`
}

// Regulation represents a compliance obligation.
type Regulation struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	ControlFamilies []string `json:"control_families"`
}

// Policy is an internal policy.
type Policy struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Contract is a contractual obligation.
type Contract struct {
	ID string `json:"id"`
}

// RiskRegister maintains the current organizational risk posture.
type RiskRegister struct {
	Version     string         `json:"version"`
	Scenarios   []RiskScenario `json:"scenarios"`
	LastUpdated int64          `json:"last_updated"`
}

// RiskScenario is a single risk entry.
type RiskScenario struct {
	ID             string   `json:"id"`
	Description    string   `json:"description"`
	Likelihood     float64  `json:"likelihood"`
	Impact         float64  `json:"impact"`
	Treatment      string   `json:"treatment"`
	LinkedControls []string `json:"linked_controls"`
}

// Objective is a concrete, measurable security goal.
type Objective struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	TargetMTTD  time.Duration `json:"target_mttd"`
	TargetMTTR  time.Duration `json:"target_mttr"`
	LinkedRisks []string      `json:"linked_risks"`
}

// Metrics defines quantitative targets for security posture.
type Metrics struct {
	CoverageTarget     float64       `json:"coverage_target"`
	MaxMTTD            time.Duration `json:"max_mttd"`
	MaxMTTR            time.Duration `json:"max_mttr"`
	MinAttributionRate float64       `json:"min_attribution_rate"`
}
