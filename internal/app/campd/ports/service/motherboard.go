package service

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
)

type Motherboard interface {
	DescribeBios(context.Context) (*motherboard.Bios, error)
	DescribeMotherboard(context.Context) (*motherboard.Motherboard, error)
}
