package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/core/domain/system"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (h *SystemHandler) GetSystem(ctx context.Context) (*system.System, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Getting system information")

	host, err := h.hostSvc.Describe(ctx)
	if err != nil {
		return nil, err
	}

	cpu, err := h.cpuSvc.Description(ctx)
	if err != nil {
		return nil, err
	}

	mem, err := h.memSvc.Description(ctx)
	if err != nil {
		return nil, err
	}

	bios, err := h.moboSvc.DescribeBios(ctx)
	if err != nil {
		return nil, err
	}

	mobo, err := h.moboSvc.DescribeMotherboard(ctx)
	if err != nil {
		return nil, err
	}

	nics, err := h.netSvc.ListNics(ctx)
	if err != nil {
		return nil, err
	}

	ips, err := h.netSvc.ListIpAddresses(ctx)
	if err != nil {
		return nil, err
	}

	disks, err := h.strSvc.ListDisks(ctx)
	if err != nil {
		return nil, err
	}

	return &system.System{
		Host:              host,
		Cpu:               cpu,
		Memory:            mem,
		Bios:              bios,
		Motherboard:       mobo,
		NetworkInterfaces: nics,
		IpAddresses:       ips,
		Disks:             disks,
	}, nil

}
