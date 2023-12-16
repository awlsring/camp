package machine

import (
	"context"
	"time"

	mach "github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) RegisterMachine(ctx context.Context, id mach.Identifier, endpoint mach.MachineEndpoint, key mach.AgentKey, class machine.MachineClass, cap mach.PowerCapabilities, host *host.Host, cpu *cpu.CPU, mem *memory.Memory, disks []*storage.Disk, nics []*network.Nic, vols []*storage.Volume, ips []*network.IpAddress) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Registering machine with identifier %s", id)

	now := time.Now().UTC()
	status := &power.Status{
		Status:    power.StatusCodeRunning,
		UpdatedAt: now,
	}
	machine := mach.NewMachine(id, endpoint, key, class, now, now, now, status, cap, host, cpu, mem, disks, nics, vols, ips)

	log.Debug().Msgf("Adding machine with identifier %s in repo", id)
	err := s.repo.Add(ctx, machine)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to add machine with identifier %s in repo", id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("Machine with identifier %s registered", id)
	return nil
}
