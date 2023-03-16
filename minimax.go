package main

import (
	"fmt"

	"github.com/cstuartroe/minimax/base"
	"github.com/cstuartroe/minimax/mancala"
	"github.com/cstuartroe/minimax/minimaxer"
)

type Player[State base.GameState] interface {
	Name() string
	ChooseMove(base.Prospect[State]) base.Move[State]
}

type HumanPlayer[State base.GameState] struct {
	name string
	game base.Game[State]
}

func NewHumanPlayer[State base.GameState](name string, game base.Game[State]) HumanPlayer[State] {
	return HumanPlayer[State]{name, game}
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
	_, err := fmt.Scan(&choice)

	if err != nil || choice >= len(sd.Moves) {
		fmt.Println("Invalid entry.")
		return p.ChooseMove(prospect)
	}

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
			State:      game.InitialState(),
			FirstAgent: true,
		},
	}
}

func (gp *Gameplay[State]) makeMove(move base.Move[State]) {
	nextAgent := !gp.currentProspect.FirstAgent
	if move.RetainControl {
		nextAgent = gp.currentProspect.FirstAgent
	}

	gp.currentProspect = base.Prospect[State]{
		State:      move.State,
		FirstAgent: nextAgent,
	}
}

func (gp Gameplay[State]) done() bool {
	return len(gp.game.Describe(gp.currentProspect).Moves) == 0
}

func (gp *Gameplay[State]) Play(verbose bool) int {
	log := func(format string, a ...any) (n int, err error) { return 0, nil }
	if verbose {
		log = fmt.Printf
	}

	for !gp.done() {
		player := gp.player2
		if gp.currentProspect.FirstAgent {
			player = gp.player1
		}

		log("%s's turn\n", player.Name())
		log("Current state:\n")
		log("%s\n", gp.currentProspect.State.String())
		move := player.ChooseMove(gp.currentProspect)
		log("%s chose %s\n", player.Name(), move.Summary)
		switch p := player.(type) {
		case *minimaxer.Minimaxer[State]:
			log("%s analyzed %d game states.\n", player.Name(), p.Size())
		}

		fmt.Println()

		gp.makeMove(move)
	}

	log("Final state:\n")
	log("%s\n", gp.currentProspect.State.String())

	score := gp.game.Describe(gp.currentProspect).Score
	log("Game score: %d\n", score)
	if score > 0 {
		log("%s wins!\n", gp.player1.Name())
	} else if score < 0 {
		log("%s wins!\n", gp.player2.Name())
	} else {
		log("It's a draw.\n")
	}

	return score
}

func main() {
	game := mancala.MancalaGame(6, 4)

	gp := NewGameplay(game)
	gp.player1 = minimaxer.NewMinimaxer(game, 8)
	gp.player2 = NewHumanPlayer("website", game)

	gp.Play(true)

	fmt.Println()
}
