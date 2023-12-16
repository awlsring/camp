package motherboard

import "github.com/awlsring/camp/internal/pkg/values"

type Motherboard struct {
	AssetTag *string
	Product  *string
	Serial   *string
	Vendor   *string
	Version  *string
}

func NewMotherboard(a, p, s, ve, ver string) *Motherboard {
	return &Motherboard{
		AssetTag: values.ParseOptional(a),
		Product:  values.ParseOptional(p),
		Serial:   values.ParseOptional(s),
		Vendor:   values.ParseOptional(ve),
		Version:  values.ParseOptional(ver),
	}
}
