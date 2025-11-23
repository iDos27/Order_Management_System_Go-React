#!/bin/bash

# Script to build all Docker images for Order Management System
# Run from project root: ./scripts/build-all.sh

set -e  # Exit on error

echo "🏗️  Building all Docker images..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to build and check
build_service() {
    local service_name=$1
    local context_path=$2
    local dockerfile_path=$3
    local image_tag=$4
    
    echo -e "${YELLOW}📦 Building ${service_name}...${NC}"
    
    if docker build -t "${image_tag}" -f "${dockerfile_path}" "${context_path}"; then
        echo -e "${GREEN}✅ ${service_name} built successfully${NC}"
        echo ""
        return 0
    else
        echo -e "${RED}❌ ${service_name} build failed${NC}"
        echo ""
        return 1
    fi
}

# Track failures
failures=0

# Build Auth Service
build_service "Auth Service" \
    "./services/auth-service" \
    "./services/auth-service/docker/dockerfile" \
    "auth-service:latest" || ((failures++))

# Build Order Service
build_service "Order Service" \
    "./services/order-service" \
    "./services/order-service/docker/dockerfile" \
    "order-service:latest" || ((failures++))

# Build Raport Service
build_service "Raport Service" \
    "./services/raport-service" \
    "./services/raport-service/docker/dockerfile" \
    "raport-service:latest" || ((failures++))

# Build Frontend
build_service "Frontend" \
    "./frontend" \
    "./frontend/dockerfile" \
    "frontend:latest" || ((failures++))

echo ""
echo "================================================"

if [ $failures -eq 0 ]; then
    echo -e "${GREEN}🎉 All images built successfully!${NC}"
    echo ""
    echo "📋 Built images:"
    docker images | grep -E "(auth-service|order-service|raport-service|frontend)" || true
    echo ""
    echo "Next steps:"
    echo "1. Run: docker compose up -d"
    echo "2. Check: docker compose ps"
    echo "3. Access: http://localhost:30080"
else
    echo -e "${RED}❌ ${failures} image(s) failed to build${NC}"
    exit 1
fi
