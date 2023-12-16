package service

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/internal/pkg/domain/network"
)

var ErrNicNotFound = fmt.Errorf("nic not found")

type Network interface {
	ListNics(ctx context.Context) ([]*network.Nic, error)
	DescribeNic(ctx context.Context, name string) (*network.Nic, error)
}
