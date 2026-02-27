# L8: Assurance & Audit

## Responsibilities

- Generate reports (executive, operational, regulatory)
- Map evidence to control attestations (NIST CSF, ISO 27001, SOC 2)
- Feed governance feedback (unmapped obligations, untested controls)
- Summarize evidence for audit

## Inputs

- Evidence records from all layers
- GovernanceToken (obligations, controls)
- Incidents, engineering changes, offensive findings

## Outputs

- AssuranceReport on topic/assurance_reports
- GovernanceFeedback on topic/governance_feedback
- Evidence summary with integrity hash

## Invariants

- No attestation without linked evidence
- Reports generated within required timeframes

## Global Requirements

- S1: Attribution
- S3: Observability & Evidence
- S6: Risk & Business Alignment
- S10: Auditability & Explainability
