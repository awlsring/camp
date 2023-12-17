package board

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/sys"
	"github.com/awlsring/camp/internal/pkg/values"
	"github.com/jaypipes/ghw"
)

type Service struct {
	bios *motherboard.Bios
	mobo *motherboard.Motherboard
}

func InitService(ctx context.Context) (service.Motherboard, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Initializing motherboard service")

	if sys.IsMacOS() {
		log.Warn().Msg("Motherboard service not implemented for apple")
		return &Service{}, nil
	}

	b, err := ghw.BIOS()
	if err != nil {
		log.Error().Err(err).Msg("failed to get BIOS info")
		return nil, err
	}

	bo, err := ghw.Baseboard()
	if err != nil {
		log.Error().Err(err).Msg("failed to get baseboard info")
		return nil, err
	}

	return &Service{
		bios: &motherboard.Bios{
			Vendor:  values.ParseOptional(b.Vendor),
			Version: values.ParseOptional(b.Version),
			Date:    values.ParseOptional(b.Date),
		},
		mobo: &motherboard.Motherboard{
			AssetTag: values.ParseOptional(bo.AssetTag),
			Product:  values.ParseOptional(bo.Product),
			Serial:   values.ParseOptional(bo.SerialNumber),
			Vendor:   values.ParseOptional(bo.Vendor),
			Version:  values.ParseOptional(bo.Version),
		},
	}, nil
}
