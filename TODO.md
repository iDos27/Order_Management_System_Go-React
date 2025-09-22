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
                            ┌────────────┼────────────┐
                            │            │            │
                    ┌───────▼───┐   ┌────▼────┐   ┌───▼─────┐
                    │   Auth    │   │ Orders  │   │Analytics│
                    │  Service  │   │ Service │   │ Service │
                    │           │   │(Obecny) │   │         │
                    └───────────┘   └────┬────┘   └─────────┘
                                         │
                                   ┌─────▼──────┐
                                   │Notification│
                                   │  Service   │
                                   └────────────┘
                            
            ┌─────────────────────────────────────────────────────────────┐
            │                     WARSTWA DANYCH                          │
            ├─────────────────┬─────────────────┬─────────────────────────┤
            │   PostgreSQL    │      Redis      │       RabbitMQ          │
            │   (Orders DB)   │   (Cache/Auth)  │    (Notifications)      │
            └─────────────────┴─────────────────┴─────────────────────────┘

                                 KUBERNETES CLUSTER
             ┌─────────────────────────────────────────────────────────┐
             │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────┐    │
             │  │  Auth   │ │ Orders  │ │Analytics│ │Notification │    │
             │  │   Pod   │ │   Pod   │ │   Pod   │ │    Pod      │    │
             │  └─────────┘ └─────────┘ └─────────┘ └─────────────┘    │
             │                                                         │
             │  ┌─────────┐ ┌─────────┐ ┌─────────────────────────┐    │
             │  │  Redis  │ │ Postgres│ │        RabbitMQ         │    │
             │  │   Pod   │ │   Pod   │ │          Pod            │    │
             │  └─────────┘ └─────────┘ └─────────────────────────┘    │
             └─────────────────────────────────────────────────────────┘
```

## PLAN ROZWOJU - LISTA ZADAŃ

### 1: PRZYGOTOWANIE PROJEKTU
- [ ] **Utworzenie branchu `feature/microservices`**
- [ ] **Restructuryzacja folderów** - podział na services/
- [ ] **Docker Compose** dla lokalnego developmentu

### 2: AUTHENTICATION SERVICE
- [ ] **Auth Service** - Autoryzacja użytkowników
- [ ] **Login endpoint** - POST /auth/login z email/password albo nr.zamówienia dla klientów
- [ ] **Register endpoint** - POST /auth/register
- [ ] **Token validation** - middleware dla innych serwisów
- [ ] **Role management** - przypisywanie ról użytkownikom
- [ ] **Password hashing** - bcrypt dla bezpieczeństwa haseł
- [ ] **Redis integration** - cache dla sesji użytkowników

### 3: NOTIFICATION SERVICE  
- [ ] **Notification Service** - Powiadomienia systemowe Linux o nowym zamówieniu
- [ ] **Desktop notifications** - libnotify dla natywnych powiadomień systemowych
- [ ] **Click to open** - kliknięcie w powiadomienie otwiera aplikację w przeglądarce
- [ ] **Webhook endpoints** - odbieranie eventów o zmianach statusu z Orders Service
- [ ] **Notification formatting** - czytelne komunikaty o statusach zamówień
- [ ] **Queue system** - Redis/RabbitMQ dla asynchronicznych powiadomień
- [ ] **Service integration** - komunikacja z Orders Service przez HTTP/WebSocket

### 4: ANALYTICS SERVICE (???)
- [ ] **Analytics Service** - Prosta analiza danych zamówień
- [ ] **Cyclical reports** - raporty cykliczne w formacie TXT
- [ ] **Order statistics** - podstawowe statystyki (liczba, statusy, trendy)
- [ ] **Data aggregation** - zbieranie danych z Orders Service przez API
- [ ] **Simple calculations** - suma zamówień, średnia wartość, najpopularniejsze statusy
- [ ] **Cron jobs** - automatyczne generowanie cyklicznych raportów
- [ ] **File storage** - zapisywanie raportów w folderze /reports

### 5: INTEGRACJA SERWISÓW
- [ ] **API Gateway** - Nginx reverse proxy dla routingu
- [ ] **Service discovery** - automatyczne wykrywanie serwisów
- [ ] **Inter-service communication** - HTTP calls między serwisami
- [ ] **Event-driven architecture** - komunikacja przez eventy
- [ ] **Error handling** - obsługa błędów między serwisami
- [ ] **Distributed logging** - centralne logowanie wszystkich serwisów

### 6: KONTENERYZACJA
- [ ] **Dockerfile dla serwisu** - Auth, Notification, Analytics
- [ ] **Docker Compose** - kompletne środowisko developerskie
- [ ] **Multi-stage builds** - optymalizacja rozmiarów obrazów
- [ ] **Health checks** - sprawdzanie stanu kontenerów
- [ ] **Environment variables** - konfiguracja przez zmienne środowiskowe

### 7: KUBERNETES
- [ ] **Kubernetes manifests** - Deployment, Service, ConfigMap
- [ ] **Namespaces** - podział na środowiska (dev/staging/prod)
- [ ] **Ingress controller** - routing ruchu zewnętrznego
- [ ] **Persistent volumes** - trwałe przechowywanie danych
- [ ] **Secrets management** - bezpieczne przechowywanie haseł/kluczy
- [ ] **Auto-scaling** - automatyczne skalowanie podów
- [ ] **Rolling updates** - bezpieczne wdrażanie nowych wersji

### 8: MONITORING I OBSERVABILITY (???)
- [ ] **Prometheus** - zbieranie metryk z aplikacji
- [ ] **Grafana** - dashboardy i wizualizacja metryk
- [ ] **Health endpoints** - /health dla każdego serwisu
- [ ] **Application metrics** - custom metryki biznesowe
- [ ] **Alerting** - powiadomienia o problemach
- [ ] **Log aggregation** - centralne zbieranie logów

### 9: ŚRODOWISKO (Nie wymagane)
- [] **Terraform** - tworzenie maszyny wirtualnej na serwerze
- [] **Ansible** - konfiguracja systemu pod kątem bezpieczeństwa

## INFORMACJE TECHNICZNE

### **Technologie:**
- **Backend**: Go (Gin framework) dla wszystkich serwisów
- **Auth**: JWT tokens + Redis cache
- **Queue**: RabbitMQ dla asynchronicznych zadań  
- **Database**: PostgreSQL (główna) + Redis (cache)
- **Containers**: Docker + Docker Compose
- **Orchestration**: Kubernetes
- **Monitoring**: Prometheus + Grafana
- **API Gateway**: Nginx

### **Porty serwisów:**
- **Orders Service**: 8080
- **Auth Service**: 8081
- **Notification Service**: 8082 (Linux desktop notifications)
- **Analytics Service**: 8083 (TXT reports generation)
- **API Gateway**: 80/443

### **Bazy danych:**
- **orders_db**: PostgreSQL dla Orders Service
- **auth_db**: PostgreSQL dla Auth Service
- **notifications_log**: SQLite/PostgreSQL dla historii powiadomień (opcjonalne)
- **analytics_cache**: Redis dla cache'owania obliczeń (opcjonalne)
- **redis**: Cache dla sesji i temp data

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