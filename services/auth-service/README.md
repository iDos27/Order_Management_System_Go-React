# Auth Service

## Opis

Mikroserwis odpowiedzialny za uwierzytelnianie i autoryzację użytkowników w systemie zarządzania zamówieniami. Realizuje mechanizmy rejestracji, logowania oraz walidacji tokenów JWT. Integruje się z innymi serwisami poprzez middleware autoryzacyjny obsługiwany przez Nginx (auth_request).

## Architektura

### Port
- **8081** - HTTP Server

### Baza danych
- **PostgreSQL** (port 5433)
- **Nazwa bazy:** `auth_service`
- **Tabela:** `users`

## Funkcjonalności

### 1. Rejestracja użytkownika (`POST /api/v1/register`)
- Hashowaniehasła za pomocą **bcrypt**
- Obsługa ról: `admin`, `employee`, `customer`
- Walidacja unikalnego adresu email
- Zwraca ID nowo utworzonego użytkownika

### 2. Logowanie (`POST /api/v1/login`)
- Weryfikacja hasła z wykorzystaniem bcrypt
- Generowanie tokenu **JWT** (JSON Web Token)
- Ważność tokenu: **24 godziny**
- Claims zawierają: `user_id`, `email`, `role`, `exp`

### 3. Walidacja tokenu (`GET /api/v1/validate`)
- Endpoint dla Nginx **auth_request**
- Weryfikacja tokenu JWT z nagłówka Authorization
- Sprawdzenie istnienia użytkownika w bazie danych
- Zwraca status 200 (OK) lub 401 (Unauthorized)

### 4. Profil użytkownika (`GET /api/v1/profile`)
- Chroniony endpoint (wymaga tokenu JWT)
- Zwraca szczegóły zalogowanego użytkownika
- Informacje: ID, email, rola, daty utworzenia/aktualizacji

### 5. Middleware autoryzacji
- **RequireAuth()** - wymaga ważnego tokenu JWT
- **RequireRole()** - wymaga określonej roli
- **RequireAdmin()** - dostęp tylko dla administratorów
- **RequireAdminOrEmployee()** - dostęp dla admin/pracowników

## Struktura projektu

```
auth-service/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikacji
├── internal/
│   ├── database/
│   │   └── connection.go        # Połączenie z PostgreSQL
│   ├── handlers/
│   │   └── handlers.go          # Obsługa endpointów (register, login, profile)
│   └── models/
│       └── user.go              # Modele danych (User, LoginRequest, etc.)
├── middleware/
│   └── auth.go                  # Middleware autoryzacji (JWT validation)
├── migrations/
│   └── create_tables.sql        # Schemat bazy danych
├── docker/
│   └── dockerfile               # Multi-stage build dla kontenera
└── go.mod                       # Zależności Go
```

## Technologie

- **Go 1.24** - Język programowania
- **Gin** - Framework HTTP
- **PostgreSQL** - Baza danych
- **JWT (golang-jwt/jwt)** - Tokeny autoryzacyjne
- **bcrypt** - Hashowanie hasła
- **Docker/Podman** - Konteneryzacja

## Zmienne środowiskowe

| Zmienna | Domyślna wartość | Opis |
|---------|-------------------|------|
| `DATABASE_URL` | `postgres://postgres:password@localhost:5433/auth_service?sslmode=disable` | URL bazy danych |
| `JWT_SECRET` | `secret-key` | Sekret do podpisywania tokenów JWT |
| `PORT` | `8081` | Port HTTP servera |

## Endpointy API

### Publiczne
- `POST /api/v1/register` - Rejestracja nowego użytkownika
- `POST /api/v1/login` - Logowanie (zwraca token JWT)
- `GET /health` - Health check

### Chronione (wymagają tokenu JWT)
- `GET /api/v1/profile` - Profil zalogowanego użytkownika
- `GET /api/v1/validate` - Walidacja tokenu (używane przez Nginx)
- `GET /api/v1/validate-token` - Szczegółowa walidacja tokenu z zwrotem claimów

### Administracyjne (tylko admin)
- `GET /api/v1/admin/users` - Lista użytkowników

## Bezpieczeństwo

1. **Hashowanie hasła** - bcrypt z domyślnym kosztem (10 rund)
2. **JWT signing** - HMAC-SHA256
3. **CORS** - Skonfigurowany middleware z obsługą preflight
4. **Non-root user** - Kontener uruchamia aplikację jako `appuser` (UID 1000)
5. **Walidacja ról** - Middleware sprawdza uprawnienia przed dostępem do zasobów

## Integracja z systemem

### Nginx auth_request
Auth-service służy jako backend dla mechanizmu `auth_request` w Nginx:

```nginx
location /api/orders {
    auth_request /validate;
    proxy_pass http://order-service;
}

location = /validate {
    internal;
    proxy_pass http://auth-service/api/v1/validate;
}
```

### Inne mikroserwisy
Każdy chroniony endpoint w innych serwisach powinien:
1. Odebrać token z nagłówka `Authorization: Bearer <token>`
2. Zweryfikować token przez wywołanie `/api/v1/validate-token`
3. Użyć claimów (user_id, role) do kontroli dostępu

## Deployment

### Docker
```bash
cd services/auth-service
docker build -t auth-service -f docker/dockerfile .
docker run -p 8081:8081 \
  -e DATABASE_URL="postgres://..." \
  -e JWT_SECRET="production-secret" \
  auth-service
```

### Health Check
```bash
curl http://localhost:8081/health
# Odpowiedź: {"status":"ok","service":"auth-service"}
```

## Role użytkowników

- **admin** - Pełny dostęp do wszystkich zasobów
- **employee** - Dostęp do zarządzania zamówieniami
- **customer** - Dostęp do własnych zamówień
