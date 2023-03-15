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
	FinalRemarks() string
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
	fmt.Scan(&choice)
	return sd.Moves[choice]
}

func (p HumanPlayer[State]) FinalRemarks() string {
	return "Thanks for playing!"
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

		log("Current state:\n")
		log("%s\n", gp.currentProspect.State.String())
		move := player.ChooseMove(gp.currentProspect)
		log("%s chose %s\n\n", player.Name(), move.Summary)
		gp.makeMove(move.State)
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

	log("\n")
	for _, player := range []Player[State]{gp.player1, gp.player2} {
		log("%s says: %q\n", player.Name(), player.FinalRemarks())
	}

	return score
}

func perfectPlay[State base.GameState](game base.Game[State]) int {
	gp := NewGameplay(game)

	gp.player1 = minimaxer.NewMinimaxer(game)
	gp.player2 = minimaxer.NewMinimaxer(game)

	return gp.Play(false)
}

func main() {
	for total := 0; total <= 12; total++ {
		for _, piles := range nim.NimStates(total, total) {
			fmt.Print(piles)

			score := perfectPlay(nim.GenerateNim(piles, 0, false))

			word := "Win"
			if score == -1 {
				word = "Lose"
			}

			fmt.Printf(" %s\n", word)

			mscore := perfectPlay(nim.GenerateNim(piles, 0, true))

			if score != mscore {
				fmt.Println("Misere is different!")
			}
		}
		fmt.Println()
	}
}
