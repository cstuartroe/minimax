package minimaxer

import (
	"fmt"

	"github.com/cstuartroe/minimax/base"
)

type Minimaxer[State base.GameState] struct {
	game           base.Game[State]
	prospectScores map[string]int
}

func NewMinimaxer[State base.GameState](game base.Game[State]) *Minimaxer[State] {
	return &Minimaxer[State]{
		game:           game,
		prospectScores: map[string]int{},
	}
}

func (m *Minimaxer[State]) Name() string {
	return fmt.Sprintf("Minimaxer @%p", m)
}

func (m *Minimaxer[State]) ChooseMove(prospect base.Prospect[State]) base.Move[State] {
	_, move := m.chooseMove(prospect)
	return *move
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
