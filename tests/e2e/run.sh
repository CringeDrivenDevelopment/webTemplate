#!/bin/bash
set -euo pipefail

cleanup() {
    echo "Cleaning up Docker containers..."
    docker compose down --remove-orphans || true
}
trap cleanup EXIT

# start backend
echo "Running tests..."
docker compose up --build --exit-code-from runner --abort-on-container-exit
