package tests

import (
	"testing"

	"github.com/codethor0/security-stack-engineering/internal/l0_govern"
	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

func TestGovernStrategyTick(t *testing.T) {
	store := l0_govern.NewGovTokenStore()
	bus := orchestrator.NewInMemoryBus(nil, nil)
	signer, err := l0_govern.NewEd25519Signer()
	if err != nil {
		t.Fatalf("NewEd25519Signer: %v", err)
	}

	strategy := l0_govern.NewGovernStrategy(store, bus, signer)
	token, err := strategy.Tick()
	if err != nil {
		t.Fatalf("Tick: %v", err)
	}

	if token.TokenID == "" {
		t.Error("TokenID empty")
	}
	if token.Issuer != "L0_GOVERN_STRATEGY" {
		t.Errorf("Issuer: got %s", token.Issuer)
	}
	if len(token.Objectives) == 0 {
		t.Error("Objectives empty")
	}
	if len(token.Objectives[0].LinkedRisks) == 0 {
		t.Error("Objective has no linked risks")
	}
	if len(token.Signature) == 0 {
		t.Error("Signature empty")
	}
}

func TestInvariantValidation(t *testing.T) {
	store := l0_govern.NewGovTokenStore()
	bus := orchestrator.NewInMemoryBus(nil, nil)
	signer, _ := l0_govern.NewEd25519Signer()
	strategy := l0_govern.NewGovernStrategy(store, bus, signer)

	token, err := strategy.Tick()
	if err != nil {
		t.Fatal(err)
	}

	for _, obj := range token.Objectives {
		if len(obj.LinkedRisks) == 0 {
			t.Errorf("objective %s has no linked risks", obj.ID)
		}
	}
	for _, reg := range token.Governance.Regulations {
		if len(reg.ControlFamilies) == 0 {
			t.Errorf("regulation %s has no control families", reg.ID)
		}
	}
}
