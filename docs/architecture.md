# SSE Architecture

## End-to-End Message Flow

```
                    +------------------+
                    |       L0         |
                    | Govern & Strategy|
                    +--------+---------+
                             | topic/governance_tokens
                             v
    +------------+   +-------+--------+   +------------+
    |     L1     |   |    ORCHESTRATOR   |   |     L2     |
    | Asset/Surf |<->|   (message bus)   |<->| Identity   |
    +------+-----+   +-------+--------+   +------+------+
           |                  |                   |
           | topic/env_snap   |                   | topic/access_*
           v                  v                   v
    +------------+    +-------+--------+   +------------+
    |     L3     |    |       L4       |   |     L6     |
    | Telemetry  |<-->| Adversary      |<->| Response   |
    +------+-----+    | Simulation     |   +------+-----+
           |         +-------+--------+          |
           |                  |                  |
           v                  v                  v
    +------------+    +-------+--------+   +------------+
    |     L5     |    |       L8       |   |     L7     |
    | Detection  |--->| Assurance       |   | Engineering|
    +------------+    +----------------+   +------------+
```

## Topics

| Topic | Producer | Consumer | Payload |
|-------|----------|----------|---------|
| topic/governance_tokens | L0 | L2, L4, L6, L8 | GovernanceToken |
| topic/environment_snapshots | L1 | L2, L5 | EnvironmentSnapshot |
| topic/offensive_plan | L4 | L5 | SignedTask[] |
| topic/offensive_findings | L4 | L5, L8 | OffensiveFindings |
| topic/alerts | L5 | L6 | Alert |
| topic/incidents | L6 | L8 | Incident |
| topic/engineering_changes | L7 | L8 | EngineeringChange |
| topic/assurance_reports | L8 | L0 (feedback) | AssuranceReport |
| topic/evidence_records | All | L8 | EvidenceRecord |
