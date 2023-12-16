package handler

import (
	mach "github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
)

func modelToSummary(m *mach.Machine) camplocal.MachineSummary {
	cap := powerCapabilitiesToSummary(m.PowerCapabilities)
	// system := systemModelToSummary(m.System)
	host := hostModelToSummary(m.Host)
	cpu := cpuModelToSummary(m.Cpu)
	memory := memoryModelToSummary(m.Memory)
	disks := diskModelsToSummaryList(m.Disks)
	networkInterfaces := networkInterfaceModelsToSummaryList(m.NetworkInterfaces)
	// volumes := volumeModelsToSummaryList(m.Volumes)
	addresses := addressModelsToSummaryList(m.Addresses)
	tags := domainToTags(m.Tags)

	return camplocal.MachineSummary{
		Identifier: m.Identifier.String(),
		Class:      camplocal.NewOptMachineClass(camplocal.MachineClass(m.Class.String())),
		Status: camplocal.StatusSummary{
			Status:      camplocal.StatusCode(m.Status.Status.String()),
			LastUpdated: float64(m.Status.UpdatedAt.Unix()),
		},
		PowerCapabilities: cap,
		Tags:              tags,
		Host:              host,
		// System:            system,
		CPU:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		// Volumes:           volumes,
		Addresses: addresses,
	}
}

func powerCapabilitiesToSummary(cap mach.PowerCapabilities) camplocal.MachinePowerCapabilitiesSummary {
	reboot := camplocal.MachinePowerCapabilityRebootSummary{
		Enabled: cap.Reboot.Enabled,
	}

	powerOff := camplocal.MachinePowerCapabilityPowerOffSummary{
		Enabled: cap.PowerOff.Enabled,
	}

	wakeOnLan := camplocal.MachinePowerCapabilityWakeOnLanSummary{
		Enabled: cap.WakeOnLan.Enabled,
	}
	if cap.WakeOnLan.MacAddress != nil {
		wakeOnLan.MacAddress = camplocal.NewOptString(cap.WakeOnLan.MacAddress.String())
	}

	return camplocal.MachinePowerCapabilitiesSummary{
		Reboot:    reboot,
		PowerOff:  powerOff,
		WakeOnLan: wakeOnLan,
	}
}

func hostModelToSummary(m *host.Host) camplocal.HostSummary {
	sum := camplocal.HostSummary{}

	if m != nil {
		if m.Hostname != nil {
			sum.Hostname = camplocal.NewOptString(*m.Hostname)
		}
		if m.HostId != nil {
			sum.HostId = camplocal.NewOptString(*m.HostId)
		}
		if m.OS != nil {
			os := camplocal.OsSummary{}
			if m.OS.Name != nil {
				os.Name = camplocal.NewOptString(*m.OS.Name)
			}
			if m.OS.Platform != nil {
				os.Platform = camplocal.NewOptString(*m.OS.Platform)
			}
			if m.OS.Version != nil {
				os.Version = camplocal.NewOptString(*m.OS.Version)
			}
			if m.OS.Kernel != nil {
				os.Kernel = camplocal.NewOptString(*m.OS.Kernel)
			}
			if m.OS.Family != nil {
				os.Family = camplocal.NewOptString(*m.OS.Family)
			}
			sum.Os = camplocal.NewOptOsSummary(os)
		}
	}

	return sum
}

func cpuModelToSummary(m *cpu.CPU) camplocal.CpuSummary {
	sum := camplocal.CpuSummary{
		TotalCores:   float64(m.TotalCores),
		TotalThreads: float64(m.TotalThreads),
		Architecture: camplocal.Architecture(m.Architecture.String()),
	}
	if m.Model != nil {
		sum.Model = camplocal.NewOptString(*m.Model)
	}
	if m.Vendor != nil {
		sum.Vendor = camplocal.NewOptString(*m.Vendor)
	}

	processors := []camplocal.ProcessorSummary{}
	for _, p := range m.Processors {
		processor := processorModelToSummary(p)
		processors = append(processors, processor)
	}
	sum.Processors = processors

	return sum
}

