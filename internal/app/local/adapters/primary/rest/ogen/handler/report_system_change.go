package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
)

func (h Handler) ReportSystemChange(ctx context.Context, req *camplocal.ReportSystemChangeRequestContent) (camplocal.ReportSystemChangeRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke ReportSystemChange")

	id, err := machine.IdentifierFromString(req.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.Identifier)
		return nil, err
	}

	host := hostSummaryToDomain(req.Summary.Host)
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

	// vol, err := volumeSummariesToDomain(req.Summary.Volumes)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to parse volume summaries")
	// 	return nil, err
	// }

	ips, err := addressSummariesToDomain(req.Summary.Addresses)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse ip address summaries")
		return nil, err
	}

	err = h.mSvc.ReportSystemChange(ctx, id, host, cpu, mem, disk, nic, nil, ips)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update machine")
		return nil, err
	}
	return &camplocal.ReportSystemChangeResponseContent{
		Success: true,
	}, nil
}
