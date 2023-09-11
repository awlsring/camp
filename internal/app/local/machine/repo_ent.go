package machine

import (
	"context"
	"time"

	"github.com/awlsring/camp/models/database/local/ent"
	"github.com/awlsring/camp/models/database/local/ent/cpu"
	"github.com/awlsring/camp/models/database/local/ent/disk"
	"github.com/awlsring/camp/models/database/local/ent/ipaddress"
	"github.com/awlsring/camp/models/database/local/ent/machine"
)

type EntRepo struct {
	client *ent.Client
}

func modelToMachine(model *Model, create *ent.MachineCreate) *ent.MachineCreate {
	now := time.Now()
	return create.
		SetIdentifier(model.Identifier).
		SetClass(machine.Class(model.Class)).
		SetState(machine.StateRunning).
		SetLastHeartbeat(now).
		SetRegisteredAt(now).
		SetUpdatedAt(now)
}

func modelToCpu(model *CpuModel, create *ent.CPUCreate) *ent.CPUCreate {
	create.
		SetCores(model.Cores).
		SetArchitecture(cpu.Architecture(model.Architecture))

	if model.Model != nil {
		create.SetModel(*model.Model)
	}
	if model.Vendor != nil {
		create.SetVendor(*model.Vendor)
	}
	return create
}

func modelToMemory(model *MemoryModel, create *ent.MemoryCreate) *ent.MemoryCreate {
	return create.SetTotal(model.Total)
}

func modelToDisk(model *DiskModel, create *ent.DiskCreate) *ent.DiskCreate {
	create.
		SetDevice(model.Device).
		SetInterface(disk.Interface(model.Interface)).
		SetDiskType(disk.DiskType(model.Type)).
		SetSectorSize(model.SectorSize).
		SetSize(model.Size)

	if model.Model != nil {
		create.SetModel(*model.Model)
	}
	if model.Vendor != nil {
		create.SetVendor(*model.Vendor)
	}
	if model.Serial != nil {
		create.SetSerial(*model.Serial)
	}
	if model.SizeRaw != nil {
		create.SetSizeRaw(*model.SizeRaw)
	}
	return create
}

func modelToNetworkInterface(model *NetworkInterfaceModel, create *ent.NetworkInterfaceCreate) *ent.NetworkInterfaceCreate {
	create.
		SetName(model.Name).
		SetVirtual(model.Virtual)

	if model.MacAddress != nil {
		create.SetMacAddress(*model.MacAddress)
	}
	if model.Vendor != nil {
		create.SetVendor(*model.Vendor)
	}
	if model.Mtu != nil {
		create.SetMtu(int(*model.Mtu))
	}
	if model.Speed != nil {
		create.SetSpeed(*model.Speed)
	}
	if model.Duplex != nil {
		create.SetDuplex(*model.Duplex)
	}
	return create
}

func modelToVolume(model *VolumeModel, create *ent.VolumeCreate) *ent.VolumeCreate {
	create.
		SetName(model.Name).
		SetMountPoint(model.MountPoint).
		SetTotal(model.Total)

	if model.FileSystem != nil {
		create.SetFileSystem(*model.FileSystem)
	}
	return create
}

func modelToIpAddress(model *AddressModel, create *ent.IpAddressCreate) *ent.IpAddressCreate {
	return create.SetAddress(model.Address).SetVersion(ipaddress.Version(model.Version))
}

func getIpFromMap(m map[string]*ent.IpAddress, addresses []AddressModel) *ent.IpAddress {
	for _, address := range addresses {
		if ip, ok := m[address.Address]; ok {
			return ip
		}
	}
	return nil
}

func (e *EntRepo) CreateMachine(ctx context.Context, model *Model) error {
	machine, err := modelToMachine(model, e.client.Machine.Create()).Save(ctx)
	if err != nil {
		return err
	}

	_, err = modelToCpu(model.Cpu, e.client.CPU.Create()).SetMachine(machine).Save(ctx)
	if err != nil {
		return err
	}

	_, err = modelToMemory(model.Memory, e.client.Memory.Create()).SetMachine(machine).Save(ctx)

	for _, disk := range model.Disks {
		_, err = modelToDisk(disk, e.client.Disk.Create()).SetMachine(machine).Save(ctx)
		if err != nil {
			return err
		}
	}

	ips := map[string]*ent.IpAddress{}
	for _, address := range model.Addresses {
		a, err := modelToIpAddress(address, e.client.IpAddress.Create()).SetMachine(machine).Save(ctx)
		if err != nil {
			return err
		}
		ips[address.Address] = a
	}

	for _, networkInterface := range model.NetworkInterfaces {
		address := getIpFromMap(ips, networkInterface.IpAddresses)
		create := modelToNetworkInterface(networkInterface, e.client.NetworkInterface.Create())

		if address != nil {
			create.AddIpAddresses()
		}

		create.SetMachine(machine).Save(ctx)

		if err != nil {
			return err
		}
	}

	for _, volume := range model.Volumes {
		_, err := modelToVolume(volume, e.client.Volume.Create()).SetMachine(machine).Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
