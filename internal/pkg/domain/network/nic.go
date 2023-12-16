package network

import "github.com/awlsring/camp/internal/pkg/values"

type Nic struct {
	Name        string
	Virtual     bool
	MacAddress  *MacAddress
	Speed       *string
	Duplex      *string
	PCIAddress  string
	Vendor      *string
	IpAddresses []*IpAddress
}

func NewNic(name string, macAddress *MacAddress, virtual bool, speed string, duplex string, pciAddress string) *Nic {
	return &Nic{
		Name:       name,
		MacAddress: macAddress,
		Virtual:    virtual,
		Speed:      values.ParseOptional(speed),
		Duplex:     values.ParseOptional(duplex),
		PCIAddress: pciAddress,
	}
}
