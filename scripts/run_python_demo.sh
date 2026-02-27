#!/usr/bin/env bash
# Run Python layers standalone with MockBus (demonstrates algorithm structure).
# Use when Python cannot connect to Go orchestrator.

set -e
cd "$(dirname "$0")/.."
export PYTHONPATH="$PWD"

echo "Starting Python layer demos (MockBus - standalone)"

python3 -c "
import asyncio
from python.common.mock_bus import MockBus
from python.l1_asset_surface.asset_surface import AssetSurfaceAlgorithm

async def main():
    bus = MockBus(verbose=True)
    bus.publish('topic/governance_tokens', {'critical_processes': ['auth']})
    l1 = AssetSurfaceAlgorithm(bus)
    print('L1: Running one cycle')
    await l1.execute_cycle()
    print('L1: Done')

asyncio.run(main())
print('Python demo complete')
"
