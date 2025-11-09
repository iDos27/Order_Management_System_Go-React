-- Baza danych dla raportów
CREATE DATABASE IF NOT EXISTS reports_db;

-- Tabela do przechowywania wygenerowanych raportów
CREATE TABLE IF NOT EXISTS reports (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,           -- 'weekly', 'daily', 'monthly'
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    total_orders INTEGER NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    file_path VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'completed', 'failed'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela do szczegółów źródeł zamówień w raporcie
CREATE TABLE IF NOT EXISTS report_sources (
    id SERIAL PRIMARY KEY,
    report_id INTEGER REFERENCES reports(id) ON DELETE CASCADE,
    source_name VARCHAR(100) NOT NULL,   -- 'website', 'źródło_dwa', 'manual', 'źródło_jeden'
    order_count INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL
);