package queues

import (
	"courses-api/domain/courses"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	QueueName string
}

type Rabbit struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewRabbit(config RabbitConfig) Rabbit {
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port))
	if err != nil {
		log.Fatalf("error getting Rabbit connection: %v", err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("error creating Rabbit channel: %v", err)
	}
	queue, err := channel.QueueDeclare(config.QueueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("error declaring Rabbit queue: %v", err)
	}

	return Rabbit{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}

}

func (queue Rabbit) Publish(courseNew courses.CourseNew) error {
	bytes, err := json.Marshal(courseNew)
	if err != nil {
		return fmt.Errorf("error marshaling Rabbit courseNew: %w", err)
	}
	if err := queue.channel.Publish(
		"",
		queue.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		}); err != nil {
		return fmt.Errorf("error publishing to Rabbit: %w", err)
	}
	return nil
}

// Close cleans up the RabbitMQ resources
func (queue Rabbit) Close() {
	if err := queue.channel.Close(); err != nil {
		log.Printf("error closing Rabbit channel: %v", err)
	}
	if err := queue.connection.Close(); err != nil {
		log.Printf("error closing Rabbit connection: %v", err)
	}
}
