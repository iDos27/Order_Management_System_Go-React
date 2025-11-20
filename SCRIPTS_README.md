# 🚀 Skrypty Uruchomieniowe Systemu

## Szybki Start

### Uruchomienie całego środowiska
```bash
./makeEnv.sh
```

Ten skrypt automatycznie:
1. ✅ Sprawdza i czyści istniejące kontenery
2. 🐘 Uruchamia 3 bazy danych PostgreSQL (porty 5432, 5433, 5434)
3. 📊 Wykonuje migracje schematów baz danych
4. 🌐 Uruchamia Nginx Gateway (port 80)
5. 🔧 Uruchamia wszystkie serwisy Go w osobnych terminalach:
   - Auth Service (port 8081)
   - Order Service (port 8080)
   - Raport Service (port 8083)

### Zatrzymanie środowiska
```bash
./stopEnv.sh
```

Ten skrypt:
1. 🛑 Zatrzymuje i usuwa wszystkie kontenery
2. 🔧 Zabija procesy serwisów Go
3. 🧹 Czyści pliki logów

---

## Wymagania

### Zainstalowane narzędzia:
- **Podman** - konteneryzacja (zamiennik Dockera)
- **Go 1.21+** - kompilacja serwisów
- **Terminal graficzny** (jeden z):
  - `konsole` (KDE) - preferowany
  - `gnome-terminal` (GNOME)
  - `xterm` (fallback)
  - Jeśli brak - serwisy działają w tle z logami w `logs/`

### Sprawdzenie instalacji:
```bash
podman --version    # Powinno zwrócić wersję Podman
go version          # Powinno zwrócić wersję Go
konsole --version   # Opcjonalnie - sprawdź terminal
```

---

## Struktura Portów

| Komponent            | Port  | Protokół | Opis                        |
|---------------------|-------|----------|-----------------------------|
| **Nginx Gateway**    | 80    | HTTP     | Reverse Proxy               |
| **Order Service**    | 8080  | HTTP     | API zamówień                |
| **Auth Service**     | 8081  | HTTP     | API autoryzacji             |
| **Raport Service**   | 8083  | HTTP     | API raportów                |
| **PostgreSQL Orders**| 5432  | TCP      | Baza zamówień               |
| **PostgreSQL Auth**  | 5433  | TCP      | Baza użytkowników           |
| **PostgreSQL Raports**| 5434 | TCP      | Baza raportów               |
| **Frontend (dev)**   | 5173  | HTTP     | React (uruchom ręcznie)     |

---

## Uruchomienie Frontendu

Po uruchomieniu `makeEnv.sh`, frontend uruchamiasz ręcznie:

```bash
cd frontend/admin-panel
npm install    # Tylko przy pierwszym uruchomieniu
npm run dev    # Uruchomienie serwera deweloperskiego
```

Frontend będzie dostępny pod: **http://localhost:5173**

---

## Debugowanie

### Sprawdzenie działających kontenerów
```bash
podman ps
```

### Logi kontenerów
```bash
podman logs postgres-orders
podman logs postgres-auth
podman logs postgres-raports
podman logs nginx-gateway
```

### Logi serwisów Go (jeśli działają w tle)
```bash
tail -f logs/Auth\ Service.log
tail -f logs/Order\ Service.log
tail -f logs/Raport\ Service.log
```

### Testowanie połączeń z bazami
```bash
# Orders
podman exec -it postgres-orders psql -U postgres -d orders_management

# Auth
podman exec -it postgres-auth psql -U postgres -d auth_service

# Raports
podman exec -it postgres-raports psql -U postgres -d raports_management
```

### Sprawdzenie czy porty są zajęte
```bash
ss -tuln | grep -E ':(80|8080|8081|8083|5432|5433|5434)'
```

---

## Rozwiązywanie Problemów

### Problem: Port już zajęty
```bash
# Znajdź proces blokujący port
sudo lsof -i :8080

# Zabij proces (zmień PID)
kill -9 <PID>
```

### Problem: Kontenery nie startują
```bash
# Sprawdź logi błędów
podman logs <nazwa_kontenera>

# Usuń wszystkie kontenery i uruchom ponownie
./stopEnv.sh
./makeEnv.sh
```

### Problem: Baza danych niedostępna
```bash
# Sprawdź czy kontener działa
podman ps | grep postgres

# Sprawdź logi bazy
podman logs postgres-orders

# Zrestartuj kontener
podman restart postgres-orders
```

### Problem: Serwis Go nie uruchamia się
```bash
# Sprawdź czy Go jest zainstalowane
go version

# Sprawdź błędy kompilacji
cd services/auth-service
go run cmd/server/main.go
```

---

## Przydatne Komendy

### Pełne wyczyszczenie środowiska
```bash
./stopEnv.sh
podman system prune -a --volumes  # UWAGA: usuwa WSZYSTKIE dane Podman!
./makeEnv.sh
```

### Backup baz danych
```bash
# Orders
podman exec postgres-orders pg_dump -U postgres orders_management > backup_orders.sql

# Auth
podman exec postgres-auth pg_dump -U postgres auth_service > backup_auth.sql

# Raports
podman exec postgres-raports pg_dump -U postgres raports_management > backup_raports.sql
```

### Restore baz danych
```bash
# Orders
podman exec -i postgres-orders psql -U postgres -d orders_management < backup_orders.sql

# Auth
podman exec -i postgres-auth psql -U postgres -d auth_service < backup_auth.sql

# Raports
podman exec -i postgres-raports psql -U postgres -d raports_management < backup_raports.sql
```

---

## Architektura Skryptu makeEnv.sh

```
1. Cleanup ━━━━━━━━━━━━━━━━━━━━━━┐
   └─ Usuwa stare kontenery      │
                                  │
2. PostgreSQL Containers ━━━━━━━━┤
   ├─ postgres-orders  (:5432)   │
   ├─ postgres-auth    (:5433)   │
   └─ postgres-raports (:5434)   │
                                  │
3. Database Migrations ━━━━━━━━━━┤
   ├─ orders: create_tables.sql  │
   ├─ auth: inline SQL           │
   └─ raports: create_tables.sql │
                                  │
4. Nginx Gateway ━━━━━━━━━━━━━━━━┤
   └─ Reverse Proxy (:80)        │
                                  │
5. Go Services ━━━━━━━━━━━━━━━━━━┤
   ├─ Auth Service    (:8081)    │
   ├─ Order Service   (:8080)    │
   └─ Raport Service  (:8083)    │
                                  │
6. Summary ━━━━━━━━━━━━━━━━━━━━━━┘
   └─ Status i instrukcje
```

---

## Dodatkowe Informacje

- **Logi**: Jeśli terminal graficzny nie jest dostępny, logi znajdziesz w katalogu `logs/`
- **Automatyczne czyszczenie**: Skrypt sam wykrywa i usuwa istniejące kontenery przed startem
- **Cross-terminal**: Działa z KDE Konsole, GNOME Terminal, Xterm
- **Kolorowy output**: Łatwe śledzenie postępu uruchamiania

---

## Autor

System Order Management  
Piotr - Praca Inżynierska 2025
