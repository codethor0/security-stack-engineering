package l2_identity

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

// AccessRequest is the normalized decision input for L2.
type AccessRequest struct {
	RequestID  string         `json:"request_id"`
	IdentityID string         `json:"identity_id"`
	ResourceID string         `json:"resource_id"`
	Action     string         `json:"action"`
	Context    RequestContext `json:"context"`
	Timestamp  int64          `json:"timestamp"`
}

// RequestContext holds request context.
type RequestContext struct {
	SourceIP    string   `json:"source_ip"`
	DeviceID    string   `json:"device_id"`
	MFAUsed     bool     `json:"mfa_used"`
	SessionRisk float64  `json:"session_risk"`
}

// AccessDecision is the result of zero trust evaluation.
type AccessDecision struct {
	RequestID     string  `json:"request_id"`
	Decision      string  `json:"decision"`
	Reason        string  `json:"reason"`
	PolicyApplied string  `json:"policy_applied"`
	IdentityRisk  float64 `json:"identity_risk"`
	Timestamp     int64   `json:"timestamp"`
	ExpiresAt     int64   `json:"expires_at"`
}

// IdentityRisk represents a detected identity risk.
type IdentityRisk struct {
	IdentityID        string   `json:"identity_id"`
	RiskType          string   `json:"risk_type"`
	Severity          string   `json:"severity"`
	Evidence          []string `json:"evidence"`
	RecommendedAction string   `json:"recommended_action"`
}

// MessageBus for L2.
type MessageBus interface {
	Publish(topic string, payload interface{}) error
	Subscribe(topic string) <-chan orchestrator.Message
	ConsumeLatest(topic string) interface{}
}

// EvidenceStore for L2.
type EvidenceStore interface {
	StoreEvidence(record orchestrator.EvidenceRecord) error
}

// ZeroTrustEngine encapsulates the L2 algorithm.
type ZeroTrustEngine struct {
	graph    *IdentityGraph
	policies []Policy
	bus      MessageBus
	store    EvidenceStore
	mu       sync.RWMutex
}

// NewZeroTrustEngine creates an L2 engine.
func NewZeroTrustEngine(bus MessageBus, store EvidenceStore) *ZeroTrustEngine {
	return &ZeroTrustEngine{
		graph:    NewIdentityGraph(),
		policies: loadDefaultPolicies(),
		bus:      bus,
		store:    store,
	}
}

// Run starts the L2 loops.
func (zte *ZeroTrustEngine) Run() {
	go zte.updateGraphLoop(5 * time.Minute)
	go zte.evaluateRequestLoop()
	go zte.riskAssessmentLoop(15 * time.Minute)
}

// SyncFromLatestSnapshot fetches the latest snapshot and updates the graph (for tests).
func (zte *ZeroTrustEngine) SyncFromLatestSnapshot() {
	zte.syncFromSnapshot(zte.bus.ConsumeLatest("topic/environment_snapshots"))
}

func (zte *ZeroTrustEngine) updateGraphLoop(interval time.Duration) {
	zte.syncFromSnapshot(zte.bus.ConsumeLatest("topic/environment_snapshots"))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		snapshot := zte.bus.ConsumeLatest("topic/environment_snapshots")
		if snapshot == nil {
			continue
		}
		zte.syncFromSnapshot(snapshot)
	}
}

func (zte *ZeroTrustEngine) syncFromSnapshot(snapshot interface{}) {
	zte.mu.Lock()
	defer zte.mu.Unlock()

	data, ok := snapshot.(map[string]interface{})
	if !ok {
		return
	}
	identities, ok := data["identities"].(map[string]interface{})
	if !ok {
		return
	}

	for id, raw := range identities {
		idMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		identity := Identity{
			ID:         id,
			Type:       getStr(idMap, "identity_type", "user"),
			Owner:      getStr(idMap, "owner", ""),
			Privileged: getBool(idMap, "privileged", false),
			Active:     true,
			Attributes: make(map[string]string),
		}
		if lastAuth := getStr(idMap, "last_auth", ""); lastAuth != "" {
			if t, err := time.Parse(time.RFC3339, lastAuth); err == nil {
				identity.LastAuth = t.Unix()
			}
		}
		zte.graph.Identities[id] = identity
	}
	zte.graph.LastUpdate = time.Now().Unix()
}

func (zte *ZeroTrustEngine) evaluateRequestLoop() {
	ch := zte.bus.Subscribe("topic/access_requests")
	for msg := range ch {
		var ar AccessRequest
		if err := json.Unmarshal(msg.Payload, &ar); err != nil {
			continue
		}
		decision := zte.EvaluateZeroTrust(ar)
		_ = zte.bus.Publish("topic/access_decisions", decision)
		zte.emitEvidence("access_decision", decision)
	}
}

