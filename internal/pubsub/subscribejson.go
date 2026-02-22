package pubsub

import (
	"fmt"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON [T any] (conn *amqp.Connection, exchange, queueName, key, queueType string, handler func(T),) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {return err}

	deliveryChan, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {return err}

	go func() {
		for data := range deliveryChan {
			var message T

			err = json.Unmarshal(data.Body, &message)
			if err != nil {
				fmt.Println("Error unmarshaling message: ", err)
				continue
			}

			handler(message)
			err = data.Ack(false)
			if err != nil {
				fmt.Println("Error acknowledging message: ", err)
				continue
			}
		}
	}()

	return nil
}