# System Obsługi Realizacji Zamówień

>  System zarządzania zamówieniami z wykorzystaniem technologii **React.js** oraz **Go**

## Architektura Systemu

```
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Admin Panel   │
│   (React.js)    │    │   (Warehouse)   │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          └──────────┬───────────┘
                     │
            ┌────────▼────────┐
            │   Backend API   │
            │    (Go/Gin)     │
            └─────────────────┘
                     │
            ┌────────▼────────┐
            │   PostgreSQL    │
            │   + WebSocket   │
            └─────────────────┘
```

## Aktualna Funkcjonalność

### **Backend (Go + Gin)**
- **REST API** - CRUD operacje na zamówieniach
- **WebSocket** - Real-time aktualizacje statusów
- **PostgreSQL** - Przechowywanie danych zamówień
- **CORS** - Obsługa żądań z frontendu

### **Frontend (React + Vite)**
- **Admin Panel** - Kanban board dla pracowników magazynu
- **Real-time Updates** - Synchronizacja między kartami przeglądarki
- **Workflow System** - Logiczne przejścia między statusami

### **Funkcje Biznesowe**
- Przeglądanie wszystkich zamówień w układzie Kanban
- Zmiana statusów zamówień (new → confirmed → shipped → delivered)
- Możliwość anulowania zamówień na każdym etapie
- Real-time powiadomienia o zmianach dla wszystkich użytkowników
- Szczegółowy widok pojedynczego zamówienia

## Stack Technologiczny

### **Backend**
```go
// Główne zależności Go
github.com/gin-gonic/gin           // Framework web REST API
github.com/lib/pq                  // Driver PostgreSQL
github.com/joho/godotenv           // Zmienne środowiskowe z .env
github.com/gorilla/websocket       // WebSocket real-time communication
github.com/gin-contrib/cors        // Cross-Origin Resource Sharing
```

### **Frontend**
```json
// Główne zależności React
"react": "^18.0.0"                 // Biblioteka UI
"vite": "^4.0.0"                   // Build tool i dev server
"axios": "^1.0.0"                  // HTTP client dla API calls
```

### **Baza Danych**
- **PostgreSQL 17** - Relacyjna baza danych
- **Tabele**: `orders`, `order_items`
- **Docker container** - Łatwe uruchomienie lokalnie

## Struktura Projektu

```
Order_Management_System/
├── backend/                       # Backend Go
│   ├── cmd/server/main.go        # Entry point aplikacji
│   ├── internal/
│   │   ├── database/             # Połączenie z bazą danych
│   │   ├── handlers/             # REST API endpoints
│   │   ├── models/               # Struktury danych
│   │   └── websocket/            # WebSocket hub i komunikacja
│   ├── migrations/               # SQL skrypty dla bazy
│   ├── go.mod                    # Zależności Go
│   └── .env                      # Konfiguracja (DATABASE_URL, PORT)
├── frontend/admin-panel/         # Frontend React
│   ├── src/
│   │   ├── components/           # Komponenty React (OrderCard)
│   │   ├── hooks/                # Custom hooks (useWebSocket)
│   │   ├── services/             # API client (axios)
│   │   ├── App.jsx               # Główny komponent
│   │   └── App.css               # Style CSS
│   ├── package.json              # Zależności npm
│   └── vite.config.js            # Konfiguracja Vite
├── kubernetes/                   # Manifesty K8s (przyszłość)
├── README.md                     # Ta dokumentacja
└── TODO.md                       # Plan rozwoju mikrousług
```

## Szczegóły Implementacji

### **Backend - Kluczowe Komponenty**

#### `models/order.go`
```go
type OrderStatus string
const (
    StatusNew       OrderStatus = "new"        // Nowe zamówienie
    StatusConfirmed OrderStatus = "confirmed"  // Potwierdzone
    StatusShipped   OrderStatus = "shipped"    // Wysłane
    StatusDelivered OrderStatus = "delivered"  // Dostarczone
    StatusCancelled OrderStatus = "cancelled"  // Anulowane
)

type Order struct {
    ID            int       `json:"id" db:"id"`
    CustomerName  string    `json:"customer_name" db:"customer_name"`
    CustomerEmail string    `json:"customer_email" db:"customer_email"`
    Status        string    `json:"status" db:"status"`
    TotalAmount   float64   `json:"total_amount" db:"total_amount"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
```

#### `websocket/websocket.go`
```go
// Hub - zarządza wszystkimi połączeniami WebSocket
type Hub struct {
    Clients    map[*Client]bool  // Aktywni klienci
    Register   chan *Client      // Kanał rejestracji
    Unregister chan *Client      // Kanał wyrejestrowania  
    Broadcast  chan Message      // Kanał broadcast
}

