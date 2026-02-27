"""L5 Detection & Analytics algorithm."""

import asyncio
import hashlib
from dataclasses import dataclass, asdict
from typing import Dict, List, Any, Optional, Tuple
from datetime import datetime, timezone


@dataclass
class DetectionRule:
    rule_id: str
    name: str
    technique_ids: List[str]
    severity: str
    enabled: bool = True


@dataclass
class Alert:
    alert_id: str
    timestamp: str
    rule_id: str
    technique_id: str
    severity: str
    entities: Dict[str, Any]
    evidence_refs: List[str]
    engagement_correlation: Optional[str]
    status: str


class DetectionAnalyticsEngine:
    """L5 algorithm: rules, alerts, coverage, offensive integration."""

    def __init__(
        self,
        message_bus: Any = None,
        evidence_store: Any = None,
        threat_intel: Optional[Dict] = None,
    ):
        self.bus = message_bus
        self.store = evidence_store
        self.threat_intel = threat_intel or {}
        self.detection_rules: Dict[str, DetectionRule] = {}
        self.active_alerts: Dict[str, Alert] = {}
        self._add_default_rules()

    def _add_default_rules(self) -> None:
        self.detection_rules["rule-rtea-simulated"] = DetectionRule(
            rule_id="rule-rtea-simulated",
            name="RTE-A simulated detection",
            technique_ids=["T1078"],
            severity="high",
            enabled=True,
        )

    async def run(self) -> None:
        """Run detection and offensive integration loops."""
        await asyncio.gather(
            self.detection_loop(),
            self.offensive_integration_loop(),
        )

    async def detection_loop(self) -> None:
        """Evaluate telemetry against rules."""
        if not self.bus or not hasattr(self.bus, "subscribe"):
            await asyncio.sleep(60)
            return

        def on_telemetry(data: Dict) -> None:
            events = data.get("events", [])
            for event in events:
                matches = self._evaluate_rules_against_event(event)
                for rule_id, confidence in matches:
                    if confidence > 0.8:
                        alert = self._create_alert(rule_id, event, confidence)
                        self.active_alerts[alert.alert_id] = alert
                        if self.bus:
                            self.bus.publish("topic/alerts", asdict(alert))

        self.bus.subscribe("topic/telemetry_view", on_telemetry)
        while True:
            await asyncio.sleep(1)

    def _evaluate_rules_against_event(self, event: Dict) -> List[Tuple[str, float]]:
        matches = []
        for rule_id, rule in self.detection_rules.items():
            if not rule.enabled:
                continue
            engagement_id = event.get("engagement_id")
            technique = event.get("enrichment", {}).get("technique_id") or (
                "T1078" if engagement_id else None
            )
            if technique and technique in rule.technique_ids:
                matches.append((rule_id, 0.9))
        return matches

    def _create_alert(self, rule_id: str, event: Dict, confidence: float) -> Alert:
        rule = self.detection_rules[rule_id]
        alert_id = hashlib.sha256(
            f"{datetime.now(timezone.utc).isoformat()}{rule_id}".encode()
        ).hexdigest()[:16]
        return Alert(
            alert_id=alert_id,
            timestamp=datetime.now(timezone.utc).isoformat(),
            rule_id=rule_id,
            technique_id=rule.technique_ids[0] if rule.technique_ids else "unknown",
            severity=rule.severity,
            entities=event.get("actor", {}),
            evidence_refs=[event.get("event_id", "")],
            engagement_correlation=event.get("engagement_id"),
            status="new",
        )

    async def offensive_integration_loop(self) -> None:
        """Pre-position detection when offensive plan arrives."""
        if not self.bus:
            await asyncio.sleep(60)
            return

        def on_plan(data: Any) -> None:
            plan = data if isinstance(data, list) else [data]
            techniques = set()
            for task in plan:
                t = task.get("task", task) if isinstance(task, dict) else task
                tid = t.get("technique_id") if isinstance(t, dict) else None
                if tid:
                    techniques.add(tid)
            for technique in techniques:
                if technique not in self._get_covered_techniques():
                    if self.bus:
                        self.bus.publish("topic/detection_gaps", {
                            "gap_id": hashlib.sha256(technique.encode()).hexdigest()[:16],
                            "technique_id": technique,
                            "tactic": "unknown",
                            "severity": "high",
                            "reason": "pre_positioning_failure",
                            "recommended_action": f"No detection for {technique}",
                        })

        if hasattr(self.bus, "subscribe"):
            self.bus.subscribe("topic/offensive_plan", on_plan)
        while True:
            await asyncio.sleep(5)

    def _get_covered_techniques(self) -> set:
        out = set()
        for r in self.detection_rules.values():
            if r.enabled:
                out.update(r.technique_ids)
        return out

    def add_rule(self, rule: DetectionRule) -> None:
        self.detection_rules[rule.rule_id] = rule

    def evaluate_event(self, event: Dict) -> Optional[Dict]:
        matches = self._evaluate_rules_against_event(event)
        for rule_id, confidence in matches:
            if confidence > 0.8:
                alert = self._create_alert(rule_id, event, confidence)
                self.active_alerts[alert.alert_id] = alert
                return asdict(alert)
        return None


class DetectionEngineStub(DetectionAnalyticsEngine):
    """Backward-compatible stub alias."""

    pass
