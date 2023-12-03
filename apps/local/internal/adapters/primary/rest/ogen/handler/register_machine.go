package handler

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) Register(ctx context.Context, req *camplocal.RegisterRequestContent) (camplocal.RegisterRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Registering new machine")

	id, err := machine.IdentifierFromString(req.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.InternalIdentifier)
		return nil, err
	}

	endpoint, err := machine.MachineEndpointFromString(req.CallbackEndpoint)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse endpoint %s", req.CallbackEndpoint)
		return nil, err
	}

	key, err := machine.AgentKeyFromString(req.CallbackKey)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse key %s", req.CallbackKey)
		return nil, err
	}

	class, err := machine.MachineClassFromString(string(req.Class.Value))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse class %s", req.Class.Value)
		return nil, err
	}

	sys := systemSummaryToDomain(req.SystemSummary.System)
	cpu := cpuSummaryToDomain(req.SystemSummary.CPU)
	mem := memorySummaryToDomain(req.SystemSummary.Memory)
	disk, err := diskSummariesToDomain(req.SystemSummary.Disks)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse disk summaries")
		return nil, err
	}

	cap, err := capabilitySummaryToDomain(req.PowerCapabilities)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse power capabilities")
		return nil, err
	}

	nic, err := networkInterfaceSummariesToDomain(req.SystemSummary.NetworkInterfaces)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse network interface summaries")
		return nil, err
	}

	vol, err := volumeSummariesToDomain(req.SystemSummary.Volumes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse volume summaries")
		return nil, err
	}

	ips, err := addressSummariesToDomain(req.SystemSummary.Addresses)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse ip address summaries")
		return nil, err
	}

	err = h.mSvc.RegisterMachine(ctx, id, endpoint, key, class, cap, sys, cpu, mem, disk, nic, vol, ips)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register machine")
		return nil, err
	}
	return &camplocal.RegisterResponseContent{
		Success: true,
	}, nil
}

func capabilitySummaryToDomain(sum camplocal.OptReportedPowerCapabilitiesSummary) (machine.PowerCapabilities, error) {
	defaultCaps := machine.PowerCapabilities{
		Reboot:    machine.PowerCapabilityReboot{Enabled: false},
		PowerOff:  machine.PowerCapabilityPowerOff{Enabled: false},
		WakeOnLan: machine.PowerCapabilityWakeOnLan{Enabled: false},
	}
	if !sum.IsSet() {
		return defaultCaps, nil
	}
	var reboot machine.PowerCapabilityReboot
	if sum.Value.Reboot.IsSet() {
		reboot = machine.PowerCapabilityReboot{
			Enabled: sum.Value.Reboot.Value.Enabled,
		}
	}

	var powerOff machine.PowerCapabilityPowerOff
	if sum.Value.PowerOff.IsSet() {
		powerOff = machine.PowerCapabilityPowerOff{
			Enabled: sum.Value.PowerOff.Value.Enabled,
		}
	}

	var wakeOnLan machine.PowerCapabilityWakeOnLan
	if sum.Value.WakeOnLan.IsSet() {
		wakeOnLan = machine.PowerCapabilityWakeOnLan{
			Enabled: sum.Value.WakeOnLan.Value.Enabled,
		}
		if sum.Value.WakeOnLan.Value.Enabled && !sum.Value.WakeOnLan.Value.MacAddress.IsSet() {
			return defaultCaps, fmt.Errorf("enabling wake on lan requires a mac address")
		}
		if sum.Value.WakeOnLan.Value.MacAddress.IsSet() {
			mac, err := machine.MacAddressFromString(sum.Value.WakeOnLan.Value.MacAddress.Value)
			if err != nil {
				return defaultCaps, err
			}
			wakeOnLan.MacAddress = &mac
		}
	}

	return machine.PowerCapabilities{
		Reboot:    reboot,
		PowerOff:  powerOff,
		WakeOnLan: wakeOnLan,
	}, nil
}

func systemSummaryToDomain(sum camplocal.MachineSystemSummary) *machine.System {
	return &machine.System{
		Os: machine.Os{
			Family:     &sum.Family.Value,
			Kernel:     &sum.KernelVersion.Value,
			Name:       &sum.Os.Value,
			Version:    &sum.OsVersion.Value,
			PrettyName: &sum.OsPretty.Value,
		},
		Hostname: &sum.Hostname.Value,
	}
}

