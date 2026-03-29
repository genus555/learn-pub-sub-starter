package pubsub

import (
	"fmt"
	"log"
	"encoding/gob"
	"bytes"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob [T any] (conn *amqp.Connection, exchange, queueName, key, queueType string, handler func(T) Acktype,) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {return err}

	deliveryChan, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {return err}

	go func() {
		for data := range deliveryChan {
			var message T
			decoded_data := bytes.NewBuffer(data.Body)
			dec := gob.NewDecoder(decoded_data)

			if err := dec.Decode(&message); err != nil {
				fmt.Println("Error unmarshaling message: ", err)
				continue
			}

			ackType := handler(message)
			switch ackType {
			case Ack:
				data.Ack(false)
				log.Println("Message Ack")
			case NackRequeue:
				data.Nack(false, true)
				log.Println("Message NackRequeue")
			case NackDiscard:
				data.Nack(false, false)
				log.Println("Message NackDiscard")
			}
		}
	}()

	return nil
}