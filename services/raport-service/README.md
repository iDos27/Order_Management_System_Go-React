# Raport Service

## Opis

Mikroserwis odpowiedzialny za generowanie raportów sprzedażowych w formacie Excel. Analizuje dane zamówień z bazy order-service, tworzy statystyki według źródeł zamówień oraz automatycznie generuje raporty okresowe za pomocą zadania cron.

## Architektura

### Port
- **8083** - HTTP Server

### Bazy danych
- **PostgreSQL orders_management** (port 5432) - READ ONLY - dane zamówień
- **PostgreSQL reports_management** (port 5434) - READ/WRITE - przechowywanie raportów

### Dostęp
- Tylko dla użytkowników z rolą **admin**
- Autoryzacja przez auth-service via Nginx

## Funkcjonalności

### 1. Generowanie raportu (`POST /api/reports/generate`)
- Parametry:
  - `type` - typ raportu (daily, weekly, monthly, custom)
  - `period_start` - data początkowa (format: 2006-01-02)
  - `period_end` - data końcowa
- **Proces:**
  1. Pobranie danych zamówień z bazy orders_management
  2. Agregacja statystyk według źródeł
  3. Generowanie pliku Excel z formatowaniem
  4. Zapis pliku w katalogu `./reports`
  5. Zapis metadanych raportu w bazie reports_management
- **Zwraca:** Pełne dane raportu z ścieżką do pliku

### 2. Szybkie statystyki (`GET /api/reports/stats`)
- Podgląd statystyk bez generowania pliku Excel
- Parametry query:
  - `start` - data początkowa (opcjonalne, domyślnie -7 dni)
  - `end` - data końcowa (opcjonalne, domyślnie dzisiaj)
- **Zwraca:** JSON z agregowanymi danymi

### 3. Automatyczne raporty (Cron Job)
- **Harmonogram:** Co 5 minut (dla testów, w produkcji: co poniedziałek)
- **Zakres:** Ostatnie 7 dni
- **Typ:** weekly
- **Proces:**
  1. Automatyczne generowanie raportu
  2. Zapis pliku Excel
  3. Zapis w bazie danych z statusem "completed"
  4. Logowanie wyniku operacji

## Struktura raportu Excel

### Arkusz: "Raport Zamówień"

**Nagłówek:**
- Tytuł: "RAPORT ZAMÓWIEŃ"
- Okres: daty od-do
- Typ raportu

**Podsumowanie:**
- Łączna liczba zamówień
- Łączna kwota (PLN)

**Szczegóły według źródeł:**
| Źródło | Liczba zamówień | Kwota |
|---------|-------------------|-------|
| website | 45 | 12,345.67 |
| źródło_jeden | 23 | 5,678.90 |
| manual | 12 | 3,456.78 |

## Struktura projektu

```
raport-service/
├── cmd/
│   └── server/
│       └── main.go              # Entry point + cron jobs
├── internal/
│   ├── database/
│   │   └── connection.go        # Dual DB connection (orders + reports)
│   ├── handlers/
│   │   └── reports.go           # HTTP handlers dla endpointów
│   ├── models/
│   │   └── report.go            # Modele Report, OrderStats, SourceStat
│   └── services/
│       └── raport_service.go    # Logika biznesowa (stats, Excel generation)
├── middleware/
│   └── auth.go                  # Middleware autoryzacji
├── migrations/
│   └── create_tables.sql        # Schemat bazy reports_management
├── reports/                     # Katalog z wygenerowanymi plikami Excel
├── docker/
│   └── dockerfile               # Multi-stage build
└── go.mod                       # Zależności
```

## Technologie

- **Go 1.24** - Język programowania
- **Gorilla Mux** - HTTP router
- **Excelize** - Generowanie plików Excel
- **Cron v3** - Harmonogram zadań
- **PostgreSQL** - Dwie bazy danych (orders, reports)
- **Docker/Podman** - Konteneryzacja

## Zmienne środowiskowe

| Zmienna | Domyślna wartość | Opis |
|---------|-------------------|------|
| `PORT` | `8083` | Port HTTP servera |
| `ORDERS_DB_HOST` | `localhost` | Host bazy orders |
| `ORDERS_DB_PORT` | `5432` | Port bazy orders |
| `ORDERS_DB_NAME` | `orders_management` | Nazwa bazy orders |
| `ORDERS_DB_USER` | `postgres` | Użytkownik bazy orders |
| `ORDERS_DB_PASSWORD` | `password` | Hasło bazy orders |
| `RAPORTS_DB_HOST` | `localhost` | Host bazy reports |
| `RAPORTS_DB_PORT` | `5434` | Port bazy reports |
| `RAPORTS_DB_NAME` | `reports_management` | Nazwa bazy reports |
| `RAPORTS_DB_USER` | `postgres` | Użytkownik bazy reports |
| `RAPORTS_DB_PASSWORD` | `password` | Hasło bazy reports |

