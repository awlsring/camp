package handler

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h Handler) ReportSystemChange(ctx context.Context, req *camplocal.ReportSystemChangeRequestContent) (camplocal.ReportSystemChangeRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke ReportSystemChange")

	id, err := machine.IdentifierFromString(req.Summary.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.Summary.InternalIdentifier)
		return nil, err
	}

	class, err := machine.MachineClassFromString(string(req.Summary.GetClass().Value))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse class %s", req.Summary.Class.Value)
		return nil, err
	}

	sys := systemSummaryToDomain(req.Summary.System)
	cpu := cpuSummaryToDomain(req.Summary.CPU)
	mem := memorySummaryToDomain(req.Summary.Memory)
	disk, err := diskSummariesToDomain(req.Summary.Disks)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse disk summaries")
		return nil, err
	}

	nic, err := networkInterfaceSummariesToDomain(req.Summary.NetworkInterfaces)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse network interface summaries")
		return nil, err
	}

	vol, err := volumeSummariesToDomain(req.Summary.Volumes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse volume summaries")
		return nil, err
	}

	ips, err := addressSummariesToDomain(req.Summary.Addresses)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse ip address summaries")
		return nil, err
	}

	err = h.mSvc.ReportSystemChange(ctx, id, class, sys, cpu, mem, disk, nic, vol, ips)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update machine")
		return nil, err
	}
	return &camplocal.ReportSystemChangeResponseContent{
		Success: true,
	}, nil
}
