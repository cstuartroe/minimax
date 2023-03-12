package main

import (
	"fmt"

	"github.com/cstuartroe/minimax/base"
	"github.com/cstuartroe/minimax/nim"
)

type Minimaxer[State base.GameState] struct {
	game           base.Game[State]
	prospectScores map[string]int
}

func (m *Minimaxer[State]) chooseMove(prospect base.Prospect[State]) (int, *base.Move[State]) {
	sd := m.game.Describe(prospect)

	if len(sd.Moves) == 0 {
		return sd.Score, nil
	}

	score := 0
	var chosenMove base.Move[State]

	for i, move := range sd.Moves {
		ps := m.getProspectScore(base.Prospect[State]{State: move.State, FirstAgent: !prospect.FirstAgent})
		if (i == 0) || (ps > score && prospect.FirstAgent) || (ps < score && !prospect.FirstAgent) {
			score = ps
			chosenMove = move
		}
	}

	return score, &chosenMove
}

func (m *Minimaxer[State]) getProspectScore(prospect base.Prospect[State]) int {
	scoreString := prospect.String()

	if _, ok := m.prospectScores[scoreString]; !ok {
		m.prospectScores[scoreString], _ = m.chooseMove(prospect)
	}

	return m.prospectScores[scoreString]
}

type Gameplay[State base.GameState] struct {
	game            base.Game[State]
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

func (gp *Gameplay[State]) playerMove() {
	sd := gp.game.Describe(gp.currentProspect)

	fmt.Println("Current state:")
	fmt.Println(gp.currentProspect.State.String())
	fmt.Println("Possible moves:")
	for i, move := range sd.Moves {
		fmt.Printf("%d: %s\n", i, move.Summary)
	}
	fmt.Print("Choose: ")
	var choice int
	fmt.Scan(&choice)
	gp.makeMove(sd.Moves[choice].State)
}

func main() {
	game := nim.GenerateNim([]int{5, 4, 3, 2}, 0, true)

	gp := NewGameplay(game)
	mx := Minimaxer[nim.NimState]{
		game:           game,
		prospectScores: map[string]int{},
	}

	for !gp.done() {
		if gp.currentProspect.FirstAgent {
			gp.playerMove()
		} else {
			_, newState := mx.chooseMove(gp.currentProspect)
			gp.makeMove(newState.State)
		}
	}

	fmt.Println("Final state:")
	fmt.Println(gp.currentProspect.State.String())
	fmt.Printf("Game score: %d\n", gp.game.Describe(gp.currentProspect).Score)
}
