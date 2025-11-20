#!/bin/bash

# Kolory dla czytelności
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_info "🛑 Zatrzymywanie systemu Order Management..."

# Zatrzymanie i usunięcie kontenerów
print_info "📦 Zatrzymywanie kontenerów..."

containers=("postgres-orders" "postgres-auth" "postgres-raports" "nginx-gateway")

for container in "${containers[@]}"; do
    if docker ps -a --format "{{.Names}}" | grep -q "^${container}$"; then
        print_info "Zatrzymywanie $container..."
        docker stop "$container" 2>/dev/null
        docker rm "$container" 2>/dev/null
        print_success "✓ $container zatrzymany i usunięty"
    else
        print_warning "$container nie jest uruchomiony"
    fi
done

# Zatrzymanie procesów Go (jeśli działają w tle)
print_info "🔧 Zatrzymywanie procesów Go..."

# Sprawdzenie czy są jakieś procesy go run
GO_PIDS=$(pgrep -f "go run cmd/server/main.go")

if [ ! -z "$GO_PIDS" ]; then
    print_info "Znaleziono procesy Go: $GO_PIDS"
    echo "$GO_PIDS" | xargs kill -15 2>/dev/null
    sleep 2
    
    # Jeśli nadal działają, wymuszam zamknięcie
    GO_PIDS=$(pgrep -f "go run cmd/server/main.go")
    if [ ! -z "$GO_PIDS" ]; then
        print_warning "Wymuszam zamknięcie procesów Go..."
        echo "$GO_PIDS" | xargs kill -9 2>/dev/null
    fi
    
    print_success "✓ Procesy Go zatrzymane"
else
    print_info "Brak procesów Go do zatrzymania"
fi

# Zatrzymanie procesów npm (frontend)
print_info "🌐 Zatrzymywanie Frontend (npm)..."

# Sprawdzenie czy są jakieś procesy npm run dev
NPM_PIDS=$(pgrep -f "npm run dev")

if [ ! -z "$NPM_PIDS" ]; then
    print_info "Znaleziono procesy npm: $NPM_PIDS"
    echo "$NPM_PIDS" | xargs kill -15 2>/dev/null
    sleep 2
    
    # Jeśli nadal działają, wymuszam zamknięcie
    NPM_PIDS=$(pgrep -f "npm run dev")
    if [ ! -z "$NPM_PIDS" ]; then
        print_warning "Wymuszam zamknięcie procesów npm..."
        echo "$NPM_PIDS" | xargs kill -9 2>/dev/null
    fi
    
    # Zatrzymaj także procesy node (Vite)
    NODE_PIDS=$(pgrep -f "vite")
    if [ ! -z "$NODE_PIDS" ]; then
        echo "$NODE_PIDS" | xargs kill -9 2>/dev/null
    fi
    
    print_success "✓ Procesy Frontend zatrzymane"
else
    print_info "Brak procesów Frontend do zatrzymania"
fi

# Czyszczenie logów (opcjonalne)
if [ -d "logs" ]; then
    print_info "🧹 Czyszczenie logów..."
    rm -rf logs/*.log 2>/dev/null
    print_success "✓ Logi wyczyszczone"
fi

echo ""
print_success "=============================================="
print_success "✅ SYSTEM ZATRZYMANY"
print_success "=============================================="
echo ""
print_info "Aby uruchomić system ponownie: ./makeEnv.sh"
echo ""
