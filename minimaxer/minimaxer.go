package minimaxer

import (
	"fmt"
	"math/rand"

	"github.com/cstuartroe/minimax/base"
)

type Minimaxer[State base.GameState] struct {
	game           base.Game[State]
	prospectScores map[string]int
	lookahead      int
}

func NewMinimaxer[State base.GameState](game base.Game[State], lookahead int) *Minimaxer[State] {
	return &Minimaxer[State]{
		game:           game,
		prospectScores: map[string]int{},
		lookahead:      lookahead,
	}
}

func (m *Minimaxer[State]) Name() string {
	return fmt.Sprintf("Minimaxer @%p", m)
}

func (m Minimaxer[State]) Size() int {
	return len(m.prospectScores)
}

func (m *Minimaxer[State]) ChooseMove(prospect base.Prospect[State]) base.Move[State] {
	m.prospectScores = map[string]int{}
	_, move := m.chooseMove(prospect, m.lookahead)
	return *move
}

func (m *Minimaxer[State]) chooseMove(prospect base.Prospect[State], searchDepth int) (int, *base.Move[State]) {
	sd := m.game.Describe(prospect)

	if len(sd.Moves) == 0 || searchDepth == 0 {
		return sd.Score, nil
	}

	score := 0
	var goodMoves []base.Move[State]

	for i, move := range sd.Moves {
		firstAgent := !prospect.FirstAgent
		if move.RetainControl {
			firstAgent = prospect.FirstAgent
		}

		ps := m.getProspectScore(base.Prospect[State]{State: move.State, FirstAgent: firstAgent}, searchDepth-1)
		if (i == 0) || (ps > score && prospect.FirstAgent) || (ps < score && !prospect.FirstAgent) {
			score = ps
			goodMoves = []base.Move[State]{move}
		} else if ps == score {
			goodMoves = append(goodMoves, move)
		}
	}

	chosenMove := goodMoves[rand.Intn(len(goodMoves))]

	return score, &chosenMove
}

func (m *Minimaxer[State]) getProspectScore(prospect base.Prospect[State], searchDepth int) int {
	scoreString := prospect.String()

	if _, ok := m.prospectScores[scoreString]; !ok {
		m.prospectScores[scoreString], _ = m.chooseMove(prospect, searchDepth)
	}

	return m.prospectScores[scoreString]
}
