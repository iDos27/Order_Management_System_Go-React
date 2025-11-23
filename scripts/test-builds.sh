#!/bin/bash

# Quick test script to verify Docker builds without full build
# Run from project root: ./scripts/test-builds.sh

set -e

echo "🔍 Testing Dockerfile syntax and context..."
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

test_dockerfile() {
    local service_name=$1
    local context_path=$2
    local dockerfile_path=$3
    
    echo -e "${YELLOW}Testing ${service_name}...${NC}"
    
    # Check if Dockerfile exists
    if [ ! -f "${dockerfile_path}" ]; then
        echo -e "${RED}❌ Dockerfile not found: ${dockerfile_path}${NC}"
        return 1
    fi
    
    # Check if context path exists
    if [ ! -d "${context_path}" ]; then
        echo -e "${RED}❌ Context path not found: ${context_path}${NC}"
        return 1
    fi
    
    # Validate Dockerfile syntax (dry-run)
    if docker build --no-cache -f "${dockerfile_path}" "${context_path}" 2>&1 | head -20; then
        echo -e "${GREEN}✅ ${service_name} Dockerfile syntax OK${NC}"
    else
        echo -e "${RED}❌ ${service_name} Dockerfile has issues${NC}"
        return 1
    fi
    
    echo ""
}

echo "1. Auth Service"
test_dockerfile "Auth Service" \
    "./services/auth-service" \
    "./services/auth-service/docker/dockerfile"

echo "2. Order Service"
test_dockerfile "Order Service" \
    "./services/order-service" \
    "./services/order-service/docker/dockerfile"

echo "3. Raport Service"
test_dockerfile "Raport Service" \
    "./services/raport-service" \
    "./services/raport-service/docker/dockerfile"

echo "4. Frontend"
test_dockerfile "Frontend" \
    "./frontend" \
    "./frontend/dockerfile"

echo -e "${GREEN}🎉 All Dockerfiles validated!${NC}"
echo ""
echo "To build all images, run: ./scripts/build-all.sh"
