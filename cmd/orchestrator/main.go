// Command orchestrator starts the SSE message bus and wires layers.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/codethor0/security-stack-engineering/internal/l0_govern"
	"github.com/codethor0/security-stack-engineering/internal/l2_identity"
	"github.com/codethor0/security-stack-engineering/internal/l4_adversary"
	"github.com/codethor0/security-stack-engineering/internal/l5_detection"
	"github.com/codethor0/security-stack-engineering/internal/l6_response"
	"github.com/codethor0/security-stack-engineering/internal/l8_assurance"
	"github.com/codethor0/security-stack-engineering/internal/orchestrator"
)

func main() {
	pathway := flag.String("pathway", "a", "Pathway: a, b, or c")
	configPath := flag.String("config", "", "Path to config JSON (overrides -pathway)")
	flag.Parse()

	persistence := orchestrator.NewInMemoryPersistence()
	bus := orchestrator.NewInMemoryBus(nil, persistence)

	path := strings.ToLower(*pathway)
	if *configPath != "" {
		cfg := orchestrator.LoadConfigOrDefault(*configPath)
		path = strings.ToLower(cfg.Pathway)
	}

	switch path {
	case "a":
		runPathwayA(bus, persistence)
	case "b":
		runPathwayB(bus, persistence)
	case "c":
		runPathwayC(bus, persistence)
	default:
		runPathwayA(bus, persistence)
	}
}

func runPathwayA(bus *orchestrator.InMemoryBus, persistence *orchestrator.InMemoryPersistence) {
	fmt.Println("SSE Pathway A: L0 + L4 + L5 + L8")

	signer, err := l0_govern.NewEd25519Signer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "signer: %v\n", err)
		os.Exit(1)
	}

	tokenStore := l0_govern.NewGovTokenStore()
	l0 := l0_govern.NewGovernStrategy(tokenStore, bus, signer)

	if _, err := l0.Tick(); err != nil {
		fmt.Fprintf(os.Stderr, "L0 Tick: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("L0: Governance token published")

	rtea, err := l4_adversary.NewRTEAEngine(bus, persistence)
	if err != nil {
		fmt.Fprintf(os.Stderr, "L4: %v\n", err)
		os.Exit(1)
	}

	l5Bridge := l5_detection.NewDetectionBridge(bus)
	go l5Bridge.Run()
	fmt.Println("L5: Detection bridge running (listens for L4 evidence)")

	scenario := l4_adversary.Scenario{
		EngagementID: "eng-demo-001",
		OperatorID:   "op-demo",
		ApproverID:   "approver-demo",
		Techniques:   []string{"T1078"},
	}

	if err := rtea.Run(scenario); err != nil {
		fmt.Fprintf(os.Stderr, "L4 Run: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("L4: Engagement completed, findings published")

	time.Sleep(200 * time.Millisecond)

	l8 := l8_assurance.NewAssuranceEngine(persistence, bus)
	report := l8.GenerateReport(time.Now().Add(-1*time.Hour), time.Now())
	fmt.Printf("L8: Assurance report generated (id=%s, evidence=%d)\n", report.ReportID, report.EvidenceSummary.TotalRecords)

	waitForShutdown()
}

func runPathwayB(bus *orchestrator.InMemoryBus, persistence *orchestrator.InMemoryPersistence) {
	fmt.Println("SSE Pathway B: L0 + L1(stub) + L2 + L3(stub) + L4 + L5 + L6 + L8")

	signer, err := l0_govern.NewEd25519Signer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "signer: %v\n", err)
		os.Exit(1)
	}

	tokenStore := l0_govern.NewGovTokenStore()
	l0 := l0_govern.NewGovernStrategy(tokenStore, bus, signer)
	if _, err := l0.Tick(); err != nil {
		fmt.Fprintf(os.Stderr, "L0 Tick: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("L0: Governance token published")

	snapshot := map[string]interface{}{
		"snapshot_id": orchestrator.GenerateID(),
		"timestamp":   time.Now().Format(time.RFC3339),
		"identities": map[string]interface{}{
			"id-1": map[string]interface{}{
				"identity_type": "user",
				"owner":         "admin",
				"privileged":    true,
				"last_auth":     time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			},
		},
	}
	_ = bus.Publish("topic/environment_snapshots", snapshot)
	fmt.Println("L1: Stub environment snapshot published")

	l2 := l2_identity.NewZeroTrustEngine(bus, persistence)
	go l2.Run()
	fmt.Println("L2: Zero trust engine running")

	l5Bridge := l5_detection.NewDetectionBridge(bus)
	go l5Bridge.Run()
	fmt.Println("L5: Detection bridge running")

	l6 := l6_response.NewResponseEngine(bus, persistence)
	go l6.Run()
	fmt.Println("L6: Response orchestration running")

	rtea, err := l4_adversary.NewRTEAEngine(bus, persistence)
	if err != nil {
		fmt.Fprintf(os.Stderr, "L4: %v\n", err)
		os.Exit(1)
	}

	scenario := l4_adversary.Scenario{
		EngagementID: "eng-demo-001",
		OperatorID:   "op-demo",
		ApproverID:   "approver-demo",
		Techniques:   []string{"T1078"},
	}

	if err := rtea.Run(scenario); err != nil {
		fmt.Fprintf(os.Stderr, "L4 Run: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("L4: Engagement completed")

	time.Sleep(500 * time.Millisecond)

	l8 := l8_assurance.NewAssuranceEngine(persistence, bus)
	report := l8.GenerateReport(time.Now().Add(-1*time.Hour), time.Now())
	fmt.Printf("L8: Assurance report generated (id=%s, evidence=%d)\n", report.ReportID, report.EvidenceSummary.TotalRecords)

	waitForShutdown()
}

func runPathwayC(bus *orchestrator.InMemoryBus, persistence *orchestrator.InMemoryPersistence) {
	fmt.Println("SSE Pathway C: Full stack L0-L8 (Python layers stubbed)")

	runPathwayB(bus, persistence)
}

func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press Ctrl+C to exit")
	<-sigCh
	fmt.Println("Shutdown")
}
