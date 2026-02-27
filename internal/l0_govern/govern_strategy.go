package l0_govern

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

// MessageBus publishes tokens and evidence.
type MessageBus interface {
	Publish(topic string, payload interface{}) error
}

// EvidenceStore stores tokens for audit.
type EvidenceStore interface {
	Store(token *GovernanceToken) error
}

// SigningKey provides token signing.
type SigningKey interface {
	Sign(data []byte) ([]byte, error)
}

// GovernStrategy implements the L0 control loop.
type GovernStrategy struct {
	evidenceStore EvidenceStore
	messageBus    MessageBus
	signingKey    SigningKey
}

// NewGovernStrategy creates an L0 strategy.
func NewGovernStrategy(store EvidenceStore, bus MessageBus, key SigningKey) *GovernStrategy {
	return &GovernStrategy{
		evidenceStore: store,
		messageBus:    bus,
		signingKey:    key,
	}
}

// Run executes the governance loop at the specified period.
func (gs *GovernStrategy) Run(period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for range ticker.C {
		token, err := gs.generateGovernanceToken()
		if err != nil {
			gs.emitErrorEvidence("token_generation_failed", err)
			continue
		}

		if err := gs.validateInvariants(token); err != nil {
			gs.emitErrorEvidence("invariant_violation", err)
			continue
		}

		_ = gs.messageBus.Publish("topic/governance_tokens", token)
		_ = gs.evidenceStore.Store(token)
	}
}

// Tick produces one token immediately (for testing and Pathway A demo).
func (gs *GovernStrategy) Tick() (*GovernanceToken, error) {
	token, err := gs.generateGovernanceToken()
	if err != nil {
		return nil, err
	}
	if err := gs.validateInvariants(token); err != nil {
		return nil, err
	}
	_ = gs.messageBus.Publish("topic/governance_tokens", token)
	_ = gs.evidenceStore.Store(token)
	return token, nil
}

func (gs *GovernStrategy) generateGovernanceToken() (*GovernanceToken, error) {
	g := gs.updateGovernance()
	r := gs.updateRiskRegister()
	o, m := gs.deriveObjectivesAndTargets(g, r)

	token := &GovernanceToken{
		TokenID:       orchestrator.GenerateID(),
		IssuedAt:      time.Now().Unix(),
		ValidUntil:    time.Now().Add(30 * 24 * time.Hour).Unix(),
		Governance:    g,
		RiskRegister:  r,
		Objectives:    o,
		TargetMetrics: m,
		Issuer:        "L0_GOVERN_STRATEGY",
	}

	sig, err := gs.signToken(token)
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}
	token.Signature = sig

	return token, nil
}

func (gs *GovernStrategy) validateInvariants(token *GovernanceToken) error {
	for _, obj := range token.Objectives {
		if len(obj.LinkedRisks) == 0 {
			return fmt.Errorf("objective %s has no linked risks", obj.ID)
		}
	}
	for _, reg := range token.Governance.Regulations {
		if len(reg.ControlFamilies) == 0 {
			return fmt.Errorf("regulation %s has no control families", reg.ID)
		}
	}
	return nil
}

func (gs *GovernStrategy) signToken(token *GovernanceToken) ([]byte, error) {
	bytes, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return gs.signingKey.Sign(bytes)
}

func (gs *GovernStrategy) updateGovernance() GovernanceSet {
	return GovernanceSet{
		Regulations: []Regulation{
			{ID: "reg-1", Name: "Internal Security Policy", ControlFamilies: []string{"AC", "IA"}},
		},
		InternalPolicies: []Policy{{ID: "pol-1", Name: "MFA Required"}},
		Contracts:        []Contract{},
	}
}

func (gs *GovernStrategy) updateRiskRegister() RiskRegister {
	return RiskRegister{
		Version:   "1",
		LastUpdated: time.Now().Unix(),
		Scenarios: []RiskScenario{
			{
				ID:             "risk-1",
				Description:    "Credential abuse leading to data breach",
				Likelihood:     0.3,
				Impact:         0.8,
				Treatment:      "mitigate",
				LinkedControls: []string{"AC-2", "IA-2"},
			},
		},
	}
}

func (gs *GovernStrategy) deriveObjectivesAndTargets(g GovernanceSet, r RiskRegister) ([]Objective, Metrics) {
	objectives := []Objective{
		{
			ID:          "obj-1",
			Description: "Contain credential abuse within 15 minutes",
			TargetMTTD:  5 * time.Minute,
			TargetMTTR:  15 * time.Minute,
			LinkedRisks: []string{"risk-1"},
		},
	}
	metrics := Metrics{
		CoverageTarget:     0.9,
		MaxMTTD:            15 * time.Minute,
		MaxMTTR:            30 * time.Minute,
		MinAttributionRate: 1.0,
	}
	return objectives, metrics
}

func (gs *GovernStrategy) emitErrorEvidence(action string, err error) {
	record := orchestrator.EvidenceRecord{
		Layer:       "L0",
		Action:      action,
		Error:       err.Error(),
		Timestamp:   time.Now().Unix(),
		Attribution: "L0_GOVERN_STRATEGY",
	}
	_ = gs.messageBus.Publish("topic/evidence_records", record)
}
