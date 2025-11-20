#!/bin/bash

# Kolory dla czytelności
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funkcja do wyświetlania kolorowych komunikatów
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

# Ścieżka główna projektu
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

print_info "🚀 Uruchamianie systemu Order Management..."
print_info "📁 Katalog projektu: $PROJECT_ROOT"

# ==========================================
# 1. SPRAWDZENIE CZY KONTENERY JUŻ ISTNIEJĄ
# ==========================================
print_info "🔍 Sprawdzanie istniejących kontenerów..."

# Funkcja do zatrzymania i usunięcia kontenera jeśli istnieje
cleanup_container() {
    local container_name=$1
    if podman ps -a --format "{{.Names}}" | grep -q "^${container_name}$"; then
        print_warning "Kontener $container_name już istnieje. Usuwam..."
        podman stop "$container_name" 2>/dev/null
        podman rm "$container_name" 2>/dev/null
        print_success "Kontener $container_name usunięty"
    fi
}

cleanup_container "postgres-orders"
cleanup_container "postgres-auth"
cleanup_container "postgres-raports"
cleanup_container "nginx-gateway"

# ==========================================
# 2. URUCHOMIENIE BAZ DANYCH
# ==========================================
print_info "🐘 Uruchamianie baz danych PostgreSQL..."

# Baza Orders (port 5432)
print_info "Starting postgres-orders..."
podman run --name postgres-orders \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=password123 \
    -e POSTGRES_DB=orders_management \
    -p 5432:5432 \
    -d postgres:17

if [ $? -eq 0 ]; then
    print_success "✓ postgres-orders uruchomiony (port 5432)"
else
    print_error "✗ Błąd uruchamiania postgres-orders"
    exit 1
fi

# Baza Auth (port 5433)
print_info "Starting postgres-auth..."
podman run --name postgres-auth \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=password \
    -e POSTGRES_DB=auth_service \
    -p 5433:5432 \
    -d postgres:17

if [ $? -eq 0 ]; then
    print_success "✓ postgres-auth uruchomiony (port 5433)"
else
    print_error "✗ Błąd uruchamiania postgres-auth"
    exit 1
fi

# Baza Raports (port 5434)
print_info "Starting postgres-raports..."
podman run --name postgres-raports \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=password123 \
    -e POSTGRES_DB=raports_management \
    -p 5434:5432 \
    -d postgres:17

if [ $? -eq 0 ]; then
    print_success "✓ postgres-raports uruchomiony (port 5434)"
else
    print_error "✗ Błąd uruchamiania postgres-raports"
    exit 1
fi

# Czekamy aż bazy będą gotowe
print_info "⏳ Oczekiwanie na gotowość baz danych (10 sekund)..."
sleep 10

# ==========================================
# 3. MIGRACJE BAZ DANYCH
# ==========================================
print_info "📊 Uruchamianie migracji baz danych..."

# Migracja Orders
print_info "Migracja bazy orders_management..."
if [ -f "$PROJECT_ROOT/services/order-service/migrations/create_tables.sql" ]; then
    podman exec -i postgres-orders psql -U postgres -d orders_management < "$PROJECT_ROOT/services/order-service/migrations/create_tables.sql"
    if [ $? -eq 0 ]; then
        print_success "✓ Migracja orders_management zakończona"
    else
        print_error "✗ Błąd migracji orders_management"
    fi
else
    print_warning "Plik migracji orders nie znaleziony"
fi

# Migracja Auth
print_info "Migracja bazy auth_service..."
if [ -f "$PROJECT_ROOT/services/auth-service/migrations/create_tables.sql" ]; then
    podman exec -i postgres-auth psql -U postgres -d auth_service < "$PROJECT_ROOT/services/auth-service/migrations/create_tables.sql"
    if [ $? -eq 0 ]; then
        print_success "✓ Migracja auth_service zakończona"
    else
        print_error "✗ Błąd migracji auth_service"
    fi
else
    print_warning "Plik migracji auth nie znaleziony"
fi

# Migracja Raports
print_info "Migracja bazy raports_management..."
if [ -f "$PROJECT_ROOT/services/raport-service/migration/create_tables.sql" ]; then
    podman exec -i postgres-raports psql -U postgres -d raports_management < "$PROJECT_ROOT/services/raport-service/migration/create_tables.sql"
    if [ $? -eq 0 ]; then
        print_success "✓ Migracja raports_management zakończona"
    else
        print_error "✗ Błąd migracji raports_management"
    fi
