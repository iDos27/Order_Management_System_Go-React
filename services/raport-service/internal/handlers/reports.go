package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/iDos27/order-management/raport-service/internal/models"
	"github.com/iDos27/order-management/raport-service/internal/services"
)

type ReportsHandler struct {
	service *services.ReportService
}

func NewReportsHandler(service *services.ReportService) *ReportsHandler {
	return &ReportsHandler{service: service}
}

// POST /reports/generate - Generowanie nowego raportu
func (h *ReportsHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req models.GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Błąd parsowania JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Pobranie danych z bazy orders
	stats, err := h.service.GetOrderStats(req.PeriodStart, req.PeriodEnd)
	if err != nil {
		log.Printf("Błąd pobierania danych: %v", err)
		http.Error(w, "Błąd pobierania danych zamówień: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Tworzenie raportu Excel
	filePath, err := h.service.GenerateExcelReport(stats, req.Type, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		log.Printf("Błąd generowania raportu: %v", err)
		http.Error(w, "Błąd generowania raportu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Struktura raportu
	report := &models.Report{
		Type:        req.Type,
		PeriodStart: req.PeriodStart,
		PeriodEnd:   req.PeriodEnd,
		TotalOrders: stats.TotalOrders,
		TotalAmount: stats.TotalAmount,
		FilePath:    &filePath,
		Status:      "completed",
		Sources:     make([]models.ReportSource, 0),
	}

	// Przekształcanie źródeł
	for _, source := range stats.Sources {
		report.Sources = append(report.Sources, models.ReportSource{
			SourceName: source.SourceName,
			OrderCount: source.Count,
			Amount:     source.Amount,
		})
	}

	// Zapis w bazie raportów
	reportID, err := h.service.SaveReport(report)
	if err != nil {
		log.Printf("Błąd zapisywania raportu: %v", err)
		http.Error(w, "Błąd zapisywania raportu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	report.ID = reportID
	report.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GET /reports/stats - Podgląd statystyk raportu bez generowania pliku
func (h *ReportsHandler) GetQuickStats(w http.ResponseWriter, r *http.Request) {
	// Domyślnie ostatnie 7 dni
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	// Opcjonalnie można przekazać daty przez query params
	if start := r.URL.Query().Get("start"); start != "" {
		if parsed, err := time.Parse("2006-01-02", start); err == nil {
			startDate = parsed
		}
	}
	if end := r.URL.Query().Get("end"); end != "" {
		if parsed, err := time.Parse("2006-01-02", end); err == nil {
			endDate = parsed
		}
	}

	stats, err := h.service.GetOrderStats(startDate, endDate)
	if err != nil {
		log.Printf("Błąd pobierania statystyk: %v", err)
		http.Error(w, "Błąd pobierania danych", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
