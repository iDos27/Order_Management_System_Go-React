# Order Service

## Opis

Główny mikroserwis systemu odpowiedzialny za zarządzanie zamówieniami. Realizuje operacje CRUD na zamówieniach, komunikację w czasie rzeczywistym przez WebSocket, oraz integruję się z systemem powiadomień poprzez RabbitMQ.

## Architektura

### Port
- **8080** - HTTP Server + WebSocket

### Baza danych
- **PostgreSQL** (port 5432)
- **Nazwa bazy:** `orders_management`
- **Tabele:** `orders`, `order_items`

### Integracje
- **RabbitMQ** (port 5672) - Publikowanie powiadomień o zamówieniach
- **WebSocket** - Real-time updates dla frontendów
- **Auth Service** - Walidacja użytkowników przez Nginx

## Funkcjonalności

### 1. Pobieranie zamówień (`GET /api/orders`)
- Lista wszystkich zamówień posortowanych po dacie utworzenia (DESC)
- Chronione przez autoryzację (wymagany token JWT)
- Zwraca pełne informacje o zamówieniu: ID, klient, źródło, status, kwota, daty

### 2. Pobieranie pojedynczego zamówienia (`GET /api/orders/:id`)
- Szczegóły konkretnego zamówienia
- Walidacja istnienia zamówienia
- Zwraca 404 jeśli nie znaleziono

### 3. Tworzenie zamówienia (`POST /api/orders`)
- Przyjmuje dane: customer_name, customer_email, source, total_amount
- Automatyczne ustawienie statusu na `new`
- **Powiadomienia:**
  - Broadcast przez WebSocket do wszystkich podłączonych klientów
  - Publikacja do RabbitMQ dla notification-service
- Zwraca utworzone zamówienie z ID i timestampami

### 4. Aktualizacja statusu (`PATCH /api/orders/:id/status`)
- Zmiana statusu zamówienia
- **Dozwolone statusy:** `new`, `confirmed`, `shipped`, `delivered`, `cancelled`
- Walidacja poprawności statusu
- **Powiadomienia:**
  - Broadcast przez WebSocket
  - Publikacja do RabbitMQ (z danymi klienta i kwotą)

### 5. WebSocket komunikacja (`/ws`)
- Real-time updates dla frontendów
- Automatyczne ponowne połączenie przy rozłączeniu
- **Typy wiadomości:**
  - `order_update` - zmiana statusu zamówienia
  - Payload zawiera: `order_id`, `new_status`, `updated_by`

### 6. RabbitMQ Publisher
- Publikacja powiadomień do kolejki `order_notifications`
- **Struktura powiadomienia:**
  ```json
  {
    "order_id": 123,
    "customer_name": "Jan Kowalski",
    "status": "new",
    "total_amount": 299.99,
    "timestamp": "2025-11-21T10:30:00Z"
  }
  ```
- Graceful degradation - serwis działa nawet jeśli RabbitMQ jest niedostępny

## Struktura projektu

```
order-service/
├── cmd/
│   └── server/
│       └── main.go              # Entry point + inicjalizacja WebSocket/RabbitMQ
├── internal/
│   ├── database/
│   │   └── connection.go        # Połączenie z PostgreSQL
│   ├── handlers/
│   │   └── orders.go            # CRUD dla zamówień
│   ├── models/
│   │   └── order.go             # Modele Order, Status, Source
│   ├── publisher/
│   │   └── publisher.go         # RabbitMQ publisher
│   └── websocket/
│       └── websocket.go         # WebSocket Hub, Client, Message handling
├── migrations/
│   └── create_tables.sql        # Schemat bazy + przykładowe dane
├── docker/
│   └── dockerfile               # Multi-stage build
└── go.mod                       # Zależności
```

## Technologie

- **Go 1.24** - Język programowania
- **Gin** - HTTP framework
- **Gorilla WebSocket** - Komunikacja real-time
- **RabbitMQ (amqp091-go)** - Message broker
- **PostgreSQL** - Baza danych
- **Docker/Podman** - Konteneryzacja

## Zmienne środowiskowe

| Zmienna | Domyślna wartość | Opis |
|---------|-------------------|------|
| `SERVER_PORT` | `8080` | Port HTTP servera |
| `DATABASE_URL` | `postgres://postgres:password@localhost:5432/orders_management?sslmode=disable` | URL bazy danych |
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672/` | URL RabbitMQ |
| `RABBITMQ_QUEUE` | `order_notifications` | Nazwa kolejki powiadomień |

## Endpointy API

### HTTP REST
- `GET /api/orders` - Lista wszystkich zamówień (chronione)
- `GET /api/orders/:id` - Pojedyncze zamówienie (chronione)
- `POST /api/orders` - Utworzenie nowego zamówienia (chronione)
- `PATCH /api/orders/:id/status` - Aktualizacja statusu (chronione)
- `GET /health` - Health check

### WebSocket
- `GET /ws` - Połączenie WebSocket dla real-time updates

## Statusy zamówień

1. **new** - Nowe zamówienie
2. **confirmed** - Potwierdzone
3. **shipped** - Wysłane
4. **delivered** - Dostarczone
5. **cancelled** - Anulowane

## Źródła zamówień

- **website** - Strona internetowa
- **źródło_jeden** - Zewnętrzny system #1
- **źródło_dwa** - Zewnętrzny system #2
- **manual** - Ręczne wprowadzenie

## WebSocket Protocol

### Struktura wiadomości
```json
{
  "type": "order_update",
  "payload": {
    "order_id": 123,
    "new_status": "confirmed",
    "updated_by": "admin"
  }
}
```

### Połączenie (JavaScript przykad)
```javascript
const ws = new WebSocket('ws://localhost/ws');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.type === 'order_update') {
    console.log(`Order #${message.payload.order_id} changed to ${message.payload.new_status}`);
  }
};
```

## RabbitMQ Integration

### Kolejka
- **Nazwa:** `order_notifications`
- **Durability:** `true` (przeżywa restart RabbitMQ)
- **Format:** JSON

### Publisher
- Połączenie przy starcie aplikacji
- Timeout publikacji: 5 sekund
- Logowanie błędów bez przerywania requesta
- Graceful close przy shutdown

## Deployment

### Docker
```bash
cd services/order-service
docker build -t order-service -f docker/dockerfile .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e RABBITMQ_URL="amqp://..." \
  order-service
```

### Health Check
```bash
curl http://localhost:8080/health
# Odpowiedź: {"status":"ok","service":"order-service"}
```

## Przykładowe użycie

### Utworzenie zamówienia
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "customer_name": "Jan Kowalski",
    "customer_email": "jan@example.com",
    "source": "website",
    "total_amount": 299.99
  }'
```

### Zmiana statusu
```bash
curl -X PATCH http://localhost:8080/api/orders/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"status": "confirmed"}'
```

## Integracja z systemem

### Frontend
- Połączenie WebSocket dla real-time updates
- REST API dla operacji CRUD
- Autoryzacja przez tokeny JWT

### Notification Service
- Subskrypcja kolejki `order_notifications` w RabbitMQ
- Automatyczne powiadomienia dla użytkowników

### Raport Service
- Dostęp read-only do bazy `orders_management`
- Generowanie raportów na podstawie danych zamówień
