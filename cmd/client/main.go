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

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

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

	queueName := fmt.Sprintf("pause.%v", username)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	//Client REPL
	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, "transient", handlerPause(gameState))

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
			_, err := gameState.CommandMove(inputs)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("Unit has successfully moved")
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
