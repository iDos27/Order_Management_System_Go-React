package notifier

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"

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

// NewNotifier tworzy nowy notifier z połączeniem D-Bus i nasłuchuje na kliknięcia
func NewNotifier() (*Notifier, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("błąd połączenia z D-Bus: %w", err)
	}

	n := &Notifier{conn: conn}

	// Rejestrujemy regułę, aby otrzymywać sygnał ActionInvoked z serwera powiadomień
	if err = conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"interface='org.freedesktop.Notifications',member='ActionInvoked'").Store(); err != nil {
		return nil, fmt.Errorf("błąd dodawania reguły match: %w", err)
	}

	// Tworzymy kanał dla sygnałów
	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	// Uruchamiamy goroutine do obsługi kliknięć w tle
	go n.handleSignals(c)

	log.Println("Połączono z D-Bus dla powiadomień systemowych i nasłuchiwania akcji")
	return n, nil
}

// handleSignals nasłuchuje na kliknięcia w powiadomienia
func (n *Notifier) handleSignals(c chan *dbus.Signal) {
	for v := range c {
		// Sprawdzamy, czy to sygnał "ActionInvoked" (kliknięcie w akcję)
		if v.Name == "org.freedesktop.Notifications.ActionInvoked" {
			// v.Body[0] to ID powiadomienia (uint32)
			// v.Body[1] to klucz akcji (string) - szukamy "default"
			if len(v.Body) >= 2 {
				actionKey, ok := v.Body[1].(string)
				if ok && actionKey == "default" {
					log.Println("Kliknięto w powiadomienie! Otwieranie localhost...")
					openBrowser("http://localhost")
				}
			}
		}
	}
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("nieobsługiwana platforma")
	}
	if err != nil {
		log.Printf("Błąd otwierania przeglądarki: %v", err)
	}
}

// HandleOrderUpdate przetwarza wiadomość o aktualizacji zamówienia
func (n *Notifier) HandleOrderUpdate(data []byte) error {
	var notification OrderNotification
	if err := json.Unmarshal(data, &notification); err != nil {
		return fmt.Errorf("błąd parsowania powiadomienia: %w", err)
	}

	log.Printf("Przetwarzanie powiadomienia dla Zamówienia #%d: %s", notification.OrderID, notification.Status)

	// Wyświetlaj powiadomienie tylko dla nowych zamówień
	if notification.Status != "new" {
		log.Printf("Pomijam powiadomienie - status to '%s' (tylko 'new' generuje powiadomienia)", notification.Status)
		return nil
	}

	title := "Nowe zamówienie"
	body := fmt.Sprintf("Zamówienie #%d\nKlient: %s\nKwota: %.2f PLN", notification.OrderID, notification.CustomerName, notification.TotalAmount)

	if err := n.SendNotification(title, body); err != nil {
		return fmt.Errorf("błąd wysyłania powiadomienia: %w", err)
	}

	log.Printf("✓ Powiadomienie wysłane pomyślnie dla Zamówienia #%d", notification.OrderID)
	return nil
}

// SendNotification wysyła natywne powiadomienie systemowe przez D-Bus
func (n *Notifier) SendNotification(title, body string) error {
	obj := n.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"System Zarządzania Zamówieniami", // app_name
		uint32(0),                         // replaces_id
		"",                                // app_icon
		title,                             // summary
		body,                              // body
		[]string{"default", "Otwórz"},     // actions
		map[string]dbus.Variant{ // hints
			"urgency": dbus.MakeVariant(byte(1)), // normal urgency
		},
		int32(5000), // timeout
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
