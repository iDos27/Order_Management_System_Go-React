# System Obsługi Realizacji Zamówień - Dokumentacja Techniczna

> **Dokumentacja Referencyjna do Pracy Inżynierskiej**
>
> Ten dokument stanowi techniczny opis implementacji systemu, architektury oraz kluczowych rozwiązań programistycznych.

---

## 1. Cel i Opis Systemu

Celem projektu jest stworzenie rozproszonego systemu obsługi zamówień magazynowych, który umożliwia:

- **Czas rzeczywisty**: Natychmiastową synchronizację statusów zamówień między pracownikami (WebSocket).
- **Skalowalność**: Podział na niezależne mikroserwisy (Auth, Order, Raport).
- **Bezpieczeństwo**: Pełne uwierzytelnianie oparte na tokenach JWT.
- **Analitykę**: Generowanie raportów w formatach biznesowych (Excel).

---

## 2. Stack Technologiczny

### Backend (Go)

Wykorzystano język Go ze względu na wysoką wydajność i natywną obsługę współbieżności.

```go
// Główne zależności (go.mod)
github.com/gin-gonic/gin           // Framework web REST API (v1.9.1)
github.com/lib/pq                  // Driver PostgreSQL (v1.10.9)
github.com/joho/godotenv           // Obsługa zmiennych środowiskowych .env (v1.5.1)
github.com/gorilla/websocket       // Implementacja protokołu WebSocket (v1.5.1)
github.com/gin-contrib/cors        // Middleware CORS (v1.5.0)
github.com/golang-jwt/jwt/v5       // Obsługa tokenów JWT (v5.2.0)
github.com/xuri/excelize/v2        // Generowanie plików Excel (v2.8.0)
```

### Frontend (React)

Aplikacja kliencka typu SPA (Single Page Application) zbudowana w oparciu o nowoczesny stack Reacta.

```json
// Główne zależności (package.json)
"react": "^18.2.0",                // Biblioteka UI
"react-dom": "^18.2.0",            // Renderowanie DOM
"react-router-dom": "^6.22.0",     // Routing po stronie klienta
"vite": "^5.1.0",                  // Build tool i serwer deweloperski (szybszy niż CRA)
"axios": "^1.6.7"                  // Klient HTTP
```

### Baza Danych i Infrastruktura

- **PostgreSQL 17**: Główny silnik bazy danych.
- **Podman**: Konteneryzacja usług i baz danych.
- **Nginx**: API Gateway i Reverse Proxy.

---

## 3. Architektura Systemu

System oparty jest na architekturze mikroserwisów, gdzie każdy moduł odpowiada za jedną domenę biznesową. Całość komunikacji zewnętrznej przechodzi przez API Gateway.

```
                                   ┌─────────────────┐
                                   │   Admin Panel   │
                                   │  (React :5173)  │
                                   └────────▼────────┘
                                            │
                                            │
                                            │
                                   ┌────────▲────────┐
                                   │    [HTTP/WS]    │
                                   └────────▼────────┘
                                            │
                                   ┌────────▼────────┐
                                   │  Nginx Gateway  │
                                   │       :80       │
                                   └─────────────────┘
                                            │
         ┌──────────────────────────────────┼───────────────────────────────────┐
         │                                  │                                   │
┌────────▼────────┐                ┌────────▼────────┐                 ┌────────▼────────┐
│  Auth Service   │                │  Order Service  │                 │  Raport Service │
│   (Go :8081)    │                │   (Go :8080)    │        ┌────────│   (Go :8083)    │
└─────────────────┘                └─────────────────┘        │        └─────────────────┘
         │                                  │                 │                 │
┌────────▼────────┐                ┌────────▼────────┐        │        ┌────────▼────────┐
│       DB        │                │       DB        │        │        │       DB        │
│  postgres_auth  │                │  postgres_order │◄───────┘        │ postgres_raport │
│      :5433      │                │      :5432      │                 │      :5434      │
└─────────────────┘                └─────────────────┘                 └─────────────────┘
                                            │
                                   ┌────────▼────────┐
                                   │    RabbitMQ     │
                                   │  (Port :5672)   │
                                   └─────────────────┘
                                            │
                                   ┌────────▼────────┐
                                   │  Notification   │
                                   │     Service     │
                                   │   (Go :8084)    │
                                   └─────────────────┘

```

---

## 4. Backend - Kluczowe Implementacje (Go)

### 4.1. Auth Service - Generowanie Tokena JWT

Implementacja bezpiecznego, bezstanowego uwierzytelniania. Token zawiera ID użytkownika, email i rolę, co pozwala na weryfikację uprawnień bez odpytywania bazy przy każdym żądaniu.

**Plik:** `services/auth-service/internal/handlers/handlers.go`

```go
// generateJWTToken tworzy podpisany token ważny przez 24 godziny
func (h *AuthHandler) generateJWTToken(user models.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    secret := os.Getenv("JWT_SECRET")
    return token.SignedString([]byte(secret))
}
```

### 4.2. Order Service - WebSocket Hub (Pub/Sub)

Hub zarządza aktywnymi połączeniami WebSocket. Wykorzystuje kanały Go (`chan`) do bezpiecznej komunikacji między wątkami (goroutines). Wzorzec ten pozwala na efektywne rozsyłanie powiadomień do wielu klientów jednocześnie.

**Plik:** `services/order-service/internal/websocket/websocket.go`

