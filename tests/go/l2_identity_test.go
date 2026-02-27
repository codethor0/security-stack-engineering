package tests

import (
	"testing"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/l2_identity"
	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

func TestZeroTrustEngineEvaluate(t *testing.T) {
	bus := orchestrator.NewInMemoryBus(nil, nil)
	persistence := orchestrator.NewInMemoryPersistence()

	snapshot := map[string]interface{}{
		"identities": map[string]interface{}{
			"id-1": map[string]interface{}{
				"identity_type": "user",
				"owner":         "admin",
				"privileged":    true,
				"last_auth":     time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			},
		},
	}
	_ = bus.Publish("topic/environment_snapshots", snapshot)

	zte := l2_identity.NewZeroTrustEngine(bus, persistence)
	zte.SyncFromLatestSnapshot()

	req := l2_identity.AccessRequest{
		RequestID:  "req-1",
		IdentityID: "id-1",
		Action:     "read",
		Context:    l2_identity.RequestContext{MFAUsed: true},
		Timestamp:  time.Now().Unix(),
	}
	decision := zte.EvaluateZeroTrust(req)

	if decision.Decision != "allow" && decision.Decision != "step_up" {
		t.Errorf("expected allow or step_up, got %s", decision.Decision)
	}
	if decision.RequestID != "req-1" {
		t.Errorf("RequestID: got %s", decision.RequestID)
	}
}

func TestZeroTrustEngineDenyUnknownIdentity(t *testing.T) {
	bus := orchestrator.NewInMemoryBus(nil, nil)
	persistence := orchestrator.NewInMemoryPersistence()
	zte := l2_identity.NewZeroTrustEngine(bus, persistence)

	req := l2_identity.AccessRequest{
		RequestID:  "req-1",
		IdentityID: "unknown-id",
		Timestamp:  time.Now().Unix(),
	}
	decision := zte.EvaluateZeroTrust(req)

	if decision.Decision != "deny" {
		t.Errorf("expected deny for unknown identity, got %s", decision.Decision)
	}
	if decision.Reason != "identity_not_found" {
		t.Errorf("Reason: got %s", decision.Reason)
	}
}

func TestZeroTrustEngineExportedEvaluate(t *testing.T) {
	bus := orchestrator.NewInMemoryBus(nil, nil)
	zte := l2_identity.NewZeroTrustEngine(bus, nil)
	req := l2_identity.AccessRequest{RequestID: "r1", IdentityID: "x", Timestamp: 1}
	_ = zte.EvaluateZeroTrust(req)
}
