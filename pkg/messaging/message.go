package messaging

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IMessage interface {
	Body() []byte
	Ack() error
	Nack() error
}

type IMessaging interface {
	DeclareQueue(name string) error
	Publish(queueName string, body []byte) error
	Consume(queueName, consumer string, autoAck bool) (<-chan IMessage, error)
	Close()
}

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

type rabbitMessage struct {
	delivery amqp.Delivery
}

func (r *rabbitMessage) Body() []byte {
	return r.delivery.Body
}

func (r *rabbitMessage) Ack() error {
	return r.delivery.Ack(false)
}

func (r *rabbitMessage) Nack() error {
	return r.delivery.Nack(false, true)
}

// NewRabbitMQConnection connects to RabbitMQ and returns *RabbitMQ
func NewRabbitMQConnection() (*RabbitMQ, error) {
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPass := os.Getenv("RABBITMQ_PASS")

	if rabbitHost == "" || rabbitPort == "" || rabbitUser == "" || rabbitPass == "" {
		return nil, fmt.Errorf("missing required RabbitMQ environment variables")
	}

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitUser, rabbitPass, rabbitHost, rabbitPort,
	)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	return &RabbitMQ{Conn: conn, Channel: ch}, nil
}

// DeclareQueue declares a durable RabbitMQ queue
func (r *RabbitMQ) DeclareQueue(name string) error {
	_, err := r.Channel.QueueDeclare(
		name,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	return err
}

// Publish sends a message to a queue
func (r *RabbitMQ) Publish(queueName string, body []byte) error {
	return r.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Consume starts consuming messages and wraps them in Message interface
func (r *RabbitMQ) Consume(queueName, consumer string, autoAck bool) (<-chan IMessage, error) {
	deliveries, err := r.Channel.Consume(
		queueName,
		consumer,
		autoAck,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan IMessage)
	go func() {
		defer close(msgChan)
		for d := range deliveries {
			msgChan <- &rabbitMessage{delivery: d}
		}
	}()

	return msgChan, nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}
