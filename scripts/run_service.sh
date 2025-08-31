#!/bin/bash

# Created by an LLM.

# scripts/run_service.sh - Run individual Go service with air hot-reload
# This script runs inside the container

set -euo pipefail

SERVICE=${1:-}

if [[ -z "$SERVICE" ]]; then
    echo "Usage: $0 <service_name>"
    echo "Available services: api, worker, aggregator, streamer"
    exit 1
fi

# Ensure we're in the workspace directory
cd /workspace

# Copy backend source if not mounted (for build consistency)
if [[ ! -d "backend" ]]; then
    echo "❌ Backend source not found. Make sure volume is mounted correctly."
    exit 1
fi

echo "🔧 Running $SERVICE with air hot-reload..."

case "$SERVICE" in
    "api")
        air -c backend/configs/build/api.air.toml backend/configs/api.yaml 0.0.0.0:8081
        ;;
    "worker")
        air -c backend/configs/build/worker.air.toml backend/configs/worker.yaml 0.0.0.0:8082
        ;;
    "aggregator")
        air -c backend/configs/build/aggregator.air.toml backend/configs/aggregator.yaml 0.0.0.0:8083
        ;;
    "streamer")
        air -c backend/configs/build/streamer.air.toml backend/configs/streamer.yaml 0.0.0.0:8084
        ;;
    *)
        echo "❌ Unknown service: $SERVICE"
        echo "Available services: api, worker, aggregator, streamer"
        exit 1
        ;;
esac





