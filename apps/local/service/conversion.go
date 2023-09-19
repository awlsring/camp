package service

import (
	"time"

	"github.com/awlsring/camp/apps/local/machine"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func reportedMachineSummaryToModel(sum *camplocal.ReportedMachineSummary) *machine.Model {
	return &machine.Model{
		Identifier:        sum.InternalIdentifier,
		Class:             machine.MachineClass(sum.GetClass().Value),
		Status:            machine.Running,
		LastHeartbeat:     time.Now(),
		System:            machineSystemSummaryToModel(sum.System),
		Cpu:               machineCpuSummaryToModel(sum.CPU),
		Memory:            machineMemorySummaryToModel(sum.Memory),
		Disks:             machineDiskSummariesToModels(sum.Disks),
		NetworkInterfaces: machineNetworkInterfaceSummariesToModels(sum.NetworkInterfaces),
		Volumes:           machineVolumeSummariesToModels(sum.Volumes),
		Addresses:         machineAddressSummariesToModels(sum.Addresses),
	}
}

func machineSystemSummaryToModel(sum camplocal.MachineSystemSummary) *machine.SystemModel {
	return &machine.SystemModel{
		Family:        &sum.Family.Value,
		KernelVersion: &sum.KernelVersion.Value,
		Os:            &sum.Os.Value,
		OsVersion:     &sum.OsVersion.Value,
		OsPretty:      &sum.OsPretty.Value,
		Hostname:      &sum.Hostname.Value,
	}
}

func machineCpuSummaryToModel(sum camplocal.MachineCpuSummary) *machine.CpuModel {
	return &machine.CpuModel{
		Cores:        int(sum.Cores),
		Architecture: machine.CpuArchitecture(sum.Architecture),
		Model:        &sum.Model.Value,
		Vendor:       &sum.Vendor.Value,
	}
}

func machineMemorySummaryToModel(sum camplocal.MachineMemorySummary) *machine.MemoryModel {
	return &machine.MemoryModel{
		Total: int64(sum.Total),
	}
}

func machineDiskSummaryToModel(sum camplocal.MachineDiskSummary) *machine.DiskModel {
	sizeRaw := int64(sum.SizeRaw.Value)
	return &machine.DiskModel{
		Device:     sum.Device,
		Model:      &sum.Model.Value,
		Vendor:     &sum.Vendor.Value,
		Interface:  machine.DiskInterface(sum.Interface),
		Type:       machine.DiskType(sum.Type),
		Serial:     &sum.Serial.Value,
		SectorSize: int(sum.SectorSize.Value),
		Size:       int64(sum.Size),
		SizeRaw:    &sizeRaw,
	}
}

func machineDiskSummariesToModels(sums []camplocal.MachineDiskSummary) []*machine.DiskModel {
	var out []*machine.DiskModel
	for _, sum := range sums {
		out = append(out, machineDiskSummaryToModel(sum))
	}
	return out
}

func machineNetworkInterfaceSummaryToModel(sum camplocal.MachineNetworkInterfaceSummary) *machine.NetworkInterfaceModel {
	addresses := machineAddressSummariesToModels(sum.Addresses)

	mtu := int(sum.Mtu.Value)
	speed := int(sum.Speed.Value)
	return &machine.NetworkInterfaceModel{
		Name:        sum.Name,
		IpAddresses: addresses,
		Virtual:     sum.Virtual,
		MacAddress:  &sum.MacAddress.Value,
		Vendor:      &sum.Vendor.Value,
		Mtu:         &mtu,
		Speed:       &speed,
		Duplex:      &sum.Duplex.Value,
	}
}

func machineNetworkInterfaceSummariesToModels(sums []camplocal.MachineNetworkInterfaceSummary) []*machine.NetworkInterfaceModel {
	var out []*machine.NetworkInterfaceModel
	for _, sum := range sums {
		out = append(out, machineNetworkInterfaceSummaryToModel(sum))
	}
	return out
}

func machineVolumeSummaryToModel(sum camplocal.MachineVolumeSummary) *machine.VolumeModel {
	return &machine.VolumeModel{
		Name:       sum.Name,
		MountPoint: sum.MountPoint,
		Total:      int64(sum.Total),
		FileSystem: &sum.FileSystem.Value,
	}
}

func machineVolumeSummariesToModels(sums []camplocal.MachineVolumeSummary) []*machine.VolumeModel {
	var out []*machine.VolumeModel
	for _, sum := range sums {
		out = append(out, machineVolumeSummaryToModel(sum))
	}
	return out
}

func machineAddressSummaryToModel(sum camplocal.IpAddressSummary) *machine.AddressModel {
	return &machine.AddressModel{
		Version: machine.IpAddressType(sum.Version),
		Address: sum.Address,
	}
}

func machineAddressSummariesToModels(sums []camplocal.IpAddressSummary) []*machine.AddressModel {
	var out []*machine.AddressModel
	for _, sum := range sums {
		out = append(out, machineAddressSummaryToModel(sum))
	}
	return out
}

func modelToSummary(m *machine.Model) camplocal.MachineSummary {
	system := systemModelToSummary(m.System)
	cpu := cpuModelToSummary(m.Cpu)
	memory := memoryModelToSummary(m.Memory)
	disks := diskModelsToSummaryList(m.Disks)
	networkInterfaces := networkInterfaceModelsToSummaryList(m.NetworkInterfaces)
	volumes := volumeModelsToSummaryList(m.Volumes)
	addresses := addressModelsToSummaryList(m.Addresses)

	return camplocal.MachineSummary{
		Identifier: m.Identifier,
		Class:      camplocal.NewOptMachineClass(camplocal.MachineClass(m.Class)),
		Status: camplocal.MachineStatusSummary{
			Status:      camplocal.MachineStatus(m.Status),
			LastChecked: float64(m.LastHeartbeat.Unix()),
		},
		Tags:              []camplocal.Tag{},
		System:            system,
		CPU:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		Volumes:           volumes,
		Addresses:         addresses,
	}
}

func systemModelToSummary(m *machine.SystemModel) camplocal.MachineSystemSummary {
	sum := camplocal.MachineSystemSummary{}
	if m.Family != nil {
		sum.Family = camplocal.NewOptString(*m.Family)
	}
	if m.KernelVersion != nil {
		sum.KernelVersion = camplocal.NewOptString(*m.KernelVersion)
	}
	if m.Os != nil {
		sum.Os = camplocal.NewOptString(*m.Os)
	}
	if m.OsVersion != nil {
		sum.OsVersion = camplocal.NewOptString(*m.OsVersion)
	}
	if m.OsPretty != nil {
		sum.OsPretty = camplocal.NewOptString(*m.OsPretty)
	}
	if m.Hostname != nil {
		sum.Hostname = camplocal.NewOptString(*m.Hostname)
	}
	return sum
}

func cpuModelToSummary(m *machine.CpuModel) camplocal.MachineCpuSummary {
	sum := camplocal.MachineCpuSummary{
		Cores:        float64(m.Cores),
		Architecture: camplocal.CpuArchitecture(m.Architecture),
	}
	if m.Model != nil {
		sum.Model = camplocal.NewOptString(*m.Model)
	}
	if m.Vendor != nil {
		sum.Vendor = camplocal.NewOptString(*m.Vendor)
	}
	return sum
}

func memoryModelToSummary(m *machine.MemoryModel) camplocal.MachineMemorySummary {
	return camplocal.MachineMemorySummary{
		Total: float64(m.Total),
	}
}

func diskModelToSummary(m *machine.DiskModel) camplocal.MachineDiskSummary {
	sum := camplocal.MachineDiskSummary{
		Device:    m.Device,
		Interface: camplocal.DiskInterface(m.Interface),
		Type:      camplocal.DiskType(m.Type),
		Serial:    camplocal.NewOptString(*m.Serial),
		Size:      float64(m.Size),
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

func diskModelsToSummaryList(ms []*machine.DiskModel) []camplocal.MachineDiskSummary {
	var out []camplocal.MachineDiskSummary
	for _, m := range ms {
		out = append(out, diskModelToSummary(m))
	}
	return out
}

func networkInterfaceModelToSummary(m *machine.NetworkInterfaceModel) camplocal.MachineNetworkInterfaceSummary {
	sum := camplocal.MachineNetworkInterfaceSummary{
		Name:       m.Name,
		Virtual:    m.Virtual,
		MacAddress: camplocal.NewOptString(*m.MacAddress),
		Vendor:     camplocal.NewOptString(*m.Vendor),
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

func networkInterfaceModelsToSummaryList(ms []*machine.NetworkInterfaceModel) []camplocal.MachineNetworkInterfaceSummary {
	var out []camplocal.MachineNetworkInterfaceSummary
	for _, m := range ms {
		out = append(out, networkInterfaceModelToSummary(m))
	}
	return out
}

func volumeModelToSummary(m *machine.VolumeModel) camplocal.MachineVolumeSummary {
	return camplocal.MachineVolumeSummary{
		Name:       m.Name,
		MountPoint: m.MountPoint,
		Total:      float64(m.Total),
	}
}

func volumeModelsToSummaryList(ms []*machine.VolumeModel) []camplocal.MachineVolumeSummary {
	var out []camplocal.MachineVolumeSummary
	for _, m := range ms {
		out = append(out, volumeModelToSummary(m))
	}
	return out
}

func addressModelToSummary(m *machine.AddressModel) camplocal.IpAddressSummary {
	return camplocal.IpAddressSummary{
		Version: camplocal.IpAddressVersion(m.Version),
		Address: m.Address,
	}
}

func addressModelsToSummaryList(ms []*machine.AddressModel) []camplocal.IpAddressSummary {
	var out []camplocal.IpAddressSummary
	for _, m := range ms {
		out = append(out, addressModelToSummary(m))
	}
	return out
}
