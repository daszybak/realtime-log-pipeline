# Created by an LLM.
# Edited by me: docker image versions, air version.

# Development build Dockerfile
# For running Go services in containers with hot-reload capability

FROM golang:1.25.0-alpine as dev

# Install development tools
RUN apk add --no-cache \
    git \
    make \
    bash \
    wget \
    curl

# Install air for hot-reload
RUN go install github.com/cosmtrek/air@v1.62.0

# Set working directory
WORKDIR /workspace

# Copy dependency files for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy Makefile for build commands
COPY Makefile ./

# Create target directories
RUN mkdir -p target/artfs target/tmp

# Set up air directories for all services
RUN mkdir -p \
    target/tmp/air_api \
    target/tmp/air_api_testdata \
    target/tmp/air_worker \
    target/tmp/air_worker_testdata \
    target/tmp/air_aggregator \
    target/tmp/air_aggregator_testdata \
    target/tmp/air_streamer \
    target/tmp/air_streamer_testdata

# Default command will be overridden by docker-compose
CMD ["bash"]
