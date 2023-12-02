package handler

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) RemoveTagFromMachine(ctx context.Context, req camplocal.RemoveTagFromMachineParams) (camplocal.RemoveTagFromMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Removing tag from machine: %s", req.Identifier)

	log.Debug().Msg("Parsing machine identifier")
	id, err := machine.IdentifierFromString(req.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.Identifier)
		return nil, errors.New(errors.ErrValidation, err)
	}

	log.Debug().Msg("Parsing tag key")
	key, err := tag.TagKeyFromString(req.TagKey)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse tag key")
		return nil, errors.New(errors.ErrValidation, err)
	}

	log.Debug().Msgf("Removing tag %s from machine %s", key, id)
	err = h.mSvc.RemoveTagFromMachine(ctx, id, key)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to add tags to machine %s", req.Identifier)
		return nil, err
	}

	log.Debug().Msgf("Successfully removed tag %s from machine %s", key, id)
	return &camplocal.RemoveTagFromMachineResponseContent{
		Sucess: true,
	}, nil
}
