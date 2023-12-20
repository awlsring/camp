package network

import (
	"context"
	"strings"

	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/sys"
	"github.com/awlsring/camp/internal/pkg/values"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/v3/net"
)

func inIgnoreList(ignored []string, name string) bool {
	for _, ignore := range ignored {
		if strings.Contains(name, ignore) {
			return true
		}
	}

	return false
}

func makeMacAddress(mac string) *network.MacAddress {
	m, err := network.MacAddressFromString(mac)
	if err != nil {
		return nil
	}
	return &m
}

func (s *Service) load(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("loading nics")
	if sys.IsMacOS() {
		log.Warn().Msg("network service not implemented for mac")
		return nil
	}

	netw, err := ghw.Network()
	if err != nil {
		log.Error().Err(err).Msg("error loading network")
		return err
	}

	log.Debug().Msg("creating network models")
	nics := map[string]*network.Nic{}
	for _, nic := range netw.NICs {
		if inIgnoreList(s.ignoredPrefix, nic.Name) {
			continue
		}
		n := &network.Nic{
			Name:       nic.Name,
			MacAddress: makeMacAddress(nic.MacAddress),
			Virtual:    nic.IsVirtual,
			Speed:      values.ParseOptional(nic.Speed),
			Duplex:     values.ParseOptional(nic.Duplex),
			PCIAddress: nic.PCIAddress,
		}
		nics[nic.Name] = n
		log.Debug().Interface("nic", n).Msg("created network model")
	}
	s.nics = nics

	addresses := []*network.IpAddress{}

	ifaces, err := net.InterfacesWithContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error loading network interfaces")
		return err
	}
	for _, iface := range ifaces {
		nic, ok := s.nics[iface.Name]
		if !ok {
			continue
		}

		if nic.MacAddress == nil {
			nic.MacAddress = makeMacAddress(iface.HardwareAddr)
		}

		for _, addr := range iface.Addrs {
			address, err := network.AddressFromString(addr.Addr)
			if err != nil {
				log.Error().Err(err).Msg("error parsing ip address")
				continue
			}
			version := network.DetermineIpAddressType(addr.Addr)
			ipa := network.NewIpAddress(version, address, nic)
			addresses = append(addresses, &ipa)
			nic.IpAddresses = append(nic.IpAddresses, &ipa)
		}
	}
	s.addresses = addresses

	log.Debug().Msg("network models created")
	return nil
}
