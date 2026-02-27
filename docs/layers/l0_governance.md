# L0: Govern & Strategy

## Responsibilities

- Produce and maintain the GovernanceToken (cryptographically signed)
- Aggregate governance state from upstream GRC/legal sources
- Derive objectives and target metrics from risk register
- Enforce structural invariants before token issuance

## Inputs

- External governance sources (regulations, policies, contracts)
- Risk register updates
- Evidence records for audit

## Outputs

- GovernanceToken on topic/governance_tokens
- Evidence records on governance failures

## Invariants

- Every objective must link to at least one risk scenario
- Every regulation must map to at least one control family

## Global Requirements

- S1: Attribution (token signed, issuer identified)
- S2: Policy & Scope Compliance (governance context)
- S6: Risk & Business Alignment (objectives derived from risk)
- S7: Multi-Team Coherence (central governance substrate)
