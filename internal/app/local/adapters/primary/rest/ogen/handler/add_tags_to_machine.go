package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
)

func (h *Handler) AddTagsToMachine(ctx context.Context, req *camplocal.AddTagsToMachineRequestContent, params camplocal.AddTagsToMachineParams) (camplocal.AddTagsToMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Adding tags to machine: %s", params.Identifier)

	id, err := machine.IdentifierFromString(params.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", params.Identifier)
		return nil, err
	}

	tags, err := tagsToDomain(req.Tags)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to parse tags")
		return nil, err
	}

	err = h.mSvc.AddTagsToMachine(ctx, id, tags)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to add tags to machine %s", params.Identifier)
		return nil, err
	}

	return &camplocal.AddTagsToMachineResponseContent{
		Sucess: true,
	}, nil
}
