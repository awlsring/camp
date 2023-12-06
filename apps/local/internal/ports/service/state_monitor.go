package service

import "context"

type StateMonitor interface {
	VerifyAndAdjustMachineStates(ctx context.Context) error
}
