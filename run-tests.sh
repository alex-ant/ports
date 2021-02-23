#!/bin/bash

if [ "$1" == "compose" ]; then
  # This part is called by Go container holding the source code.

  # Run tests.
  go test -cover ./...

  exit $?
fi

# Exit on error.
set -e

# Start MySQL, Go containers and run tests.
docker-compose -f docker-compose-test.yml down
docker-compose -f docker-compose-test.yml rm -f
docker-compose -f docker-compose-test.yml up --no-build -d
docker-compose -f docker-compose-test.yml run ports-test ./run-tests.sh compose
TEST_EXIT_CODE=$?
docker-compose -f docker-compose-test.yml down

# Delete unused volumes.
docker volume prune -f > /dev/null

exit $TEST_EXIT_CODE
