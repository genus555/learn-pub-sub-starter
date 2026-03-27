package main

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	gamelogic "github.com/genus555/learn-pub-sub-starter/internal/gamelogic"
	routing "github.com/genus555/learn-pub-sub-starter/internal/routing"
	pubsub "github.com/genus555/learn-pub-sub-starter/internal/pubsub"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			war_key := routing.WarRecognitionsPrefix + "." + gs.GetUsername()
			war := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			}
			if err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, war_key, war); err != nil {
				fmt.Printf("error publish war: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			fallthrough
		default:
			return pubsub.NackDiscard
		}
	}
}

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		msg := ""
		outcome, winner, loser := gs.HandleWar(rw)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			fallthrough
		case gamelogic.WarOutcomeYouWon:
			msg = fmt.Sprintf("%s won a war against %s", winner, loser)
			fallthrough
		case gamelogic.WarOutcomeDraw:
			if msg == "" {
				msg = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			}
			ack := publishGameLog(gs, msg, ch)
			return ack
		default:
			fmt.Println("error determining war outcome")
			return pubsub.NackDiscard
		}
	}
}

func publishGameLog (gs *gamelogic.GameState, msg string, ch *amqp.Channel) pubsub.Acktype {
	log := routing.GameLog {
		CurrentTime:	time.Now(),
		Message:		msg,
		Username:		gs.GetUsername(),
	}
	key := routing.GameLogSlug + "." + gs.GetUsername()
	if err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, key, log); err != nil {
		fmt.Println("error publishing war gob")
		return pubsub.NackDiscard
	}
	return pubsub.Ack
}