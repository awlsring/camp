package handler

import (
	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func modelToSummary(m *machine.Machine) camplocal.MachineSummary {
	cap := powerCapabilitiesToSummary(m.PowerCapabilities)
	system := systemModelToSummary(m.System)
	cpu := cpuModelToSummary(m.Cpu)
	memory := memoryModelToSummary(m.Memory)
	disks := diskModelsToSummaryList(m.Disks)
	networkInterfaces := networkInterfaceModelsToSummaryList(m.NetworkInterfaces)
	volumes := volumeModelsToSummaryList(m.Volumes)
	addresses := addressModelsToSummaryList(m.Addresses)
	tags := domainToTags(m.Tags)

	return camplocal.MachineSummary{
		Identifier: m.Identifier.String(),
		Class:      camplocal.NewOptMachineClass(camplocal.MachineClass(m.Class.String())),
		Status: camplocal.MachineStatusSummary{
			Status:      camplocal.MachineStatus(m.Status.String()),
			LastChecked: float64(m.LastHeartbeat.Unix()),
		},
		PowerCapabilities: cap,
		Tags:              tags,
		System:            system,
		CPU:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		Volumes:           volumes,
		Addresses:         addresses,
	}
}

func powerCapabilitiesToSummary(cap machine.PowerCapabilities) camplocal.MachinePowerCapabilitiesSummary {
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

func systemModelToSummary(m *machine.System) camplocal.MachineSystemSummary {
	sum := camplocal.MachineSystemSummary{}
	if m.Os.Family != nil {
		sum.Family = camplocal.NewOptString(*m.Os.Family)
	}
	if m.Os.Kernel != nil {
		sum.KernelVersion = camplocal.NewOptString(*m.Os.Kernel)
	}
	if m.Os.Name != nil {
		sum.Os = camplocal.NewOptString(*m.Os.Name)
	}
	if m.Os.Version != nil {
		sum.OsVersion = camplocal.NewOptString(*m.Os.Version)
	}
	if m.Os.PrettyName != nil {
		sum.OsPretty = camplocal.NewOptString(*m.Os.PrettyName)
	}
	if m.Hostname != nil {
		sum.Hostname = camplocal.NewOptString(*m.Hostname)
	}
	return sum
}

func cpuModelToSummary(m *machine.Cpu) camplocal.MachineCpuSummary {
	sum := camplocal.MachineCpuSummary{
		Cores:        float64(m.Cores),
		Architecture: camplocal.CpuArchitecture(m.Architecture.String()),
	}
	if m.Model != nil {
		sum.Model = camplocal.NewOptString(*m.Model)
	}
	if m.Vendor != nil {
		sum.Vendor = camplocal.NewOptString(*m.Vendor)
	}
	return sum
}

func memoryModelToSummary(m *machine.Memory) camplocal.MachineMemorySummary {
	return camplocal.MachineMemorySummary{
		Total: float64(m.Total),
	}
}

func diskModelToSummary(m *machine.Disk) camplocal.MachineDiskSummary {
	sum := camplocal.MachineDiskSummary{
		Device:    m.Device.String(),
		Interface: camplocal.DiskInterface(m.Interface.String()),
		Type:      camplocal.DiskType(m.Type.String()),
		Size:      float64(m.Size),
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
	if m.SizeRaw != nil {
		sum.SizeRaw = camplocal.NewOptFloat64(float64(*m.SizeRaw))
	}
	if m.SectorSize != 0 {
		sum.SectorSize = camplocal.NewOptFloat64(float64(m.SectorSize))
	}
	return sum
}

func diskModelsToSummaryList(ms []*machine.Disk) []camplocal.MachineDiskSummary {
	var out []camplocal.MachineDiskSummary
	for _, m := range ms {
		out = append(out, diskModelToSummary(m))
	}
	return out
}

func networkInterfaceModelToSummary(m *machine.NetworkInterface) camplocal.MachineNetworkInterfaceSummary {
	sum := camplocal.MachineNetworkInterfaceSummary{
		Name:    m.Name.String(),
		Virtual: m.Virtual,
	}
	if m.MacAddress != nil {
		sum.MacAddress = camplocal.NewOptString(m.MacAddress.String())
	}
	if m.Vendor != nil {
		sum.Vendor = camplocal.NewOptString(*m.Vendor)
	}
	if m.Mtu != nil {
		sum.Mtu = camplocal.NewOptFloat64(float64(*m.Mtu))
	}
	if m.Speed != nil {
		sum.Speed = camplocal.NewOptFloat64(float64(*m.Speed))
	}
	if m.Duplex != nil {
		sum.Duplex = camplocal.NewOptString(*m.Duplex)
	}
	return sum
}

func networkInterfaceModelsToSummaryList(ms []*machine.NetworkInterface) []camplocal.MachineNetworkInterfaceSummary {
	var out []camplocal.MachineNetworkInterfaceSummary
	for _, m := range ms {
		out = append(out, networkInterfaceModelToSummary(m))
	}
	return out
}

func volumeModelToSummary(m *machine.Volume) camplocal.MachineVolumeSummary {
	vol := camplocal.MachineVolumeSummary{
		Name:       m.Name.String(),
		MountPoint: m.MountPoint.String(),
		Total:      float64(m.Total),
	}
	if m.FileSystem != nil {
		vol.FileSystem = camplocal.NewOptString(*m.FileSystem)
	}
	return vol
}

func volumeModelsToSummaryList(ms []*machine.Volume) []camplocal.MachineVolumeSummary {
	var out []camplocal.MachineVolumeSummary
	for _, m := range ms {
		out = append(out, volumeModelToSummary(m))
	}
	return out
}

func addressModelToSummary(m *machine.IpAddress) camplocal.IpAddressSummary {
	return camplocal.IpAddressSummary{
		Version: camplocal.IpAddressVersion(m.Version.String()),
		Address: m.Address.String(),
	}
}

func addressModelsToSummaryList(ms []*machine.IpAddress) []camplocal.IpAddressSummary {
	var out []camplocal.IpAddressSummary
	for _, m := range ms {
		out = append(out, addressModelToSummary(m))
	}
	return out
}