else
    print_warning "Plik migracji raports nie znaleziony"
fi

# ==========================================
# 4. URUCHOMIENIE NGINX GATEWAY
# ==========================================
print_info "🌐 Uruchamianie Nginx Gateway..."

if [ -f "$PROJECT_ROOT/nginx/nginx.conf" ]; then
    podman run -d \
        --name nginx-gateway \
        --network host \
        -v "$PROJECT_ROOT/nginx/nginx.conf:/etc/nginx/nginx.conf:ro" \
        nginx:alpine
    
    if [ $? -eq 0 ]; then
        print_success "✓ Nginx Gateway uruchomiony (port 80)"
    else
        print_error "✗ Błąd uruchamiania Nginx"
        exit 1
    fi
else
    print_error "Plik nginx.conf nie znaleziony!"
    exit 1
fi

# ==========================================
# 5. URUCHOMIENIE SERWISÓW GO W TERMINALACH
# ==========================================
print_info "🚀 Uruchamianie serwisów Go w osobnych terminalach..."

# Funkcja do uruchamiania serwisu w Ghostty
launch_service() {
    local service_name=$1
    local service_path=$2
    local port=$3
    
    print_info "Uruchamianie $service_name (port $port)..."
    
    if [ -d "$PROJECT_ROOT/$service_path" ]; then
        ghostty --working-directory="$PROJECT_ROOT/$service_path" \
            --title="$service_name" \
            -e bash -c "echo '🚀 Uruchamianie $service_name na porcie $port...'; go run cmd/server/main.go; exec bash" &
        print_success "✓ $service_name uruchomiony w nowym oknie Ghostty"
    else
        print_error "✗ Katalog $service_path nie istnieje!"
    fi
}

# Sprawdzenie czy Ghostty jest dostępne
if ! command -v ghostty &> /dev/null; then
    print_warning "Ghostty nie jest zainstalowane. Sprawdzam alternatywy..."
    
    # Próba użycia Konsole (KDE)
    if command -v konsole &> /dev/null; then
        print_info "Używam konsole..."
        
        launch_service() {
            local service_name=$1
            local service_path=$2
            local port=$3
            
            print_info "Uruchamianie $service_name (port $port)..."
            
            if [ -d "$PROJECT_ROOT/$service_path" ]; then
                konsole --new-tab \
                    --workdir "$PROJECT_ROOT/$service_path" \
                    -e bash -c "echo '🚀 Uruchamianie $service_name na porcie $port...'; go run cmd/server/main.go; exec bash" &
                print_success "✓ $service_name uruchomiony w nowej karcie Konsole"
            else
                print_error "✗ Katalog $service_path nie istnieje!"
            fi
        }
    
    # Próba użycia gnome-terminal
    elif command -v gnome-terminal &> /dev/null; then
        print_info "Używam gnome-terminal..."
        
        launch_service() {
            local service_name=$1
            local service_path=$2
            local port=$3
            
            print_info "Uruchamianie $service_name (port $port)..."
            
            if [ -d "$PROJECT_ROOT/$service_path" ]; then
                gnome-terminal --tab --title="$service_name" --working-directory="$PROJECT_ROOT/$service_path" -- bash -c "echo '🚀 Uruchamianie $service_name na porcie $port...'; go run cmd/server/main.go; exec bash" &
                print_success "✓ $service_name uruchomiony w nowej karcie gnome-terminal"
            else
                print_error "✗ Katalog $service_path nie istnieje!"
            fi
        }
    # Próba użycia xterm
    elif command -v xterm &> /dev/null; then
        print_info "Używam xterm..."
        
        launch_service() {
            local service_name=$1
            local service_path=$2
            local port=$3
            
            print_info "Uruchamianie $service_name (port $port)..."
            
            if [ -d "$PROJECT_ROOT/$service_path" ]; then
                xterm -T "$service_name" -e "cd '$PROJECT_ROOT/$service_path' && echo '🚀 Uruchamianie $service_name na porcie $port...' && go run cmd/server/main.go; bash" &
                print_success "✓ $service_name uruchomiony w nowym oknie xterm"
            else
                print_error "✗ Katalog $service_path nie istnieje!"
            fi
        }
    else
        print_error "Nie znaleziono terminala graficznego (ghostty, konsole, gnome-terminal, xterm)"
        print_warning "Uruchamiam serwisy w tle..."
        
        launch_service() {
            local service_name=$1
            local service_path=$2
            local port=$3
            
            print_info "Uruchamianie $service_name (port $port) w tle..."
            
            if [ -d "$PROJECT_ROOT/$service_path" ]; then
                cd "$PROJECT_ROOT/$service_path"
                nohup go run cmd/server/main.go > "$PROJECT_ROOT/logs/${service_name}.log" 2>&1 &
                print_success "✓ $service_name uruchomiony w tle (logi: logs/${service_name}.log)"
                cd "$PROJECT_ROOT"
            else
                print_error "✗ Katalog $service_path nie istnieje!"
            fi
        }
        
        mkdir -p "$PROJECT_ROOT/logs"
    fi
