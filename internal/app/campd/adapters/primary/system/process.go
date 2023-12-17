package system

import "context"

type Process struct {
}

func (p *Process) Start(ctx context.Context) {
	// startup and register

	// report status as up

	// start a loop that will continue to heartbeat

	// if context is broken, report status as down and close
}
