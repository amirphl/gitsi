package component

import "github.com/amirphl/gitsi/gibri/component/state"

type TransitionHandler func(from, to state.State)

type NotifyingStateMachine interface {
	Notify(from, to state.State)
	OnTransition(TransitionHandler)
}
