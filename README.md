# Order_Management_System_Go-React


Użyte zależności:
go get github.com/gin-gonic/gin         - framework web do tworzenia REST API
go get github.com/lib/pq                - driver PostgreSQL
go get github.com/joho/godotenv         - ładowanie zmiennych środowiskowych z pliku .env
go get github.com/gorilla/websocket     - framework web do tworzenia websocketów
go get github.com/gin-contrib/cors




WYJAŚNIENIE     order.go
OrderStatus - enum dla statusów zamówienia
Order - główna struktura zamówienia
OrderItem - pozycje w zamówieniu
tagi JSON i DB - dla serializacji JSON i mapowania bazy danych

WYJAŚNIENIE     connection.go
Łączy z PostgreSQL używając zmiennych środowiskowych
Sprawdza połączenie funkcją Ping()



docker run --name postgres-orders -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password123 -e POSTGRES_DB=orders_management -p 5432:5432 -d postgres:17


docker exec -i postgres-orders psql -U postgres -d orders_management < backend/migration/create_tables.sql




INSTALACJA REACT
npm create vite@latest admin-panel -- --template react
npm install
npm install axios