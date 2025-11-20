#!/bin/bash

# remove old venv
rm -rf ./.venv

# install uv and deps
curl -LsSf https://astral.sh/uv/install.sh | sh
uv venv && uv sync

# start backend
docker compose up -d --build

# wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# run tests with explicit failure handling
echo "Running Tavern tests..."
TEST_EXIT_CODE=0

# run tests and capture exit code
uv run pytest -n auto --maxfail=1 --exitfirst || TEST_EXIT_CODE=$?

# check if any tests were actually run
if [ $TEST_EXIT_CODE -eq 5 ]; then
    # stop backend
    # docker compose down

    echo "ERROR: No tests were collected or run!"
    echo "This could mean:"
    echo "1. Test files not found in the expected location"
    echo "2. No test cases matched the selection criteria"
    echo "3. Import errors in test files"
    exit 1
fi

# check if all tests failed
if [ $TEST_EXIT_CODE -ne 0 ]; then
    # stop backend
    # docker compose down

    echo "ERROR: Tests failed with exit code $TEST_EXIT_CODE"
    echo "Pytest exit codes:"
    echo "  0 = all tests passed"
    echo "  1 = some tests failed"
    echo "  2 = test execution interrupted"
    echo "  3 = internal error"
    echo "  4 = pytest command line usage error"
    echo "  5 = no tests collected"
    exit $TEST_EXIT_CODE
fi

# stop backend
docker compose down
