package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/iDos27/order-management/raport-service/internal/database"
	"github.com/iDos27/order-management/raport-service/internal/models"

	"github.com/xuri/excelize/v2"
)

type ReportService struct {
	db *database.DB
}

func NewReportService(db *database.DB) *ReportService {
	return &ReportService{db: db}
}

// Pobieramy dane o zamówieniach z bazy orders
func (rs *ReportService) GetOrderStats(periodStart, periodEnd time.Time) (*models.OrderStats, error) {
	query := `
 		SELECT 
			source,
			COUNT(*) as count,
			COALESCE(SUM(total_amount), 0) as amount
		FROM orders 
		WHERE created_at BETWEEN $1 AND $2 
		GROUP BY source
		ORDER BY source
	`

	rows, err := rs.db.OrdersDB.Query(query, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("błąd pobierania danych zamówień: %v", err)
	}
	defer rows.Close()

	stats := &models.OrderStats{
		Sources: make([]models.SourceStat, 0),
	}

	// Skanuj wyniki dla każdego źródła
	for rows.Next() {
		var source models.SourceStat
		err := rows.Scan(&source.SourceName, &source.Count, &source.Amount)
		if err != nil {
			return nil, fmt.Errorf("błąd skanowania danych: %v", err)
		}

		stats.Sources = append(stats.Sources, source)
		stats.TotalOrders += source.Count
		stats.TotalAmount += source.Amount
	}

	return stats, nil
}

// Generuje plik Excel z raportem
func (rs *ReportService) GenerateExcelReport(stats *models.OrderStats, reportType string, periodStart, periodEnd time.Time) (string, error) {
	f := excelize.NewFile()

	// Utwórz nowy arkusz
	sheetName := "Raport Zamówień"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("błąd tworzenia arkusza: %v", err)
	}

	// Ustaw aktywny arkusz
	f.SetActiveSheet(index)

	// Nagłówki raportu
	f.SetCellValue(sheetName, "A1", "RAPORT ZAMÓWIEŃ")
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Okres: %s - %s",
		periodStart.Format("2006-01-02"),
		periodEnd.Format("2006-01-02")))
	f.SetCellValue(sheetName, "A3", fmt.Sprintf("Typ: %s", reportType))

	// Podsumowanie
	f.SetCellValue(sheetName, "A5", "PODSUMOWANIE")
	f.SetCellValue(sheetName, "A6", "Łączna liczba zamówień:")
	f.SetCellValue(sheetName, "B6", stats.TotalOrders)
	f.SetCellValue(sheetName, "A7", "Łączna kwota:")
	f.SetCellValue(sheetName, "B7", stats.TotalAmount)

	// Szczegóły źródeł
	f.SetCellValue(sheetName, "A9", "SZCZEGÓŁY WEDŁUG ŹRÓDEŁ")
	f.SetCellValue(sheetName, "A10", "Źródło")
	f.SetCellValue(sheetName, "B10", "Liczba zamówień")
	f.SetCellValue(sheetName, "C10", "Kwota")

	for i, source := range stats.Sources {
		row := 11 + i
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), source.SourceName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), source.Count)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), source.Amount)
	}

	// Utwórz katalog dla plików
	reportsDir := "./reports"
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return "", fmt.Errorf("błąd tworzenia katalogu: %v", err)
	}

	// Nazwa pliku
	filename := fmt.Sprintf("%s_raport_%s.xlsx",
		reportType,
		time.Now().Format("2006_01_02_15_04_05"))
	filePath := filepath.Join(reportsDir, filename)

	// Zapisz plik
	if err := f.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("błąd zapisywania pliku: %v", err)
	}

	return filePath, nil
}

// Zapisuje raport w bazie danych
func (rs *ReportService) SaveReport(report *models.Report) (int, error) {
	query := `
		INSERT INTO reports (type, period_start, period_end, total_orders, total_amount, file_path, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var reportID int
	err := rs.db.RaportsDB.QueryRow(query,
		report.Type,
		report.PeriodStart,
		report.PeriodEnd,
		report.TotalOrders,
		report.TotalAmount,
		report.FilePath,
		report.Status,
	).Scan(&reportID)

	if err != nil {
		return 0, fmt.Errorf("błąd zapisywania raportu: %v", err)
	}

	// Zapisz źródła
	for _, source := range report.Sources {
		_, err := rs.db.RaportsDB.Exec(`
			INSERT INTO report_sources (report_id, source_name, order_count, amount)
			VALUES ($1, $2, $3, $4)
		`, reportID, source.SourceName, source.OrderCount, source.Amount)

		if err != nil {
			return 0, fmt.Errorf("błąd zapisywania źródeł raportu: %v", err)
		}
	}

	return reportID, nil
}
