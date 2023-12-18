package system

import (
	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	cpuSvc  service.CPU
	hostSvc service.Host
	memSvc  service.Memory
	moboSvc service.Motherboard
	netSvc  service.Network
	strSvc  service.Storage
}

func New(cpuSvc service.CPU, hostSvc service.Host, memSvc service.Memory, moboSvc service.Motherboard, netSvc service.Network, strSvc service.Storage) *Handler {
	return &Handler{
		cpuSvc:  cpuSvc,
		hostSvc: hostSvc,
		memSvc:  memSvc,
		moboSvc: moboSvc,
		netSvc:  netSvc,
		strSvc:  strSvc,
	}
}

func (h *Handler) Mount(r chi.Router) {
	r.Get("/", h.GetSystem)
}
