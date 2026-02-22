package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	gamelogic "github.com/genus555/learn-pub-sub-starter/internal/gamelogic"
	pubsub "github.com/genus555/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/genus555/learn-pub-sub-starter/internal/routing"
)

func main() {
	connection_string := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connection_string)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	} else {
		log.Println("Connection was successful")
	}
	defer connection.Close()

	_, q, err := pubsub.DeclareAndBind(connection, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", "durable")
	if err != nil {log.Fatalf("Something went wrong with creating queue: %v", err)}
	log.Println("Queue ", q.Name, " was declared and bound")

	publishCh, err := connection.Channel()
	if err != nil {log.Fatalf("Something went wrong with channel: %v", err)}

	gamelogic.PrintServerHelp()
	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}
		state := routing.PlayingState{}
		switch inputs[0] {
		case "pause":
			log.Println("Game state paused")
			state.IsPaused = true
		case "resume":
			log.Println("Game state resumed")
			state.IsPaused = false
		case "quit":
			log.Println("Closing Connection")
			return
		default:
			log.Printf("\"%s\" is not a valid command", inputs[0])
			continue
		}
		err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, state)
		if err != nil {log.Fatalf("Something went wrong with publishing JSON: %v", err)}
	}
}