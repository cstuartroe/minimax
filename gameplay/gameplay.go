package gameplay

import (
	"fmt"

	"github.com/cstuartroe/minimax/games"
	"github.com/cstuartroe/minimax/minimaxer"
)

type Player[State games.GameState] interface {
	Name() string
	ChooseMove(games.Prospect[State]) games.Move[State]
}

type HumanPlayer[State games.GameState] struct {
	name string
	game games.Game[State]
}

func NewHumanPlayer[State games.GameState](name string, game games.Game[State]) HumanPlayer[State] {
	return HumanPlayer[State]{name, game}
}

func (p HumanPlayer[State]) Name() string {
	return p.name
}

func (p HumanPlayer[State]) ChooseMove(prospect games.Prospect[State]) games.Move[State] {
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

type Gameplay[State games.GameState] struct {
	game            games.Game[State]
	player1         Player[State]
	player2         Player[State]
	currentProspect games.Prospect[State]
}

func NewGameplay[State games.GameState](game games.Game[State], player1 Player[State], player2 Player[State]) Gameplay[State] {
	return Gameplay[State]{
		game:    game,
		player1: player1,
		player2: player2,
		currentProspect: games.Prospect[State]{
			State:      game.InitialState(),
			FirstAgent: true,
		},
	}
}

func (gp *Gameplay[State]) makeMove(move games.Move[State]) {
	nextAgent := !gp.currentProspect.FirstAgent
	if move.RetainControl {
		nextAgent = gp.currentProspect.FirstAgent
	}

	gp.currentProspect = games.Prospect[State]{
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
