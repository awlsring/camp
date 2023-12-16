package motherboard

import "github.com/awlsring/camp/internal/pkg/values"

type Bios struct {
	Date    *string
	Vendor  *string
	Version *string
}

func NewBios(d, v, ve string) *Bios {
	return &Bios{
		Date:    values.ParseOptional(d),
		Vendor:  values.ParseOptional(v),
		Version: values.ParseOptional(ve),
	}
}