```go
type Hub struct {
    Clients    map[*Client]bool  // Mapa aktywnych klientów (Set)
    Register   chan *Client      // Kanał rejestracji nowych połączeń
    Unregister chan *Client      // Kanał zamykania połączeń
    Broadcast  chan Message      // Kanał rozgłoszeniowy
}

// Główna pętla obsługująca zdarzenia w Hubie
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            h.Clients[client] = true
        case client := <-h.Unregister:
            if _, ok := h.Clients[client]; ok {
                delete(h.Clients, client)
                close(client.Send)
            }
        case message := <-h.Broadcast:
            // Rozsyłanie wiadomości do wszystkich połączonych klientów
            for client := range h.Clients {
                select {
                case client.Send <- message:
                default:
                    close(client.Send)
                    delete(h.Clients, client)
                }
            }
        }
    }
}
```

### 4.3. Order Service - Tworzenie Zamówienia z Powiadomieniem

Handler REST API, który po poprawnym zapisaniu zamówienia w bazie danych, natychmiast wysyła powiadomienie przez WebSocket do wszystkich podłączonych klientów.

**Plik:** `services/order-service/internal/handlers/orders.go`

```go
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var order models.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // 1. Zapis do bazy danych (transakcyjność zapewniona przez PostgreSQL)
    err := h.db.QueryRow(`
        INSERT INTO orders (customer_name, customer_email, source, status, total_amount, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at
        `, order.CustomerName, order.CustomerEmail, order.Source, models.StatusNew, order.TotalAmount).
        Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
        return
    }

    // 2. Real-time Broadcast: Powiadomienie wszystkich klientów o nowym zamówieniu
    h.hub.BroadcastOrderUpdate(order.ID, string(models.StatusNew), "system")

    c.JSON(http.StatusCreated, order)
}
```

---

## 5. Frontend - Implementacja React

### 5.1. Globalny Stan Autoryzacji (Context API)

Wykorzystanie `React Context` do przechowywania stanu zalogowanego użytkownika w całej aplikacji. Pozwala to uniknąć "prop drilling" (przekazywania propsów przez wiele poziomów).

**Plik:** `frontend/admin-panel/src/context/AuthContext.jsx`

```javascript
export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);

  // Automatyczne przywracanie sesji po odświeżeniu strony
  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    const storedUser = localStorage.getItem("user");

    if (storedToken && storedUser) {
      setToken(storedToken);
      setUser(JSON.parse(storedUser));
    }
  }, []);

  const login = async (credentials) => {
    const response = await api.login(credentials);
    if (response.token) {
      // Zapisanie tokena w localStorage dla persystencji
      localStorage.setItem("token", response.token);
      localStorage.setItem("user", JSON.stringify(response.user));
      setToken(response.token);
      setUser(response.user);
      return { success: true };
    }
  };

  return (
    <AuthContext.Provider value={{ user, token, login }}>
      {children}
    </AuthContext.Provider>
  );
};
```

### 5.2. Obsługa WebSocket (Custom Hook)

Wydzielenie logiki WebSocket do osobnego hooka `useWebSocket` zapewnia czystość kodu komponentów i łatwość ponownego użycia.

**Plik:** `frontend/admin-panel/src/hooks/useWebSocket.js`

```javascript
const useWebSocket = (url) => {
  const [lastMessage, setLastMessage] = useState(null);
  const ws = useRef(null);

  useEffect(() => {
    ws.current = new WebSocket(url);

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("Otrzymano wiadomość:", data);
      setLastMessage(data); // Aktualizacja stanu wymusi re-render komponentu używającego hooka
    };

    return () => {
      if (ws.current) ws.current.close(); // Sprzątanie połączenia przy odmontowaniu
    };
  }, [url]);

  return { lastMessage };
};
```

### 5.3. Aktualizacja UI w Czasie Rzeczywistym

Komponent `OrderManagement` reaguje na zmiany w `lastMessage` z hooka WebSocket, aktualizując listę zamówień bez konieczności odświeżania strony.

**Plik:** `frontend/admin-panel/src/components/OrderManagement.jsx`

```javascript
const OrderManagement = () => {
  const [orders, setOrders] = useState([]);
  const { lastMessage } = useWebSocket("ws://localhost/ws");

  // Efekt nasłuchujący na nowe wiadomości WebSocket
  useEffect(() => {
    if (lastMessage && lastMessage.type === "order_update") {
      const { order_id, new_status } = lastMessage.payload;

      setOrders((prevOrders) => {
        const orderExists = prevOrders.find((o) => o.id === order_id);

        if (orderExists) {
          // Aktualizacja statusu istniejącego zamówienia (optymistyczna aktualizacja UI)
          return prevOrders.map((order) =>
            order.id === order_id ? { ...order, status: new_status } : order
          );
        } else {
          // Nowe zamówienie - pobranie świeżej listy
          fetchOrders();
          return prevOrders;
        }
      });

      alert(`Aktualizacja zamówienia #${order_id}: ${new_status}`);
    }
  }, [lastMessage]);

  // ... renderowanie Kanban Board
};
```

---

## 6. Uruchomienie Projektu (Podman)

```bash
# 1. Uruchomienie baz danych i Nginx
podman run --name postgres-orders -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password123 -e POSTGRES_DB=orders_management -p 5432:5432 -d postgres:17

podman run --name postgres-auth -e POSTGRES_PASSWORD=password -e POSTGRES_DB=auth_service -p 5433:5432 -d postgres:17

podman run -d --name nginx-gateway --network host -v $(pwd)/nginx/nginx.conf:/etc/nginx/nginx.conf:ro nginx:alpine

# 2. Uruchomienie serwisów Go (w osobnych terminalach)
cd services/order-service && go run cmd/server/main.go
cd services/auth-service && go run cmd/server/main.go
cd services/raport-service && go run cmd/server/main.go

# 3. Uruchomienie Frontendu
cd frontend/admin-panel && npm run dev
```
