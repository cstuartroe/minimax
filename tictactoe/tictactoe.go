package tictactoe

import "github.com/cstuartroe/minimax/base"

type TicTacToeSquare byte

const (
	X     TicTacToeSquare = 'X'
	O     TicTacToeSquare = 'O'
	Space TicTacToeSquare = ' '
)

type TicTacToeBoard [3][3]TicTacToeSquare

func (board TicTacToeBoard) String() string {
	out := []TicTacToeSquare{}
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			out = append(out, board[x][y])
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
	return b[i.x][i.y]
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
			out[x][y] = board[x][y]
		}
	}
	return out
}

func getMoves(prospect base.Prospect[TicTacToeBoard]) []TicTacToeBoard {
	agentByte := O
	if prospect.FirstAgent {
		agentByte = X
	}

	out := []TicTacToeBoard{}

	if getWinner(prospect.State) != Space {
		return out
	}

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if prospect.State[x][y] == Space {
				newBoard := copy(prospect.State)
				newBoard[x][y] = agentByte
				out = append(out, newBoard)
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
