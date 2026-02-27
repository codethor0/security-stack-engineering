# L4: Adversary Simulation (RTE-A)

## Responsibilities

- Governed adversary simulation with typed tasking
- Validate scenarios against governance before execution
- Sign tasks with Ed25519 for attribution
- Enforce TTL and lifecycle safety
- Emit evidence and audit records for every execution

## Inputs

- Scenario (techniques, targets, operator, approver)
- GovernanceToken (scope, allowed techniques)
- Data plane abstraction for execution

## Outputs

- OffensivePlan on topic/offensive_plan (before execution)
- OffensiveFindings on topic/offensive_findings (after)
- Evidence records with task_id, engagement_id

## Invariants

- Every task has operator and approver
- TTL within bounds (1–3600 seconds)
- Invalid signature blocks execution

## Global Requirements

- S1: Attribution (signed tasks)
- S3: Observability & Evidence
- S4: Lifecycle Safety (TTL, cancellation)
- S5: Coverage & Depth
