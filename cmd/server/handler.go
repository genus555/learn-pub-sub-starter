package main

import (
	"fmt"
	
	gamelogic "github.com/genus555/learn-pub-sub-starter/internal/gamelogic"
	routing "github.com/genus555/learn-pub-sub-starter/internal/routing"
	pubsub "github.com/genus555/learn-pub-sub-starter/internal/pubsub"
)

func handleLog () func(routing.GameLog) pubsub.Acktype {
	return func(gl routing.GameLog) pubsub.Acktype {
		defer fmt.Print("> ")
		if err := gamelogic.WriteLog(gl); err != nil {
			fmt.Println(err)
			return pubsub.NackRequeue
		}

		return pubsub.Ack
	}
}