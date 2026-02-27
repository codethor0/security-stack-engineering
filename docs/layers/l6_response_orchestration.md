# L6: Response & Orchestration

## Responsibilities

- Aggregate alerts into incidents
- Select and execute playbooks (NIST 800-61 lifecycle)
- Orchestrate automated and manual response actions
- Enforce SLAs and escalation
- Record every action with evidence

## Inputs

- Alert from L5
- GovernanceToken (SLAs, automation policy)
- Playbook definitions

## Outputs

- Incident on topic/incidents
- ResponseAction records
- Evidence records

## Invariants

- Critical incidents must have at least one action within 15 minutes
- SLA breach triggers escalation

## Global Requirements

- S3: Observability & Evidence
- S8: Learning & Adaptation (playbook updates)
- S9: Guardrailed Automation (approval for high-impact actions)
