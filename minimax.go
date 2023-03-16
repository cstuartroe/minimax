package main

import (
	"github.com/cstuartroe/minimax/gameplay"
	"github.com/cstuartroe/minimax/mancala"
	"github.com/cstuartroe/minimax/minimaxer"
)

func main() {
	game := mancala.MancalaGame(6, 4)

	var player1, player2 gameplay.Player[mancala.MancalaState]

	player1 = gameplay.NewHumanPlayer("Conor", game)
	player2 = minimaxer.NewMinimaxer(game, 8)

	gp := gameplay.NewGameplay(game, player1, player2)

	gp.Play(true)
}
