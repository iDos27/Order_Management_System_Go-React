package publisher

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

type OrderNotification struct {
	OrderID      int       `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	Status       string    `json:"status"`
	TotalAmount  float64   `json:"total_amount"`
	Timestamp    time.Time `json:"timestamp"`
}

func NewPublisher(rabbitMQURL, queueName string) (*Publisher, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Deklaracja kolejki
	_, err = channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	log.Printf("Połączono z RabbitMQ i utworzono kolejkę: %s", queueName)

	return &Publisher{
		conn:    conn,
		channel: channel,
		queue:   queueName,
	}, nil
}

func (p *Publisher) PublishOrderNotification(notification OrderNotification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		"",      // exchange
		p.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		log.Printf("Błąd publikacji wiadomości do RabbitMQ: %v", err)
		return err
	}

	log.Printf("Opublikowano powiadomienie o zamówieniu #%d (status: %s)", notification.OrderID, notification.Status)
	return nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
