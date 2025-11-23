# Kubernetes/Podman Deployment

## Opis

Kompletna konfiguracja wdrożenia systemu zarządzania zamówieniami w środowisku Kubernetes/Podman. Plik YAML definiuje wszystkie komponenty systemu jako Deployment z automatycznym zarządzaniem replikacją i aktualizacjami.

## Architektura

### Bazy danych (PostgreSQL)
- **postgres-orders** - Baza zamówień (port 5432, 2Gi storage)
- **postgres-auth** - Baza użytkowników (port 5432, 1Gi storage)
- **postgres-reports** - Baza raportów (port 5432, 1Gi storage)

### Message Broker
- **rabbitmq** - RabbitMQ (AMQP 5672, Management 15672, 1Gi storage)

### Mikroserwisy
- **auth-service** - Autoryzacja (port 8081)
- **order-service** - Zarządzanie zamówieniami (port 8080)
- **raport-service** - Generowanie raportów (port 8083, 5Gi storage dla plików Excel)

### Frontend i Proxy
- **frontend** - React/Vite (port 5173)
- **nginx** - Reverse proxy (port 80, NodePort 30080)

## Przed uruchomieniem

### 1. Zbuduj obrazy Docker

```bash
# Auth Service
cd services/auth-service
podman build -t auth-service:latest -f docker/dockerfile .

# Order Service
cd ../order-service
podman build -t order-service:latest -f docker/dockerfile .

# Raport Service
cd ../raport-service
podman build -t raport-service:latest -f docker/dockerfile .

cd ../..
```

### 2. Zainstaluj zależności frontendu

```bash
cd frontend/admin-panel
npm install
cd ../..
```

## Uruchomienie

### Podman Play Kube

```bash
# Uruchom cały system
podman play kube kubernetes/order-management-system.yaml

# Sprawdź status podów
podman pod ps

# Sprawdź wszystkie kontenery
podman ps

# Logi z konkretnego kontenera
podman logs <container_id>
```

### Zatrzymanie

```bash
# Zatrzymaj i usuń wszystkie pody
podman play kube --down kubernetes/order-management-system.yaml
```

### Restart (Replace)

```bash
# Zatrzymaj stare i uruchom nowe w jednej komendzie
podman play kube --replace kubernetes/order-management-system.yaml
```

## Dostęp do aplikacji

### Główne endpointy
- **Frontend**: http://localhost:30080
- **API (przez Nginx)**: http://localhost:30080/api
- **RabbitMQ Management**: http://localhost:30672

### Bezpośrednie porty serwisów (ClusterIP)
Te porty są dostępne tylko wewnątrz sieci Podman:
- Auth Service: http://auth-service:8081
- Order Service: http://order-service:8080
- Raport Service: http://raport-service:8083

## Inicjalizacja baz danych

Po pierwszym uruchomieniu musisz uruchomić migracje dla każdego serwisu.

### 1. Znajdź ID kontenerów

```bash
podman ps | grep postgres
podman ps | grep service
```

### 2. Wykonaj migracje

```bash
# Auth Service
AUTH_CONTAINER=$(podman ps --filter name=auth-service --format "{{.ID}}")
podman exec -i $(podman ps --filter name=postgres-auth --format "{{.ID}}") \
  psql -U postgres -d auth_service < services/auth-service/migrations/create_tables.sql

# Order Service
podman exec -i $(podman ps --filter name=postgres-orders --format "{{.ID}}") \
  psql -U postgres -d orders_management < services/order-service/migrations/create_tables.sql

# Raport Service
podman exec -i $(podman ps --filter name=postgres-reports --format "{{.ID}}") \
  psql -U postgres -d reports_management < services/raport-service/migrations/create_tables.sql
```

### 3. Utwórz użytkowników

```bash
# Admin
curl -X POST http://localhost:30080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123","role":"admin"}'

# Employee
curl -X POST http://localhost:30080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"gosc@test.com","password":"gosc123","role":"employee"}'
```

## Persistent Storage

System używa PersistentVolumeClaim dla trwałego przechowywania danych:

- **postgres-orders-pvc** - 2Gi dla zamówień
- **postgres-auth-pvc** - 1Gi dla użytkowników
- **postgres-reports-pvc** - 1Gi dla metadanych raportów
- **rabbitmq-pvc** - 1Gi dla kolejki wiadomości
- **reports-files-pvc** - 5Gi dla wygenerowanych plików Excel

Podman automatycznie tworzy volumes w `~/.local/share/containers/storage/volumes/`

## ConfigMap

Nginx konfiguracja jest przechowywana w ConfigMap `nginx-config`, co pozwala na łatwą aktualizację bez rebuildu obrazu:

```bash
# Edytuj konfigurację
kubectl edit configmap nginx-config  # lub ręcznie w YAML

# Przeładuj nginx
podman exec <nginx_container_id> nginx -s reload
```

## Health Checks

Wszystkie mikroserwisy mają skonfigurowane:
- **Liveness Probe** - sprawdza czy kontener żyje (restart przy błędzie)
- **Readiness Probe** - sprawdza czy kontener gotowy przyjmować ruch

