package machine

import (
	"context"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) RegisterMachine(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, key machine.AgentKey, class machine.MachineClass, cap machine.PowerCapabilities, sys *machine.System, cpu *machine.Cpu, mem *machine.Memory, disks []*machine.Disk, nics []*machine.NetworkInterface, vols []*machine.Volume, ips []*machine.IpAddress) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Registering machine with identifier %s", id)

	now := time.Now().UTC()
	machine := machine.NewMachine(id, endpoint, key, class, now, now, now, machine.MachineStatusRunning, cap, sys, cpu, mem, disks, nics, vols, ips)

	log.Debug().Msgf("Adding machine with identifier %s in repo", id)
	err := s.repo.Add(ctx, machine)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to add machine with identifier %s in repo", id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("Machine with identifier %s registered", id)
	return nil
}
