package camp_reporting

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/secondary/reporting/camp/conversion"
	"github.com/awlsring/camp/internal/app/campd/core/domain/system"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
)

func (c *CampLocalReporting) Register(ctx context.Context, id machine.Identifier, class machine.MachineClass, system *system.System, power *machine.PowerCapabilities, callback machine.MachineEndpoint) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("registering system")

	in := &local.RegisterInput{
		InternalIdentifier: conversion.InternalIdentifierFromDomain(id),
		Class:              conversion.ClassFromDomain(class),
		SystemSummary:      conversion.SystemFromDomain(system),
		PowerCapabilities:  conversion.PowerCapabilitiesFromDomain(power),
		CallbackEndpoint:   callback.String(),
	}

	log.Debug().Msg("sending register request")
	_, err := c.client.Register(ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to send register request")
		return err
	}

	log.Debug().Msg("register request sent")
	return nil
}
