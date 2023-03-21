package main

import (
	"math/rand"

	"github.com/cstuartroe/minimax/gameplay"
	"github.com/cstuartroe/minimax/mancala"
	"github.com/cstuartroe/minimax/minimaxer"
)

func main() {
	game := mancala.MancalaGame(6, 4)

	mx1 := minimaxer.NewMinimaxer(game, 4)
	mx2 := minimaxer.NewMinimaxer(game, 8)
	me := gameplay.NewHumanPlayer("Conor", game)

	var player1, player2 gameplay.Player[mancala.MancalaState]

	if rand.Float32() > .5 {
		player1 = mx1
		player2 = minimaxer.NewAssistedHumanPlayer(me, mx2)
	} else {
		player1 = minimaxer.NewAssistedHumanPlayer(me, mx1)
		player2 = mx2
	}

	gp := gameplay.NewGameplay(game, player1, player2)

	gp.Play(true)
}
