package service_discovery

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
)

type state int

const (
	stateInitial state = iota
	stateLookup
	stateServer
	stateClient
)

func (s state) String() string {
	switch s {
	case stateClient:
		return "Client"

	case stateServer:
		return "Server"

	case stateLookup:
		return "Lookup"

	case stateInitial:
		return "Initial"
	}
	return ""
}

type action struct {
	state state
	data  any
}

type actionHandler func(a action, ch chan action, logger *slog.Logger)

var handlers = map[state]actionHandler{
	stateLookup: lookupHandler,
	stateClient: clientHandler,
	stateServer: serverHandler,
}

type StateMachine struct {
	logger   *slog.Logger
	state    state
	actionch chan action
}

func NewStateMachine(logger *slog.Logger) *StateMachine {
	sm := &StateMachine{
		logger:   logger,
		state:    stateInitial,
		actionch: make(chan action, 1),
	}

	sm.lookup()

	return sm
}

func (sm *StateMachine) Run(ctx context.Context) error {
	sm.logger.Info("started")
	defer sm.logger.Info("stopped")
	for {
		select {
		case a := <-sm.actionch:
			sm.logger.Info("state changed", "old_state", sm.state.String(), "new_state", a.state.String())
			sm.state = a.state
			h, ok := handlers[sm.state]
			if !ok {
				return fmt.Errorf("no handlers for state %q", sm.state)
			}
			go h(a, sm.actionch, sm.logger)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (sm *StateMachine) lookup() {
	sm.actionch <- action{state: stateLookup}
}

func lookupHandler(a action, ch chan action, logger *slog.Logger) {
	services := LookupService(logger)

	if len(services) > 0 {
		ch <- action{state: stateClient, data: services}
	} else {
		ch <- action{state: stateServer, data: nil}
	}
}

func serverHandler(a action, ch chan action, logger *slog.Logger) {
	time.Sleep(5 * time.Second)

	ch <- action{state: stateClient, data: nil}
}

func clientHandler(a action, ch chan action, logger *slog.Logger) {
	time.Sleep(5 * time.Second)

	ch <- action{state: stateLookup, data: nil}
}
