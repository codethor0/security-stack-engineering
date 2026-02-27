# L2: Identity & Access Plane

## Responsibilities

- Maintain IdentityGraph from L1 snapshots
- Evaluate AccessRequest, produce AccessDecision (allow, deny, challenge, step_up)
- Periodically assess identity risk (orphaned, stale, excessive privilege)
- Enforce zero trust policies (e.g., MFA for privileged)

## Inputs

- EnvironmentSnapshot, GovernanceToken
- AccessRequest on topic/access_requests

## Outputs

- AccessDecision on topic/access_decisions
- IdentityRisk on topic/identity_risks
- Evidence records

## Invariants

- Default deny when no policy matches
- Privileged identities require MFA when governance mandates

## Global Requirements

- S1: Attribution (decisions attributed)
- S2: Policy & Scope Compliance
- S9: Guardrailed Automation (approval flows)
