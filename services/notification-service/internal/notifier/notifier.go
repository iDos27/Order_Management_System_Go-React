package notifier

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/godbus/dbus/v5"
)

type OrderNotification struct {
	OrderID      int     `json:"order_id"`
	CustomerName string  `json:"customer_name"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
}

type Notifier struct {
	conn *dbus.Conn
}

// NewNotifier tworzy nowy notifier z połączeniem D-Bus
func NewNotifier() (*Notifier, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("błąd połączenia z D-Bus: %w", err)
	}

	log.Println("Połączono z D-Bus dla powiadomień systemowych")
	return &Notifier{conn: conn}, nil
}

// HandleOrderUpdate przetwarza wiadomość o aktualizacji zamówienia
func (n *Notifier) HandleOrderUpdate(data []byte) error {
	var notification OrderNotification
	if err := json.Unmarshal(data, &notification); err != nil {
		return fmt.Errorf("błąd parsowania powiadomienia: %w", err)
	}

	log.Printf("Przetwarzanie powiadomienia dla Zamówienia #%d: %s", notification.OrderID, notification.Status)

	// Wysłanie natywnego powiadomienia przez D-Bus
	if err := n.SendNotification(); err != nil {
		return fmt.Errorf("błąd wysyłania powiadomienia: %w", err)
	}

	log.Printf("✓ Powiadomienie wysłane pomyślnie dla Zamówienia #%d", notification.OrderID)
	return nil
}

// SendNotification wysyła natywne powiadomienie systemowe przez D-Bus
func (n *Notifier) SendNotification() error {
	obj := n.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"Order Management",             // app_name
		uint32(0),                      // replaces_id
		"package-new",                  // app_icon
		"Nowe zamówienie",              // summary (tytuł)
		"Kliknij by przejść do panelu", // body
		[]string{},                     // actions
		map[string]dbus.Variant{ // hints
			"urgency":    dbus.MakeVariant(byte(1)),            // normal urgency
			"x-kde-urls": dbus.MakeVariant("http://localhost"), // KDE: URL po kliknięciu
		},
		int32(5000), // timeout (5 sekund)
	)

	if call.Err != nil {
		return fmt.Errorf("błąd wywołania D-Bus: %w", call.Err)
	}

	return nil
}

// translateStatus tłumaczy status na polski
func (n *Notifier) translateStatus(status string) string {
	translations := map[string]string{
		"new":       "Nowe",
		"confirmed": "Potwierdzone",
		"shipped":   "Wysłane",
		"delivered": "Dostarczone",
		"cancelled": "Anulowane",
	}

	if translated, ok := translations[status]; ok {
		return translated
	}
	return status
}

// Close zamyka połączenie D-Bus
func (n *Notifier) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
