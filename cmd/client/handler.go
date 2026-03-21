package main

import (
	"fmt"

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

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			fallthrough
		case gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			fallthrough
		default:
			return pubsub.NackDiscard
		}
	}
}