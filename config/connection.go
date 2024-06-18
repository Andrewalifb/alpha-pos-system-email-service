// config/connection.go
package config

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

func ConnectToRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return conn, nil
}
