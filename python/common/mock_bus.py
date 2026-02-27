"""Mock message bus for standalone Python layer demos."""

import asyncio
import json
from typing import Any, Dict, Optional
from collections import deque


class MockBus:
    """In-memory bus that logs all actions. Used when Python runs standalone."""

    def __init__(self, verbose: bool = True):
        self.verbose = verbose
        self._topics: Dict[str, deque] = {}
        self._subscribers: Dict[str, list] = {}
        self._latest: Dict[str, Any] = {}

    def publish(self, topic: str, payload: Any) -> None:
        data = payload if isinstance(payload, dict) else _to_dict(payload)
        self._latest[topic] = data
        if topic not in self._topics:
            self._topics[topic] = deque(maxlen=100)
        self._topics[topic].append(data)
        if self.verbose:
            print(f"[MockBus] PUB {topic}: {_short(data)}")
        for cb in self._subscribers.get(topic, []):
            try:
                cb(data)
            except Exception:
                pass

    def subscribe(self, topic: str, callback=None):
        if topic not in self._subscribers:
            self._subscribers[topic] = []
        if callback:
            self._subscribers[topic].append(callback)
        return self._topics.get(topic, deque())

    async def consume_latest(self, topic: str, timeout: float = 5.0) -> Optional[Dict]:
        await asyncio.sleep(0)
        return self._latest.get(topic)

    def get_latest(self, topic: str) -> Optional[Dict]:
        return self._latest.get(topic)


def _to_dict(obj: Any) -> Dict:
    if hasattr(obj, "__dict__"):
        return {k: v for k, v in obj.__dict__.items() if not k.startswith("_")}
    if hasattr(obj, "_asdict"):
        return obj._asdict()
    return dict(obj) if isinstance(obj, (dict, list)) else {"value": obj}


def _short(data: Any, max_len: int = 80) -> str:
    s = json.dumps(data, default=str)[:max_len]
    return s + "..." if len(s) >= max_len else s
