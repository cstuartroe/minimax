package main

import (
	"github.com/cstuartroe/minimax/gameplay"
	"github.com/cstuartroe/minimax/minimaxer"
	"github.com/cstuartroe/minimax/peg_solitaire"
)

func main() {
	game := peg_solitaire.TrianglePegSolitaire()

	mx1 := minimaxer.NewMinimaxer(game, 20)
	// me := gameplay.NewHumanPlayer("Conor", game)

	var player1 gameplay.Player[peg_solitaire.TrianglePegSolitaireState] = mx1

	gp := gameplay.NewGameplay(game, player1, nil)

	gp.Play(true)
}
