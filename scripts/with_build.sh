#!/bin/bash

# Created by an LLM.

# scripts/with_build.sh - Run Go services in containers with build consistency
# Usage: ./scripts/with_build.sh <service_name> [command]
#
# Examples:
#   ./scripts/with_build.sh api           # Run api service in container
#   ./scripts/with_build.sh worker        # Run worker service in container  
#   ./scripts/with_build.sh api bash      # Get shell in api container
#   ./scripts/with_build.sh all           # Run all services in containers

set -euo pipefail

SERVICE=${1:-}
COMMAND=${2:-}

if [[ -z "$SERVICE" ]]; then
    echo "Usage: $0 <service_name> [command]"
    echo ""
    echo "Available services: api, worker, aggregator, streamer, all"
    echo ""
    echo "Examples:"
    echo "  $0 api           # Run api service"
    echo "  $0 worker        # Run worker service"  
    echo "  $0 api bash      # Get shell in api container"
    echo "  $0 all           # Run all services"
    exit 1
fi

# Ensure infrastructure is running
echo "🚀 Starting infrastructure services..."
docker-compose -f docker-compose.dev.yml up -d postgres rabbitmq prometheus grafana

# Wait for services to be healthy
echo "⏳ Waiting for infrastructure to be ready..."
docker-compose -f docker-compose.dev.yml exec postgres pg_isready -U user -d db || {
    echo "Waiting for PostgreSQL..."
    sleep 5
    docker-compose -f docker-compose.dev.yml exec postgres pg_isready -U user -d db
}

echo "✅ Infrastructure ready!"

case "$SERVICE" in
    "api")
        echo "🏃 Running API service in container..."
        if [[ -n "$COMMAND" ]]; then
            docker-compose -f docker-compose.dev.yml run --rm --service-ports api $COMMAND
        else
            docker-compose -f docker-compose.dev.yml run --rm --service-ports api ./scripts/run_service.sh api
        fi
        ;;
    "worker")
        echo "🏃 Running Worker service in container..."
        if [[ -n "$COMMAND" ]]; then
            docker-compose -f docker-compose.dev.yml run --rm worker $COMMAND
        else
            docker-compose -f docker-compose.dev.yml run --rm worker ./scripts/run_service.sh worker
        fi
        ;;
    "aggregator")
        echo "🏃 Running Aggregator service in container..."
        if [[ -n "$COMMAND" ]]; then
            docker-compose -f docker-compose.dev.yml run --rm --service-ports aggregator $COMMAND
        else
            docker-compose -f docker-compose.dev.yml run --rm --service-ports aggregator ./scripts/run_service.sh aggregator
        fi
        ;;
    "streamer")
        echo "🏃 Running streamer service in container..."
        if [[ -n "$COMMAND" ]]; then
            docker-compose -f docker-compose.dev.yml run --rm --service-ports streamer $COMMAND
        else
            docker-compose -f docker-compose.dev.yml run --rm --service-ports streamer ./scripts/run_service.sh streamer
        fi
        ;;
    "all")
        echo "🏃 Running all Go services in containers..."
        # Uncomment services in docker-compose.dev.yml and start them
        sed -i.bak 's/^  # \(api\|worker\|aggregator\|streamer\):/  \1:/' docker-compose.dev.yml
        sed -i.bak 's/^  #   /    /' docker-compose.dev.yml
        docker-compose -f docker-compose.dev.yml up api worker aggregator streamer
        ;;
    *)
        echo "❌ Unknown service: $SERVICE"
        echo "Available services: api, worker, aggregator, streamer, all"
        exit 1
        ;;
esac


