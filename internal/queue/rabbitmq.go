package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"router-cloud-platform/internal/config"
)

const HeartbeatQueue = "heartbeat_queue"

var Conn *amqp.Connection
var Channel *amqp.Channel

func Connect() {
	var err error
	for i := 1; i <= 5; i++ {
		Conn, err = amqp.Dial(config.App.RabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connection attempt %d failed: %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	_, err = Channel.QueueDeclare(HeartbeatQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	log.Println("RabbitMQ connected successfully")
}

func Publish(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return Channel.PublishWithContext(
		context.Background(), "", HeartbeatQueue, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func Close() {
	if Channel != nil {
		Channel.Close()
	}
	if Conn != nil {
		Conn.Close()
	}
}