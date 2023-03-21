package minimaxer

import (
	"fmt"
	"math/rand"

	"github.com/cstuartroe/minimax/gameplay"
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

type RatedMove[State games.GameState] struct {
	Move  games.Move[State]
	Score int
}

func (m *Minimaxer[State]) RateChoices(prospect games.Prospect[State]) []RatedMove[State] {
	m.prospectScores = map[string]int{}
	out := []RatedMove[State]{}

	for _, move := range m.game.Describe(prospect).Moves {
		out = append(out, RatedMove[State]{
			Score: m.getProspectScore(moveToProspect(move, prospect.FirstAgent), m.lookahead),
			Move:  move,
		})
	}

	return out
}

func (m Minimaxer[State]) Comment() string {
	return fmt.Sprintf("I analyzed %d game states!", len(m.prospectScores))
}

func moveToProspect[State games.GameState](move games.Move[State], firstAgent bool) games.Prospect[State] {
	if !move.RetainControl {
		firstAgent = !firstAgent
	}

	return games.Prospect[State]{State: move.State, FirstAgent: firstAgent}
}

func (m *Minimaxer[State]) chooseMove(prospect games.Prospect[State], searchDepth int) (int, *games.Move[State]) {
	sd := m.game.Describe(prospect)

	if len(sd.Moves) == 0 || searchDepth == 0 {
		return sd.Score, nil
	}

	score := 0
	var goodMoves []games.Move[State]

	for i, move := range sd.Moves {
		ps := m.getProspectScore(moveToProspect(move, prospect.FirstAgent), searchDepth-1)
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

type AssistedHumanPlayer[State games.GameState] struct {
	human     gameplay.HumanPlayer[State]
	minimaxer *Minimaxer[State]
}

func NewAssistedHumanPlayer[State games.GameState](human gameplay.HumanPlayer[State], minimaxer *Minimaxer[State]) AssistedHumanPlayer[State] {
	return AssistedHumanPlayer[State]{human, minimaxer}
}

func (p AssistedHumanPlayer[State]) Name() string {
	return p.human.Name() + " with some help from " + p.minimaxer.Name()
}

func (p AssistedHumanPlayer[State]) ChooseMove(prospect games.Prospect[State]) games.Move[State] {
	move := p.human.ChooseMove(prospect)

	bestScore := -100000
	bestMoves := []games.Move[State]{}

	fmt.Println("Good thought! Here's how the minimaxer rates the moves:")
	for _, ratedMove := range p.minimaxer.RateChoices(prospect) {
		fmt.Printf("%s: %d\n", ratedMove.Move.Summary, ratedMove.Score)
		score := ratedMove.Score
		if !prospect.FirstAgent {
			score = -ratedMove.Score
		}
		if score > bestScore {
			bestScore = score
			bestMoves = []games.Move[State]{ratedMove.Move}
		} else if score == bestScore {
			bestMoves = append(bestMoves, ratedMove.Move)
		}
	}

	recommendation := ""
	agreed := false
	for i, bestMove := range bestMoves {
		if i > 0 {
			recommendation += " or "
		}
		recommendation += bestMove.Summary
		agreed = agreed || bestMove.Summary == move.Summary
	}
	fmt.Printf("%s recommends %s\n", p.minimaxer.Name(), recommendation)
	if agreed {
		fmt.Println("Looks like you two agree!")
	}

	fmt.Print("Change choice? ")
	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil && err.Error() != "unexpected newline" {
		panic(err)
	}
	if answer != "" && (answer[0] == 'y' || answer[0] == 'Y') {
		move = p.human.ChooseMove(prospect)
	}

	return move
}

func (p AssistedHumanPlayer[State]) Comment() string {
	return fmt.Sprintf("%s says %q, %s says %q", p.human.Name(), p.human.Comment(), p.minimaxer.Name(), p.minimaxer.Comment())
}
