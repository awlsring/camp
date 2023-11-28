package machine

import (
	"context"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/ports/repository"
	"github.com/awlsring/camp/apps/local/internal/core/ports/service"
	"github.com/rs/zerolog/log"
)

var _ service.Machine = &machineService{}

type machineService struct {
	repo repository.Machine
}

func NewMachineService(repo repository.Machine) service.Machine {
	return &machineService{
		repo: repo,
	}
}

func (s *machineService) RegisterMachine(ctx context.Context, id machine.Identifier, class machine.MachineClass, sys *machine.System, cpu *machine.Cpu, mem *machine.Memory, disks []*machine.Disk, nics []*machine.NetworkInterface, vols []*machine.Volume, ips []*machine.IpAddress) error {
	log.Debug().Msg("Invoke Controller.RegisterMachine")

	now := time.Now()
	machine := machine.NewMachine(id, class, now, now, now, machine.MachineStatusRunning, sys, cpu, mem, disks, nics, vols, ips)

	return s.repo.Add(ctx, machine)
}

func (s *machineService) DescribeMachine(ctx context.Context, identifier machine.Identifier) (*machine.Machine, error) {
	log.Debug().Msg("Invoke Controller.DescribeMachine")
	return s.repo.Get(ctx, identifier)
}

func (s *machineService) ListMachines(ctx context.Context) ([]*machine.Machine, error) {
	log.Debug().Msg("Invoke Controller.ListMachines")
	return s.repo.List(ctx, nil)
}

func (s *machineService) AcknowledgeHeartbeat(ctx context.Context, id machine.Identifier) error {
	log.Debug().Msg("Invoke Controller.AcknowledgeHeartbeat")
	return s.repo.UpdateHeartbeat(ctx, id)
}

func (s *machineService) UpdateStatus(ctx context.Context, id machine.Identifier, status machine.MachineStatus) error {
	log.Debug().Msg("Invoke Controller.UpdateStatus")
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *machineService) ReportSystemChange(ctx context.Context, id machine.Identifier, class machine.MachineClass, sys *machine.System, cpu *machine.Cpu, mem *machine.Memory, disks []*machine.Disk, nics []*machine.NetworkInterface, vols []*machine.Volume, ips []*machine.IpAddress) error {
	panic("not implemented")
}