Endpointy health check:
- `/health` - zwraca `{"status":"ok","service":"<nazwa>"}`

## Zmienne środowiskowe

### Produkcja - zmień przed wdrożeniem!

```yaml
# Auth Service
- name: JWT_SECRET
  value: "production-secret-change-me"  # ⚠️ ZMIEŃ NA LOSOWY SEKRET!

# Wszystkie bazy danych
- name: POSTGRES_PASSWORD
  value: "password"  # ⚠️ UŻYJ SILNEGO HASŁA!
```

### Użycie Secrets (zalecane dla produkcji)

Zamiast plaintext w Deployment, użyj Kubernetes Secrets:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: database-secrets
type: Opaque
data:
  postgres-password: cGFzc3dvcmQ=  # base64 encoded
  jwt-secret: c2VjcmV0LWtleQ==

---
# W Deployment:
env:
- name: POSTGRES_PASSWORD
  valueFrom:
    secretKeyRef:
      name: database-secrets
      key: postgres-password
```

## Skalowanie

Podman play kube wspiera pole `replicas` w Deployment:

```yaml
spec:
  replicas: 3  # 3 instancje auth-service
```

**Uwaga dla baz danych:** PostgreSQL nie wspiera natywnie replikacji read-write w Kubernetes. Dla produkcji użyj operatorów jak Zalando Postgres Operator.

## Troubleshooting

### Sprawdzenie logów

```bash
# Wszystkie pody
podman pod logs <pod_name>

# Konkretny kontener
podman logs <container_id>

# Follow logs (real-time)
podman logs -f <container_id>
```

### Sprawdzenie sieci

```bash
# Lista sieci
podman network ls

# Inspekcja sieci
podman network inspect <network_name>
```

### Dostęp do bazy danych

```bash
# PostgreSQL Orders
podman exec -it $(podman ps --filter name=postgres-orders --format "{{.ID}}") \
  psql -U postgres -d orders_management

# PostgreSQL Auth
podman exec -it $(podman ps --filter name=postgres-auth --format "{{.ID}}") \
  psql -U postgres -d auth_service
```

### Restart pojedynczego serwisu

```bash
# Znajdź deployment
podman ps | grep auth-service

# Usuń kontener (Deployment automatycznie utworzy nowy)
podman stop <container_id>
podman rm <container_id>
```

## Migracja do Kubernetes

Ten sam plik YAML działa na prawdziwym klastrze Kubernetes:

```bash
# Minikube
minikube start
kubectl apply -f kubernetes/order-management-system.yaml

# K3s
kubectl apply -f kubernetes/order-management-system.yaml

# GKE, EKS, AKS
kubectl apply -f kubernetes/order-management-system.yaml
```

**Różnice:**
- NodePort zamiast LoadBalancer (zmień `type: NodePort` na `type: LoadBalancer`)
- PersistentVolume może wymagać StorageClass
- Secrets zarządzane przez Kubernetes
- Ingress zamiast NodePort dla produkcji

## Backup i Restore

### Backup baz danych

```bash
# Orders
podman exec $(podman ps --filter name=postgres-orders --format "{{.ID}}") \
  pg_dump -U postgres orders_management > backup_orders.sql

# Auth
podman exec $(podman ps --filter name=postgres-auth --format "{{.ID}}") \
  pg_dump -U postgres auth_service > backup_auth.sql

# Reports
podman exec $(podman ps --filter name=postgres-reports --format "{{.ID}}") \
  pg_dump -U postgres reports_management > backup_reports.sql
```

### Restore

```bash
podman exec -i $(podman ps --filter name=postgres-orders --format "{{.ID}}") \
  psql -U postgres -d orders_management < backup_orders.sql
```

## Monitoring

### Zasoby

```bash
# Zużycie CPU/RAM
podman stats

# Szczegóły konkretnego kontenera
podman stats <container_id>
```

### Health status

```bash
# Sprawdź wszystkie health endpoints
curl http://localhost:30080/api/v1/health  # Auth przez Nginx
curl http://auth-service:8081/health       # Bezpośrednio (jeśli w sieci)
curl http://order-service:8080/health
curl http://raport-service:8083/health
```

## Bezpieczeństwo

### Checklist produkcyjny

- [ ] Zmień `POSTGRES_PASSWORD` na silne hasła
- [ ] Zmień `JWT_SECRET` na losowy sekret (min. 32 znaki)
- [ ] Użyj Kubernetes Secrets zamiast plaintext
- [ ] Skonfiguruj SSL/TLS dla Nginx
- [ ] Ogranicz dostęp do RabbitMQ Management (15672)
- [ ] Włącz Network Policies dla izolacji serwisów
- [ ] Skonfiguruj RBAC dla dostępu do klastra
- [ ] Regularny backup baz danych
- [ ] Monitoring i alerting (Prometheus + Grafana)

## Autorzy

- Piotr - System zarządzania zamówieniami
- Data utworzenia: 23 listopada 2025

## Licencja

Projekt studencki - Praca Inżynierska
