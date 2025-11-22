package consumer

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// Nowe połaczenie z RabbitMQ
func NewConsumer(amqpURL string) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("błąd łączenia z RabbitMQ: %v", err)
	}

	chanel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("błąd tworzenia kanału RabbitMQ: %v", err)
	}

	log.Println("Utworzono połączenie z RabbitMQ")

	return &Consumer{
		conn:    conn,
		channel: chanel,
		url:     amqpURL,
	}, nil
}

// Rozpocznij nasłuchiwanie na wiadomości z kolejki
func (c *Consumer) StartConsuming(queueName string, handler func([]byte)) error {
	// Deklaracja kolejki
	queue, err := c.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("BŁĄD deklaracji kolejki: %w", err)
	}

	log.Printf("Nasłuchiwanie na kolejce: %s", queue.Name)

	// Ustawienie prefetch - ile wiadomości może być przetwarzanych jednocześnie
	err = c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("BŁĄD ustawiania QoS: %w", err)
	}

	messages, err := c.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("BŁĄD rozpoczęcia konsumpcji: %w", err)
	}

	log.Println("Oczekiwanie na wiadomosci...")

	// Przetwarzanie wiadomości
	for msg := range messages {
		log.Printf("Otrzymano wiadomość: %s", string(msg.Body))

		// Obsłuż wiadomość
		handler(msg.Body)

		// Potwierdź przetworzenie wiadomości
		if err := msg.Ack(false); err != nil {
			log.Printf("BŁĄD potwierdzania wiadomości: %v", err)
		}
	}
	return nil
}

// Reconnect próbuje ponownie połączyć się z RabbitMQ
func (c *Consumer) Reconnect() error {
	log.Println("Próba ponownego połączenia z RabbitMQ...")

	// Zamknij istniejące połączenie i kanał
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}

	// Spróbuj ponownie połączyć się
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		conn, err := amqp.Dial(c.url)
		if err != nil {
			waitTime := time.Duration(i+1) * 2 * time.Second
			log.Printf("Reconnect nieudany (próba %d/%d): %v. Ponawiam w %v...", i+1, maxRetries, err, waitTime)
			time.Sleep(waitTime)
			continue
		}

		channel, err := conn.Channel()
		if err != nil {
			conn.Close()
			log.Printf("Nieudane otwarcie kanału (próba %d/%d): %v", i+1, maxRetries, err)
			continue
		}

		c.conn = conn
		c.channel = channel
		log.Println("Skuteczne ponowne połączenie z RabbitMQ")
		return nil
	}

	return fmt.Errorf("Nie udało się połączyć w %d attempts", maxRetries)
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	log.Println("Zamknięto połączenie z RabbitMQ")
}
