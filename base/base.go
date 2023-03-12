package base

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

type StateDescriptor[State GameState] struct {
	Score int
	Moves []State
}

type Game[State GameState] struct {
	InitialState State
	Describe     func(Prospect[State]) StateDescriptor[State]
}
