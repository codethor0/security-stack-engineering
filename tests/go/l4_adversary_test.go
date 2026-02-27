package tests

import (
	"strings"
	"testing"

	"github.com/codethor0/security-stack-engineering/internal/l4_adversary"
	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

func TestRTEAEngineRun(t *testing.T) {
	bus := orchestrator.NewInMemoryBus(nil, nil)
	persistence := orchestrator.NewInMemoryPersistence()

	// Publish governance token first
	_ = bus.Publish("topic/governance_tokens", map[string]string{"token_id": "tok-1"})

	rtea, err := l4_adversary.NewRTEAEngine(bus, persistence)
	if err != nil {
		t.Fatalf("NewRTEAEngine: %v", err)
	}

	scenario := l4_adversary.Scenario{
		EngagementID: "eng-test-001",
		OperatorID:   "op-test",
		ApproverID:   "approver-test",
		Techniques:   []string{"T1078"},
	}

	if err := rtea.Run(scenario); err != nil {
		t.Fatalf("Run: %v", err)
	}

	latest := bus.ConsumeLatest("topic/offensive_findings")
	if latest == nil {
		t.Fatal("offensive_findings not published")
	}
}

func TestRTEAEngineRequiresGovernance(t *testing.T) {
	bus := orchestrator.NewInMemoryBus(nil, nil)
	persistence := orchestrator.NewInMemoryPersistence()
	// Do NOT publish governance token

	rtea, err := l4_adversary.NewRTEAEngine(bus, persistence)
	if err != nil {
		t.Fatalf("NewRTEAEngine: %v", err)
	}

	scenario := l4_adversary.Scenario{
		EngagementID: "eng-test",
		OperatorID:   "op",
		ApproverID:   "approver",
		Techniques:  []string{"T1078"},
	}

	err = rtea.Run(scenario)
	if err == nil {
		t.Error("expected error when no governance token")
	}
	if err != nil && !strings.Contains(err.Error(), "governance") {
		t.Errorf("unexpected error: %v", err)
	}
}
