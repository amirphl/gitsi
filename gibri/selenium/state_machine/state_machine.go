package statemachine

import (
	"context"
	"time"

	"github.com/amirphl/gitsi/gibri/component"
	"github.com/amirphl/gitsi/gibri/component/event"
	"github.com/amirphl/gitsi/gibri/component/state"
	"github.com/looplab/fsm"
)

type StateMachine interface {
	component.NotifyingStateMachine
}

type TransitionTimeTracker interface {
	Timestamp() time.Time
	MaybeUpdate(eventOccured bool)
	ExceededTimeout(timeout time.Duration) bool
}

type stateMachine struct {
	f        *fsm.FSM
	handlers []component.TransitionHandler
}

type transitionTimeTracker struct {
	timestamp time.Time
}

func (s *stateMachine) Notify(from, to state.State) {
	for _, f := range s.handlers {
		f(from, to)
	}
}

func (s *stateMachine) OnTransition(h component.TransitionHandler) {
	s.handlers = append(s.handlers, h)
}

func (s *transitionTimeTracker) Timestamp() time.Time {
	return s.timestamp
}

func (s *transitionTimeTracker) MaybeUpdate(eventOccured bool) {
	if eventOccured && s.timestamp.IsZero() {
		s.timestamp = time.Now()
	} else if !eventOccured {
		s.timestamp = time.Time{}
	}
}

func (s *transitionTimeTracker) ExceededTimeout(timeout time.Duration) bool {
	if s.timestamp.IsZero() {
		return false
	}

	return time.Since(s.timestamp) > timeout
}

func New() StateMachine {
	startingUp := state.StartingUp.String()
	running := state.Running.String()
	finished := state.Finished.String()

	sm := &stateMachine{
		handlers: []component.TransitionHandler{},
	}

	// TODO https://github.com/Tinder/StateMachine/blob/main/src/main/kotlin/com/tinder/StateMachine.kt
	// TODO https://pkg.go.dev/github.com/looplab/fsm#section-readme

	f := fsm.NewFSM(
		startingUp,
		fsm.Events{
			{
				Name: event.CallJoined.String(),
				Src:  []string{startingUp},
				Dst:  running,
			},
			{
				Name: event.FailedToJoinCall.String(),
				Src:  []string{startingUp},
				Dst:  state.FailedToJoinCall.String(),
			},
			{
				Name: event.ChromeHung.String(),
				Src:  []string{startingUp},
				Dst:  state.ChromeHung.String(),
			},
			{
				Name: event.CallEmpty.String(),
				Src:  []string{running},
				Dst:  finished,
			},
			{
				Name: event.NoMediaReceived.String(),
				Src:  []string{running},
				Dst:  state.NoMediaReceived.String(),
			},
			{
				Name: event.ICEFailed.String(),
				Src:  []string{running},
				Dst:  state.ICEFailed.String(),
			},
			{
				Name: event.LocalParticipantKicked.String(),
				Src:  []string{running},
				Dst:  finished,
			},
			{
				Name: event.ChromeHung.String(),
				Src:  []string{running},
				Dst:  state.ChromeHung.String(),
			},
			// TODO dontTransition()
		},
		fsm.Callbacks{
			"after_event": func(_ context.Context, e *fsm.Event) {
				// TODO throw Exception("Invalid state transition: $it")
				if e.Src != e.Dst {
					sm.Notify(state.FromString(e.Src), state.FromString(e.Dst))
				}
			},
		},
	)

	sm.f = f

	return sm
}

func NewTransitionTimeTracker() TransitionTimeTracker {
	return &transitionTimeTracker{}
}
