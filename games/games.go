package games

import "fmt"

type GameState interface {
	String() string
}

type Prospect[State GameState] struct {
	State      State
	FirstAgent bool
}

func (p Prospect[State]) String() string {
	return fmt.Sprintf("%t%s", p.FirstAgent, p.State.String())
}

type Move[State GameState] struct {
	Summary       string
	State         State
	RetainControl bool
}

type StateDescriptor[State GameState] struct {
	Score int
	Moves []Move[State]
}

type Game[State GameState] interface {
	InitialState() State
	Describe(Prospect[State]) StateDescriptor[State]
}
