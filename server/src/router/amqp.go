package router

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)


// CreateProducer creates a new producer and sends a message to the queue
func CreateProducer() (*amqp.Channel, error) {
	queueName := os.Getenv("CLOUDAMQP_QUEUE")
	queueUrl := os.Getenv("CLOUDAMQP_URL")
    conn, err := amqp.Dial(queueUrl)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

	channel, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open a channel: %v", err)
    }
    defer channel.Close()

    _, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to declare a queue: %v", err)
    }

	return channel, nil
}