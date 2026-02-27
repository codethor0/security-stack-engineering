# L7: Engineering & Hardening

## Responsibilities

- Aggregate findings from all layers
- Prioritize by risk, asset criticality, exploitability
- Create EngineeringChange (patch, config, IAM, IaC, code)
- Execute via CI/CD and infrastructure tools
- Validate that changes fix findings
- Report hardening progress

## Inputs

- Asset gaps, identity risks, detection gaps, telemetry gaps
- Offensive findings, incident lessons
- Vuln scans, AppSec, CSPM
- GovernanceToken (remediation SLAs)

## Outputs

- EngineeringChange on topic/engineering_changes
- HardeningProgress on topic/hardening_progress
- Evidence records

## Invariants

- Critical findings must have owners and due dates
- Failed changes trigger rollback when available

## Global Requirements

- S4: Lifecycle Safety (rollback plans)
- S6: Risk & Business Alignment (prioritization)
- S8: Learning & Adaptation
- S9: Guardrailed Automation