func processorModelToSummary(p *cpu.Processor) camplocal.ProcessorSummary {
	pro := camplocal.ProcessorSummary{
		Identifier:  float64(p.Id),
		CoreCount:   float64(p.CoreCount),
		ThreadCount: float64(p.ThreadCount),
	}

	if p.Model != nil {
		pro.Model = camplocal.NewOptString(*p.Model)
	}
	if p.Vendor != nil {
		pro.Vendor = camplocal.NewOptString(*p.Vendor)
	}

	cores := []camplocal.CoreSummary{}
	for _, c := range p.Cores {
		core := coreModelToSummary(c)
		cores = append(cores, core)
	}

	pro.Cores = cores
	return pro
}

func coreModelToSummary(c *cpu.Core) camplocal.CoreSummary {
	return camplocal.CoreSummary{
		Identifier:  float64(c.Id),
		ThreadCount: float64(c.Threads),
	}
}

func memoryModelToSummary(m *memory.Memory) camplocal.MemorySummary {
	return camplocal.MemorySummary{
		Total: float64(m.Total),
	}
}

func diskModelToSummary(m *storage.Disk) camplocal.DiskSummary {
	sum := camplocal.DiskSummary{
		Name:              m.Name,
		StorageController: camplocal.StorageController(m.StorageController.String()),
		Type:              camplocal.DiskType(m.DriveType.String()),
		Size:              float64(m.Size),
		Removable:         m.Removable,
	}
	if m.Serial != nil {
		sum.Serial = camplocal.NewOptString(*m.Serial)
	}
	if m.Model != nil {
		sum.Model = camplocal.NewOptString(*m.Model)
	}
	if m.Vendor != nil {
		sum.Vendor = camplocal.NewOptString(*m.Vendor)
	}
	return sum
}

func diskModelsToSummaryList(ms []*storage.Disk) []camplocal.DiskSummary {
	var out []camplocal.DiskSummary
	for _, m := range ms {
		out = append(out, diskModelToSummary(m))
	}
	return out
}

func networkInterfaceModelToSummary(m *network.Nic) camplocal.NetworkInterfaceSummary {
	sum := camplocal.NetworkInterfaceSummary{
		Name:    m.Name,
		Virtual: m.Virtual,
	}
	if m.MacAddress != nil {
		sum.MacAddress = camplocal.NewOptString(m.MacAddress.String())
	}
	if m.Vendor != nil {
		sum.Vendor = camplocal.NewOptString(*m.Vendor)
	}
	// if m.Speed != nil {
	// 	sum.Speed = camplocal.NewOptFloat64(float64(*m.Speed))
	// }
	if m.Duplex != nil {
		sum.Duplex = camplocal.NewOptString(*m.Duplex)
	}
	return sum
}

func networkInterfaceModelsToSummaryList(ms []*network.Nic) []camplocal.NetworkInterfaceSummary {
	var out []camplocal.NetworkInterfaceSummary
	for _, m := range ms {
		out = append(out, networkInterfaceModelToSummary(m))
	}
	return out
}

// func volumeModelToSummary(m *machine.Volume) camplocal.MachineVolumeSummary {
// 	vol := camplocal.MachineVolumeSummary{
// 		Name:       m.Name.String(),
// 		MountPoint: m.MountPoint.String(),
// 		Total:      float64(m.Total),
// 	}
// 	if m.FileSystem != nil {
// 		vol.FileSystem = camplocal.NewOptString(*m.FileSystem)
// 	}
// 	return vol
// }

// func volumeModelsToSummaryList(ms []*machine.Volume) []camplocal.MachineVolumeSummary {
// 	var out []camplocal.MachineVolumeSummary
// 	for _, m := range ms {
// 		out = append(out, volumeModelToSummary(m))
// 	}
// 	return out
// }

func addressModelToSummary(m *network.IpAddress) camplocal.IpAddressSummary {
	return camplocal.IpAddressSummary{
		Version: camplocal.IpAddressVersion(m.Version.String()),
		Address: m.Address.String(),
	}
}

func addressModelsToSummaryList(ms []*network.IpAddress) []camplocal.IpAddressSummary {
	var out []camplocal.IpAddressSummary
	for _, m := range ms {
		out = append(out, addressModelToSummary(m))
	}
	return out
}
