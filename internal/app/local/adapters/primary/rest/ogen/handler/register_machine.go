package handler

import (
	"context"
	"fmt"

	mach "github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
)

func (h *Handler) Register(ctx context.Context, req *camplocal.RegisterRequestContent) (camplocal.RegisterRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Registering new machine")

	id, err := mach.IdentifierFromString(req.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.InternalIdentifier)
		return nil, err
	}

	endpoint, err := mach.MachineEndpointFromString(req.CallbackEndpoint)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse endpoint %s", req.CallbackEndpoint)
		return nil, err
	}

	key, err := mach.AgentKeyFromString("a")
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse key")
		return nil, err
	}

	class, err := machine.MachineClassFromString(string(req.Class))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse class %s", string(req.Class))
		return nil, err
	}

	host := hostSummaryToDomain(req.SystemSummary.Host)
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

	// vol, err := volumeSummariesToDomain(req.SystemSummary.Volumes)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to parse volume summaries")
	// 	return nil, err
	// }

	ips, err := addressSummariesToDomain(req.SystemSummary.Addresses)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse ip address summaries")
		return nil, err
	}

	err = h.mSvc.RegisterMachine(ctx, id, endpoint, key, class, cap, host, cpu, mem, disk, nic, nil, ips)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register machine")
		return nil, err
	}
	return &camplocal.RegisterResponseContent{
		AccessKey: key.String(),
	}, nil
}

func capabilitySummaryToDomain(sum camplocal.OptReportedPowerCapabilitiesSummary) (mach.PowerCapabilities, error) {
	defaultCaps := mach.PowerCapabilities{
		Reboot:    mach.PowerCapabilityReboot{Enabled: false},
		PowerOff:  mach.PowerCapabilityPowerOff{Enabled: false},
		WakeOnLan: mach.PowerCapabilityWakeOnLan{Enabled: false},
	}
	if !sum.IsSet() {
		return defaultCaps, nil
	}
	var reboot mach.PowerCapabilityReboot
	if sum.Value.Reboot.IsSet() {
		reboot = mach.PowerCapabilityReboot{
			Enabled: sum.Value.Reboot.Value.Enabled,
		}
	}

	var powerOff mach.PowerCapabilityPowerOff
	if sum.Value.PowerOff.IsSet() {
		powerOff = mach.PowerCapabilityPowerOff{
			Enabled: sum.Value.PowerOff.Value.Enabled,
		}
	}

	var wakeOnLan mach.PowerCapabilityWakeOnLan
	if sum.Value.WakeOnLan.IsSet() {
		wakeOnLan = mach.PowerCapabilityWakeOnLan{
			Enabled: sum.Value.WakeOnLan.Value.Enabled,
		}
		if sum.Value.WakeOnLan.Value.Enabled && !sum.Value.WakeOnLan.Value.MacAddress.IsSet() {
			return defaultCaps, fmt.Errorf("enabling wake on lan requires a mac address")
		}
		if sum.Value.WakeOnLan.Value.MacAddress.IsSet() {
			mac, err := network.MacAddressFromString(sum.Value.WakeOnLan.Value.MacAddress.Value)
			if err != nil {
				return defaultCaps, err
			}
			wakeOnLan.MacAddress = &mac
		}
	}

	return mach.PowerCapabilities{
		Reboot:    reboot,
		PowerOff:  powerOff,
		WakeOnLan: wakeOnLan,
	}, nil
}

func hostSummaryToDomain(sum camplocal.HostSummary) *host.Host {
	s := &host.Host{}

	if sum.Hostname.IsSet() {
		s.Hostname = &sum.Hostname.Value
	}

	if sum.Os.IsSet() {
		os := &host.OS{}
		if sum.Os.Value.Name.IsSet() {
			os.Name = &sum.Os.Value.Name.Value
		}
		if sum.Os.Value.Platform.IsSet() {
			os.Platform = &sum.Os.Value.Platform.Value
		}
		if sum.Os.Value.Version.IsSet() {
			os.Version = &sum.Os.Value.Version.Value
		}
		if sum.Os.Value.Kernel.IsSet() {
			os.Kernel = &sum.Os.Value.Kernel.Value
		}
		if sum.Os.Value.Family.IsSet() {
			os.Family = &sum.Os.Value.Family.Value
		}
		s.OS = os
	}

	return s
}

func cpuSummaryToDomain(sum camplocal.CpuSummary) *cpu.CPU {
	c := &cpu.CPU{
		TotalCores:   uint32(sum.TotalCores),
		TotalThreads: uint32(sum.TotalThreads),
		Architecture: cpu.ArchitectureFromString(string(sum.Architecture)),
	}
	if sum.Model.IsSet() {
		c.Model = &sum.Model.Value
	}
	if sum.Vendor.IsSet() {
		c.Vendor = &sum.Vendor.Value
	}

	processors := []*cpu.Processor{}
	for _, p := range sum.Processors {
		processor := processorSummaryToDomain(p)
		processors = append(processors, processor)
	}
	c.Processors = processors

	return c
}

