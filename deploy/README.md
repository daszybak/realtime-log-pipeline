# Docker Deployment Guide

This directory contains the Docker Compose configuration for the Realtime Log Pipeline project.

## Services

The docker-compose.yml defines the following services:

### Application Services

- **app** (port 8080) - React/Vite frontend application
- **api** (port 8081) - Main API service with health checks
- **worker** (port 8082) - Background job processor (scalable)
- **aggregator** (port 8083) - Data aggregation service with metrics
- **streamer** (port 8084) - Data polling/ingestion service

### Infrastructure Services

- **rabbitmq** (ports 5672, 15672) - Message broker with management UI
- **postgres** (port 5432) - Database with persistent storage
- **prometheus** (port 9090) - Metrics collection
- **grafana** (port 3000) - Monitoring dashboard (admin/admin)

## Usage

### Start all services:

```bash
docker-compose -f deploy/docker-compose.yml up -d
```

### Start with logs:

```bash
docker-compose -f deploy/docker-compose.yml up
```

### Scale workers:

```bash
docker-compose -f deploy/docker-compose.yml up --scale worker=3 -d
```

### Stop all services:

```bash
docker-compose -f deploy/docker-compose.yml down
```

### Stop and remove volumes:

```bash
docker-compose -f deploy/docker-compose.yml down -v
```

## Accessing Services

- **Frontend**: http://localhost:8080
- **API**: http://localhost:8081
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **PostgreSQL**: localhost:5432 (user/pass/db)

## Configuration

- Backend services use YAML configuration files in `backend/configs/`
- Environment variables are set in docker-compose.yml
- Prometheus configuration in `deploy/prometheus/prometheus.yml`
- Grafana provisioning in `deploy/grafana/provisioning/`

## Health Checks

The following services include health checks:

- API, Aggregator, streamer: HTTP health endpoints
- RabbitMQ: Built-in diagnostics
- PostgreSQL: Connection checks

## Networks and Volumes

- All services run on the `app-network` bridge network
- Persistent volumes: postgres-data, grafana-data, prometheus-data, rabbitmq-data

## Development

For development with hot-reload, use the Justfile commands instead:

```bash
just dev_all  # Start all services in development mode
```

