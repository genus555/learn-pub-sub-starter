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

	publishCh, err := connection.Channel()

	username, err := gamelogic.ClientWelcome()
	if err != nil {log.Fatalf("Welcome went wrong: %v", err)}

	queueName := fmt.Sprintf("pause.%v", username)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	//Client REPL
	gameState := gamelogic.NewGameState(username)

	move_key := routing.ArmyMovesPrefix + "." + username
	war_key := routing.WarRecognitionsPrefix + ".*"

	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, "transient", handlerPause(gameState))
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, move_key, routing.ArmyMovesPrefix+".*", "transient", handlerMove(gameState, publishCh))
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, war_key, "durable", handlerWar(gameState, publishCh))

	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}
		switch inputs[0] {
		case "spawn":
			err := gameState.CommandSpawn(inputs)
			if err != nil {log.Println(err)}
		case "move":
			move, err := gameState.CommandMove(inputs)
			if err != nil {
				log.Println(err)
			} else {
				err = pubsub.PublishJSON(publishCh, routing.ExchangePerilTopic, move_key, move)
				log.Printf("%s moved successfully", username)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			log.Println("Spamming is not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			log.Printf("\"%s\" is not a valid command", inputs[0])
			continue
		}
	}
}