func processorSummaryToDomain(sum camplocal.ProcessorSummary) *cpu.Processor {
	processor := &cpu.Processor{
		Id:          int(sum.Identifier),
		CoreCount:   uint32(sum.CoreCount),
		ThreadCount: uint32(sum.ThreadCount),
	}
	if sum.Model.IsSet() {
		processor.Model = &sum.Model.Value
	}
	if sum.Vendor.IsSet() {
		processor.Vendor = &sum.Vendor.Value
	}

	cores := []*cpu.Core{}
	for _, c := range sum.Cores {
		core := coreSummaryToDomain(c)
		cores = append(cores, core)
	}
	processor.Cores = cores

	return processor
}

func coreSummaryToDomain(sum camplocal.CoreSummary) *cpu.Core {
	return &cpu.Core{
		Id:      int(sum.Identifier),
		Threads: uint32(sum.ThreadCount),
	}
}

func memorySummaryToDomain(sum camplocal.MemorySummary) *memory.Memory {
	return &memory.Memory{
		Total: uint64(sum.Total),
	}
}

func diskSummaryToDomain(sum camplocal.DiskSummary) (*storage.Disk, error) {
	storageCtl := storage.StorageControllerFromString(string(sum.StorageController))
	driveType := storage.DriveTypeFromString(string(sum.Type))
	disk := &storage.Disk{
		Name:              sum.Name,
		StorageController: storageCtl,
		DriveType:         driveType,
		Size:              uint64(sum.Size),
		Removable:         sum.Removable,
	}
	if sum.Model.IsSet() {
		disk.Model = &sum.Model.Value
	}

	if sum.Vendor.IsSet() {
		disk.Vendor = &sum.Vendor.Value
	}

	if sum.Serial.IsSet() {
		disk.Serial = &sum.Serial.Value
	}

	if sum.Wwn.IsSet() {
		disk.WWN = &sum.Wwn.Value
	}

	return disk, nil
}

func diskSummariesToDomain(sums []camplocal.DiskSummary) ([]*storage.Disk, error) {
	var out []*storage.Disk
	for _, sum := range sums {
		disk, err := diskSummaryToDomain(sum)
		if err != nil {
			return nil, err
		}
		out = append(out, disk)
	}
	return out, nil
}

func networkInterfaceSummaryToDomain(sum camplocal.NetworkInterfaceSummary) (*network.Nic, error) {
	addresses, err := addressSummariesToDomain(sum.IpAddresses)
	if err != nil {
		return nil, err
	}

	nic := &network.Nic{
		Name:        sum.Name,
		IpAddresses: addresses,
		Virtual:     sum.Virtual,
	}

	if sum.Vendor.IsSet() {
		nic.Vendor = &sum.Vendor.Value
	}

	if sum.Duplex.IsSet() {
		nic.Duplex = &sum.Duplex.Value
	}

	// if sum.Speed.IsSet() {
	// 	nic.Speed = &sum.Speed.Value
	// }

	if sum.MacAddress.IsSet() {
		mac, err := network.MacAddressFromString(sum.MacAddress.Value)
		if err != nil {
			return nil, err
		}
		nic.MacAddress = &mac
	}

	return nic, nil
}

func networkInterfaceSummariesToDomain(sums []camplocal.NetworkInterfaceSummary) ([]*network.Nic, error) {
	var out []*network.Nic
	for _, sum := range sums {
		nic, err := networkInterfaceSummaryToDomain(sum)
		if err != nil {
			return nil, err
		}
		out = append(out, nic)
	}
	return out, nil
}

// func volumeSummaryToDomain(sum camplocal.V) (*machine.Volume, error) {
// 	vol, err := machine.VolumeIdentifierFromString(sum.Name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mp, err := machine.MountPointFromString(sum.MountPoint)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &machine.Volume{
// 		Name:       vol,
// 		MountPoint: mp,
// 		Total:      int64(sum.Total),
// 		FileSystem: &sum.FileSystem.Value,
// 	}, nil
// }

// func volumeSummariesToDomain(sums []camplocal.MachineVolumeSummary) ([]*machine.Volume, error) {
// 	var out []*machine.Volume
// 	for _, sum := range sums {
// 		vol, err := volumeSummaryToDomain(sum)
// 		if err != nil {
// 			return nil, err
// 		}
// 		out = append(out, vol)
// 	}
// 	return out, nil
// }

func addressSummaryToDomain(sum camplocal.IpAddressSummary) (*network.IpAddress, error) {
	addr, err := network.AddressFromString(sum.Address)
	if err != nil {
		return nil, err
	}
	version := network.IpAddressTypeFromString(string(sum.Version))

	return &network.IpAddress{
		Version: version,
		Address: addr,
	}, nil
}

func addressSummariesToDomain(sums []camplocal.IpAddressSummary) ([]*network.IpAddress, error) {
	var out []*network.IpAddress
	for _, sum := range sums {
		addr, err := addressSummaryToDomain(sum)
		if err != nil {
			return nil, err
		}
		out = append(out, addr)
	}
	return out, nil
}
