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

-- Przykładowe wygenerowane raporty z poprzednich miesięcy
INSERT INTO reports (type, period_start, period_end, total_orders, total_amount, file_path, status, created_at) VALUES
('monthly', '2025-08-01', '2025-08-31', 4, 946.50, 'reports/2025-08_monthly.xlsx', 'completed', '2025-09-01 08:00:00'),
('monthly', '2025-09-01', '2025-09-30', 5, 1281.54, 'reports/2025-09_monthly.xlsx', 'completed', '2025-10-01 08:00:00'),
('monthly', '2025-10-01', '2025-10-31', 5, 1337.40, 'reports/2025-10_monthly.xlsx', 'completed', '2025-11-01 08:00:00'),
('daily', '2025-11-18', '2025-11-18', 1, 410.00, 'reports/2025-11-18_daily.xlsx', 'completed', '2025-11-18 23:00:00');

-- Szczegóły źródeł dla raportów
INSERT INTO report_sources (report_id, source_name, order_count, amount) VALUES
-- Sierpień 2025
(1, 'website', 2, 345.75),
(1, 'źródło_jeden', 1, 180.00),
(1, 'manual', 1, 420.75),
-- Wrzesień 2025
(2, 'website', 3, 695.74),
(2, 'źródło_dwa', 1, 310.00),
(2, 'źródło_jeden', 1, 199.99),
(2, 'manual', 1, 275.80),
-- Październik 2025
(3, 'website', 2, 509.65),
(3, 'źródło_dwa', 1, 445.00),
(3, 'źródło_jeden', 1, 167.25),
(3, 'manual', 1, 215.50),
-- Listopad 18
(4, 'website', 1, 410.00);