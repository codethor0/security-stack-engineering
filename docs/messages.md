# Message Types and Topics

All messages use JSON. Schemas are in `schemas/`.

| Topic | Schema | Description |
|-------|--------|--------------|
| topic/governance_tokens | governance_token.schema.json | Signed governance context |
| topic/environment_snapshots | environment_snapshot.schema.json | Asset/identity/data model |
| topic/identity_graph | identity_graph.schema.json | L2 identity state |
| topic/telemetry_event | telemetry_event.schema.json | Normalized event |
| topic/offensive_plan | offensive_plan.schema.json | Signed tasks before execution |
| topic/offensive_findings | offensive_findings.schema.json | Post-engagement results |
| topic/detection_rule | detection_rule.schema.json | Detection content |
| topic/alerts | alert.schema.json | SOC alert |
| topic/incidents | incident.schema.json | IR incident |
| topic/engineering_change | engineering_change.schema.json | Remediation action |
| topic/assurance_reports | assurance_report.schema.json | Report output |
| topic/evidence_records | evidence_record.schema.json | Audit evidence |

## Common Fields

- timestamp: Unix or ISO8601
- producer: Layer identifier (e.g., L0_GOVERN_STRATEGY)
- message_id: Unique per message
