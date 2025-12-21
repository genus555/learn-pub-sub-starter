package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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
	fmt.Println("Starting Peril client...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {log.Fatalf("Welcome went wrong: %v", err)}

	queueName := fmt.Sprintf("%v.%v", routing.PauseKey, username)
	ch, _, err := pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, "transient")
	if err != nil {log.Fatalf("DeclareAndBind went wrong: %v", err)}

	state := routing.PlayingState{
		IsPaused:	true,
	}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, state)
	if err != nil { log.Fatalf("Something went wrong with publishing JSON: %v", err) }

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Connection is closed")
}
