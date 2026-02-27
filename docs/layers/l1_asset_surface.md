# L1: Asset & Attack Surface

## Responsibilities

- Maintain high-fidelity model of assets, identities, data flows
- Produce EnvironmentSnapshot objects
- Identify asset gaps (orphaned identities, missing coverage)
- Derive attack surface view (paths from exposed to critical assets)

## Inputs

- Raw inventories from cloud, CMDB, EASM
- GovernanceToken (for critical processes, scope)

## Outputs

- EnvironmentSnapshot on topic/environment_snapshots
- AssetGap on topic/asset_gaps

## Invariants

- No critical asset unseen for more than seven days
- All critical business processes have at least one mapped asset

## Global Requirements

- S5: Coverage & Depth
- S2: Policy & Scope Compliance (scope for L4)
