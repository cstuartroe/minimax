package connect_four

import (
	"fmt"

	"github.com/cstuartroe/minimax/games"
)

type ConnectFourPiece byte

const (
	CFBlank ConnectFourPiece = iota
	CFRed
	CFYellow
)

func (p ConnectFourPiece) String() string {
	if p == CFBlank {
		return " "
	} else if p == CFRed {
		return "X"
	} else if p == CFYellow {
		return "O"
	}
	panic("?")
}

type ConnectFourState [6][7]ConnectFourPiece

func (s ConnectFourState) String() string {
	out := ""
	for y := 5; y >= 0; y-- {
		row := s[y]
		out += "| "
		for _, piece := range row {
			out += piece.String() + " "
		}
		out += "|\n"
	}
	return out
}

func (s ConnectFourState) copy() ConnectFourState {
	out := ConnectFourState{}

	for x := 0; x < 7; x++ {
		for y := 0; y < 6; y++ {
			out[y][x] = s[y][x]
		}
	}

	return out
}

type _ConnectFour struct{}

func (cf _ConnectFour) InitialState() ConnectFourState {
	return ConnectFourState{}
}

type Position struct {
	x int8
	y int8
}

func (s ConnectFourState) at(p Position) ConnectFourPiece {
	return s[p.y][p.x]
}

func getStreak(start, delta, size Position, length int8) *[]Position {
	out := []Position{}

	pos := start

	for i := int8(0); i < length; i++ {
		if pos.x < 0 || pos.x >= size.x || pos.y < 0 || pos.y >= size.y {
			return nil
		}
		out = append(out, pos)
		pos.x += delta.x
		pos.y += delta.y
	}

	return &out
}

var boardSize Position = Position{7, 6}

func getAllStreaks() [][]Position {
	out := [][]Position{}

	for x := int8(0); x < boardSize.x; x++ {
		for y := int8(0); y < boardSize.y; y++ {
			for _, delta := range []Position{{0, 1}, {-1, 1}, {1, 1}, {1, 0}} {
				streak := getStreak(Position{x, y}, delta, boardSize, 4)
				if streak != nil {
					out = append(out, *streak)
				}
			}
		}
	}

	return out
}

var allStreaks [][]Position = getAllStreaks()

func getScore(s ConnectFourState) int {
	for _, streak := range allStreaks {
		piece := s.at(streak[0])
		if piece == CFBlank {
			continue
		}
		won := true
		for _, pos := range streak {
			if s.at(pos) != piece {
				won = false
				break
			}
		}
		if won {
			if piece == CFRed {
				return 1
			} else {
				return -1
			}
		}
	}
	return 0
}

func getMoves(prospect games.Prospect[ConnectFourState]) []games.Move[ConnectFourState] {
	out := []games.Move[ConnectFourState]{}
	var piece ConnectFourPiece
	if prospect.FirstAgent {
		piece = CFRed
	} else {
		piece = CFYellow
	}

	for x := int8(0); x < boardSize.x; x++ {
		row := 5
		if prospect.State[row][x] == CFBlank {
			for row > 0 && prospect.State[row-1][x] == CFBlank {
				row--
			}

			newState := prospect.State.copy()
			newState[row][x] = piece

			out = append(out, games.Move[ConnectFourState]{
				Summary:       fmt.Sprintf("Go in column #%d", x),
				State:         newState,
				RetainControl: false,
			})
		}
	}

	return out
}

func (cg _ConnectFour) Describe(prospect games.Prospect[ConnectFourState]) games.StateDescriptor[ConnectFourState] {
	score := getScore(prospect.State)

	moves := []games.Move[ConnectFourState]{}
	if score == 0 {
		moves = getMoves(prospect)
	}

	return games.StateDescriptor[ConnectFourState]{
		Score: score,
		Moves: moves,
	}
}

func ConnectFour() games.Game[ConnectFourState] {
	return _ConnectFour{}
}
