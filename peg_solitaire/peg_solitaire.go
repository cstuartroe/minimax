package peg_solitaire

import (
	"fmt"
	"strings"

	"github.com/cstuartroe/minimax/games"
)

type TrianglePegSolitaireState struct {
	pegs [15]bool
}

func pegSymbol(b bool) string {
	if b {
		return "."
	}
	return "o"
}

func (s TrianglePegSolitaireState) copy() [15]bool {
	out := [15]bool{}
	copy(out[:], s.pegs[:])
	return out
}

func (s TrianglePegSolitaireState) String() string {
	out := ""

	i := 0
	for y := 0; y < 5; y++ {
		out += strings.Repeat(" ", 4-y)
		for x := 0; x <= y; x++ {
			out += pegSymbol(s.pegs[i]) + " "
			i += 1
		}
		out += "\n"
	}

	return out
}

var horizontal_peg_dimension [][]int = [][]int{{0}, {1, 2}, {3, 4, 5}, {6, 7, 8, 9}, {10, 11, 12, 13, 14}}

func rotate(dimensions [][]int) [][]int {
	out := [][]int{}

	for i := range dimensions {
		new_d := []int{}
		for _, d := range dimensions {
			if len(d) > i {
				new_d = append(new_d, d[i])
			}
		}
		out = append(out, new_d)
	}

	return out
}

var peg_dimensions [][][]int = [][][]int{
	horizontal_peg_dimension,
	rotate(horizontal_peg_dimension),
	rotate(rotate(horizontal_peg_dimension)),
}

func (s TrianglePegSolitaireState) jump(from int, over int, to int) games.Move[TrianglePegSolitaireState] {
	pegs := s.copy()

	if !pegs[from] || !pegs[over] || pegs[to] {
		panic("Pegs in the wrong holes")
	}

	pegs[from] = false
	pegs[over] = false
	pegs[to] = true

	return games.Move[TrianglePegSolitaireState]{
		Summary:       fmt.Sprintf("Jump peg #%d over peg #%d", from, over),
		State:         TrianglePegSolitaireState{pegs},
		RetainControl: true,
	}
}

func (s TrianglePegSolitaireState) findJumps() []games.Move[TrianglePegSolitaireState] {
	out := []games.Move[TrianglePegSolitaireState]{}

	for _, dim := range peg_dimensions {
		for _, row := range dim {
			for i := 0; i < len(row)-2; i++ {
				a, b, c := row[i], row[i+1], row[i+2]

				if !s.pegs[b] {
					continue
				}

				if s.pegs[a] && !s.pegs[c] {
					out = append(out, s.jump(a, b, c))
				}

				if s.pegs[c] && !s.pegs[a] {
					out = append(out, s.jump(c, b, a))
				}
			}
		}
	}

	return out
}

var initialRemovals []int = []int{0, 1, 3, 4}

type _TrianglePegSolitaire struct{}

func (s _TrianglePegSolitaire) InitialState() TrianglePegSolitaireState {
	return TrianglePegSolitaireState{
		[15]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
	}
}

func (s _TrianglePegSolitaire) Describe(prospect games.Prospect[TrianglePegSolitaireState]) games.StateDescriptor[TrianglePegSolitaireState] {
	score := len(prospect.State.pegs)
	for _, p := range prospect.State.pegs {
		if p {
			score -= 1
		}
	}

	moves := []games.Move[TrianglePegSolitaireState]{}

	if score == 0 {
		for _, i := range initialRemovals {
			pegs := prospect.State.copy()
			pegs[i] = false
			moves = append(moves, games.Move[TrianglePegSolitaireState]{
				Summary:       fmt.Sprintf("Remove peg #%d", i),
				State:         TrianglePegSolitaireState{pegs},
				RetainControl: true,
			})
		}
	} else {
		moves = prospect.State.findJumps()
	}

	return games.StateDescriptor[TrianglePegSolitaireState]{
		Score: score,
		Moves: moves,
	}
}

func TrianglePegSolitaire() games.Game[TrianglePegSolitaireState] {
	return _TrianglePegSolitaire{}
}
