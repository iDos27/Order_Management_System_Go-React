CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    source VARCHAR(50) NOT NULL DEFAULT 'website',
    status VARCHAR(50) NOT NULL DEFAULT 'new',
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    product_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL
);

-- Wstawienie przykładowych zamówień z różnych miesięcy (2025)
-- Równomierny rozkład po statusach: new(4), confirmed(4), shipped(4), delivered(4), cancelled(3)
-- Sierpień 2025
INSERT INTO orders (customer_name, customer_email, source, status, total_amount, created_at, updated_at) VALUES
('Jan Kowalski', 'jan.kowalski@example.com', 'website', 'new', 250.50, '2025-08-05 10:30:00', '2025-08-05 10:30:00'),
('Anna Nowak', 'anna.nowak@example.com', 'źródło_jeden', 'confirmed', 180.00, '2025-08-12 14:15:00', '2025-08-12 14:15:00'),
('Piotr Wiśniewski', 'piotr.wisniewski@example.com', 'manual', 'shipped', 420.75, '2025-08-18 09:45:00', '2025-08-18 09:45:00'),
('Maria Lewandowska', 'maria.lewandowska@example.com', 'website', 'delivered', 95.25, '2025-08-25 16:20:00', '2025-08-25 16:20:00'),

-- Wrzesień 2025
('Tomasz Kamiński', 'tomasz.kaminski@example.com', 'źródło_dwa', 'cancelled', 310.00, '2025-09-03 11:00:00', '2025-09-03 11:00:00'),
('Katarzyna Zielińska', 'katarzyna.zielinska@example.com', 'website', 'new', 155.50, '2025-09-10 13:30:00', '2025-09-10 13:30:00'),
('Michał Woźniak', 'michal.wozniak@example.com', 'manual', 'confirmed', 275.80, '2025-09-15 10:15:00', '2025-09-15 10:15:00'),
('Agnieszka Dąbrowska', 'agnieszka.dabrowska@example.com', 'źródło_jeden', 'shipped', 199.99, '2025-09-22 15:45:00', '2025-09-22 15:45:00'),
('Paweł Wojciechowski', 'pawel.wojciechowski@example.com', 'website', 'delivered', 340.25, '2025-09-28 12:00:00', '2025-09-28 12:00:00'),

-- Październik 2025
('Ewa Kowalczyk', 'ewa.kowalczyk@example.com', 'website', 'cancelled', 128.75, '2025-10-04 09:30:00', '2025-10-04 09:30:00'),
('Krzysztof Mazur', 'krzysztof.mazur@example.com', 'źródło_dwa', 'new', 445.00, '2025-10-09 14:00:00', '2025-10-09 14:00:00'),
('Magdalena Krawczyk', 'magdalena.krawczyk@example.com', 'manual', 'confirmed', 215.50, '2025-10-14 11:20:00', '2025-10-14 11:20:00'),
('Adam Piotrowski', 'adam.piotrowski@example.com', 'website', 'shipped', 380.90, '2025-10-20 16:10:00', '2025-10-20 16:10:00'),
('Joanna Grabowska', 'joanna.grabowska@example.com', 'źródło_jeden', 'delivered', 167.25, '2025-10-26 10:50:00', '2025-10-26 10:50:00'),

-- Listopad 2025 (do 20-go)
('Marek Pawłowski', 'marek.pawlowski@example.com', 'website', 'cancelled', 295.00, '2025-11-02 13:15:00', '2025-11-02 13:15:00'),
('Barbara Michalska', 'barbara.michalska@example.com', 'źródło_dwa', 'new', 220.50, '2025-11-08 09:00:00', '2025-11-08 09:00:00'),
('Grzegorz Nowakowski', 'grzegorz.nowakowski@example.com', 'manual', 'confirmed', 175.75, '2025-11-12 15:30:00', '2025-11-12 15:30:00'),
('Monika Adamczyk', 'monika.adamczyk@example.com', 'website', 'shipped', 410.00, '2025-11-18 11:45:00', '2025-11-18 11:45:00'),
('Robert Dudek', 'robert.dudek@example.com', 'źródło_jeden', 'delivered', 185.25, '2025-11-20 14:20:00', '2025-11-20 14:20:00');