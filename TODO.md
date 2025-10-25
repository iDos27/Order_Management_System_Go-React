# Rozwój do mikrousług

## SCHEMAT INFRASTRUKTURY

```
                    ┌─────────────────┐    ┌─────────────────┐
                    │   Frontend      │    │   Admin Panel   │
                    │   (Customer)    │    │   (Warehouse)   │
                    └─────────┬───────┘    └─────────┬───────┘
                              │                      │
                              └──────────┬───────────┘
                                         │
                                ┌────────▼────────┐
                                │  API Gateway    │ 
                                │   (Nginx)       │
                                └─────────────────┘
                                         │
                            ┌────────────┼────────────┬────────────┐
                            │            │            │            │
                    ┌───────▼───┐   ┌────▼────┐   ┌───▼─────┐  ┌───▼─────┐
                    │   Auth    │   │ Orders  │   │Analytics│  │  Mail   │
                    │  Service  │   │ Service │   │ Service │  │ Service │
                    │           │   │(Obecny) │   │         │  │         │
                    └───────────┘   └────┬────┘   └─────────┘  └─────────┘
                                         │
                                   ┌─────▼──────┐
                                   │Notification│
                                   │  Service   │
                                   └────────────┘
                            
            ┌─────────────────────────────────────────────────────────────┐
            │                     WARSTWA DANYCH                          │
            ├─────────────────┬─────────────────┬─────────────────────────┤
            │   PostgreSQL    │    PostgreSQL   │       RabbitMQ          │
            │   (Orders DB)   │   (Cache/Auth)  │    (Notifications)      │
            │    Port 5432    │    Port 5433    │       Port 5434         │
            └─────────────────┴─────────────────┴─────────────────────────┘

                                 KUBERNETES CLUSTER
             ┌─────────────────────────────────────────────────────────┐
             │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────┐    │
             │  │  Auth   │ │ Orders  │ │Analytics│ │Notification │    │
             │  │   Pod   │ │   Pod   │ │   Pod   │ │    Pod      │    │
             │  └─────────┘ └─────────┘ └─────────┘ └─────────────┘    │
             │                                                         │
             │  ┌─────────┐ ┌─────────┐ ┌─────────────────────────┐    │
             │  │ Postgres│ │ Postgres│ │        RabbitMQ         │    │
             │  │Order DB │ │ Auth DB │ │          Pod            │    │
             │  └─────────┘ └─────────┘ └─────────────────────────┘    │
             └─────────────────────────────────────────────────────────┘
```

##  WAŻNE DECYZJE DO PODJĘCIA

### ŚRODOWISKO

#### Opcja 1: Jedna VM + Docker Compose
- **Terraform** tworzy VM KVM (jako zastępstwo Proxmox) [openSUSE]
- **Ansible** instaluje Docker i uruchamia Docker Compose
#### Opcja 2: Mały klaster Kubernetes
- **Terraform** stawia kilka VM KVM
- Instalowanie **microKube**
- Każdy mikroserwis i baza jako Deployment itd.


## PLAN ROZWOJU - LISTA ZADAŃ

### 1: PRZYGOTOWANIE PROJEKTU
- [x] **Utworzenie branchu `feature/microservices`**
- [x] **Restructuryzacja folderów** - podział na services/

### 2: AUTHENTICATION SERVICE
- [x] **Auth Service** - Autoryzacja użytkowników
- [x] **Login endpoint** - POST /auth/login
- [x] **Register endpoint** - POST /auth/register
- [x] **Password hashing** - bcrypt
- [x] **JWT tokens** - generowanie i walidacja
- [x] **PostgreSQL integration** - baza użytkowników (port 5433)
- [x] **CORS middleware** - obsługa żądań cross-origin
- [x] **Health endpoint** - GET /api/v1/health
- [ ] **Token validation middleware** - middleware dla innych serwisów
- [ ] **Role management** - przypisywanie ról użytkownikom (admin/employee/customer)
- [ ] **Session management** - tracking aktywnych sesji (PostgreSQL)
- [ ] **Token blacklist** - możliwość unieważnienia tokenów
- [ ] **Database migrations** - skrypt SQL dla nowych tabel
- [ ] **Rate limiting** - ograniczenie prób logowania
- [ ] **Input validation** - walidacja email, siła hasła
- [ ] **Logout endpoint** - POST /logout z blacklist tokenów