// BroadcastOrderUpdate - wysyła aktualizację do wszystkich klientów
func (h *Hub) BroadcastOrderUpdate(orderID int, newStatus string, updatedBy string)
```

#### `handlers/orders.go`
```go
// REST API endpoints
GET    /api/orders           // Lista wszystkich zamówień
GET    /api/orders/:id       // Pojedyncze zamówienie
POST   /api/orders           // Nowe zamówienie
PATCH  /api/orders/:id/status // Zmiana statusu
```

### **Frontend - Kluczowe Komponenty**

#### `hooks/useWebSocket.js`
```javascript
// Custom hook dla WebSocket connection
const useWebSocket = (url) => {
  const [lastMessage, setLastMessage] = useState(null);
  // Automatyczne połączenie, obsługa wiadomości, cleanup
}
```

#### `App.jsx`
```javascript
// Główny komponent z Kanban board
const App = () => {
  // Real-time synchronizacja przez WebSocket
  useEffect(() => {
    if (lastMessage && lastMessage.type === 'order_update') {
      // Aktualizacja lokalnego stanu bez przeładowania
    }
  }, [lastMessage]);
}
```

## Uruchomienie Projektu

### **1. Uruchomienie Bazy Danych**
```bash
# PostgreSQL w Docker
docker run --name postgres-orders \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password123 \
  -e POSTGRES_DB=orders_management \
  -p 5432:5432 -d postgres:17

# Wykonanie migracji
docker exec -i postgres-orders psql -U postgres -d orders_management < backend/migrations/create_tables.sql
```

### **2. Uruchomienie Backend**
```bash
cd backend
go mod tidy
go run cmd/server/main.go
# Server dostępny na http://localhost:8080
```

### **3. Uruchomienie Frontend**
```bash
cd frontend/admin-panel
npm install
npm run dev
# Aplikacja dostępna na http://localhost:5173
```

## Testowanie API

### **Tworzenie nowego zamówienia**
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Jan Kowalski",
    "customer_email": "jan@example.com", 
    "total_amount": 299.99
  }'
```

### **Zmiana statusu zamówienia**
```bash
curl -X PATCH http://localhost:8080/api/orders/1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "confirmed"}'
```

### **Lista wszystkich zamówień**
```bash
curl http://localhost:8080/api/orders
```

## WebSocket Real-Time

### **Połączenie**
```javascript
// Frontend automatycznie łączy się z WebSocket
const socket = new WebSocket('ws://localhost:8080/ws');
```

### **Format wiadomości**
```json
{
  "type": "order_update",
  "payload": {
    "order_id": 1,
    "new_status": "shipped",
    "updated_by": "admin"
  }
}
```

## Konfiguracja Środowiska

### **Backend (.env)**
```env
DATABASE_URL=postgres://postgres:password123@localhost:5432/orders_management?sslmode=disable
SERVER_PORT=8080
```

### **Frontend (Vite)**
```javascript
// vite.config.js
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173
  }
})
```

## Planowany Rozwój

> **Szczegóły w [TODO.md](./TODO.md)**

### **Najbliższe Cele:**
- **Auth Service** - System logowania z rolami użytkowników
- **Notification Service** - Email/SMS/Desktop powiadomienia  
- **Analytics Service** - Raporty i statystyki
- **Kubernetes** - Orkiestracja kontenerów
- **Docker Compose** - Kompletne środowisko deweloperskie

### **Architektura Docelowa:**
- **Mikrousługi** - Podział na niezależne serwisy
- **API Gateway** - Centralne zarządzanie ruchem
- **Service Mesh** - Komunikacja między serwisami
- **Monitoring** - Prometheus + Grafana
- **CI/CD** - Automatyczne wdrażanie

## Informacje Techniczne

### **Porty:**
- **Backend**: 8080
- **Frontend**: 5173  
- **PostgreSQL**: 5432
- **WebSocket**: ws://localhost:8080/ws

### **Statusy Zamówień:**
1. **new** → Nowe zamówienie
2. **confirmed** → Potwierdzone przez magazyn
3. **shipped** → Wysłane do klienta
4. **delivered** → Dostarczone
5. **cancelled** → Anulowane (możliwe na każdym etapie)

### **Workflow Biznesowy:**
```
new → confirmed → shipped → delivered
 ↓       ↓          ↓
cancelled ← cancelled ← cancelled
```
