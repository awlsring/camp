package handler

import (
	"context"
	"errors"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	camperror "github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) reportSystemChangeErrorHandler(err error) (camplocal.ReportSystemChangeRes, error) {
	var campErr camperror.Error
	if errors.As(err, &campErr) {
		e := campErr.CampError()
		switch e {
		case camperror.ErrResourceNotFound:
			return &camplocal.ResourceNotFoundExceptionResponseContent{
				Message: err.Error(),
			}, nil
		case camperror.ErrValidation:
			return &camplocal.ValidationExceptionResponseContent{
				Message: err.Error(),
			}, nil
		}
	}
	return nil, err
}

func (h Handler) ReportSystemChange(ctx context.Context, req *camplocal.ReportSystemChangeRequestContent) (camplocal.ReportSystemChangeRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke ReportSystemChange")

	iid, err := machine.InternalIdentifierFromString(req.Summary.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.Summary.InternalIdentifier)
		return h.reportSystemChangeErrorHandler(err)
	}

	class, err := machine.MachineClassFromString(string(req.Summary.GetClass().Value))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse class %s", req.Summary.Class)
		return h.reportSystemChangeErrorHandler(err)
	}

	sys := systemSummaryToDomain(req.Summary.System)
	cpu := cpuSummaryToDomain(req.Summary.CPU)
	mem := memorySummaryToDomain(req.Summary.Memory)
	disk, err := diskSummariesToDomain(req.Summary.Disks)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse disk summaries")
		return h.reportSystemChangeErrorHandler(err)
	}

	nic, err := networkInterfaceSummariesToDomain(req.Summary.NetworkInterfaces)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse network interface summaries")
		return h.reportSystemChangeErrorHandler(err)
	}

	vol, err := volumeSummariesToDomain(req.Summary.Volumes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse volume summaries")
		return h.reportSystemChangeErrorHandler(err)
	}

	ips, err := addressSummariesToDomain(req.Summary.Addresses)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse ip address summaries")
		return h.reportSystemChangeErrorHandler(err)
	}

	err = h.mSvc.ReportSystemChange(ctx, iid, class, sys, cpu, mem, disk, nic, vol, ips)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update machine")
		return nil, err
	}
	return &camplocal.ReportSystemChangeResponseContent{
		Success: true,
	}, nil
}
