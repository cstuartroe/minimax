package main

import (
	"fmt"

	"github.com/cstuartroe/minimax/base"
	"github.com/cstuartroe/minimax/minimaxer"
	"github.com/cstuartroe/minimax/nim"
)

type Player[State base.GameState] interface {
	Name() string
	ChooseMove(base.Prospect[State]) base.Move[State]
}

type HumanPlayer[State base.GameState] struct {
	name string
	game base.Game[State]
}

func NewHumanPlayer[State base.GameState](name string, game base.Game[State]) *HumanPlayer[State] {
	return &HumanPlayer[State]{name, game}
}

func (p HumanPlayer[State]) Name() string {
	return p.name
}

func (p HumanPlayer[State]) ChooseMove(prospect base.Prospect[State]) base.Move[State] {
	sd := p.game.Describe(prospect)

	fmt.Println("Possible moves:")
	for i, move := range sd.Moves {
		fmt.Printf("%d: %s\n", i, move.Summary)
	}
	fmt.Print("Choose: ")
	var choice int
	fmt.Scan(&choice)
	return sd.Moves[choice]
}

type Gameplay[State base.GameState] struct {
	game            base.Game[State]
	player1         Player[State]
	player2         Player[State]
	currentProspect base.Prospect[State]
}

func NewGameplay[State base.GameState](game base.Game[State]) Gameplay[State] {
	return Gameplay[State]{
		game: game,
		currentProspect: base.Prospect[State]{
			State:      game.InitialState,
			FirstAgent: true,
		},
	}
}

func (gp *Gameplay[State]) makeMove(newState State) {
	gp.currentProspect = base.Prospect[State]{
		State:      newState,
		FirstAgent: !gp.currentProspect.FirstAgent,
	}
}

func (gp Gameplay[State]) done() bool {
	return len(gp.game.Describe(gp.currentProspect).Moves) == 0
}

func (gp *Gameplay[State]) Play() {
	for !gp.done() {
		player := gp.player2
		if gp.currentProspect.FirstAgent {
			player = gp.player1
		}

		fmt.Println("Current state:")
		fmt.Println(gp.currentProspect.State.String())
		move := player.ChooseMove(gp.currentProspect)
		fmt.Printf("%s chose %s\n\n", player.Name(), move.Summary)
		gp.makeMove(move.State)
	}

	fmt.Println("Final state:")
	fmt.Println(gp.currentProspect.State.String())

	score := gp.game.Describe(gp.currentProspect).Score
	fmt.Printf("Game score: %d\n", score)
	if score > 0 {
		fmt.Printf("%s wins!\n", gp.player1.Name())
	} else if score < 0 {
		fmt.Printf("%s wins!\n", gp.player2.Name())
	} else {
		fmt.Println("It's a draw.")
	}
}

func main() {
	game := nim.GenerateNim([]int{5, 4, 3, 2}, 0, true)

	gp := NewGameplay(game)

	// me := NewHumanPlayer("Conor", game)
	mx1 := minimaxer.NewMinimaxer(game)
	mx2 := minimaxer.NewMinimaxer(game)

	gp.player1 = mx1
	gp.player2 = mx2

	gp.Play()
}
