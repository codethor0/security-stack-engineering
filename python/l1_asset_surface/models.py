"""L1 data models for assets, identities, and environment snapshots."""

from dataclasses import dataclass, asdict
from typing import List, Dict, Optional


@dataclass
class Asset:
    asset_id: str
    asset_type: str
    owner: str
    criticality: str
    business_processes: List[str]
    data_classification: str
    exposure: str
    last_seen: str


@dataclass
class Identity:
    identity_id: str
    identity_type: str
    owner: str
    privileged: bool
    last_auth: Optional[str]
    orphaned: bool


@dataclass
class EnvironmentSnapshot:
    snapshot_id: str
    timestamp: str
    assets: Dict[str, Asset]
    identities: Dict[str, Identity]
    data_flows: List[Dict]
    attack_surface: List[Dict]


@dataclass
class AssetGap:
    gap_type: str
    severity: str
    description: str
    remediation_owner: str
