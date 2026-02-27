# PoC Threat Model

This document describes the threat model for the reference implementation itself, not the defended environment.

## Assets

- Source code and configuration
- In-memory message bus (no persistent queue in PoC)
- Governance tokens and signing keys
- Evidence records

## Threats

| Threat | Mitigation |
|--------|------------|
| Message tampering | Schema validation; L4 uses cryptographic signing |
| Key compromise | Keys are process-local in PoC; production would use HSM/KMS |
| Information disclosure via logs | No secrets in code; evidence records omit sensitive payloads |
| Denial of service (bus overload) | Bounded buffers; slow subscriber drop |
| Supply chain | Minimal dependencies; vendoring recommended |

## Assumptions

- PoC runs in trusted environment (no network exposure by default)
- Operators have legitimate access to orchestration
- External integrations are behind mockable interfaces
