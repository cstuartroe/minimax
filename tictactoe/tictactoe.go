package tictactoe

import (
	"fmt"

	"github.com/cstuartroe/minimax/base"
)

type TicTacToeSquare byte

const (
	X     TicTacToeSquare = 'X'
	O     TicTacToeSquare = 'O'
	Space TicTacToeSquare = '_'
)

type TicTacToeBoard [3][3]TicTacToeSquare

func (board TicTacToeBoard) String() string {
	out := []TicTacToeSquare{}
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			out = append(out, board[y][x])
		}
		out = append(out, '\n')
	}
	return string(out)
}

type TicTacToeIndex struct {
	x int
	y int
}

func (b TicTacToeBoard) at(i TicTacToeIndex) TicTacToeSquare {
	return b[i.y][i.x]
}

var ticTacToeRuns [][3]TicTacToeIndex = [][3]TicTacToeIndex{
	{{0, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},

	{{0, 0}, {1, 0}, {2, 0}},
	{{0, 1}, {1, 1}, {2, 1}},
	{{0, 2}, {1, 2}, {2, 2}},

	{{0, 0}, {1, 1}, {2, 2}},
	{{0, 2}, {1, 1}, {2, 0}},
}

func getWinner(board TicTacToeBoard) TicTacToeSquare {
	for _, run := range ticTacToeRuns {
		c := board.at(run[0])
		if c != Space && c == board.at(run[1]) && c == board.at(run[2]) {
			return c
		}
	}

	return Space
}

func copy(board TicTacToeBoard) TicTacToeBoard {
	out := TicTacToeBoard{}
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			out[y][x] = board[y][x]
		}
	}
	return out
}

var rowNames [3]string = [3]string{"top", "middle", "bottom"}
var columnNames [3]string = [3]string{"left", "center", "right"}

func getMoves(prospect base.Prospect[TicTacToeBoard]) []base.Move[TicTacToeBoard] {
	agentByte := O
	if prospect.FirstAgent {
		agentByte = X
	}

	out := []base.Move[TicTacToeBoard]{}

	if getWinner(prospect.State) != Space {
		return out
	}

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if prospect.State[y][x] == Space {
				newBoard := copy(prospect.State)
				newBoard[y][x] = agentByte

				summary := fmt.Sprintf("%c to %s %s", agentByte, rowNames[y], columnNames[x])

				out = append(out, base.Move[TicTacToeBoard]{Summary: summary, State: newBoard})
			}
		}
	}

	return out
}

func DescribeTicTacToe(prospect base.Prospect[TicTacToeBoard]) base.StateDescriptor[TicTacToeBoard] {
	out := base.StateDescriptor[TicTacToeBoard]{
		Moves: getMoves(prospect),
	}

	winner := getWinner(prospect.State)

	if winner == X {
		out.Score = 1
	} else if winner == O {
		out.Score = -1
	}

	return out
}

var TicTacToe base.Game[TicTacToeBoard] = base.Game[TicTacToeBoard]{
	InitialState: [3][3]TicTacToeSquare{{Space, Space, Space}, {Space, Space, Space}, {Space, Space, Space}},
	Describe:     DescribeTicTacToe,
}
