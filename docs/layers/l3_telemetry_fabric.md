# L3: Telemetry & Data Fabric

## Responsibilities

- Collect, normalize, route security telemetry
- Produce NormalizedEvent with common schema
- Link telemetry to red team engagements (engagement_id)
- Detect telemetry gaps (missing logs for critical assets)
- Correlate L4 evidence with raw events

## Inputs

- Raw events from SIEM, XDR, cloud, endpoints
- EnvironmentSnapshot for asset context
- Evidence records from L4

## Outputs

- NormalizedEvent in TelemetryView
- TelemetryGap on topic/telemetry_gaps
- Correlated evidence on topic/correlated_evidence

## Invariants

- Critical assets have telemetry within SLO window

## Global Requirements

- S3: Observability & Evidence
- S5: Coverage & Depth
