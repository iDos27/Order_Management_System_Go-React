package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"

	"github.com/iDos27/order-management/raport-service/internal/database"
	"github.com/iDos27/order-management/raport-service/internal/handlers"
	"github.com/iDos27/order-management/raport-service/internal/models"
	"github.com/iDos27/order-management/raport-service/internal/services"
	"github.com/iDos27/order-management/raport-service/middleware"
)

func main() {
	// Połączenie z bazami danych
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Błąd połączenia z bazą danych:", err)
	}
	defer db.Close()

	// Inicjalizacja serwisów
	reportService := services.NewReportService(db)

	// Inicjalizacja handlerów
	reportHandler := handlers.NewReportsHandler(reportService)

	// Konfiguracja routingu
	router := mux.NewRouter()

	// API endpoints z middleware autoryzacji
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/reports/generate", middleware.AuthMiddleware(reportHandler.GenerateReport)).Methods("POST")
	api.HandleFunc("/reports/stats", middleware.AuthMiddleware(reportHandler.GetQuickStats)).Methods("GET")

	// Health check endpoint (bez autoryzacji)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Raport Service działa"))
	}).Methods("GET")

	// Uruchomienie cron job (co 5 minut dla testów)
	startCronJobs(reportService)

	// Konfiguracja portu
	port := getEnv("PORT", "8083")

	log.Printf("Raport Service uruchomiony na porcie %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func startCronJobs(reportService *services.ReportService) {
	c := cron.New()

	// Co 5 minut generuj tygodniowy raport (dla testów)
	// Później zmienisz na: "0 0 * * 1" (co poniedziałek o północy)
	c.AddFunc("*/5 * * * *", func() {
		log.Println("Uruchamiam automatyczne generowanie raportu...")

		// Ostatnie 7 dni
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -7)

		stats, err := reportService.GetOrderStats(startDate, endDate)
		if err != nil {
			log.Printf("Błąd pobierania danych dla automatycznego raportu: %v", err)
			return
		}

		filePath, err := reportService.GenerateExcelReport(stats, "weekly", startDate, endDate)
		if err != nil {
			log.Printf("Błąd generowania automatycznego raportu: %v", err)
			return
		}

		// Zapisz raport w bazie
		report := &models.Report{
			Type:        "weekly",
			PeriodStart: startDate,
			PeriodEnd:   endDate,
			TotalOrders: stats.TotalOrders,
			TotalAmount: stats.TotalAmount,
			FilePath:    &filePath,
			Status:      "completed",
			Sources:     make([]models.ReportSource, 0),
		}

		for _, source := range stats.Sources {
			report.Sources = append(report.Sources, models.ReportSource{
				SourceName: source.SourceName,
				OrderCount: source.Count,
				Amount:     source.Amount,
			})
		}

		reportID, err := reportService.SaveReport(report)
		if err != nil {
			log.Printf("Błąd zapisywania automatycznego raportu: %v", err)
			return
		}

		log.Printf("Automatyczny raport wygenerowany pomyślnie! ID: %d, Plik: %s", reportID, filePath)
	})

	c.Start()
	log.Println("Cron jobs uruchomione - automatyczne raporty co 5 minut")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
