package pubsub

import (
	"encoding/gob"
	"bytes"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob [T any] (ch *amqp.Channel, exchange, key string, val T) error {
	var data bytes.Buffer

	enc := gob.NewEncoder(&data)
	if err := enc.Encode(val); err != nil {return err}

	pub := amqp.Publishing {
		ContentType:	"application/gob",
		Body:			data.Bytes(),
	}
	if err := ch.PublishWithContext(context.Background(), exchange, key, false, false, pub); err != nil {return err}

	return nil
}