package minimaxer

import (
	"fmt"
	"math/rand"

	"github.com/cstuartroe/minimax/games"
)

type Minimaxer[State games.GameState] struct {
	game           games.Game[State]
	prospectScores map[string]int
	lookahead      int
}

func NewMinimaxer[State games.GameState](game games.Game[State], lookahead int) *Minimaxer[State] {
	return &Minimaxer[State]{
		game:      game,
		lookahead: lookahead,
	}
}

func (m *Minimaxer[State]) Name() string {
	return fmt.Sprintf("Minimaxer @%p", m)
}

func (m Minimaxer[State]) Size() int {
	return len(m.prospectScores)
}

func (m *Minimaxer[State]) ChooseMove(prospect games.Prospect[State]) games.Move[State] {
	m.prospectScores = map[string]int{}
	_, move := m.chooseMove(prospect, m.lookahead)
	return *move
}

func (m Minimaxer[State]) Comment() string {
	return fmt.Sprintf("I analyzed %d game states!", len(m.prospectScores))
}

func (m *Minimaxer[State]) chooseMove(prospect games.Prospect[State], searchDepth int) (int, *games.Move[State]) {
	sd := m.game.Describe(prospect)

	if len(sd.Moves) == 0 || searchDepth == 0 {
		return sd.Score, nil
	}

	score := 0
	var goodMoves []games.Move[State]

	for i, move := range sd.Moves {
		firstAgent := !prospect.FirstAgent
		if move.RetainControl {
			firstAgent = prospect.FirstAgent
		}

		ps := m.getProspectScore(games.Prospect[State]{State: move.State, FirstAgent: firstAgent}, searchDepth-1)
		if (i == 0) || (ps > score && prospect.FirstAgent) || (ps < score && !prospect.FirstAgent) {
			score = ps
			goodMoves = []games.Move[State]{move}
		} else if ps == score {
			goodMoves = append(goodMoves, move)
		}
	}

	chosenMove := goodMoves[rand.Intn(len(goodMoves))]

	return score, &chosenMove
}

func (m *Minimaxer[State]) getProspectScore(prospect games.Prospect[State], searchDepth int) int {
	scoreString := prospect.String()

	if _, ok := m.prospectScores[scoreString]; !ok {
		m.prospectScores[scoreString], _ = m.chooseMove(prospect, searchDepth)
	}

	return m.prospectScores[scoreString]
}
