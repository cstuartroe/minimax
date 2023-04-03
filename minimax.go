package main

import (
	"github.com/cstuartroe/minimax/connect_four"
	"github.com/cstuartroe/minimax/gameplay"
	"github.com/cstuartroe/minimax/games"
	"github.com/cstuartroe/minimax/minimaxer"
)

func p[State games.GameState](game games.Game[State], player gameplay.Player[State]) gameplay.Player[State] {
	return player
}

func main() {
	game := connect_four.ConnectFour()

	mx1 := minimaxer.NewMinimaxer(game, 10)
	mx2 := minimaxer.NewMinimaxer(game, 10)
	// me := gameplay.NewHumanPlayer("Conor", game)

	var player1 gameplay.Player[connect_four.ConnectFourState] = mx1
	var player2 gameplay.Player[connect_four.ConnectFourState] = mx2

	gp := gameplay.NewGameplay(game, player1, player2)

	gp.Play(true)
}
