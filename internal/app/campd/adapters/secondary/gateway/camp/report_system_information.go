package camp_reporting

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/secondary/gateway/camp/conversion"
	"github.com/awlsring/camp/internal/app/campd/core/domain/system"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
)

func (c *CampLocalReporting) ReportSystemInformation(ctx context.Context, id machine.Identifier, system *system.System) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("sending system change system")

	in := &local.ReportSystemChangeInput{
		Identifier: conversion.InternalIdentifierFromDomain(id),
		Summary:    conversion.SystemFromDomain(system),
	}

	log.Debug().Msg("sending register request")
	_, err := c.client.ReportSystemChange(ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to send system change request")
		return err
	}

	log.Debug().Msg("system change request sent")
	return nil
}