// EvaluateZeroTrust runs zero trust policy evaluation (exported for tests).
func (zte *ZeroTrustEngine) EvaluateZeroTrust(req AccessRequest) AccessDecision {
	zte.mu.RLock()
	identity, exists := zte.graph.Identities[req.IdentityID]
	zte.mu.RUnlock()

	if !exists {
		return AccessDecision{
			RequestID: req.RequestID,
			Decision:  "deny",
			Reason:    "identity_not_found",
			Timestamp: time.Now().Unix(),
		}
	}

	if identity.Privileged && !req.Context.MFAUsed {
		return AccessDecision{
			RequestID:     req.RequestID,
			Decision:      "step_up",
			Reason:        "privileged_access_requires_mfa",
			PolicyApplied: "governance_mfa_mandate",
			IdentityRisk:  identity.RiskScore,
			Timestamp:     time.Now().Unix(),
		}
	}

	for _, policy := range zte.policies {
		if zte.matchesPolicy(req, identity, policy) {
			expiresAt := int64(0)
			if policy.Action == "allow" {
				expiresAt = time.Now().Add(1 * time.Hour).Unix()
			}
			return AccessDecision{
				RequestID:     req.RequestID,
				Decision:      policy.Action,
				Reason:        fmt.Sprintf("policy_match: %s", policy.Name),
				PolicyApplied: policy.ID,
				IdentityRisk:  identity.RiskScore,
				Timestamp:     time.Now().Unix(),
				ExpiresAt:     expiresAt,
			}
		}
	}

	return AccessDecision{
		RequestID: req.RequestID,
		Decision:  "deny",
		Reason:    "no_matching_policy",
		Timestamp: time.Now().Unix(),
	}
}

func (zte *ZeroTrustEngine) matchesPolicy(req AccessRequest, identity Identity, policy Policy) bool {
	for _, c := range policy.Conditions {
		switch c.Attribute {
		case "privileged":
			if c.Value == true && !identity.Privileged {
				return false
			}
		case "mfa_used":
			if c.Value == true && !req.Context.MFAUsed {
				return false
			}
		case "owner":
			if c.Value == "" && identity.Owner != "" {
				return false
			}
		case "identity_found":
			return true
		}
	}
	return true
}

func (zte *ZeroTrustEngine) riskAssessmentLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		risks := zte.assessIdentityRisks()
		for _, risk := range risks {
			_ = zte.bus.Publish("topic/identity_risks", risk)
			zte.emitEvidence("identity_risk_detected", risk)
		}
	}
}

func (zte *ZeroTrustEngine) assessIdentityRisks() []IdentityRisk {
	zte.mu.RLock()
	defer zte.mu.RUnlock()

	var risks []IdentityRisk
	for id, identity := range zte.graph.Identities {
		if identity.Privileged && identity.Owner == "" {
			risks = append(risks, IdentityRisk{
				IdentityID:        id,
				RiskType:          "orphaned_privileged_identity",
				Severity:          "critical",
				Evidence:          []string{"no_owner_attribute", "privileged=true"},
				RecommendedAction: "disable_account_pending_investigation",
			})
		}
		if identity.LastAuth > 0 {
			lastAuth := time.Unix(identity.LastAuth, 0)
			if time.Since(lastAuth) > 90*24*time.Hour && identity.Active {
				risks = append(risks, IdentityRisk{
					IdentityID:        id,
					RiskType:          "stale_credentials",
					Severity:          "medium",
					Evidence:          []string{fmt.Sprintf("last_auth: %s", lastAuth.Format(time.RFC3339))},
					RecommendedAction: "credential_rotation",
				})
			}
		}
	}
	return risks
}

func (zte *ZeroTrustEngine) emitEvidence(action string, data interface{}) {
	record := orchestrator.EvidenceRecord{
		Layer:       "L2",
		Action:      action,
		Data:        data,
		Timestamp:   time.Now().Unix(),
		Attribution: "L2_IDENTITY_ACCESS",
	}
	_ = zte.bus.Publish("topic/evidence_records", record)
	if zte.store != nil {
		_ = zte.store.StoreEvidence(record)
	}
}

func getStr(m map[string]interface{}, k, def string) string {
	if v, ok := m[k]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func getBool(m map[string]interface{}, k string, def bool) bool {
	if v, ok := m[k]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}
