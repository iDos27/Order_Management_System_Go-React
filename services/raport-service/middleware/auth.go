package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

// AuthMiddleware sprawdza czy użytkownik jest zalogowany i ma rolę admin
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Pobierz token z nagłówka Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Brak nagłówka Authorization")
			http.Error(w, "Brak tokena autoryzacji", http.StatusUnauthorized)
			return
		}

		// Sprawdź format "Bearer TOKEN"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Nieprawidłowy format tokena: %s", authHeader)
			http.Error(w, "Nieprawidłowy format tokena", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Sprawdź token przez auth-service
		if !validateTokenWithAuthService(token) {
			log.Printf("Token nieważny lub użytkownik nie ma uprawnień admin")
			http.Error(w, "Nieważny token lub brak uprawnień administratora", http.StatusForbidden)
			return
		}

		// Jeśli wszystko OK, kontynuuj
		next(w, r)
	}
}

// validateTokenWithAuthService sprawdza token przez wywołanie auth-service
func validateTokenWithAuthService(token string) bool {
	authServiceURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8081")

	// Wywołaj endpoint /api/v1/validate-token z auth-service (zwraca dane o użytkowniku)
	req, err := http.NewRequest("GET", authServiceURL+"/api/v1/validate-token", nil)
	if err != nil {
		log.Printf("Błąd tworzenia request do auth-service: %v", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Błąd połączenia z auth-service: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Sprawdź czy auth-service zatwierdził token
	if resp.StatusCode != http.StatusOK {
		log.Printf("Auth-service odrzucił token. Status: %d", resp.StatusCode)
		return false
	}

	// Sprawdź rolę użytkownika
	return validateUserRole(resp)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// validateUserRole sprawdza czy użytkownik ma rolę admin
func validateUserRole(resp *http.Response) bool {
	var result struct {
		Valid  bool   `json:"valid"`
		Role   string `json:"role"`
		Email  string `json:"email"`
		UserID int    `json:"user_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Błąd parsowania odpowiedzi z auth-service: %v", err)
		return false
	}

	if !result.Valid {
		log.Printf("Token nieważny według auth-service")
		return false
	}

	// Sprawdź czy użytkownik ma rolę admin
	if result.Role != "admin" {
		log.Printf("Użytkownik %s (ID: %d) próbował uzyskać dostęp do raportów, ale ma rolę: %s",
			result.Email, result.UserID, result.Role)
		return false
	}

	log.Printf("Autoryzacja udana: admin %s (ID: %d) uzyskał dostęp do raportów",
		result.Email, result.UserID)
	return true
}
