"""Tests for L5 Detection stub."""

from python.l5_detection_analytics.detection_engine import (
    DetectionEngineStub,
    DetectionRule,
    Alert,
)


def test_detection_engine_add_rule():
    engine = DetectionEngineStub()
    rule = DetectionRule(
        rule_id="r1",
        name="Credential abuse",
        technique_ids=["T1078"],
        severity="high",
    )
    engine.add_rule(rule)
    assert "r1" in engine.detection_rules


def test_detection_engine_evaluate():
    engine = DetectionEngineStub()
    engine.add_rule(
        DetectionRule(
            rule_id="r1",
            name="Brute force",
            technique_ids=["T1110"],
            severity="high",
        )
    )
    # T1110 is not covered by the engine's default rule, so only r1 can match.
    event = {
        "event_id": "e1",
        "enrichment": {"technique_id": "T1110"},
        "engagement_id": "eng-1",
        "actor": {"id": "u1"},
    }
    alert = engine.evaluate_event(event)
    assert alert is not None
    assert alert["rule_id"] == "r1"
    assert alert["technique_id"] == "T1110"
    assert alert["status"] == "new"
