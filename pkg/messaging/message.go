package messaging

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IMessaging interface {
	DeclareQueue(name string) error
	Publish(queueName string, body []byte) error
	Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
	Close()
}

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQConnection() (*RabbitMQ, error) {
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPass := os.Getenv("RABBITMQ_PASS")

	if rabbitHost == "" || rabbitPort == "" || rabbitUser == "" || rabbitPass == "" {
		return nil, fmt.Errorf("missing required RabbitMQ environment variables")
	}

	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitUser,
		rabbitPass,
		rabbitHost,
		rabbitPort,
	)

	conn, err := amqp.Dial(connectionString)
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

func (r *RabbitMQ) DeclareQueue(name string) error {
	_, err := r.Channel.QueueDeclare(
		name,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	return err
}

func (r *RabbitMQ) Publish(queueName string, body []byte) error {
	return r.Channel.Publish(
		"",        // default exchange
		queueName, // routing key = queue name
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(
		queueName,
		consumer, // consumer name (can be empty string)
		autoAck,  // auto acknowledge messages?
		false,    // exclusive
		false,    // no-local (deprecated)
		false,    // no-wait
		nil,      // args
	)
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}
