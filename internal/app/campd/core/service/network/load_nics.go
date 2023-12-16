package network

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/values"
	"github.com/jaypipes/ghw"
)

func (s *Service) loadNics(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("loading nics")

	net, err := ghw.Network()
	if err != nil {
		log.Error().Err(err).Msg("error loading network")
		return err
	}

	log.Debug().Msg("creating network models")
	nics := map[string]*network.Nic{}
	for _, nic := range net.NICs {
		var mac *network.MacAddress
		m, err := network.MacAddressFromString(nic.MacAddress)
		if err != nil {
			mac = nil
		}
		mac = &m

		nics[nic.Name] = &network.Nic{
			Name:       nic.Name,
			MacAddress: mac,
			Virtual:    nic.IsVirtual,
			Speed:      values.ParseOptional(nic.Speed),
			Duplex:     values.ParseOptional(nic.Duplex),
			PCIAddress: nic.PCIAddress,
		}
	}
	s.nics = nics

	log.Debug().Msg("network models created")
	return nil
}
