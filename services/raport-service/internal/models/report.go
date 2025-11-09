package models

import (
	"time"
)

// Główna struktura raportu
type Report struct {
	ID          int            `json:"id" db:"id"`
	Type        string         `json:"type" db:"type"`
	PeriodStart time.Time      `json:"period_start" db:"period_start"`
	PeriodEnd   time.Time      `json:"period_end" db:"period_end"`
	TotalOrders int            `json:"total_orders" db:"total_orders"`
	TotalAmount float64        `json:"total_amount" db:"total_amount"`
	FilePath    *string        `json:"file_path,omitempty" db:"file_path"`
	Status      string         `json:"status" db:"status"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	Sources     []ReportSource `json:"sources,omitempty"`
}

// Szczegóły źródła raportu
type ReportSource struct {
	ID         int     `json:"id" db:"id"`
	ReportID   int     `json:"report_id" db:"report_id"`
	SourceName string  `json:"source_name" db:"source_name"`
	OrderCount int     `json:"order_count" db:"order_count"`
	Amount     float64 `json:"amount" db:"amount"`
}

// Struktury do generowania raportów
type OrderStats struct {
	TotalOrders int          `json:"total_orders"`
	TotalAmount float64      `json:"total_amount"`
	Sources     []SourceStat `json:"sources"`
}

// Statystyki dla pojedynczego źródła
type SourceStat struct {
	SourceName string  `json:"source_name"`
	Count      int     `json:"count"`
	Amount     float64 `json:"amount"`
}

// Żądanie generowania raportu
type GenerateReportRequest struct {
	Type        string    `json:"type"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}