## Endpointy API

### Chronione (tylko admin)
- `POST /api/reports/generate` - Generowanie nowego raportu
- `GET /api/reports/stats` - Szybkie statystyki (bez pliku)

### Publiczne
- `GET /health` - Health check

## Baza danych

### Tabela: reports
```sql
CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,           -- daily, weekly, monthly, custom
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    total_orders INTEGER NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    file_path TEXT,
    status VARCHAR(50) DEFAULT 'pending', -- pending, completed, failed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Tabela: report_sources
```sql
CREATE TABLE report_sources (
    id SERIAL PRIMARY KEY,
    report_id INTEGER REFERENCES reports(id) ON DELETE CASCADE,
    source_name VARCHAR(100) NOT NULL,
    order_count INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL
);
```

## Cron Jobs

### Konfiguracja
```go
// Testowe: co 5 minut
c.AddFunc("*/5 * * * *", generateWeeklyReport)

// Produkcyjne: co poniedziałek o północy
c.AddFunc("0 0 * * 1", generateWeeklyReport)
```

### Format cron
```
┌───────── minuta (0 - 59)
│ ┌───────── godzina (0 - 23)
│ │ ┌───────── dzień miesiąca (1 - 31)
│ │ │ ┌───────── miesiąc (1 - 12)
│ │ │ │ ┌───────── dzień tygodnia (0 - 6) (0=niedziela)
│ │ │ │ │
* * * * *
```

## Deployment

### Docker
```bash
cd services/raport-service
docker build -t raport-service -f docker/dockerfile .
docker run -p 8083:8083 \
  -e ORDERS_DB_HOST="postgres-orders" \
  -e RAPORTS_DB_HOST="postgres-reports" \
  -v $(pwd)/reports:/app/reports \
  raport-service
```

### Volume dla raportów
Ważne: Katalog `./reports` powinien być zmontowany jako volume, aby pliki Excel były dostępne po restarcie kontenera.

## Przykładowe użycie

### Generowanie raportu miesięcznego
```bash
curl -X POST http://localhost:8083/api/reports/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "type": "monthly",
    "period_start": "2025-11-01T00:00:00Z",
    "period_end": "2025-11-30T23:59:59Z"
  }'
```

### Szybkie statystyki (ostatnie 7 dni)
```bash
curl http://localhost:8083/api/reports/stats \
  -H "Authorization: Bearer <admin_token>"
```

### Szybkie statystyki (niestandardowy zakres)
```bash
curl "http://localhost:8083/api/reports/stats?start=2025-11-01&end=2025-11-15" \
  -H "Authorization: Bearer <admin_token>"
```

## Odpowiedź API

### Przykładowa odpowiedź z /generate
```json
{
  "id": 42,
  "type": "weekly",
  "period_start": "2025-11-14T00:00:00Z",
  "period_end": "2025-11-21T23:59:59Z",
  "total_orders": 87,
  "total_amount": 23456.78,
  "file_path": "./reports/weekly_raport_2025_11_21_10_30_45.xlsx",
  "status": "completed",
  "created_at": "2025-11-21T10:30:45Z",
  "sources": [
    {
      "source_name": "website",
      "order_count": 45,
      "amount": 12345.67
    },
    {
      "source_name": "źródło_jeden",
      "order_count": 23,
      "amount": 5678.90
    },
    {
      "source_name": "manual",
      "order_count": 19,
      "amount": 5432.21
    }
  ]
}
```

## Bezpieczeństwo

1. **Dostęp tylko dla adminów** - middleware sprawdza rolę użytkownika
2. **Read-only access** - baza orders_management dostępna tylko do odczytu
3. **Walidacja dat** - sprawdzenie poprawności zakresów czasowych
4. **Sanityzacja nazw plików** - timestamp w nazwie zapobiega konfliktom

## Integracja z systemem

### Order Service
- Raport-service zczytuje dane z bazy `orders_management`
- Dostęp READ ONLY
- Automatyczna synchronizacja danych (real-time)

### Auth Service
- Autoryzacja przez Nginx auth_request
- Wymaga roli `admin`

### Frontend
- Endpoint do pobierania raportów
- Wyświetlanie statystyk
- Download wygenerowanych plików Excel
