package host

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/values"
	pshost "github.com/shirou/gopsutil/v3/host"
)

type Service struct {
	host *host.Host
}

func InitService(ctx context.Context) (service.Host, error) {
	info, err := pshost.Info()
	if err != nil {
		return nil, err
	}

	host := &host.Host{
		Hostname: values.ParseOptional(info.Hostname),
		HostId:   values.ParseOptional(info.HostID),
		OS: &host.OS{
			Platform: values.ParseOptional(info.OS),
			Name:     values.ParseOptional(info.Platform),
			Version:  values.ParseOptional(info.PlatformVersion),
			Kernel:   values.ParseOptional(info.KernelVersion),
			Family:   values.ParseOptional(info.PlatformFamily),
		},
	}
	return &Service{
		host: host,
	}, nil
}
