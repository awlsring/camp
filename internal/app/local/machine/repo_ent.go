package machine

import (
	"context"
	"time"

	"entgo.io/ent/dialect"
	"github.com/rs/zerolog/log"

	"github.com/awlsring/camp/models/database/local/ent"
	"github.com/awlsring/camp/models/database/local/ent/cpu"
	"github.com/awlsring/camp/models/database/local/ent/disk"
	"github.com/awlsring/camp/models/database/local/ent/ipaddress"
	"github.com/awlsring/camp/models/database/local/ent/machine"
	"github.com/awlsring/camp/models/database/local/ent/memory"
	"github.com/awlsring/camp/models/database/local/ent/networkinterface"
	"github.com/awlsring/camp/models/database/local/ent/volume"

	_ "github.com/lib/pq"
)

type EntRepo struct {
	client *ent.Client
}

func NewRepo(cfg RepoConfig) Repo {
	conn := createPostgresConnectionString(cfg)
	client, err := ent.Open(dialect.Postgres, conn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to db.")
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed creating schema resources")
	}

	return &EntRepo{
		client: client,
	}
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
	if err != nil {
		return err
	}

	for _, disk := range model.Disks {
		_, err = modelToDisk(disk, e.client.Disk.Create()).SetMachine(machine).Save(ctx)
		if err != nil {
			return err
		}
	}

	for _, networkInterface := range model.NetworkInterfaces {
		n, err := modelToNetworkInterface(networkInterface, e.client.NetworkInterface.Create()).SetMachine(machine).Save(ctx)
		if err != nil {
			return err
		}
		for _, address := range networkInterface.IpAddresses {
			_, err := modelToIpAddress(address, e.client.IpAddress.Create()).SetNetworkInterface(n).Save(ctx)
			if err != nil {
				return err
			}
		}

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

func (e *EntRepo) UpdateMachine(ctx context.Context, m *Model) error {
	_, err := e.client.Machine.Update().
		Where(machine.Identifier(m.Identifier)).
		SetClass(machine.Class(m.Class)).
		SetState(machine.State(m.Status)).
		SetLastHeartbeat(m.LastHeartbeat).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = e.client.CPU.Update().
		Where(cpu.HasMachineWith(machine.Identifier(m.Identifier))).
		SetCores(m.Cpu.Cores).
		SetArchitecture(cpu.Architecture(m.Cpu.Architecture)).
		SetModel(*m.Cpu.Model).
		SetVendor(*m.Cpu.Vendor).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = e.client.Memory.Update().
		Where(memory.HasMachineWith(machine.Identifier(m.Identifier))).SetTotal(m.Memory.Total).Save(ctx)
	if err != nil {
		return err
	}

	for _, d := range m.Disks {
		_, err = e.client.Disk.Update().
			Where(disk.HasMachineWith(machine.Identifier(m.Identifier))).
			SetDevice(d.Device).
			SetModel(*d.Model).
			SetVendor(*d.Vendor).
			SetInterface(disk.Interface(d.Interface)).
			SetDiskType(disk.DiskType(d.Type)).
			SetSerial(*d.Serial).
			SetSectorSize(d.SectorSize).
			SetSize(d.Size).
			SetSizeRaw(*d.SizeRaw).
			Save(ctx)
	}

	for _, nic := range m.NetworkInterfaces {
		_, err = e.client.NetworkInterface.Update().
			Where(networkinterface.HasMachineWith(machine.Identifier(m.Identifier))).
			SetName(nic.Name).
			SetVirtual(nic.Virtual).
			SetMacAddress(*nic.MacAddress).
			SetVendor(*nic.Vendor).
			SetMtu(int(*nic.Mtu)).
			SetSpeed(*nic.Speed).
			SetDuplex(*nic.Duplex).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	for _, v := range m.Volumes {
		_, err = e.client.Volume.Update().
			Where(volume.HasMachineWith(machine.Identifier(m.Identifier))).
			SetName(v.Name).
			SetMountPoint(v.MountPoint).
			SetTotal(v.Total).
			SetFileSystem(*v.FileSystem).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return err
}

func (e *EntRepo) DeleteMachine(ctx context.Context, id string) error {
	m, err := e.client.Machine.Query().Where(machine.Identifier(id)).Only(ctx)
	if err != nil {
		return err
	}

	err = e.client.Machine.DeleteOne(m).Exec(ctx)
	if err != nil {
		return err
	}

	return err
}

func (e *EntRepo) GetMachineById(ctx context.Context, id string) (*Model, error) {
	log.Debug().Msg("Invoke DescribeMachine")
	disks := []*DiskModel{}
	networkInterfaces := []*NetworkInterfaceModel{}
	volumes := []*VolumeModel{}
	addresses := []*AddressModel{}

	m, err := e.client.Machine.Query().
		Where(machine.Identifier(id)).
		WithSystem().
		WithCPU().
		WithMemory().
		Only(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to query machine")
		return nil, err
	}

	d, err := e.client.Disk.Query().Where(disk.HasMachineWith(machine.Identifier(id))).All(ctx)
	if err != nil {
		return nil, err
	}
	for _, disk := range d {
		d, err := entDisktoModel(disk)
		if err != nil {
			return nil, err
		}
		disks = append(disks, d)
	}

	n, err := e.client.NetworkInterface.Query().Where(networkinterface.HasMachineWith(machine.Identifier(id))).All(ctx)
	if err != nil {
		return nil, err
	}
	for _, networkInterface := range n {
		a, err := e.client.IpAddress.Query().Where(ipaddress.HasNetworkInterfaceWith(networkinterface.ID(networkInterface.ID))).All(ctx)
		if err != nil {
			return nil, err
		}
		for _, address := range a {
			a, err := entIpAddresstoModel(address)
			if err != nil {
				return nil, err
			}
			addresses = append(addresses, a)
		}
		n, err := entNetworkInterfacetoModel(networkInterface, addresses)
		if err != nil {
			return nil, err
		}
		networkInterfaces = append(networkInterfaces, n)
	}

	v, err := e.client.Volume.Query().Where(volume.HasMachineWith(machine.Identifier(id))).All(ctx)
	if err != nil {
		return nil, err
	}
	for _, volume := range v {
		v, err := entVolumetoModel(volume)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, v)
	}

	return entMachineToModel(m, disks, networkInterfaces, volumes)
}

func (e *EntRepo) GetMachines(ctx context.Context, filters *GetMachinesFilters) ([]*Model, error) {
	models := []*Model{}
	m, err := e.client.Machine.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	for _, machine := range m {
		d, err := machine.QueryDisks().All(ctx)
		if err != nil {
			return nil, err
		}
		disks, err := entDisksToModels(d)
		v, err := machine.QueryVolumes().All(ctx)
		if err != nil {
			return nil, err
		}
		volumes, err := entVolumesToModels(v)
		n, err := machine.QueryNetworkInterfaces().All(ctx)
		if err != nil {
			return nil, err
		}
		networkInterfaces, err := entNetworkInterfacesToModels(ctx, n)
		if err != nil {
			return nil, err
		}
		model, err := entMachineToModel(machine, disks, networkInterfaces, volumes)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

func (e *EntRepo) AcknowledgeHeartbeat(ctx context.Context, id string) error {
	_, err := e.client.Machine.Update().
		Where(machine.Identifier(id)).
		SetLastHeartbeat(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	return err
}

func (e *EntRepo) UpdateStatus(ctx context.Context, id string, status MachineStatus) error {
	_, err := e.client.Machine.Update().
		Where(machine.Identifier(id)).
		SetState(machine.State(status)).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	return err
}

func entDisksToModels(disks []*ent.Disk) ([]*DiskModel, error) {
	models := []*DiskModel{}
	for _, disk := range disks {
		model, err := entDisktoModel(disk)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

func entDisktoModel(m *ent.Disk) (*DiskModel, error) {
	return &DiskModel{
		Device:     m.Device,
		Model:      &m.Model,
		Vendor:     &m.Vendor,
		Interface:  DiskInterface(m.Interface),
		Type:       DiskType(m.DiskType),
		Serial:     &m.Serial,
		SectorSize: m.SectorSize,
		Size:       m.Size,
		SizeRaw:    &m.SizeRaw,
	}, nil
}

func entNetworkInterfacetoModel(m *ent.NetworkInterface, addresses []*AddressModel) (*NetworkInterfaceModel, error) {
	return &NetworkInterfaceModel{
		Name:        m.Name,
		MacAddress:  &m.MacAddress,
		Vendor:      &m.Vendor,
		Mtu:         &m.Mtu,
		Speed:       &m.Speed,
		Duplex:      &m.Duplex,
		IpAddresses: addresses,
	}, nil
}

func entNetworkInterfacesToModels(ctx context.Context, networkInterfaces []*ent.NetworkInterface) ([]*NetworkInterfaceModel, error) {
	models := []*NetworkInterfaceModel{}
	for _, networkInterface := range networkInterfaces {
		addresses := []*AddressModel{}
		a, err := networkInterface.QueryIpAddresses().All(ctx)
		if err != nil {
			return nil, err
		}
		for _, address := range a {
			a, err := entIpAddresstoModel(address)
			if err != nil {
				return nil, err
			}
			addresses = append(addresses, a)
		}
		model, err := entNetworkInterfacetoModel(networkInterface, addresses)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

func entVolumesToModels(volumes []*ent.Volume) ([]*VolumeModel, error) {
	models := []*VolumeModel{}
	for _, volume := range volumes {
		model, err := entVolumetoModel(volume)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

func entVolumetoModel(m *ent.Volume) (*VolumeModel, error) {
	return &VolumeModel{
		Name:       m.Name,
		MountPoint: m.MountPoint,
		Total:      m.Total,
		FileSystem: &m.FileSystem,
	}, nil
}

func entIpAddresstoModel(m *ent.IpAddress) (*AddressModel, error) {
	return &AddressModel{
		Address: m.Address,
		Version: IpAddressType(m.Version),
	}, nil
}

func entMachineToModel(m *ent.Machine, disks []*DiskModel, nics []*NetworkInterfaceModel, vols []*VolumeModel) (*Model, error) {
	addrs := []*AddressModel{}
	for _, nic := range nics {
		addrs = append(addrs, nic.IpAddresses...)
	}

	return &Model{
		Identifier:    string(m.Identifier),
		Class:         MachineClass(m.Class),
		Status:        MachineStatus(m.State),
		LastHeartbeat: m.LastHeartbeat,
		System: &SystemModel{
			Family:        &m.Edges.System.Family,
			KernelVersion: &m.Edges.System.KernelVersion,
			Os:            &m.Edges.System.Os,
			OsVersion:     &m.Edges.System.OsVersion,
			OsPretty:      &m.Edges.System.OsPretty,
			Hostname:      &m.Edges.System.Hostname,
		},
		Cpu: &CpuModel{
			Cores:        m.Edges.CPU.Cores,
			Architecture: CpuArchitecture(m.Edges.CPU.Architecture),
			Model:        &m.Edges.CPU.Model,
			Vendor:       &m.Edges.CPU.Vendor,
		},
		Memory: &MemoryModel{
			Total: m.Edges.Memory.Total,
		},
		Disks:             disks,
		NetworkInterfaces: nics,
		Volumes:           vols,
		Addresses:         addrs,
	}, nil
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
