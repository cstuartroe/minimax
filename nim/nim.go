package nim

import (
	"fmt"

	"github.com/cstuartroe/minimax/games"
)

type NimState []int

func (s NimState) String() string {
	out := ""
	for i, pile := range s {
		out += fmt.Sprintf("%d:%d ", i, pile)
	}
	return out
}

type nimGame struct {
	initialState NimState
	maxTake      int
	misere       bool
}

func NimGame(initialState NimState, maxTake int, misere bool) games.Game[NimState] {
	return nimGame{initialState, maxTake, misere}
}

func (g nimGame) InitialState() NimState {
	return g.initialState
}

func (g nimGame) Describe(prospect games.Prospect[NimState]) games.StateDescriptor[NimState] {
	moves := []games.Move[NimState]{}

	for i, pile := range prospect.State {
		turnMaxTake := g.maxTake
		if g.maxTake == 0 || g.maxTake > pile {
			turnMaxTake = pile
		}

		for take := 1; take <= turnMaxTake; take++ {
			newState := NimState{}
			for _, n := range prospect.State {
				newState = append(newState, n)
			}
			newState[i] = pile - take

			summary := fmt.Sprintf("Take %d from pile #%d", take, i)

			moves = append(moves, games.Move[NimState]{
				State:   newState,
				Summary: summary,
			})
		}
	}

	score := 0
	if len(moves) == 0 {
		if prospect.FirstAgent {
			score = -1
		} else {
			score = 1
		}

		if g.misere {
			score = -score
		}
	}

	return games.StateDescriptor[NimState]{
		Moves: moves,
		Score: score,
	}
}

func NimStates(total int, maxPile int) [][]int {
	if total == 0 {
		return [][]int{{}}
	}

	if maxPile > total {
		maxPile = total
	}

	out := [][]int{}

	for lastPile := 1; lastPile <= maxPile; lastPile++ {
		for _, piles := range NimStates(total-lastPile, lastPile) {
			out = append(out, append(piles, lastPile))
		}
	}

	return out
}
