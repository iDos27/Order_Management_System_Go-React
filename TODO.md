# Order Management System

> System zarządzania zamówieniami oparty na mikroserwisach (Go) oraz frontendzie w React.
> Plik ten pełni rolę dokumentacji, statusu projektu oraz listy zadań (TODO).

---

## Aktualna Architektura

System składa się z niezależnych mikroserwisów komunikujących się przez REST API (docelowo również RabbitMQ). Całość jest schowana za API Gateway (Nginx).

```
                    ┌─────────────────┐
                    │   Admin Panel   │
                    │  (React :5173)  │
                    └─────────┬───────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │  Nginx Gateway  │
                    │   (Port :80)    │
                    └─────────┬───────┘
                              │
          ┌───────────────────┼───────────────────┐
          │                   │                   │
┌─────────▼───────┐ ┌─────────▼───────┐ ┌─────────▼───────┐
│  Auth Service   │ │  Order Service  │ │ Raport Service  │
│   (Go :8081)    │ │   (Go :8080)    │ │   (Go :8083)    │
└─────────┬───────┘ └─────────┬───────┘ └──────┬──┬───────┘
          │                   │                │  │
┌─────────▼───────┐ ┌─────────▼───────┐        │  │
│ Postgres Auth   │ │ Postgres Orders │◄───────┘  │
│  (Port :5433)   │ │  (Port :5432)   │           │
└─────────────────┘ └─────────┬───────┘           │
                              │                   │
                              │ (Events)          │
                              ▼                   │
                    ┌─────────────────┐           │
                    │    RabbitMQ     │           │
                    │  (Port :5672)   │           │
                    └─────────┬───────┘           │
                              │                   │
                              ▼                   │
                    ┌─────────────────┐           │
                    │  Notification   │           │
                    │     Service     │           │
                    └─────────────────┘           ▼
                                          ┌───────────────┐
                                          │ Raports DB    │
                                          │ (Port :5432)  │
                                          └───────────────┘
```

### Komponenty Systemu

| Serwis                   | Port    | Status       | Opis                                                                     |
| ------------------------ | ------- | ------------ | ------------------------------------------------------------------------ |
| **Nginx Gateway**        | `:80`   | ✅ Działa    | Reverse proxy, kieruje ruch do odpowiednich serwisów.                    |
| **Frontend (Admin)**     | `:5173` | ✅ Działa    | Panel React dla pracowników magazynu.                                    |
| **Auth Service**         | `:8081` | ✅ Działa    | Rejestracja, logowanie, JWT. Baza: `auth_db`.                            |
| **Order Service**        | `:8080` | ✅ Działa    | Zarządzanie zamówieniami. Baza: `orders_db`.                             |
| **Raport Service**       | `:8083` | ✅ Działa    | Generowanie raportów. Bazy: `orders_db` (odczyt) i `raports_db` (zapis). |
| **Notification Service** | -       | 📅 Planowany | Powiadomienia systemowe (Linux Native).                                  |

---

## Instrukcja Uruchomienia

### 1. Infrastruktura (Kontenery Podman)

Uruchom bazy danych i Nginx:

```bash
# Baza Orders
podman run --name postgres-orders -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password123 -e POSTGRES_DB=orders_management -p 5432:5432 -d postgres:17

# Baza Auth
podman run --name postgres-auth -e POSTGRES_PASSWORD=password -e POSTGRES_DB=auth_service -p 5433:5432 -d postgres:17

# Nginx Gateway (uruchom z głównego katalogu)
podman run -d --name nginx-gateway --network host -v $(pwd)/nginx/nginx.conf:/etc/nginx/nginx.conf:ro nginx:alpine
```

#### Inicjalizacja Baz Danych (Migracje)

Po uruchomieniu kontenerów należy utworzyć tabele:

**1. Orders Service DB:**

```bash
podman exec -i postgres-orders psql -U postgres -d orders_management < services/order-service/migrations/create_tables.sql
```

**2. Auth Service DB:**

```bash
podman exec -it postgres-auth psql -U postgres -d auth_service -c "
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);"
```

**3. Raport Service DB:**

```bash
# Najpierw utwórz bazę (jeśli nie istnieje w kontenerze orders lub auth - Raport Service korzysta z osobnej logicznej bazy, ale tutaj zakładamy osobny kontener lub tę samą instancję.
# W kodzie Raport Service łączy się z 'raports_management'. Utwórzmy ją w kontenerze postgres-orders dla uproszczenia, lub jeśli masz osobny kontener, użyj go.
# Zakładając, że używamy postgres-orders (port 5432) jako hosta również dla tej bazy:

podman exec -it postgres-orders psql -U postgres -c "CREATE DATABASE raports_management;"
podman exec -i postgres-orders psql -U postgres -d raports_management < services/raport-service/migration/create_tables.sql
```

### 2. Uruchomienie Serwisów (Go)

Otwórz osobne terminale dla każdego serwisu:

```bash
# Terminal 1: Auth Service
cd services/auth-service && go run cmd/server/main.go

# Terminal 2: Order Service
cd services/order-service && go run cmd/server/main.go

# Terminal 3: Raport Service
cd services/raport-service && go run cmd/server/main.go
```

### 3. Frontend (React)

```bash
cd frontend/admin-panel
npm install
npm run dev
```

---

## TODO / Roadmapa

### System Powiadomień (Linux Native)

Cel: Wyświetlanie natywnych dymków powiadomień na pulpicie Linuxa, gdy wpłynie nowe zamówienie.

- [ ] **Infrastruktura RabbitMQ**

  - [ ] Uruchomienie kontenera RabbitMQ (Port 5672/15672).
  - [ ] Konfiguracja Exchange `orders_exchange` i kolejki `notifications_queue`.

- [ ] **Notification Service**

  - [ ] Inicjalizacja projektu w `services/notification-service`.
  - [ ] Implementacja konsumenta AMQP w Go.
  - [ ] Integracja z systemem powiadomień (np. `libnotify` / `notify-send`).
  - [ ] Obsługa kliknięcia w powiadomienie (otwarcie przeglądarki).

- [ ] **Integracja Order Service**
  - [ ] Dodanie publikowania zdarzeń do RabbitMQ przy tworzeniu/edycji zamówienia.

### Konteneryzacja i Orkiestracja

Cel: Pełna konteneryzacja środowiska deweloperskiego przy użyciu Podmana.

- [ ] **Konteneryzacja Aplikacji**

  - [ ] Stworzenie `Containerfile` dla każdego serwisu (Auth, Order, Raport, Frontend).
  - [ ] Budowa obrazów lokalnych: `podman build ...`

- [ ] **Podman Play Kube**
  - [ ] Przygotowanie definicji Podów (YAML).
  - [ ] Uruchamianie całego stacka jedną komendą: `podman play kube system.yaml`.