### 3: NOTIFICATION SERVICE  
- [ ] **Notification Service** - Powiadomienia systemowe Linux o nowym zamówieniu
- [ ] **Desktop notifications** - libnotify dla natywnych powiadomień systemowych
- [ ] **Click to open** - kliknięcie w powiadomienie otwiera aplikację w przeglądarce
- [ ] **Webhook endpoints** - odbieranie eventów o zmianach statusu z Orders Service
- [ ] **Notification formatting** - czytelne komunikaty o statusach zamówień
- [ ] **Queue system** - RabbitMQ dla asynchronicznych powiadomień
- [ ] **Service integration** - komunikacja z Orders Service przez HTTP/WebSocket

### 4: INTEGRACJA SERWISÓW
- [ ] **API Gateway** - Nginx reverse proxy dla routingu
- [ ] **Service discovery** - automatyczne wykrywanie serwisów
- [ ] **Inter-service communication** - HTTP calls między serwisami
- [ ] **Event-driven architecture** - komunikacja przez eventy
- [ ] **Error handling** - obsługa błędów między serwisami
- [ ] **Distributed logging** - centralne logowanie wszystkich serwisów

### 5: KONTENERYZACJA
- [ ] **Dockerfile dla każdego serwisu** - Auth, Orders, Notification, Analytics
- [ ] **Docker Compose** - kompletne środowisko developerskie
- [ ] **Multi-stage builds** - optymalizacja rozmiarów obrazów

### 6: KUBERNETES
- [ ] **Kubernetes manifests** - Deployment, Service, ConfigMap
- [ ] **Namespaces** - podział na środowiska (dev/staging/prod)
- [ ] **Ingress controller** - routing ruchu zewnętrznego
- [ ] **Persistent volumes** - trwałe przechowywanie danych
- [ ] **Secrets management** - bezpieczne przechowywanie haseł/kluczy
- [ ] **Auto-scaling** - automatyczne skalowanie podów
- [ ] **Rolling updates** - bezpieczne wdrażanie nowych wersji


## INFORMACJE TECHNICZNE

### **Technologie:**
- **Backend**: Go (Gin framework) dla wszystkich serwisów
- **Auth**: JWT tokens 
- **Queue**: RabbitMQ dla asynchronicznych zadań  
- **Database**: PostgreSQL (główna)
- **Containers**: Docker + Docker Compose
- **Orchestration**: Kubernetes
- **Monitoring**: Prometheus + Grafana
- **API Gateway**: Nginx

### **Porty serwisów:**
- **Orders Service**: 8080
- **Auth Service**: 8081
- **Notification Service**: 8082
- **Analytics Service**: 8083
- **Mail Service** 8084
- **API Gateway**: 80/443

### **Bazy danych:**
- **orders_db**: PostgreSQL dla Orders Service
- **auth_db**: PostgreSQL dla Auth Service
- **notifications_log**: SQLite/PostgreSQL dla historii powiadomień (opcjonalne)
- **rabbitmq**: Queue dla powiadomień

## ZALETY

**Mikrousług:**
- **Separation of Concerns** - każdy serwis ma jedną odpowiedzialność
- **Scalability** - niezależne skalowanie komponentów
- **Fault Tolerance** - odporność na awarie serwisów
- **Technology Diversity** - różne technologie dla różnych problemów
- **DevOps** - CI/CD, monitoring, konteneryzacja
- **Cloud-Native** - gotowość na chmure

**Kubernetes:**
- **Container Orchestration** - zarządzanie kontenerami w skali
- **Service Mesh** - komunikacja między serwisami
- **Auto-scaling** - dostosowywanie zasobów  
- **Rolling Deployments** - bezpieczne wdrażanie aktualizacji
- **Infrastructure as Code** - definicja infrastruktury w kodzie