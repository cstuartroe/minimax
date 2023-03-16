package mancala

import (
	"fmt"

	"github.com/cstuartroe/minimax/base"
)

type mancalaPit struct {
	tokens int
	name   string
	store  bool
}

type MancalaState []mancalaPit

func mancalaPitString(pit mancalaPit) string {
	if pit.store {
		return fmt.Sprintf("(%2d) ", pit.tokens)
	} else {
		return fmt.Sprintf("%2d ", pit.tokens)
	}
}

func (s MancalaState) String() string {
	out := ""
	for i := len(s) - 1; i >= len(s)/2; i-- {
		out += mancalaPitString(s[i])
	}
	out += "\n     "
	for i := 0; i < len(s)/2; i++ {
		out += mancalaPitString(s[i])
	}
	return out
}

type mancalaGame struct {
	runLength  int
	startCount int
}

func MancalaGame(runLength int, startCount int) base.Game[MancalaState] {
	return mancalaGame{runLength, startCount}
}

func repeat[T interface{}](e T, times int) []T {
	out := []T{}
	for i := 0; i < times; i++ {
		out = append(out, e)
	}
	return out
}

func (g mancalaGame) InitialState() MancalaState {
	pits := []mancalaPit{}
	for _, name := range []string{"First player's", "Second player's"} {
		for i := 0; i < g.runLength; i++ {
			pits = append(pits, mancalaPit{
				tokens: g.startCount,
				name:   fmt.Sprintf("%s #%d pit", name, i),
				store:  false,
			})
		}
		pits = append(pits, mancalaPit{
			tokens: 0,
			name:   fmt.Sprintf("%s store", name),
			store:  true,
		})
	}

	return pits
}

func mancalaMove(prospect base.Prospect[MancalaState], i int) base.Move[MancalaState] {
	summary := fmt.Sprintf("pick up from %s", prospect.State[i].name)
	retainControl := false

	pits := []mancalaPit{}
	for _, pit := range prospect.State {
		pits = append(pits, pit)
	}

	tokensInHand := pits[i].tokens
	pits[i].tokens = 0

	for tokensInHand > 0 {
		i = (i + 1) % len(pits)

		if (pits[i].store) && ((i == len(pits)-1) == prospect.FirstAgent) {
			continue
		}

		pits[i].tokens++
		tokensInHand--
	}

	if pits[i].store {
		retainControl = true
	} else if pits[i].tokens == 1 && ((i < len(pits)/2) == prospect.FirstAgent) {
		oppositeIndex := len(pits) - i - 2

		if pits[oppositeIndex].tokens > 0 {
			myStoreIndex := len(pits) - 1
			if prospect.FirstAgent {
				myStoreIndex = (len(pits) / 2) - 1
			}

			pits[myStoreIndex].tokens += 1 + pits[oppositeIndex].tokens
			pits[i].tokens = 0
			pits[oppositeIndex].tokens = 0
		}
	}

	return base.Move[MancalaState]{
		State:         pits,
		Summary:       summary,
		RetainControl: retainControl,
	}
}

func (g mancalaGame) Describe(prospect base.Prospect[MancalaState]) base.StateDescriptor[MancalaState] {
	score := prospect.State[g.runLength].tokens - prospect.State[2*g.runLength+1].tokens
	moves := []base.Move[MancalaState]{}

	offset := 0
	if !prospect.FirstAgent {
		offset = g.runLength + 1
	}
	for i := 0; i < g.runLength; i++ {
		if prospect.State[i+offset].tokens == 0 {
			continue
		}

		moves = append(moves, mancalaMove(prospect, i+offset))
	}

	if len(moves) == 0 {
		opponentStones := 0
		offset := 0
		if prospect.FirstAgent {
			offset = g.runLength + 1
		}
		for i := 0; i < g.runLength; i++ {
			opponentStones += prospect.State[i+offset].tokens
		}
		if opponentStones > 0 {
			moves = append(moves, base.Move[MancalaState]{
				Summary:       "pass",
				State:         prospect.State,
				RetainControl: false,
			})
		}
	}

	return base.StateDescriptor[MancalaState]{
		Score: score,
		Moves: moves,
	}
}
