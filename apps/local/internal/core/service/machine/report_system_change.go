package machine

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) ReportSystemChange(ctx context.Context, id machine.Identifier, sys *machine.System, cpu *machine.Cpu, mem *machine.Memory, disks []*machine.Disk, nics []*machine.NetworkInterface, vols []*machine.Volume, ips []*machine.IpAddress) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Reporting system change of machine %s", id)

	log.Debug().Msgf("Getting original machine %s from repo", id)
	original, err := s.DescribeMachine(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get original machine %s from repo", id)
		return err
	}

	log.Debug().Msg("Rebuilding machine description")
	original.System = sys
	original.Cpu = cpu
	original.Memory = mem
	original.Disks = disks
	original.NetworkInterfaces = nics
	original.Volumes = vols
	original.Addresses = ips

	log.Debug().Msgf("Updating machine %s in repo", id)
	err = s.repo.Update(ctx, original)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to update machine %s in repo", id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("System change of machine %s reported", id)
	return nil
}