func cpuSummaryToDomain(sum camplocal.MachineCpuSummary) *machine.Cpu {
	return &machine.Cpu{
		Cores:        int(sum.Cores),
		Architecture: machine.CpuArchitectureFromString(string(sum.Architecture)),
		Model:        &sum.Model.Value,
		Vendor:       &sum.Vendor.Value,
	}
}

func memorySummaryToDomain(sum camplocal.MachineMemorySummary) *machine.Memory {
	return &machine.Memory{
		Total: int64(sum.Total),
	}
}

func diskSummaryToDomain(sum camplocal.MachineDiskSummary) (*machine.Disk, error) {
	dev, err := machine.DiskIdentifierFromString(sum.Device)
	if err != nil {
		return nil, err
	}
	sizeRaw := int64(sum.SizeRaw.Value)
	return &machine.Disk{
		Device:     dev,
		Model:      &sum.Model.Value,
		Vendor:     &sum.Vendor.Value,
		Interface:  machine.DiskInterfaceFromString(string(sum.Interface)),
		Type:       machine.DiskClassFromString(string(sum.Type)),
		Serial:     &sum.Serial.Value,
		SectorSize: int(sum.SectorSize.Value),
		Size:       int64(sum.Size),
		SizeRaw:    &sizeRaw,
	}, nil
}

func diskSummariesToDomain(sums []camplocal.MachineDiskSummary) ([]*machine.Disk, error) {
	var out []*machine.Disk
	for _, sum := range sums {
		disk, err := diskSummaryToDomain(sum)
		if err != nil {
			return nil, err
		}
		out = append(out, disk)
	}
	return out, nil
}

func networkInterfaceSummaryToDomain(sum camplocal.MachineNetworkInterfaceSummary) (*machine.NetworkInterface, error) {
	id, err := machine.NetworkInterfaceIdentifierFromString(sum.Name)
	if err != nil {
		return nil, err
	}

	addresses, err := addressSummariesToDomain(sum.Addresses)
	if err != nil {
		return nil, err
	}

	mtu := int(sum.Mtu.Value)
	speed := int(sum.Speed.Value)

	nic := &machine.NetworkInterface{
		Name:        id,
		IpAddresses: addresses,
		Virtual:     sum.Virtual,
		Vendor:      &sum.Vendor.Value,
		Mtu:         &mtu,
		Speed:       &speed,
		Duplex:      &sum.Duplex.Value,
	}

	if sum.MacAddress.IsSet() {
		mac, err := machine.MacAddressFromString(sum.MacAddress.Value)
		if err != nil {
			return nil, err
		}
		nic.MacAddress = &mac
	}

	return nic, nil
}

func networkInterfaceSummariesToDomain(sums []camplocal.MachineNetworkInterfaceSummary) ([]*machine.NetworkInterface, error) {
	var out []*machine.NetworkInterface
	for _, sum := range sums {
		nic, err := networkInterfaceSummaryToDomain(sum)
		if err != nil {
			return nil, err
		}
		out = append(out, nic)
	}
	return out, nil
}

func volumeSummaryToDomain(sum camplocal.MachineVolumeSummary) (*machine.Volume, error) {
	vol, err := machine.VolumeIdentifierFromString(sum.Name)
	if err != nil {
		return nil, err
	}

	mp, err := machine.MountPointFromString(sum.MountPoint)
	if err != nil {
		return nil, err
	}

	return &machine.Volume{
		Name:       vol,
		MountPoint: mp,
		Total:      int64(sum.Total),
		FileSystem: &sum.FileSystem.Value,
	}, nil
}

func volumeSummariesToDomain(sums []camplocal.MachineVolumeSummary) ([]*machine.Volume, error) {
	var out []*machine.Volume
	for _, sum := range sums {
		vol, err := volumeSummaryToDomain(sum)
		if err != nil {
			return nil, err
		}
		out = append(out, vol)
	}
	return out, nil
}

func addressSummaryToDomain(sum camplocal.IpAddressSummary) (*machine.IpAddress, error) {
	addr, err := machine.AddressFromString(sum.Address)
	if err != nil {
		return nil, err
	}
	version := machine.IpAddressTypeFromString(string(sum.Version))

	return &machine.IpAddress{
		Version: version,
		Address: addr,
	}, nil
}

func addressSummariesToDomain(sums []camplocal.IpAddressSummary) ([]*machine.IpAddress, error) {
	var out []*machine.IpAddress
	for _, sum := range sums {
		addr, err := addressSummaryToDomain(sum)
		if err != nil {
			return nil, err
		}
		out = append(out, addr)
	}
	return out, nil
}
