package board

import (
	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
	"github.com/awlsring/camp/internal/pkg/values"
	"github.com/jaypipes/ghw"
)

type Service struct {
	bios *motherboard.Bios
	mobo *motherboard.Motherboard
}

func InitService() (service.Motherboard, error) {
	b, err := ghw.BIOS()
	if err != nil {
		return nil, err
	}

	bo, err := ghw.Baseboard()
	if err != nil {
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
