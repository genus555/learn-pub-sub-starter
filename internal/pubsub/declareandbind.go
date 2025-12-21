package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key, queueType string) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {return ch, amqp.Queue{}, err}

	durable := false
	autoDelete := false
	exclusive := false
	switch queueType {
	case "durable":
		durable = true
		break
	case "transient":
		autoDelete = true
		exclusive = true
		break
	default:
		return ch, amqp.Queue{}, fmt.Errorf("Wrong QueueType")
	}
	q, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {return ch, q, err}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {return ch, q, err}

	return ch, q, nil
}