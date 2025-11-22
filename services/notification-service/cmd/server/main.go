package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iDos27/order-management/notification-service/internal/consumer"
	"github.com/iDos27/order-management/notification-service/internal/notifier"
)

func main() {
	// Konfiguracja z zmiennych środowiskowych
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	queueName := getEnv("QUEUE_NAME", "order_notifications")

	log.Println("Uruchamianie Notification Service...")
	log.Printf("RabbitMQ URL: %s", rabbitmqURL)
	log.Printf("Nazwa kolejki: %s", queueName)

	// Inicjalizacja notifiera
	notif, err := notifier.NewNotifier()
	if err != nil {
		log.Fatalf("Błąd inicjalizacji notifiera: %v", err)
	}
	defer notif.Close()

	// Inicjalizacja konsumenta RabbitMQ
	cons, err := consumer.NewConsumer(rabbitmqURL)
	if err != nil {
		log.Fatalf("Błąd tworzenia konsumenta: %v", err)
	}
	defer cons.Close()

	// Handler dla wiadomości z RabbitMQ
	handler := func(data []byte) {
		if err := notif.HandleOrderUpdate(data); err != nil {
			log.Printf("Błąd obsługi powiadomienia: %v", err)
		}
	}

	// Uruchomienie konsumowania w goroutine
	go func() {
		if err := cons.StartConsuming(queueName, handler); err != nil {
			log.Fatalf("Błąd uruchomienia konsumowania: %v", err)
		}
	}()

	log.Println("Notification Service działa. Naciśnij Ctrl+C aby zakończyć.")

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Zamykanie Notification Service...")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
