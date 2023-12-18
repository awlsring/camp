package system

import (
	"net/http"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/web/pages"
	"github.com/awlsring/camp/internal/app/campd/core/domain/system"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (h *Handler) GetSystem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	log.Debug().Msg("Getting system information")

	log.Debug().Msg("Getting host information")
	host, err := h.hostSvc.Describe(ctx)
	if err != nil {
		internalError(w, err)
		return
	}

	log.Debug().Msg("Getting cpu information")
	cpu, err := h.cpuSvc.Description(ctx)
	if err != nil {
		internalError(w, err)
		return
	}

	log.Debug().Msg("Getting memory information")
	mem, err := h.memSvc.Description(ctx)
	if err != nil {
		internalError(w, err)
		return
	}

	log.Debug().Msg("Getting bios information")
	bios, err := h.moboSvc.DescribeBios(ctx)
	if err != nil {
		internalError(w, err)
		return
	}

	log.Debug().Msg("Getting motherboard information")
	mobo, err := h.moboSvc.DescribeMotherboard(ctx)
	if err != nil {
		internalError(w, err)
		return
	}

	log.Debug().Msg("Getting network information")
	nics, err := h.netSvc.ListNics(ctx)
	if err != nil {
		internalError(w, err)
		return
	}

	log.Debug().Msg("Getting ip address information")
	ips, err := h.netSvc.ListIpAddresses(ctx)
	if err != nil {
		internalError(w, err)
		return
	}

	log.Debug().Msg("Getting disk information")
	disks, err := h.strSvc.ListDisks(ctx)
	if err != nil {
		internalError(w, err)
		return
	}

	log.Debug().Msg("Rendering summary page")
	sum := &system.System{
		Host:              host,
		Cpu:               cpu,
		Memory:            mem,
		Bios:              bios,
		Motherboard:       mobo,
		NetworkInterfaces: nics,
		IpAddresses:       ips,
		Disks:             disks,
	}

	pages.Summary(sum).Render(ctx, w)
}
