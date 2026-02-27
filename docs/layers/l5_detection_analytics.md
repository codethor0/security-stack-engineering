# L5: Detection & Analytics

## Responsibilities

- Maintain detection rules (ATT&CK-mapped)
- Evaluate telemetry against rules, emit alerts
- Pre-position detection when offensive plan arrives
- Build coverage matrix, identify detection gaps
- Threat hunting and purple-team integration

## Inputs

- TelemetryView from L3
- OffensivePlan, OffensiveFindings from L4
- Threat scenarios
- GovernanceToken

## Outputs

- Alert on topic/alerts
- CoverageMatrix on topic/coverage_matrix
- DetectionGap on topic/detection_gaps
- DefensivePlan readiness status

## Invariants

- High-priority techniques in threat model should not remain uncovered long
- Pre-positioning failure emits gap before test

## Global Requirements

- S3: Observability & Evidence
- S5: Coverage & Depth
- S8: Learning & Adaptation (rule updates from findings)
