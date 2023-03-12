package nim

import (
	"fmt"

	"github.com/cstuartroe/minimax/base"
)

type NimState []int

func (s NimState) String() string {
	out := ""
	for i, pile := range s {
		out += fmt.Sprintf("%d:%d ", i, pile)
	}
	return out
}

func nimDescriber(maxTake int, misere bool) func(base.Prospect[NimState]) base.StateDescriptor[NimState] {
	return func(prospect base.Prospect[NimState]) base.StateDescriptor[NimState] {
		moves := []base.Move[NimState]{}

		for i, pile := range prospect.State {
			turnMaxTake := maxTake
			if maxTake == 0 || maxTake > pile {
				turnMaxTake = pile
			}

			for take := 1; take <= turnMaxTake; take++ {
				newState := NimState{}
				for _, n := range prospect.State {
					newState = append(newState, n)
				}
				newState[i] = pile - take

				summary := fmt.Sprintf("Take %d from pile #%d", take, i)

				moves = append(moves, base.Move[NimState]{
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

			if misere {
				score = -score
			}
		}

		return base.StateDescriptor[NimState]{
			Moves: moves,
			Score: score,
		}
	}
}

func GenerateNim(initialState []int, maxTake int, misere bool) base.Game[NimState] {
	return base.Game[NimState]{
		InitialState: initialState,
		Describe:     nimDescriber(maxTake, misere),
	}
}
