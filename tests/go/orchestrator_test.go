package tests

import (
	"testing"

	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

func TestBusPublishSubscribe(t *testing.T) {
	bus := orchestrator.NewInMemoryBus(nil, nil)

	ch := bus.Subscribe("topic/test")
	payload := map[string]string{"key": "value"}

	if err := bus.Publish("topic/test", payload); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	msg := <-ch
	if msg.Topic != "topic/test" {
		t.Errorf("topic: got %s", msg.Topic)
	}
	if len(msg.Payload) == 0 {
		t.Error("payload empty")
	}
}

func TestConsumeLatest(t *testing.T) {
	bus := orchestrator.NewInMemoryBus(nil, nil)

	payload := map[string]string{"id": "tok-1"}
	_ = bus.Publish("topic/governance_tokens", payload)

	latest := bus.ConsumeLatest("topic/governance_tokens")
	if latest == nil {
		t.Fatal("ConsumeLatest returned nil")
	}
}

func TestPersistence(t *testing.T) {
	p := orchestrator.NewInMemoryPersistence()

	record := orchestrator.EvidenceRecord{
		Layer:     "L0",
		Action:    "test",
		Timestamp: 12345,
	}
	if err := p.StoreEvidence(record); err != nil {
		t.Fatalf("StoreEvidence: %v", err)
	}

	all := p.AllEvidence()
	if len(all) != 1 {
		t.Errorf("AllEvidence: got %d records", len(all))
	}
	if all[0].Layer != "L0" {
		t.Errorf("Layer: got %s", all[0].Layer)
	}
}
