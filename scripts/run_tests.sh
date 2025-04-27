#!/bin/bash

# Exit on any error
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Running Unit Tests...${NC}"

# Run all unit tests
echo -e "${BLUE}Running Service Tests${NC}"
go test -v ./tests/unit/services/...

echo -e "${BLUE}Running Store Tests${NC}"
go test -v ./tests/unit/store/...

echo -e "${GREEN}✓ Unit Tests Passed${NC}"

# Run E2E tests only if the 'full' argument is provided
if [ "$1" == "full" ]; then
  echo -e "${BLUE}Running E2E Tests...${NC}"
  
  echo -e "${BLUE}Running Validator E2E Tests${NC}"
  go test -v ./tests/e2e/validators/...
  
  echo -e "${BLUE}Running Delegations E2E Tests${NC}"
  go test -v ./tests/e2e/delegations/...
  
  echo -e "${GREEN}✓ E2E Tests Passed${NC}"
else
  echo -e "${BLUE}Skipping E2E tests. Run with 'full' parameter to include E2E tests.${NC}"
fi

echo -e "${GREEN}All tests completed successfully!${NC}" 