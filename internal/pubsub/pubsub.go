package pubsub

import (
	"encoding/json"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON [T any] (ch *amqp.Channel, exchange, key string, val T) error {
	jsonData, err := json.Marshal(val)
	if err != nil {return err}

	pub := amqp.Publishing{
		ContentType:	"application/json",
		Body:			jsonData,
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, pub)
	if err != nil {return err}

	return nil
}