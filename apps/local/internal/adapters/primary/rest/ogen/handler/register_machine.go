package handler

import (
	"context"
	"errors"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	camperror "github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) registerMachineErrorHandler(err error) (camplocal.RegisterRes, error) {
	var campErr *camperror.Error
	if errors.As(err, &campErr) {
		e := campErr.CampError()
		switch e {
		case camperror.ErrValidation:
			return &camplocal.ValidationExceptionResponseContent{
				Message: err.Error(),
			}, nil
		}
	}
	return nil, err
}

func (h *Handler) Register(ctx context.Context, req *camplocal.RegisterRequestContent) (camplocal.RegisterRes, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msg("Invoke Register")
	log.Debug().Msgf("Summary: %+v", req.Summary)

	id, err := machine.IdentifierFromString(req.Summary.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.Summary.InternalIdentifier)
		return h.registerMachineErrorHandler(err)
	}

	class, err := machine.MachineClassFromString(string(req.Summary.GetClass().Value))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse class %s", req.Summary.Class.Value)
		return h.registerMachineErrorHandler(err)
	}

	sys := systemSummaryToDomain(req.Summary.System)
	cpu := cpuSummaryToDomain(req.Summary.CPU)
	mem := memorySummaryToDomain(req.Summary.Memory)
	disk, err := diskSummariesToDomain(req.Summary.Disks)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse disk summaries")
		return h.registerMachineErrorHandler(err)
	}

	nic, err := networkInterfaceSummariesToDomain(req.Summary.NetworkInterfaces)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse network interface summaries")
		return h.registerMachineErrorHandler(err)
	}

	vol, err := volumeSummariesToDomain(req.Summary.Volumes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse volume summaries")
		return h.registerMachineErrorHandler(err)
	}

	ips, err := addressSummariesToDomain(req.Summary.Addresses)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse ip address summaries")
		return h.registerMachineErrorHandler(err)
	}

	err = h.mSvc.RegisterMachine(ctx, id, class, sys, cpu, mem, disk, nic, vol, ips)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register machine")
		return h.registerMachineErrorHandler(err)
	}
	return &camplocal.RegisterResponseContent{
		Success: true,
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
