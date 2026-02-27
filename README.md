# Security Stack Engineering (SSE)

[![Go Tests](https://github.com/codethor0/security-stack-engineering/actions/workflows/go-test.yml/badge.svg)](https://github.com/codethor0/security-stack-engineering/actions/workflows/go-test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A reference implementation of the Security Stack Engineering (SSE) architecture, a multi-algorithm system for governed security operations. The concept and architecture are authored by **Thor Thor** ([codethor0](https://github.com/codethor0)).

This repository provides a proof-of-concept implementation with nine cooperating layers (L0–L8) that form a governed, measurable security system. The design may be subject to future patent filings; please consult the repository owner for licensing or commercial use questions.

**Documentation:** [Substack Article](https://uniqueviolation.substack.com/p/security-stack-engineering)

## Architecture Overview

SSE splits security operations into nine algorithms, each with explicit inputs, outputs, and invariants:

| Layer | Name | Language | Description |
|-------|------|----------|-------------|
| L0 | Govern & Strategy | Go | GovernanceToken, risk register, objectives, target metrics |
| L1 | Asset & Attack Surface | Python | EnvironmentSnapshot, asset gaps |
| L2 | Identity & Access | Go | IdentityGraph, zero trust policies |
| L3 | Telemetry & Data Fabric | Python | NormalizedEvent, TelemetryView, TelemetryGap |
| L4 | Adversary Simulation (RTE-A) | Go | OffensivePlan, OffensiveFindings, signed tasks |
| L5 | Detection & Analytics | Python | DetectionRule, Alert, CoverageMatrix |
| L6 | Response & Orchestration | Go | Incident, Playbook, ResponseAction |
| L7 | Engineering & Hardening | Python | EngineeringChange, HardeningProgress |
| L8 | Assurance & Audit | Go | AssuranceReport, Attestation, GovernanceFeedback |

Layers communicate via a shared message bus with typed topics and JSON schemas.

## Adoption Pathways

**Pathway A (Minimum Viable):** L0 + L4 + L5 + L8. Governed adversary simulation with detection bridge and audit trail. No Python required.

**Pathway B (Core Defense):** Pathway A + L2 (Identity) + L6 (Response). Adds zero trust engine and incident response orchestration. L1/L3 stubbed via environment snapshot.

**Pathway C (Full Stack):** Same as Pathway B. Python layers (L1, L3, L5, L7) run standalone with MockBus via `scripts/run_python_demo.sh`.

## Quick Start

### Prerequisites

- Go 1.21+
- Python 3.11+
- (Optional) Python virtual environment

### Pathway A (L0, L4, L5, L8)

```bash
./scripts/run_pathway_a.sh
```

Or manually:

```bash
go run ./cmd/orchestrator/main.go -pathway=a
```

### Pathway B (adds L2, L6)

```bash
./scripts/run_pathway_b.sh
```

### Pathway C (full L0–L8)

```bash
./scripts/run_pathway_c.sh
```

## Hello World Scenario

A minimal end-to-end flow demonstrates:

1. **L0** produces a GovernanceToken with one objective ("contain credential abuse within 15 min")
2. **L4** runs an RTE-A engagement with one technique (T1078) and one signed task
3. **L5** listens for offensive plans and emits a detection when the task executes
4. **L8** generates an assurance summary with evidence and posture metrics

See `examples/minimal_pathway_a/` for configuration and `docs/architecture.md` for message flow.

## Repository Structure

```
security-stack-engineering/
  cmd/           # Entry points (orchestrator, standalone layer runners)
  internal/      # Go control-plane packages
  python/        # Python data-plane packages
  schemas/       # JSON schemas for message types
  docs/          # Architecture and layer documentation
  tests/         # Go and Python tests
  examples/      # Pathway configs and runners
  scripts/       # Bootstrap and run scripts
```

## Python Layer Demo

Run Python layers standalone (with MockBus) to exercise the algorithm logic:

```bash
PYTHONPATH=. ./scripts/run_python_demo.sh
```

## License

MIT License. See [LICENSE](LICENSE).

## Contact

- GitHub: [codethor0](https://github.com/codethor0)
- Email: codethor@gmail.com