fi

# Uruchomienie wszystkich serwisów
launch_service "Auth Service" "services/auth-service" "8081"
sleep 2
launch_service "Order Service" "services/order-service" "8080"
sleep 2
launch_service "Raport Service" "services/raport-service" "8083"
sleep 2

# Uruchomienie Frontendu
print_info "🌐 Uruchamianie Frontend (React)..."
if [ -d "$PROJECT_ROOT/frontend/admin-panel" ]; then
    # Sprawdź czy node_modules istnieją
    if [ ! -d "$PROJECT_ROOT/frontend/admin-panel/node_modules" ]; then
        print_warning "node_modules nie znalezione. Instaluję zależności..."
        cd "$PROJECT_ROOT/frontend/admin-panel"
        npm install
        cd "$PROJECT_ROOT"
    fi
    
    # Uruchom frontend w terminalu
    if command -v ghostty &> /dev/null; then
        ghostty --working-directory="$PROJECT_ROOT/frontend/admin-panel" \
            --title="Frontend (React)" \
            -e bash -c "echo '🌐 Uruchamianie Frontend na porcie 5173...'; npm run dev; exec bash" &
        print_success "✓ Frontend uruchomiony w nowym oknie Ghostty"
    elif command -v konsole &> /dev/null; then
        konsole --new-tab \
            --workdir "$PROJECT_ROOT/frontend/admin-panel" \
            -e bash -c "echo '🌐 Uruchamianie Frontend na porcie 5173...'; npm run dev; exec bash" &
        print_success "✓ Frontend uruchomiony w nowej karcie Konsole"
    elif command -v gnome-terminal &> /dev/null; then
        gnome-terminal --tab --title="Frontend (React)" --working-directory="$PROJECT_ROOT/frontend/admin-panel" -- bash -c "echo '🌐 Uruchamianie Frontend na porcie 5173...'; npm run dev; exec bash" &
        print_success "✓ Frontend uruchomiony w nowej karcie gnome-terminal"
    elif command -v xterm &> /dev/null; then
        xterm -T "Frontend (React)" -e "cd '$PROJECT_ROOT/frontend/admin-panel' && echo '🌐 Uruchamianie Frontend na porcie 5173...' && npm run dev; bash" &
        print_success "✓ Frontend uruchomiony w nowym oknie xterm"
    else
        cd "$PROJECT_ROOT/frontend/admin-panel"
        nohup npm run dev > "$PROJECT_ROOT/logs/Frontend.log" 2>&1 &
        print_success "✓ Frontend uruchomiony w tle (logi: logs/Frontend.log)"
        cd "$PROJECT_ROOT"
    fi
else
    print_error "✗ Katalog frontend/admin-panel nie istnieje!"
fi

# ==========================================
# 6. PODSUMOWANIE
# ==========================================
sleep 3
echo ""
print_success "=============================================="
print_success "✅ SYSTEM URUCHOMIONY POMYŚLNIE!"
print_success "=============================================="
echo ""
print_info "📦 Kontenery:"
print_info "  • postgres-orders   → localhost:5432"
print_info "  • postgres-auth     → localhost:5433"
print_info "  • postgres-raports  → localhost:5434"
print_info "  • nginx-gateway     → localhost:80"
echo ""
print_info "🔧 Serwisy Go:"
print_info "  • Auth Service      → localhost:8081"
print_info "  • Order Service     → localhost:8080"
print_info "  • Raport Service    → localhost:8083"
echo ""
print_info "🌐 Frontend:"
print_info "  • React Admin Panel → localhost:5173"
echo ""
print_info "💻 Dostęp do aplikacji:"
print_info "  • Przez Nginx       → http://localhost"
print_info "  • Bezpośrednio      → http://localhost:5173"
echo ""
print_warning "📝 Aby zatrzymać system, uruchom: ./PodmanStopEnv.sh"
echo ""
