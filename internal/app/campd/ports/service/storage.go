package service

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/internal/pkg/domain/storage"
)

var (
	ErrDiskNotFound = fmt.Errorf("disk not found")
)

type Storage interface {
	ListDisks(ctx context.Context) ([]*storage.Disk, error)
	DescribeDisk(ctx context.Context, name string) (*storage.Disk, error)
}
